package api

import (
	"context"
	"net/http"
	"time"

	"github.com/aaron2024s/dockhelm/internal/autoupdate"
	"github.com/aaron2024s/dockhelm/internal/updater"
)

// AutoRunItem 一轮自动更新里单台容器的去向。
type AutoRunItem struct {
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"` // updated | up_to_date | failed | skipped | pending
	Reason string `json:"reason,omitempty"`
}

// AutoRunSummary 描述「一轮巡检 + 自动更新」的完整结果。
//
// 刻意把「检测到什么」与「做了什么」放在同一份结构里：自动更新最容易引起的
// 疑问是「它到底动没动我的容器」，把 24 台容器各自的去向逐台列出来，
// 才能一眼看出被跳过的是哪几台、为什么跳过。
type AutoRunSummary struct {
	StartedAt  string       `json:"startedAt"`
	FinishedAt string       `json:"finishedAt"`
	Trigger    string       `json:"trigger"` // schedule | manual
	// DryRun 只巡检不执行（总开关关着，或手动勾了「仅预览」）。
	DryRun bool `json:"dryRun"`
	// Checked 本轮巡检覆盖的容器数。
	Checked int `json:"checked"`
	// Available 巡检查出的「有可用更新」的台数。
	Available int `json:"available"`
	// Updated / Failed 真正重建成功的台数 / 失败的台数。
	Updated int `json:"updated"`
	Failed  int `json:"failed"`
	// ReclaimedMB 清理旧镜像回收的空间。
	ReclaimedMB float64 `json:"reclaimedMB"`
	// DurationMs 整轮耗时。
	DurationMs int64        `json:"durationMs"`
	Containers []AutoRunItem `json:"containers"`
	Error      string        `json:"error,omitempty"`
}

// hGetAutoUpdate 返回自动更新的配置、下一次巡检时间、候选预览与上一轮结果。
//
// 「候选预览」是本接口的重点：它把「如果现在执行，会动哪几台、跳过哪几台、
// 为什么跳过」算好给前端，用户不必等到容器真的被重启才发现自己被排除列表漏掉了。
func (s *Server) hGetAutoUpdate(w http.ResponseWriter, r *http.Request) {
	cfg := s.readSettings()
	ctx, cancel := s.ctx(r)
	defer cancel()

	cands := s.candidatePlan(ctx, cfg.Exclude)

	s.autoMu.Lock()
	busy := s.autoBusy
	last := s.autoLast
	lastCheck := s.autoLastCheck
	s.autoMu.Unlock()

	next := ""
	if cfg.CheckIntervalHours > 0 {
		base := lastCheck
		if base.IsZero() {
			s.checkMu.RLock()
			base = s.checkAt
			s.checkMu.RUnlock()
		}
		if base.IsZero() {
			base = time.Now()
		}
		next = base.Add(time.Duration(cfg.CheckIntervalHours) * time.Hour).UTC().Format(time.RFC3339)
	}

	s.checkMu.RLock()
	checkedAt := s.checkAt
	s.checkMu.RUnlock()

	writeOK(w, map[string]any{
		"enabled":            cfg.AutoApply,
		"checkIntervalHours": cfg.CheckIntervalHours,
		"notifyOnCheck":      cfg.NotifyOnCheck,
		"intervalChoices":    checkIntervalChoices,
		"policy": map[string]any{
			"concurrency":  cfg.Concurrency,
			"pullOnce":     cfg.PullOnce,
			"backupBefore": cfg.BackupBefore,
			"cleanupAfter": cfg.CleanupAfter,
			"directFirst":  cfg.DirectFirst,
		},
		"running":     busy,
		"nextCheckAt": next,
		"lastCheckAt": timeOrEmpty(checkedAt),
		"candidates":  cands,
		"willUpdate":  autoupdate.Count(cands),
		"exclude":     cfg.Exclude,
		"selfName":    s.up.SelfName(),
		"lastRun":     last,
	})
}

// hRunAutoUpdate 手动立即跑一轮。「仅预览」时只巡检不执行。
func (s *Server) hRunAutoUpdate(w http.ResponseWriter, r *http.Request) {
	var in struct {
		DryRun bool `json:"dryRun"`
	}
	_ = decodeBody(r, &in)

	s.autoMu.Lock()
	busy := s.autoBusy
	s.autoMu.Unlock()
	if busy {
		writeErr(w, http.StatusConflict, "上一轮自动更新还在跑，稍候再看")
		return
	}

	go s.runAutoCycle(context.Background(), "manual", !in.DryRun)
	writeOK(w, map[string]any{"accepted": true, "dryRun": in.DryRun})
}

// ---------- 后台循环 ----------

// StartAutoLoop 启动「周期检测（可选自动更新）」后台循环。
//
// 它替代了原来 main.go 里那个「启动后只跑一次」的巡检：
// 二者只能留一个，否则启动瞬间会有两轮巡检同时打镜像仓库，
// 白白触发限流。
func (s *Server) StartAutoLoop(ctx context.Context) {
	go s.autoLoop(ctx)
}

// wakeAutoLoop 通知后台循环「设置刚改过，按新间隔重新计时」。
// 非阻塞：调用方（HTTP handler）不该因为循环正忙而卡住。
func (s *Server) wakeAutoLoop() {
	select {
	case s.autoWake <- struct{}{}:
	default:
	}
}

// autoLoop 按 update.checkInterval 周期跑一轮「检测 + 可选自动更新」。
//
// 关掉周期检测时定时器仍会走，但不干活 —— 这样「打开开关」这个动作
// 只需要唤醒一次就能马上开始计时，而不必重建 goroutine。
func (s *Server) autoLoop(ctx context.Context) {
	// 启动后先等一会儿：让面板先起来，也避开与首次巡检的竞态。
	timer := time.NewTimer(2 * time.Minute)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.autoWake:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			timer.Reset(s.nextInterval())
		case <-timer.C:
			cfg := s.readSettings()
			if cfg.CheckIntervalHours > 0 {
				s.runAutoCycle(ctx, "schedule", true)
			}
			timer.Reset(s.nextInterval())
		}
	}
}

// nextInterval 返回下一次巡检的间隔。关闭时给一个长兜底，
// 避免 0 时长定时器把循环打成忙等。
func (s *Server) nextInterval() time.Duration {
	h := s.readSettings().CheckIntervalHours
	if h <= 0 {
		return 24 * time.Hour
	}
	return time.Duration(h) * time.Hour
}

// runAutoCycle 跑一整轮：巡检 → 决策 → （可选）执行。
//
// execute 为 false 时只巡检并更新候选，绝不动任何容器 —— 这就是
// 「总开关关掉时它只是看着，永远不会擅自动手」的落点。
func (s *Server) runAutoCycle(ctx context.Context, trigger string, execute bool) {
	s.autoMu.Lock()
	if s.autoBusy {
		s.autoMu.Unlock()
		return
	}
	s.autoBusy = true
	s.autoMu.Unlock()
	defer func() {
		s.autoMu.Lock()
		s.autoBusy = false
		s.autoMu.Unlock()
	}()

	cfg := s.readSettings()
	started := time.Now()
	summary := &AutoRunSummary{
		StartedAt: started.UTC().Format(time.RFC3339),
		Trigger:   trigger,
		DryRun:    !execute,
	}

	cctx, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()

	// ---- 1. 巡检（只读，走 /distribution，绝不动容器）----
	results := s.up.CheckAll(cctx, s.excluded(), false)
	s.checkMu.Lock()
	s.checkCache = results
	s.checkAt = started
	s.checkMu.Unlock()
	summary.Checked = len(results)

	s.autoMu.Lock()
	s.autoLastCheck = started
	s.autoMu.Unlock()

	// ---- 2. 决策：哪些会被更新、哪些被跳过、为什么 ----
	cands := s.candidatePlan(cctx, cfg.Exclude)
	summary.Available = autoupdate.Count(cands)
	for _, c := range cands {
		item := AutoRunItem{Name: c.Name, Image: c.Image}
		switch {
		case c.WillUpdate:
			if execute {
				item.Status = "pending"
			} else {
				item.Status = "pending"
				item.Reason = "总开关关闭，仅预览"
			}
		case c.Protected:
			item.Status = "skipped"
			item.Reason = c.Reason
		case c.Excluded:
			item.Status = "skipped"
			item.Reason = c.Reason
		default:
			item.Status = "up_to_date"
		}
		summary.Containers = append(summary.Containers, item)
	}

	s.bus.Publish("update", "auto_check_done", "success", map[string]any{
		"checked": summary.Checked, "available": summary.Available, "trigger": trigger,
	})
	if cfg.NotifyOnCheck && summary.Available > 0 {
		s.nt.Emit("update_available", map[string]string{
			"container": "",
			"result":    "检测到新版本",
			"message":   "本轮巡检发现 " + itoa(summary.Available) + " 台容器有可用更新",
		})
	}

	// ---- 3. 总开关关着，或没有候选 ⇒ 到此为止 ----
	if !execute || summary.Available == 0 {
		s.finishCycle(summary, started, "仅巡检，未改动任何容器")
		return
	}

	// ---- 4. 执行 ----
	names := autoupdate.Names(cands)
	rows := s.up.UpdateMany(cctx, names, false)

	byName := map[string]AutoRunItem{}
	for _, it := range summary.Containers {
		byName[it.Name] = it
	}
	var cacheRows []updater.CheckResult
	for _, res := range rows {
		item := byName[res.Container]
		item.Name = res.Container
		item.Image = res.Image
		switch res.Status {
		case updater.ResultUpdated:
			item.Status = "updated"
			item.Reason = ""
			summary.Updated++
			cacheRows = append(cacheRows, freshCheckRow(res, "刚刚自动更新到最新"))
		case updater.ResultUpToDate:
			item.Status = "up_to_date"
			item.Reason = "拉取后镜像未变化"
			cacheRows = append(cacheRows, freshCheckRow(res, "拉取后镜像未变化"))
		default:
			item.Status = "failed"
			item.Reason = res.Message
			summary.Failed++
		}
		summary.ReclaimedMB += float64(res.ReclaimedBytes) / 1024 / 1024
		byName[res.Container] = item
	}
	// 按原顺序写回，保持与预览一致
	for i := range summary.Containers {
		if it, ok := byName[summary.Containers[i].Name]; ok {
			summary.Containers[i] = it
		}
	}
	s.mergeCheckResults(cacheRows)

	s.finishCycle(summary, started, "")
}

func (s *Server) finishCycle(summary *AutoRunSummary, started time.Time, note string) {
	summary.FinishedAt = time.Now().UTC().Format(time.RFC3339)
	summary.DurationMs = time.Since(started).Milliseconds()
	s.autoMu.Lock()
	s.autoLast = summary
	s.autoMu.Unlock()

	msg := "自动更新：更新 " + itoa(summary.Updated) + " 台，失败 " + itoa(summary.Failed) + " 台"
	if note != "" {
		msg = "自动巡检：" + itoa(summary.Available) + " 台有可用更新（" + note + "）"
	}
	s.st.AddRunLog("auto_update", summary.Trigger, "done", msg, "")
	s.bus.Publish("update", "auto_done", "success", map[string]any{
		"trigger": summary.Trigger,
		"updated": summary.Updated,
		"failed":  summary.Failed,
	})
	if summary.Updated > 0 || summary.Failed > 0 {
		s.nt.Emit("auto_update_done", map[string]string{
			"container": itoa(summary.Updated + summary.Failed) + " 台",
			"result":    "自动更新完成",
			"message":   msg,
		})
	}
}

// candidatePlan 把「巡检缓存 + 当前容器列表」合成候选预览。
//
// 用 ListContainers 而不是直接吃 checkCache：只有这里才能拿到容器的运行状态，
// 也才能把「在跑但还没被巡检过」的容器一并列进预览（显示成未检测到更新）。
func (s *Server) candidatePlan(ctx context.Context, exclude []string) []autoupdate.Candidate {
	s.checkMu.RLock()
	byName := map[string]updater.CheckResult{}
	for _, c := range s.checkCache {
		byName[c.Container] = c
	}
	s.checkMu.RUnlock()

	list, err := s.dc.ListContainers(ctx)
	if err != nil {
		return []autoupdate.Candidate{}
	}
	conts := make([]autoupdate.Container, 0, len(list))
	for _, c := range list {
		name := c.Name()
		cr := byName[name]
		conts = append(conts, autoupdate.Container{
			Name:      name,
			Image:     c.Image,
			ImageID:   c.ImageID,
			Running:   c.State == "running",
			HasUpdate: cr.Status == updater.StatusUpdateAvailable,
		})
	}
	return autoupdate.Plan(conts, autoupdate.Options{
		SelfName: s.up.SelfName(),
		Exclude:  exclude,
	})
}

// freshCheckRow 把一次更新结果翻译成巡检缓存里的一行。
func freshCheckRow(res updater.Result, reason string) updater.CheckResult {
	return updater.CheckResult{
		Container: res.Container,
		Image:     res.Image,
		Status:    updater.StatusUpToDate,
		Reason:    reason,
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
	}
}
