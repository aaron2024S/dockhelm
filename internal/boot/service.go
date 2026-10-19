package boot

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/aaron2024s/dockhelm/internal/dockerx"
	"github.com/aaron2024s/dockhelm/internal/notify"
	"github.com/aaron2024s/dockhelm/internal/store"
)

// 触发方式。只影响「要不要等 Docker 就绪」以及要不要看 enabled 总开关。
const (
	TriggerBoot   = "boot"   // 进程启动时自动跑
	TriggerManual = "manual" // 页面上的「立即执行一次」
)

// ErrBusy 已经有一轮在执行（方案：互斥，不并发）。
var ErrBusy = errors.New("顺序启动正在进行中")

// 执行阶段，页面据此显示不同的进度文案。
const (
	PhaseIdle         = "idle"
	PhaseWaitingDockr = "waiting-docker"
	PhaseInitialDelay = "initial-delay"
	PhaseStarting     = "starting"
)

// Status 当前进度（页面 2s 轮询）。
type Status struct {
	Running    bool   `json:"running"`
	Trigger    string `json:"trigger"`
	Phase      string `json:"phase"`
	Index      int    `json:"index"` // 第几个（按统计口径，从 1 开始）
	Total      int    `json:"total"` // 本次要处理几个（统计口径）
	Current    string `json:"current"`
	Success    int    `json:"success"`
	Skipped    int    `json:"skipped"`
	Failed     int    `json:"failed"`
	StartedAt  string `json:"startedAt"`
	FinishedAt string `json:"finishedAt"`
}

// Service 顺序启动的常驻对象。
type Service struct {
	dc *dockerx.Client
	st *store.Store
	nt *notify.Manager
	// selfName 返回 Dockhelm 自身容器名（页面要提示它当前的策略）。
	selfName func() string

	mu      sync.Mutex
	running bool
	status  Status
}

// New 创建服务。nt 可以为 nil（单测场景）。
func New(dc *dockerx.Client, st *store.Store, nt *notify.Manager, selfName func() string) *Service {
	if selfName == nil {
		selfName = func() string { return "" }
	}
	return &Service{dc: dc, st: st, nt: nt, selfName: selfName, status: Status{Phase: PhaseIdle}}
}

// LastResult 上一次执行的结构化回执（页面顶部那一行）。只是展示，不参与决策。
type LastResult struct {
	TS      string `json:"ts"`
	Status  string `json:"status"`
	Total   int    `json:"total"`
	Success int    `json:"success"`
	Skipped int    `json:"skipped"`
	Failed  int    `json:"failed"`
	Message string `json:"message"`
}

// LastResult 读上一次执行结果；没执行过返回 false。
func (s *Service) LastResult() (LastResult, bool) {
	var out LastResult
	if s.st == nil {
		return out, false
	}
	s.st.GetJSON(lastKey, &out)
	if out.TS == "" {
		return out, false
	}
	return out, true
}

// ---------- 清单读写 ----------

// LoadPlan 读清单；没有存过就用出厂默认（开关全开、延时 3 秒）。
func (s *Service) LoadPlan() Plan {
	p := DefaultPlan()
	if s.st != nil {
		s.st.GetJSON(planKey, &p)
	}
	if p.Items == nil {
		p.Items = []Item{}
	}
	return p.Normalize()
}

// SavePlan 落盘（先 Normalize 再写，保证存进去的永远是合法值）。
func (s *Service) SavePlan(p Plan) error {
	if s.st == nil {
		return nil
	}
	return s.st.SetJSON(planKey, p.Normalize())
}

// PlanMerged 读清单并与当前容器对齐（剔已删除、追加新容器）。
// 对齐结果与库里不同才回写 —— 避免每次打开页面都写一次盘。
func (s *Service) PlanMerged(infos []ContainerInfo) Plan {
	p := s.LoadPlan()
	merged := p.Merge(infos)
	if !sameItems(p.Items, merged.Items) {
		_ = s.SavePlan(merged)
	}
	return merged
}

func sameItems(a, b []Item) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Snapshot 取一次容器快照（状态 + 重启策略 + compose 项目）。
//
// 重启策略只能从 inspect 拿（/containers/json 不带它），50 多个容器串行 inspect
// 会让页面等上半秒，所以这里开 8 个并发。
func (s *Service) Snapshot(ctx context.Context) ([]ContainerInfo, error) {
	list, err := s.dc.ListContainers(ctx)
	if err != nil {
		return nil, fmt.Errorf("读取容器列表失败：%w", err)
	}
	out := make([]ContainerInfo, len(list))
	sem := make(chan struct{}, 8)
	var wg sync.WaitGroup
	for i, c := range list {
		i, c := i, c
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()
			info := ContainerInfo{Name: c.Name(), ID: c.ID, State: c.State, Project: c.ComposeProject()}
			if f, err := s.dc.Facts(ctx, c.ID); err == nil {
				info.Policy = f.Policy
				if f.State != "" {
					info.State = f.State
				}
				if f.Project != "" {
					info.Project = f.Project
				}
			}
			out[i] = info
		}()
	}
	wg.Wait()
	return out, nil
}

// SelfPolicy 读 Dockhelm 自身容器的重启策略（读不到返回空串）。
func (s *Service) SelfPolicy(ctx context.Context) string {
	name := ""
	if s.selfName != nil {
		name = s.selfName()
	}
	if name == "" {
		return ""
	}
	if f, err := s.dc.Facts(ctx, name); err == nil {
		return f.Policy
	}
	return ""
}

// SelfName 自身容器名（给页面显示用）。
func (s *Service) SelfName() string {
	if s.selfName == nil {
		return ""
	}
	return s.selfName()
}

// ---------- 状态 ----------

// Status 返回当前进度快照。
func (s *Service) Status() Status {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.status
}

func (s *Service) setStatus(fn func(*Status)) {
	s.mu.Lock()
	fn(&s.status)
	s.mu.Unlock()
}

// begin 抢占执行权；已有任务在跑时返回 false。
func (s *Service) begin(trigger string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.running {
		return false
	}
	s.running = true
	s.status = Status{
		Running:   true,
		Trigger:   trigger,
		Phase:     PhaseWaitingDockr,
		StartedAt: time.Now().UTC().Format(time.RFC3339),
	}
	return true
}

func (s *Service) end() {
	s.mu.Lock()
	s.running = false
	s.status.Running = false
	s.status.Phase = PhaseIdle
	s.status.Current = ""
	s.status.FinishedAt = time.Now().UTC().Format(time.RFC3339)
	s.mu.Unlock()
}

// ---------- 启动入口 ----------

// Startup 进程启动时调用：等 Docker 就绪 → 等初始等待 → 跑一轮。
// 这里是 goroutine 里跑的，不阻塞 HTTP 服务。
func (s *Service) Startup(ctx context.Context) {
	p := s.LoadPlan()
	if !p.Enabled {
		log.Printf("顺序启动已关闭（boot.plan.enabled=false），本次跳过")
		return
	}
	if err := s.Run(ctx, TriggerBoot, false); err != nil && !errors.Is(err, ErrBusy) && ctx.Err() == nil {
		log.Printf("顺序启动未完成：%v", err)
	}
}

// Run 同步跑一轮（Startup 与测试用）。
func (s *Service) Run(ctx context.Context, trigger string, skipInitialDelay bool) error {
	if !s.begin(trigger) {
		return ErrBusy
	}
	defer s.end()
	return s.run(ctx, trigger, skipInitialDelay)
}

// Start 异步跑一轮（页面上的「立即执行一次」）。返回 ErrBusy 表示已有一轮在跑。
func (s *Service) Start(trigger string, skipInitialDelay bool) error {
	if !s.begin(trigger) {
		return ErrBusy
	}
	go func() {
		defer s.end()
		// 不挂请求的 ctx：请求秒回，这一轮可能还要跑好几分钟
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Hour)
		defer cancel()
		if err := s.run(ctx, trigger, skipInitialDelay); err != nil && ctx.Err() == nil {
			log.Printf("顺序启动未完成：%v", err)
		}
	}()
	return nil
}

// ---------- 引擎 ----------

func (s *Service) run(ctx context.Context, trigger string, skipInitialDelay bool) error {
	p := s.LoadPlan()
	if trigger == TriggerBoot && !p.Enabled {
		return nil
	}

	// 开机场景 Docker 可能还没起好（daemon 比容器晚几十秒很常见）
	if trigger == TriggerBoot {
		if !s.waitDocker(ctx) {
			log.Printf("顺序启动放弃：等了 5 分钟 Docker 仍未就绪")
			return errors.New("Docker 守护进程未就绪")
		}
	}

	if !skipInitialDelay && p.InitialDelaySec > 0 {
		s.setStatus(func(st *Status) { st.Phase = PhaseInitialDelay })
		if err := sleepCtx(ctx, time.Duration(p.InitialDelaySec)*time.Second); err != nil {
			return err
		}
	}

	infos, err := s.Snapshot(ctx)
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		msg := "顺序启动失败：" + err.Error()
		if s.st != nil {
			s.st.AddRunLog("boot", "", "failed", msg, err.Error())
		}
		if s.nt != nil {
			s.nt.Emit("boot_failed", map[string]string{"container": "", "message": err.Error()})
		}
		return err
	}

	p = p.Merge(infos).Normalize()
	steps := Resolve(p, infos)

	total := 0
	for _, st := range steps {
		if st.Counted() {
			total++
		}
	}
	s.setStatus(func(st *Status) { st.Total = total })

	timeout := time.Duration(p.WaitTimeoutSec) * time.Second
	success, skipped := 0, 0
	failedNames := []string{}
	details := []string{}
	n := 0

	for _, st := range steps {
		if !st.Counted() {
			continue
		}
		if err := ctx.Err(); err != nil {
			return err
		}
		n++
		step := st
		s.setStatus(func(s2 *Status) {
			s2.Index = n
			s2.Current = step.Name
			s2.Phase = PhaseStarting
		})

		switch step.Kind {
		case StepRunning:
			skipped++
		case StepMissing:
			failedNames = append(failedNames, step.Name)
			details = append(details, step.Name+"：容器不存在（可能已被删除或改名）")
		case StepStart:
			if err := s.startOne(ctx, step.Name, timeout); err != nil {
				if ctx.Err() != nil {
					return ctx.Err()
				}
				failedNames = append(failedNames, step.Name)
				details = append(details, step.Name+"："+err.Error())
				log.Printf("顺序启动：%s 启动失败（%v），继续下一个", step.Name, err)
			} else {
				success++
				// 「启动成功后等 N 秒再起下一个」—— 方案 §2.2 定死的语义
				if step.DelaySec > 0 {
					if err := sleepCtx(ctx, time.Duration(step.DelaySec)*time.Second); err != nil {
						return err
					}
				}
			}
		}

		su, sk, fa := success, skipped, len(failedNames)
		s.setStatus(func(s2 *Status) {
			s2.Success = su
			s2.Skipped = sk
			s2.Failed = fa
		})
	}

	msg := Summary(total, success, skipped, failedNames)
	detail := strings.Join(details, "\n")
	status := "success"
	if len(failedNames) > 0 {
		status = "failed"
	}
	if s.st != nil {
		s.st.AddRunLog("boot", "", status, msg, detail)
		_ = s.st.SetJSON(lastKey, LastResult{
			TS:      time.Now().UTC().Format(time.RFC3339),
			Status:  status,
			Total:   total,
			Success: success,
			Skipped: skipped,
			Failed:  len(failedNames),
			Message: msg,
		})
	}
	s.notify(failedNames, details, msg, total)
	log.Printf("%s", msg)
	return nil
}

// notify 有失败发紧急推送并点名原因；全部正常发一条普通回执（默认不订阅）。
func (s *Service) notify(failedNames, details []string, brief string, total int) {
	if s.nt == nil {
		return
	}
	if len(failedNames) > 0 {
		s.nt.Emit("boot_failed", map[string]string{
			"container": fmt.Sprintf("%d 个容器", len(failedNames)),
			"message":   strings.Join(details, "\n"),
		})
		return
	}
	s.nt.Emit("boot_done", map[string]string{
		"container": fmt.Sprintf("%d 个容器", total),
		"message":   brief,
	})
}

// startOne 启动一个容器并等它进入 running（上限 timeout）。
//
// 早退判定：启动后如果看到 exited / dead，直接记失败 —— 它不会自己好起来
// （on-failure 会先进入 restarting，不是 exited）。这比干等 60 秒更有用。
func (s *Service) startOne(ctx context.Context, name string, timeout time.Duration) error {
	f, err := s.dc.Facts(ctx, name)
	if err != nil {
		return fmt.Errorf("读取容器状态失败：%s", dockerErrText(err))
	}
	if f.State == "running" {
		return nil
	}
	if err := s.dc.ContainerAction(ctx, f.ID, "start", nil); err != nil {
		return fmt.Errorf("启动命令被拒绝：%s", dockerErrText(err))
	}

	deadline := time.Now().Add(timeout)
	t := time.NewTicker(time.Second)
	defer t.Stop()
	for {
		// 先查一次再等：容器常常在 start 返回时就已经 running 了，
		// 白等这 1 秒会让 50 个容器的总耗时凭空多出一分钟。
		cur, err := s.dc.Facts(ctx, name)
		if err == nil {
			switch cur.State {
			case "running":
				return nil
			case "exited", "dead":
				return errors.New("启动后立即退出")
			}
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("%d 秒内没有进入运行状态", int(timeout.Seconds()))
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
		}
	}
}

// waitDocker Ping 重试，最多 5 分钟。
func (s *Service) waitDocker(ctx context.Context) bool {
	deadline := time.Now().Add(5 * time.Minute)
	for {
		if err := s.dc.Ping(ctx); err == nil {
			return true
		}
		if time.Now().After(deadline) || ctx.Err() != nil {
			return false
		}
		if err := sleepCtx(ctx, 2*time.Second); err != nil {
			return false
		}
	}
}

func sleepCtx(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// dockerErrText 只留守护进程给的那句话，去掉 "docker api 500:" 前缀。
func dockerErrText(err error) string {
	var se *dockerx.StatusError
	if errors.As(err, &se) {
		return strings.TrimSpace(se.Message)
	}
	return err.Error()
}
