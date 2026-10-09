// Package updater 是 Dockhelm 的核心：镜像更新。
//
// 它与 dockerCopilot 最本质的区别只有一条，但决定了「没更新的容器会不会被停掉」：
//
//	dockerCopilot：pull → stop → rename → create → start   （无条件停容器）
//	Dockhelm     ：pull → 比对镜像 ID → 没变就【直接结束，容器一个字节都不碰】
//	                                → 变了才 stop → rename → create → start
//
// 第二条同样重要：**检测（有没有新版本）与拉取永远同源**。检测侧走的是守护进程的
// /distribution 接口（与 docker pull 同一套仓库端点解析），异常时一律标记为「未知」，
// 绝不退化成「有新版本」—— 这正是 dockerCopilot 永久误报的成因。
package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aaron2024s/dockhelm/internal/bus"
	"github.com/aaron2024s/dockhelm/internal/dockerx"
)

// Notifier 通知钩子（用接口避免与 notify 包循环依赖）。
type Notifier interface {
	// Emit 投递一个事件；实现必须保证非阻塞、失败只记日志。
	Emit(event string, vars map[string]string)
}

// noopNotifier 空实现。
type noopNotifier struct{}

func (noopNotifier) Emit(string, map[string]string) {}

// Options 更新引擎参数。
type Options struct {
	// Concurrency 批量更新时的并发度（默认 2）。注册表限流时调小。
	Concurrency int
	// HealthTimeout 重建后等待健康的超时（默认 90s）。
	HealthTimeout time.Duration
	// StableSeconds 没有 healthcheck 时，启动后需保持运行多久才算成功（默认 6s）。
	StableSeconds int
	// KeepBackups 每个容器保留多少个 __bak_ 快照容器（默认 0 = 成功即删）。
	KeepBackups int
	// BackupDir 配置快照目录。
	BackupDir string
}

// DefaultOptions 返回默认参数。
func DefaultOptions(backupDir string) Options {
	return Options{
		Concurrency:   2,
		HealthTimeout: 90 * time.Second,
		StableSeconds: 6,
		KeepBackups:   0,
		BackupDir:     backupDir,
	}
}

// Updater 更新引擎。
type Updater struct {
	dc       *dockerx.Client
	bus      *bus.Bus
	notify   Notifier
	opts     Options
	selfName string
	selfID   string

	mu     sync.Mutex
	locks  map[string]*sync.Mutex // 逐容器串行锁
	onLog  func(kind, ref, status, message, detail string)
	nowStr func() string
	// mirrorFn 返回用户配置的「拉取加速源」；空字符串表示完全交给守护进程。
	mirrorFn func() string
}

// New 创建更新引擎。
func New(dc *dockerx.Client, b *bus.Bus, opts Options) *Updater {
	return &Updater{
		dc:     dc,
		bus:    b,
		notify: noopNotifier{},
		opts:   opts,
		locks:  map[string]*sync.Mutex{},
		nowStr: func() string { return time.Now().Format("2006-01-02 15:04:05") },
	}
}

// SetNotifier 注入通知器。
func (u *Updater) SetNotifier(n Notifier) {
	if n == nil {
		n = noopNotifier{}
	}
	u.notify = n
}

// SetLogSink 注入运行记录写入器。
func (u *Updater) SetLogSink(fn func(kind, ref, status, message, detail string)) {
	u.onLog = fn
}

// SetSelf 告知引擎「哪个容器是 Dockhelm 自己」，永不更新。
func (u *Updater) SetSelf(name, id string) {
	u.selfName = name
	u.selfID = id
}

// SetMirror 注入「拉取加速源」解析函数。
func (u *Updater) SetMirror(fn func() string) { u.mirrorFn = fn }

// pullImage 拉取镜像。
//
// 若用户在设置里指定了加速源、且镜像来自 Docker Hub、且引用里没写死 digest，
// 就用「<加速站>/<仓库>:<标签>」去拉，拉完再打回原始标签 —— 这样即使用户
// 没法改 daemon.json，也能让 Dockhelm 走自己选的加速站，而 compose 与其它工具
// 仍然按原来的名字找得到镜像。
//
// 无论走哪条路，**最终判定「有没有变化」用的都是内容寻址的镜像 ID**，
// 所以检测与拉取永远同源，不会出现 dockerCopilot 那种永久误报。
func (u *Updater) pullImage(ctx context.Context, imageRef string, onEvent func(dockerx.PullEvent)) (*dockerx.PullResult, error) {
	mirror := ""
	if u.mirrorFn != nil {
		mirror = strings.TrimSpace(u.mirrorFn())
	}
	ref := dockerx.ParseRef(imageRef)
	if mirror == "" || ref.Registry != "docker.io" || ref.Digest != "" {
		return u.dc.Pull(ctx, imageRef, onEvent)
	}
	host := mirrorHost(mirror)
	if host == "" {
		return u.dc.Pull(ctx, imageRef, onEvent)
	}
	target := fmt.Sprintf("%s/%s:%s", host, ref.LocalName, ref.Tag)
	res, err := u.dc.Pull(ctx, target, onEvent)
	if err != nil {
		return res, err
	}
	if res.ImageID != "" {
		if terr := u.dc.TagImage(ctx, target, ref.LocalName, ref.Tag); terr != nil {
			return res, fmt.Errorf("镜像已拉取，但回写原始标签 %s 失败：%w", imageRef, terr)
		}
	}
	res.Ref = imageRef
	return res, nil
}

// mirrorHost 去掉协议前缀与尾部斜杠，得到可以直接拼进镜像名的域名。
func mirrorHost(mirror string) string {
	m := strings.TrimSpace(strings.TrimSuffix(mirror, "/"))
	m = strings.TrimPrefix(m, "https://")
	m = strings.TrimPrefix(m, "http://")
	return m
}

// SelfName 返回自身容器名。
func (u *Updater) SelfName() string { return u.selfName }

func (u *Updater) log(kind, ref, status, message, detail string) {
	if u.onLog != nil {
		u.onLog(kind, ref, status, message, detail)
	}
}

// lockFor 取某容器的串行锁（同一容器不会被两个任务同时更新）。
func (u *Updater) lockFor(name string) *sync.Mutex {
	u.mu.Lock()
	defer u.mu.Unlock()
	m, ok := u.locks[name]
	if !ok {
		m = &sync.Mutex{}
		u.locks[name] = m
	}
	return m
}

// IsSelf 判断是否是 Dockhelm 自身。
func (u *Updater) IsSelf(nameOrID string) bool {
	if nameOrID == "" {
		return false
	}
	if u.selfName != "" && (nameOrID == u.selfName) {
		return true
	}
	if u.selfID != "" && strings.HasPrefix(u.selfID, nameOrID) {
		return true
	}
	if u.selfName != "" && strings.HasPrefix(u.selfName, strings.TrimPrefix(nameOrID, "/")) {
		return true
	}
	return false
}

// ---------- 检测 ----------

// CheckStatus 检测结论。
type CheckStatus string

const (
	StatusUpToDate        CheckStatus = "up_to_date"        // 已是最新
	StatusUpdateAvailable CheckStatus = "update_available"  // 有新版本
	StatusUnknown         CheckStatus = "unknown"           // 无法判定（网络/认证/接口不可用）
	StatusNoUpstream      CheckStatus = "no_upstream"       // 本地构建的镜像，无远端可比
)

// CheckResult 单容器检测结果。
type CheckResult struct {
	Container   string      `json:"container"`
	ID          string      `json:"id"`
	Image       string      `json:"image"`
	LocalDigest string      `json:"localDigest"`
	RemoteDigest string     `json:"remoteDigest"`
	Status      CheckStatus `json:"status"`
	Reason      string      `json:"reason"`
	CheckedAt   string      `json:"checkedAt"`
}

// Check 检测单个容器是否有新版本（只读，不拉取、不动容器）。
func (u *Updater) Check(ctx context.Context, nameOrID string) (*CheckResult, error) {
	insp, err := u.dc.Inspect(ctx, nameOrID)
	if err != nil {
		return nil, err
	}
	name := inspectName(insp)
	res := &CheckResult{
		Container: name,
		ID:        inspectID(insp),
		Image:     inspectImageRef(insp),
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
	}

	if res.Image == "" {
		res.Status = StatusNoUpstream
		res.Reason = "容器没有记录镜像引用"
		return res, nil
	}
	ref := dockerx.ParseRef(res.Image)

	// 本地 digest 集合
	local := map[string]bool{}
	for _, d := range containerRepoDigests(insp) {
		if d != "" {
			local[d] = true
			if res.LocalDigest == "" {
				res.LocalDigest = d
			}
		}
	}
	if len(local) == 0 {
		// 镜像 inspect 一次拿到 RepoDigests（有的容器 inspect 里带了 Image 的 RepoDigests）
		if imgInsp, err := u.dc.ImageInspect(ctx, ref.String()); err == nil {
			for _, d := range dockerx.LocalDigests(imgInsp) {
				local[d] = true
				if res.LocalDigest == "" {
					res.LocalDigest = d
				}
			}
		}
	}
	if len(local) == 0 {
		res.Status = StatusNoUpstream
		res.Reason = "该镜像没有仓库摘要（本地构建或已丢失 tag），无法与远端比较"
		return res, nil
	}

	remote, err := u.dc.DistributionInspect(ctx, ref.String())
	if err != nil {
		// 🚨 关键：检测失败绝不能当成「有更新」。
		res.Status = StatusUnknown
		res.Reason = "无法向仓库查询最新摘要：" + err.Error()
		return res, nil
	}
	res.RemoteDigest = remote
	if remote == "" {
		res.Status = StatusUnknown
		res.Reason = "仓库未返回摘要"
		return res, nil
	}
	if local[remote] {
		res.Status = StatusUpToDate
		res.Reason = "本地摘要与仓库一致"
		return res, nil
	}
	res.Status = StatusUpdateAvailable
	res.Reason = "仓库摘要与本地不一致"
	return res, nil
}

// DeepCheck 权威检测：真的拉一次镜像再比对镜像 ID。
//
// 因为「拉取」走的必然是守护进程（含 registry-mirrors），这是唯一 100% 同源的判定。
// 镜像已最新时，守护进程只会下载 manifest 与 config（几 KB），不会下载层。
func (u *Updater) DeepCheck(ctx context.Context, nameOrID string) (*CheckResult, error) {
	lk := u.lockFor(nameOrID)
	lk.Lock()
	defer lk.Unlock()

	insp, err := u.dc.Inspect(ctx, nameOrID)
	if err != nil {
		return nil, err
	}
	name := inspectName(insp)
	res := &CheckResult{
		Container: name,
		ID:        inspectID(insp),
		Image:     inspectImageRef(insp),
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
	}
	oldID, _ := insp["Image"].(string)
	res.LocalDigest = dockerx.ShortID(oldID)
	if res.Image == "" {
		res.Status = StatusNoUpstream
		res.Reason = "容器没有记录镜像引用"
		return res, nil
	}

	pr, err := u.pullImage(ctx, res.Image, nil)
	if err != nil {
		res.Status = StatusUnknown
		res.Reason = "拉取失败：" + err.Error()
		return res, nil
	}
	res.RemoteDigest = pr.Digest
	if pr.ImageID == "" {
		res.Status = StatusUnknown
		res.Reason = "拉取后无法读取镜像 ID"
		return res, nil
	}
	if pr.ImageID == oldID || pr.UpToDate {
		res.Status = StatusUpToDate
		res.Reason = "拉取后镜像 ID 未变化"
		return res, nil
	}
	res.Status = StatusUpdateAvailable
	res.Reason = "拉取到新镜像 " + dockerx.ShortID(pr.ImageID)
	return res, nil
}

// CheckAll 并发检测全部容器（跳过自身与被排除的）。
func (u *Updater) CheckAll(ctx context.Context, exclude map[string]bool, deep bool) []CheckResult {
	list, err := u.dc.ListContainers(ctx)
	if err != nil {
		return []CheckResult{}
	}
	type job struct{ id, name string }
	jobs := []job{}
	for _, c := range list {
		n := c.Name()
		if u.IsSelf(n) || exclude[n] {
			continue
		}
		jobs = append(jobs, job{c.ID, n})
	}

	sem := make(chan struct{}, max(1, u.opts.Concurrency*2))
	var wg sync.WaitGroup
	var mu sync.Mutex
	out := []CheckResult{}
	for _, j := range jobs {
		wg.Add(1)
		go func(j job) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			sub, cancel := context.WithTimeout(ctx, 60*time.Second)
			defer cancel()
			var r *CheckResult
			var err error
			if deep {
				r, err = u.DeepCheck(sub, j.id)
			} else {
				r, err = u.Check(sub, j.id)
			}
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				out = append(out, CheckResult{
					Container: j.name, ID: j.id, Status: StatusUnknown,
					Reason: err.Error(), CheckedAt: time.Now().UTC().Format(time.RFC3339),
				})
				return
			}
			out = append(out, *r)
		}(j)
	}
	wg.Wait()
	sort.Slice(out, func(i, j int) bool {
		rank := map[CheckStatus]int{StatusUpdateAvailable: 0, StatusUnknown: 1, StatusUpToDate: 2, StatusNoUpstream: 3}
		if rank[out[i].Status] != rank[out[j].Status] {
			return rank[out[i].Status] < rank[out[j].Status]
		}
		return out[i].Container < out[j].Container
	})
	return out
}

// ---------- 更新 ----------

// ResultStatus 更新动作的结果状态。
type ResultStatus string

const (
	ResultUpToDate ResultStatus = "up_to_date" // 拉取后镜像未变，已跳过（容器未被动过）
	ResultUpdated  ResultStatus = "updated"    // 真的更新了
	ResultFailed   ResultStatus = "failed"     // 失败且回滚成功
	ResultBroken   ResultStatus = "broken"     // 失败且回滚也失败（需人工介入）
	ResultSkipped  ResultStatus = "skipped"    // 被排除/自身
	ResultPulling  ResultStatus = "pulling"
	ResultStopped  ResultStatus = "stopped"
	ResultCreated  ResultStatus = "created"
)

// Result 单容器更新结果。
type Result struct {
	Container  string       `json:"container"`
	Image      string       `json:"image"`
	OldImageID string       `json:"oldImageId"`
	NewImageID string       `json:"newImageId"`
	Status     ResultStatus `json:"status"`
	Message    string       `json:"message"`
	RolledBack bool         `json:"rolledBack"`
	Duration   int64        `json:"durationMs"`
	Steps      []string     `json:"steps"`
}

func (u *Updater) step(res *Result, format string, a ...any) {
	line := fmt.Sprintf(format, a...)
	res.Steps = append(res.Steps, line)
	u.bus.Publish("update", "step", "running", map[string]any{
		"container": res.Container,
		"message":   line,
	})
}

func (u *Updater) status(res *Result, st ResultStatus, msg string) {
	res.Status = st
	res.Message = msg
	u.bus.Publish("update", "container_status", string(st), map[string]any{
		"container": res.Container,
		"image":     res.Image,
		"status":    string(st),
		"message":   msg,
		"oldImage":  dockerx.ShortID(res.OldImageID),
		"newImage":  dockerx.ShortID(res.NewImageID),
	})
}

// Update 更新单个容器。
//
// 语义：**镜像没变就完全不动容器**。这是对 dockerCopilot 那个
// 「一键更新 10 台，没更新的也被停掉」缺陷的直接修复。
func (u *Updater) Update(ctx context.Context, nameOrID string, force bool) *Result {
	started := time.Now()
	res := &Result{Container: nameOrID}

	lk := u.lockFor(nameOrID)
	lk.Lock()
	defer lk.Unlock()

	insp, err := u.dc.Inspect(ctx, nameOrID)
	if err != nil {
		u.status(res, ResultFailed, "读取容器信息失败："+err.Error())
		u.log("update", nameOrID, "failed", res.Message, "")
		return res
	}
	name := inspectName(insp)
	res.Container = name
	res.Image = inspectImageRef(insp)
	res.OldImageID, _ = insp["Image"].(string)

	if u.IsSelf(name) {
		u.status(res, ResultSkipped, "这是 Dockhelm 自身容器，永不自动更新")
		return res
	}
	if res.Image == "" {
		u.status(res, ResultSkipped, "容器没有镜像引用，跳过")
		return res
	}

	wasRunning := inspectRunning(insp)
	u.step(res, "开始处理 %s（镜像 %s，旧 ID %s）", name, res.Image, dockerx.ShortID(res.OldImageID))

	// ---- 1. 拉取（走守护进程，registry-mirrors 自动生效）----
	u.status(res, ResultPulling, "正在拉取镜像…")
	pr, err := u.pullImage(ctx, res.Image, func(ev dockerx.PullEvent) {
		if ev.Status != "" {
			u.bus.Publish("update", "pull_progress", "running", map[string]any{
				"container": name,
				"status":    ev.Status,
				"progress":  ev.Progress,
			})
		}
	})
	if err != nil {
		u.status(res, ResultFailed, "拉取失败："+err.Error())
		u.log("update", name, "failed", res.Message, strings.Join(res.Steps, "\n"))
		u.notify.Emit("update_failed", map[string]string{
			"container": name, "image": res.Image, "result": "拉取失败", "message": err.Error(),
		})
		res.Duration = time.Since(started).Milliseconds()
		return res
	}
	res.NewImageID = pr.ImageID
	u.step(res, "拉取完成，新 ID %s（digest %s）", dockerx.ShortID(pr.ImageID), shortDigest(pr.Digest))

	// ---- 2. 核心判定：镜像没变 ⇒ 到此为止 ----
	if res.NewImageID != "" && res.NewImageID == res.OldImageID && !force {
		u.status(res, ResultUpToDate, "镜像已是最新，容器保持原样（未停止、未重建）")
		u.step(res, "镜像 ID 未变化（%s），跳过重建", dockerx.ShortID(res.OldImageID))
		u.log("update", name, "up_to_date", res.Message, strings.Join(res.Steps, "\n"))
		res.Duration = time.Since(started).Milliseconds()
		return res
	}
	if res.NewImageID == "" {
		u.status(res, ResultFailed, "拉取后无法读取镜像 ID，拒绝在不确定状态下重建容器")
		u.log("update", name, "failed", res.Message, strings.Join(res.Steps, "\n"))
		res.Duration = time.Since(started).Milliseconds()
		return res
	}
	if pr.UpToDate && !force {
		u.status(res, ResultUpToDate, "守护进程报告镜像已是最新，容器保持原样")
		u.log("update", name, "up_to_date", res.Message, strings.Join(res.Steps, "\n"))
		res.Duration = time.Since(started).Milliseconds()
		return res
	}
	u.step(res, "检测到镜像变化，开始重建容器")

	// ---- 3. 写配置快照（回滚与还原的依据）----
	snapPath, snapErr := u.saveSnapshot(name, insp)
	if snapErr != nil {
		u.step(res, "⚠ 配置快照写入失败：%v（继续，但不写快照不影响回滚）", snapErr)
	} else {
		u.step(res, "已写入配置快照 %s", filepath.Base(snapPath))
	}

	// ---- 4~8. 停止 / 改名 / 重建 / 启动 / 自检 ----
	newID, failMsg, rb := u.recreate(ctx, insp, name, res, wasRunning)
	if failMsg != "" {
		res.RolledBack = rb
		if rb {
			u.status(res, ResultFailed, failMsg+"（已成功回滚到旧容器）")
		} else {
			u.status(res, ResultBroken, failMsg+"（回滚失败，请手工处理）")
		}
		u.log("update", name, "failed", res.Message, strings.Join(res.Steps, "\n"))
		u.notify.Emit("update_failed", map[string]string{
			"container": name, "image": res.Image,
			"result": map[bool]string{true: "已回滚", false: "回滚失败"}[rb],
			"message": failMsg,
		})
		res.Duration = time.Since(started).Milliseconds()
		return res
	}
	res.NewImageID = newID
	u.status(res, ResultUpdated, "更新成功"+map[bool]string{true: "", false: "（原容器为停止状态，保持停止）"}[wasRunning])
	u.log("update", name, "success", "已更新到 "+dockerx.ShortID(newID), strings.Join(res.Steps, "\n"))
	u.notify.Emit("update_success", map[string]string{
		"container": name, "image": res.Image,
		"result": "更新成功", "message": "新镜像 " + dockerx.ShortID(newID),
	})
	res.Duration = time.Since(started).Milliseconds()
	return res
}

// recreate 执行「停旧 → 改名保留 → 建新 → 启动 → 自检」，失败自动回滚。
// 返回 (新容器ID, 失败信息, 是否已回滚)。失败信息为空表示成功。
func (u *Updater) recreate(ctx context.Context, insp map[string]any, name string, res *Result, wasRunning bool) (string, string, bool) {
	oldID := inspectID(insp)
	bakName := fmt.Sprintf("%s__bak_%s", truncName(name, 40), time.Now().Format("20060102-150405"))

	// AutoRemove 的容器一停就自动删除，改名会失败 —— 先关掉它（新容器仍按原配置创建）。
	if hc, ok := insp["HostConfig"].(map[string]any); ok {
		if ar, _ := hc["AutoRemove"].(bool); ar {
			if err := u.dc.UpdateContainer(ctx, oldID, map[string]any{"AutoRemove": false}); err == nil {
				u.step(res, "原容器开启了自动删除，已临时关闭以保证可回滚")
			}
		}
	}

	// 5. 停止
	if wasRunning {
		u.status(res, ResultStopped, "正在停止旧容器…")
		timeout := 30
		if err := u.dc.ContainerAction(ctx, oldID, "stop", &timeout); err != nil {
			return "", "停止旧容器失败：" + err.Error(), true
		}
		u.step(res, "旧容器已停止")
	} else {
		u.step(res, "旧容器原本就是停止状态，无需停止")
	}

	// 6. 改名保留（绝不删除 —— 旧容器的可写层与运行期身份因此保留下来）
	if err := u.dc.RenameContainer(ctx, oldID, bakName); err != nil {
		if wasRunning {
			_ = u.dc.ContainerAction(ctx, oldID, "start", nil)
		}
		return "", "重命名旧容器失败：" + err.Error(), true
	}
	u.step(res, "旧容器已改名为 %s（保留以便回滚）", bakName)

	// 7. 用原配置创建同名新容器
	// 先取一次现存网络，避免把已被删除的网络别名传回去导致 404
	var existingNets map[string]bool
	if nets, err := u.dc.ListNetworks(ctx); err == nil {
		existingNets = NetworkNameSet(nets)
	}
	cfg, hostCfg, netCfg := BuildCreateSpec(insp, existingNets)
	u.status(res, ResultCreated, "正在创建新容器…")
	newID, err := u.dc.CreateContainer(ctx, name, cfg, hostCfg, netCfg)
	if err != nil {
		u.step(res, "创建新容器失败：%v，开始回滚", err)
		return "", "创建新容器失败：" + err.Error(), u.rollback(ctx, oldID, bakName, name, wasRunning, res)
	}
	u.step(res, "新容器已创建（%s）", dockerx.ShortID(newID))

	// 8. 启动（原本停止的容器保持停止）
	if wasRunning {
		if err := u.dc.ContainerAction(ctx, newID, "start", nil); err != nil {
			_ = u.dc.RemoveContainer(ctx, newID, true, false)
			return "", "启动新容器失败：" + err.Error(), u.rollback(ctx, oldID, bakName, name, wasRunning, res)
		}
		u.step(res, "新容器已启动，开始健康检查")
		// 9. 健康判定
		if err := u.waitHealthy(ctx, newID, insp); err != nil {
			_ = u.dc.ContainerAction(ctx, newID, "stop", nil)
			_ = u.dc.RemoveContainer(ctx, newID, true, false)
			return "", "新容器健康检查未通过：" + err.Error(), u.rollback(ctx, oldID, bakName, name, wasRunning, res)
		}
		u.step(res, "健康检查通过")
	}

	// 10. 清理备份容器（保留策略由 KeepBackups 决定）
	if u.opts.KeepBackups <= 0 {
		if err := u.dc.RemoveContainer(ctx, oldID, false, false); err != nil {
			u.step(res, "⚠ 备份容器 %s 删除失败（可稍后手动清理）：%v", bakName, err)
		} else {
			u.step(res, "已清理备份容器 %s", bakName)
		}
	} else {
		u.step(res, "按保留策略保留备份容器 %s", bakName)
	}
	return newID, "", false
}

// rollback 回滚：删新容器 → 把 __bak_ 改回原名 → 按原状态启动。
func (u *Updater) rollback(ctx context.Context, oldID, bakName, name string, wasRunning bool, res *Result) bool {
	u.step(res, "开始回滚：删除新容器并把旧容器改回原名")
	if err := u.dc.RenameContainer(ctx, oldID, name); err != nil {
		u.step(res, "✗ 回滚失败（改名）：%v", err)
		return false
	}
	if wasRunning {
		if err := u.dc.ContainerAction(ctx, oldID, "start", nil); err != nil {
			u.step(res, "✗ 回滚失败（启动）：%v", err)
			return false
		}
	}
	u.step(res, "✓ 已回滚，容器 %s 恢复为更新前的状态", name)
	return true
}

// waitHealthy 等待容器健康。
func (u *Updater) waitHealthy(ctx context.Context, id string, before map[string]any) error {
	hasHealth := false
	if c, ok := before["Config"].(map[string]any); ok {
		if _, ok := c["Healthcheck"]; ok {
			hasHealth = true
		}
	}
	deadline := time.Now().Add(u.opts.HealthTimeout)
	stableSince := time.Time{}
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("超时（%s）", u.opts.HealthTimeout)
		}
		insp, err := u.dc.Inspect(ctx, id)
		if err != nil {
			return err
		}
		st, _ := insp["State"].(map[string]any)
		if st == nil {
			return fmt.Errorf("无法读取容器状态")
		}
		running, _ := st["Running"].(bool)
		if !running {
			exit, _ := st["ExitCode"].(float64)
			return fmt.Errorf("容器已退出（退出码 %d）", int(exit))
		}
		if hasHealth {
			h, _ := st["Health"].(map[string]any)
			if h != nil {
				hs, _ := h["Status"].(string)
				switch hs {
				case "healthy":
					return nil
				case "unhealthy":
					return fmt.Errorf("健康检查状态为 unhealthy")
				}
			}
		} else {
			if stableSince.IsZero() {
				stableSince = time.Now()
			} else if time.Since(stableSince) >= time.Duration(u.opts.StableSeconds)*time.Second {
				return nil
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1500 * time.Millisecond):
		}
	}
}

// UpdateMany 批量更新（有界并发），并汇总成一条通知。
func (u *Updater) UpdateMany(ctx context.Context, names []string, force bool) []Result {
	if len(names) == 0 {
		return nil
	}
	conc := max(1, u.opts.Concurrency)
	sem := make(chan struct{}, conc)
	var wg sync.WaitGroup
	var mu sync.Mutex
	out := []Result{}

	batch := u.bus.Publish("update", "batch_start", "running", map[string]any{
		"total": len(names),
	})
	_ = batch

	for _, n := range names {
		wg.Add(1)
		go func(n string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			r := u.Update(ctx, n, force)
			mu.Lock()
			out = append(out, *r)
			mu.Unlock()
		}(n)
	}
	wg.Wait()

	updated, skipped, failed := 0, 0, 0
	for _, r := range out {
		switch r.Status {
		case ResultUpdated:
			updated++
		case ResultUpToDate, ResultSkipped:
			skipped++
		default:
			failed++
		}
	}
	u.bus.Publish("update", "batch_done", "success", map[string]any{
		"total": len(names), "updated": updated, "skipped": skipped, "failed": failed,
	})
	u.log("update", "batch", "done",
		fmt.Sprintf("共 %d 个容器：更新 %d、已是最新/跳过 %d、失败 %d", len(names), updated, skipped, failed), "")
	// 批量更新合并成一条汇总，而不是每台一条
	u.notify.Emit("batch_update_done", map[string]string{
		"container": fmt.Sprintf("%d 个容器", len(names)),
		"result":    "批量更新完成",
		"message":   fmt.Sprintf("更新 %d 个，已是最新 %d 个，失败 %d 个", updated, skipped, failed),
	})
	return out
}

// saveSnapshot 把容器配置快照写到 /data/backups/containers/<name>/<ts>.json。
func (u *Updater) saveSnapshot(name string, insp map[string]any) (string, error) {
	if u.opts.BackupDir == "" {
		return "", fmt.Errorf("未配置备份目录")
	}
	dir := filepath.Join(u.opts.BackupDir, name)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	b, err := json.MarshalIndent(insp, "", "  ")
	if err != nil {
		return "", err
	}
	p := filepath.Join(dir, time.Now().Format("20060102-150405")+".json")
	// 0600：配置里可能含密码（Env）
	if err := os.WriteFile(p, b, 0o600); err != nil {
		return "", err
	}
	return p, nil
}

// ---------- inspect 取值助手 ----------

func inspectName(insp map[string]any) string {
	if s, ok := insp["Name"].(string); ok {
		return strings.TrimPrefix(s, "/")
	}
	return ""
}

func inspectID(insp map[string]any) string {
	s, _ := insp["Id"].(string)
	return s
}

func inspectImageRef(insp map[string]any) string {
	if c, ok := insp["Config"].(map[string]any); ok {
		if s, ok := c["Image"].(string); ok {
			return s
		}
	}
	return ""
}

func inspectRunning(insp map[string]any) bool {
	if st, ok := insp["State"].(map[string]any); ok {
		b, _ := st["Running"].(bool)
		return b
	}
	return false
}

// containerRepoDigests 从容器 inspect 里尽力取镜像摘要。
func containerRepoDigests(insp map[string]any) []string {
	out := []string{}
	// 有的版本把 Image 的 RepoDigests 放在 ImageManifestDescriptor 里
	if imd, ok := insp["ImageManifestDescriptor"].(map[string]any); ok {
		if d, ok := imd["digest"].(string); ok && d != "" {
			out = append(out, d)
		}
	}
	return out
}

func shortDigest(d string) string {
	d = strings.TrimPrefix(d, "sha256:")
	if len(d) > 12 {
		return d[:12]
	}
	if d == "" {
		return "-"
	}
	return d
}

func truncName(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
