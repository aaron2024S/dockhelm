// Package scheduler 实现计划任务：定时启动 / 停止 / 重启 / 更新 / 备份容器。
//
// 任务定义的语法沿用标准 5 段 cron（分 时 日 月 周），并支持 @every 1h / @daily 这类描述符。
// 每条任务的执行结果会写进运行记录，并在失败时发通知。
package scheduler

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/aaron2024s/dockhelm/internal/backup"
	"github.com/aaron2024s/dockhelm/internal/bus"
	"github.com/aaron2024s/dockhelm/internal/dockerx"
	"github.com/aaron2024s/dockhelm/internal/intent"
	"github.com/aaron2024s/dockhelm/internal/notify"
	"github.com/aaron2024s/dockhelm/internal/store"
	"github.com/aaron2024s/dockhelm/internal/updater"
)

// Action 支持的任务动作。
type Action struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description"`
	// NeedsTargets 是否必须明确指定容器。
	// true 的动作在创建/编辑时至少要勾一个容器 —— 空目标**不再**表示「全部容器」。
	NeedsTargets bool `json:"needsTargets"`
}

// Actions 全部动作。
var Actions = []Action{
	{Key: "start", Label: "启动容器", Description: "把容器启起来（已在运行的会被跳过）", NeedsTargets: true},
	{Key: "stop", Label: "停止容器", Description: "停止容器，不删除", NeedsTargets: true},
	{Key: "restart", Label: "重启容器", Description: "等于 docker restart，容器不存在时跳过", NeedsTargets: true},
	{Key: "update", Label: "检查并更新", Description: "拉取镜像；只有镜像真的变了才会停容器重建", NeedsTargets: true},
	{Key: "backup", Label: "备份配置快照", Description: "为容器写一份 docker inspect 配置快照", NeedsTargets: true},
	{Key: "prune_images", Label: "清理未使用镜像", Description: "删除没有被任何容器引用的未使用镜像层", NeedsTargets: false},
}

// ActionMap 索引。
var ActionMap = func() map[string]Action {
	m := map[string]Action{}
	for _, a := range Actions {
		m[a.Key] = a
	}
	return m
}()

// Runner 定时任务的执行者。
type Runner struct {
	dc      *dockerx.Client
	up      *updater.Updater
	bk      *backup.Service
	st      *store.Store
	nt      *notify.Manager
	bus     *bus.Bus
	exclude func() map[string]bool
}

// Scheduler 调度器。
type Scheduler struct {
	c       *cron.Cron
	st      *store.Store
	runner  *Runner
	bus     *bus.Bus
	nt      *notify.Manager
	mu      sync.Mutex
	entries map[int64]cron.EntryID
	step    sync.Mutex // 多任务并发触发时串行执行
}

// New 创建调度器。
func New(st *store.Store, up *updater.Updater, dc *dockerx.Client, bk *backup.Service, nt *notify.Manager, b *bus.Bus,
	exclude func() map[string]bool) *Scheduler {
	runner := &Runner{dc: dc, up: up, bk: bk, st: st, nt: nt, bus: b, exclude: exclude}
	s := &Scheduler{
		st:      st,
		runner:  runner,
		bus:     b,
		nt:      nt,
		entries: map[int64]cron.EntryID{},
	}
	// 用「分 时 日 月 周」五段 + 描述符；秒级精度对容器任务没有意义
	s.c = cron.New(cron.WithParser(parser()))
	return s
}

func parser() cron.Parser {
	return cron.NewParser(
		cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
}

// ValidateCron 校验表达式。
func ValidateCron(expr string) error {
	if strings.TrimSpace(expr) == "" {
		return fmt.Errorf("cron 表达式不能为空")
	}
	_, err := parser().Parse(strings.TrimSpace(expr))
	if err != nil {
		return fmt.Errorf("cron 表达式无效：%v", err)
	}
	return nil
}

// Start 启动调度器并载入已有任务。
func (s *Scheduler) Start() error {
	if err := s.Reload(); err != nil {
		return err
	}
	s.c.Start()
	return nil
}

// Stop 停止调度器。
func (s *Scheduler) Stop() {
	ctx := s.c.Stop()
	<-ctx.Done()
}

// Reload 重新载入全部任务（增删改后调用）。
func (s *Scheduler) Reload() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, entry := range s.entries {
		s.c.Remove(entry)
		delete(s.entries, id)
	}
	list, err := s.st.ListSchedules()
	if err != nil {
		return err
	}
	for i := range list {
		item := list[i]
		if !item.Enabled {
			continue
		}
		schedule := item
		entryID, err := s.c.AddFunc(strings.TrimSpace(schedule.Cron), func() {
			s.runJob(schedule.ID)
		})
		if err != nil {
			s.st.AddRunLog("schedule", schedule.Name, "failed", "cron 表达式无法注册："+err.Error(), "")
			continue
		}
		s.entries[schedule.ID] = entryID
	}
	return nil
}

// NextRuns 返回每个任务的下次执行时间。
func (s *Scheduler) NextRuns() map[int64]string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := map[int64]string{}
	for id, entry := range s.entries {
		e := s.c.Entry(entry)
		if !e.Next.IsZero() {
			out[id] = e.Next.UTC().Format(time.RFC3339)
		}
	}
	return out
}

// RunNow 立即执行一次（界面上的「立即运行」按钮）。
func (s *Scheduler) RunNow(id int64) (string, bool) {
	return s.runJob(id)
}

func (s *Scheduler) runJob(id int64) (string, bool) {
	// 多个任务同时到点时串行执行，避免同一批容器被并发重启
	s.step.Lock()
	defer s.step.Unlock()

	sc, err := s.st.GetSchedule(id)
	if err != nil {
		return "任务不存在", false
	}
	started := time.Now()
	s.bus.Publish("schedule", "start", "running", map[string]any{
		"id": sc.ID, "name": sc.Name, "action": sc.Action,
	})

	out := s.runner.Execute(context.Background(), *sc)
	status := "success"
	if !out.OK {
		status = "failed"
	}
	// 运行记录/上次结果只放单行简报；完整正文（含名单）留给推送，原始报错留给记录详情。
	_ = s.st.SetScheduleResult(id, status, out.Brief)
	s.st.AddRunLog("schedule", sc.Name, status, out.Brief, out.Detail)
	s.bus.Publish("schedule", "done", status, map[string]any{
		"id": sc.ID, "name": sc.Name, "action": sc.Action,
		"message": out.Brief, "durationMs": time.Since(started).Milliseconds(),
	})

	event := "schedule_success"
	if !out.OK {
		event = "schedule_failed"
	}
	// 标题用**任务名**（compose 把 container 拼在事件名后面）：
	//   ① 定时任务一多，光看「计划任务执行成功」分不清是哪个任务；
	//   ② 去重键是 event+":"+container，原来 container 恒为空 ⇒ 同一窗口内
	//      两个不同任务的通知会互相顶掉（10 分钟内只发得出第一条）。
	label := ActionMap[sc.Action].Label
	if label == "" {
		label = "计划任务"
	}
	s.nt.Emit(event, map[string]string{
		"container": sc.Name,
		"result":    label,
		"message":   out.Full,
	})
	return out.Brief, out.OK
}

// Outcome 一次任务执行的产出。三个字段各司其职，别再挤进一个字符串：
//
//   - Brief  单行简报 → 运行记录、总览的「上次结果」（那里只显示一行）
//   - Full   多行完整正文 → 推送的「详情」（点名容器、逐项失败）
//   - Detail 原始错误 → 只进运行记录详情，**不进推送**
//
// 2026-10-12：以前只有一句「执行完成：成功 N 个，跳过 M 个」——看不出是哪几台，
// 失败还把 `Error response from daemon: …` 原样甩进推送（用户 10-11 明确要求去掉）。
type Outcome struct {
	Brief  string
	Full   string
	Detail string
	OK     bool
}

// outcome 简报与正文相同的结果（没有名单可列的场合）。
func outcome(msg string, ok bool) Outcome { return Outcome{Brief: msg, Full: msg, OK: ok} }

// maxScheduleNames 回执里最多点名几个容器，超出折叠成「等 N 个」。
const maxScheduleNames = 12

// foldNames 输出容器名单：先排序（遍历/完成顺序不定，排序后才稳定），
// 超过 maxScheduleNames 个折叠。
func foldNames(names []string) string {
	sorted := append([]string(nil), names...)
	sort.Strings(sorted)
	if len(sorted) > maxScheduleNames {
		return strings.Join(sorted[:maxScheduleNames], "、") + fmt.Sprintf(" 等 %d 个", len(sorted))
	}
	return strings.Join(sorted, "、")
}

// summarize 把一次任务的结果拼成推送正文：统计行 + 已动作名单 + 跳过名单 + 失败名单。
func summarize(action string, total, done, skipped int, doneNames, skippedNames, failedNames []string) string {
	verb := map[string]string{
		"start": "已启动", "stop": "已停止", "restart": "已重启",
		"update": "已更新", "backup": "已备份",
	}[action]
	if verb == "" {
		verb = "已完成"
	}
	skipWhy := map[string]string{
		"start": "本就已运行", "stop": "本就已停止", "update": "已是最新",
	}[action]

	head := fmt.Sprintf("共 %d 个容器：成功 %d 个、跳过 %d 个", total, done, skipped)
	if len(failedNames) > 0 {
		head = fmt.Sprintf("共 %d 个容器：成功 %d 个、跳过 %d 个、失败 %d 个",
			total, done, skipped, len(failedNames))
	}
	lines := []string{head}
	if len(doneNames) > 0 {
		lines = append(lines, verb+"："+foldNames(doneNames))
	}
	if len(skippedNames) > 0 {
		note := "跳过"
		if skipWhy != "" {
			note = "跳过（" + skipWhy + "）"
		}
		lines = append(lines, note+"："+foldNames(skippedNames))
	}
	if len(failedNames) > 0 {
		lines = append(lines, "失败："+foldNames(failedNames))
	}
	return strings.Join(lines, "\n")
}

// countLine 单行简报（统计部分），进运行记录与总览。
func countLine(total, done, skipped int, failed int) string {
	if failed > 0 {
		return fmt.Sprintf("共 %d 个容器：成功 %d 个、跳过 %d 个、失败 %d 个", total, done, skipped, failed)
	}
	return fmt.Sprintf("共 %d 个容器：成功 %d 个、跳过 %d 个", total, done, skipped)
}

// Execute 执行一次任务定义（供调度器与「立即运行」共用）。
func (r *Runner) Execute(ctx context.Context, sc store.Schedule) Outcome {
	switch sc.Action {
	case "prune_images":
		freed, err := r.dc.PruneImages(ctx)
		if err != nil {
			return outcome("清理未使用镜像失败："+err.Error(), false)
		}
		return outcome(fmt.Sprintf("已清理未使用镜像，释放 %.1f MB", float64(freed)/1024/1024), true)
	}

	// 其余动作都必须有明确目标。空目标曾经等于「全部容器」，已于 0.2.x 取消 ——
	// 老版本存下来的任务会走到这里，给一句说得清的话，而不是含糊的「没有匹配到任何容器」。
	if len(sc.Targets) == 0 {
		return outcome("任务没有选择任何容器，请检查任务配置", false)
	}

	if sc.Action == "backup" {
		targets, err := r.resolveTargets(ctx, sc.Targets)
		if err != nil {
			return outcome(err.Error(), false)
		}
		if len(targets) == 0 {
			return outcome("没有匹配到任何容器（可能已被删除或改名）", false)
		}
		done := 0
		doneNames, failedNames, raw := []string{}, []string{}, []string{}
		for _, t := range targets {
			if _, err := r.bk.Snapshot(ctx, t.name, "schedule:"+sc.Name); err != nil {
				failedNames = append(failedNames, t.name)
				raw = append(raw, fmt.Sprintf("%s（%v）", t.name, err))
			} else {
				done++
				doneNames = append(doneNames, t.name)
			}
		}
		return Outcome{
			Brief:  countLine(len(targets), done, 0, len(failedNames)),
			Full:   summarize("backup", len(targets), done, 0, doneNames, nil, failedNames),
			Detail: strings.Join(raw, "\n"),
			OK:     len(failedNames) == 0,
		}
	}

	targets, err := r.resolveTargets(ctx, sc.Targets)
	if err != nil {
		return outcome(err.Error(), false)
	}
	if len(targets) == 0 {
		return outcome("没有匹配到任何容器（可能已被删除或改名）", false)
	}

	done, skipped := 0, 0
	doneNames, skippedNames, failedNames, raw := []string{}, []string{}, []string{}, []string{}
	for _, t := range targets {
		switch sc.Action {
		case "start":
			if t.running {
				skipped++
				skippedNames = append(skippedNames, t.name)
				continue
			}
			if err := r.dc.ContainerAction(ctx, t.id, "start", nil); err != nil {
				failedNames = append(failedNames, t.name)
				raw = append(raw, fmt.Sprintf("%s（%v）", t.name, err))
			} else {
				done++
				doneNames = append(doneNames, t.name)
			}
		case "stop":
			if !t.running {
				skipped++
				skippedNames = append(skippedNames, t.name)
				continue
			}
			timeout := 30
			// 任务发起的停止不是「容器意外退出」，别让事件观察记一条失败
			intent.Mark(t.name)
			if err := r.dc.ContainerAction(ctx, t.id, "stop", &timeout); err != nil {
				intent.Release(t.name) // 没停成、不会有 die 事件，撤销豁免
				failedNames = append(failedNames, t.name)
				raw = append(raw, fmt.Sprintf("%s（%v）", t.name, err))
			} else {
				done++
				doneNames = append(doneNames, t.name)
			}
		case "restart":
			timeout := 30
			intent.Mark(t.name)
			if err := r.dc.ContainerAction(ctx, t.id, "restart", &timeout); err != nil {
				intent.Release(t.name)
				failedNames = append(failedNames, t.name)
				raw = append(raw, fmt.Sprintf("%s（%v）", t.name, err))
			} else {
				done++
				doneNames = append(doneNames, t.name)
			}
		case "update":
			res := r.up.Update(ctx, t.id, false)
			switch res.Status {
			case updater.ResultUpdated:
				done++
				doneNames = append(doneNames, t.name)
			case updater.ResultUpToDate, updater.ResultSkipped:
				skipped++
				skippedNames = append(skippedNames, t.name)
			default:
				// 推送里只写「容器名：动作短语」；底层报错留给运行记录详情。
				failedNames = append(failedNames, t.name+"："+updater.ShortPhase(res.Message))
				raw = append(raw, fmt.Sprintf("%s（%s）", t.name, res.Message))
			}
		default:
			return outcome("未知的动作："+sc.Action, false)
		}
	}

	return Outcome{
		Brief:  countLine(len(targets), done, skipped, len(failedNames)),
		Full:   summarize(sc.Action, len(targets), done, skipped, doneNames, skippedNames, failedNames),
		Detail: strings.Join(raw, "\n"),
		OK:     len(failedNames) == 0,
	}
}

type target struct {
	id      string
	name    string
	running bool
}

// resolveTargets 把「容器名列表」解析成实际存在的目标。
//
// 空列表 = 没有目标（返回空切片）。**曾经**空列表等于「全部容器（排除自身与排除项）」，
// 但那样一次误操作就会作用到整机，界面上还看不出来，已取消这个语义；调用方在更上层
// 直接拦住空目标。
func (r *Runner) resolveTargets(ctx context.Context, names []string) ([]target, error) {
	list, err := r.dc.ListContainers(ctx)
	if err != nil {
		return nil, fmt.Errorf("读取容器列表失败：%w", err)
	}
	ex := map[string]bool{}
	if r.exclude != nil {
		ex = r.exclude()
	}
	if len(names) == 0 {
		return nil, nil
	}
	byName := map[string]target{}
	for _, c := range list {
		n := c.Name()
		if r.up.IsSelf(n) || ex[n] {
			continue
		}
		t := target{id: c.ID, name: n, running: c.State == "running"}
		byName[n] = t
		byName[c.ID] = t
		byName[c.ID[:12]] = t
	}
	out := []target{}
	missing := []string{}
	for _, n := range names {
		if t, ok := byName[n]; ok {
			out = append(out, t)
		} else {
			missing = append(missing, n)
		}
	}
	if len(out) == 0 && len(missing) > 0 {
		return nil, fmt.Errorf("指定的容器都不存在：%s", strings.Join(missing, ", "))
	}
	return out, nil
}
