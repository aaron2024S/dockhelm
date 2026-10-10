package api

import (
	"context"
	"net/http"
	"time"

	"github.com/aaron2024s/dockhelm/internal/autoupdate"
	"github.com/aaron2024s/dockhelm/internal/updater"
)

// AutoRunItem 一轮自动更新里单个容器的去向。
type AutoRunItem struct {
	Name   string `json:"name"`
	Image  string `json:"image"`
	Status string `json:"status"` // updated | up_to_date | failed | skipped | pending
	Reason string `json:"reason,omitempty"`
}

// AutoRunSummary 描述「一轮巡检 + 自动更新」的完整结果。
//
// 刻意把「检测到什么」与「做了什么」放在同一份结构里：自动更新最容易引起的
// 疑问是「它到底动没动我的容器」，把 24 个容器各自的去向逐个列出来，
// 才能一眼看出被跳过的是哪几个、为什么跳过。
type AutoRunSummary struct {
	StartedAt  string `json:"startedAt"`
	FinishedAt string `json:"finishedAt"`
	Trigger    string `json:"trigger"` // schedule | manual
	// DryRun 只巡检不执行（总开关关着时的定时轮次，或手动 dry-run）。
	DryRun bool `json:"dryRun"`
	// Checked 本轮巡检覆盖的容器数。
	Checked int `json:"checked"`
	// Available 巡检查出的「有可用更新」的个数。
	Available int `json:"available"`
	// Updated / Failed 真正重建成功的个数 / 失败的个数。
	Updated int `json:"updated"`
	Failed  int `json:"failed"`
	// ReclaimedMB 清理旧镜像回收的空间。
	ReclaimedMB float64 `json:"reclaimedMB"`
	// DurationMs 整轮耗时。
	DurationMs int64         `json:"durationMs"`
	Containers []AutoRunItem `json:"containers"`
	Error      string        `json:"error,omitempty"`
}

// hGetAutoUpdate 返回自动更新的配置、下一次巡检时间、候选预览与上一轮结果。
//
// 「候选预览」是本接口的重点：它把「如果现在执行，会动哪几个、跳过哪几个、
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
	s.autoMu.Lock()
	nextAt := s.autoNextAt
	s.autoMu.Unlock()
	if !nextAt.IsZero() {
		next = nextAt.UTC().Format(time.RFC3339)
	} else if cfg.CheckIntervalHours > 0 {
		// 循环还没排过第一次（极端情况：刚起进程接口就进来了），按上次巡检推算兜底。
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

// hRunAutoUpdate 手动立即跑一轮自动更新。dryRun 为真时只巡检不执行。
//
// 全站唯一会真正动容器的手动入口。界面上不再有「立即执行自动更新」按钮（页面已删）——
// 手动更新走容器页的「更新 N 个容器」，按周期自己跑的那轮由这里承载。
// 检测一律走顶栏的只读接口（POST /api/updates/check）。
// 保留 dryRun：它是「完整跑一轮决策、但一个容器都不动」的唯一入口，
// 验证与排障都要用，去掉它等于把安全语义也删了。
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

// autoBootCooldown 是「启动冷却期」：进程起来后，第一次**可执行**的巡检
// 最早排在启动后这一刻。
//
// 为什么需要它：启动首轮巡检（8 秒）刻意是只读的 —— 那时 Docker 刚起来，
// 一堆容器在同时创建/拉镜像，此刻去重建容器既慢又容易踩到竞态。
// 但只读那一轮之后若直接按检测周期排下一轮，急着更新的人就要干等一个完整
// 周期（默认 1 小时）。折中：把第一轮可执行巡检排到启动后 15 分钟 ——
// 容器启动峰值基本过去，又远比等满一个周期快。
const autoBootCooldown = 15 * time.Minute

// wakeAction 是「设置刚保存」递给后台循环的意图。
type wakeAction int

const (
	// wakeRestart 按新间隔重新计时（绝大多数设置改动的行为）。
	wakeRestart wakeAction = iota
	// wakeSooner 把下一轮提前（见 wakeDecision）。
	wakeSooner
)

// wakeDecision 判断这次保存设置是否属于「用户明显在等着看动作」的两类改动。
//
// ① 自动更新总开关由关打到开 —— 刚打开就干等一个周期最反直觉；
// ② 检测间隔被调小 —— 用户嫌太慢，改完却还按旧的长间隔等，等于这次改动没生效。
// 其余改动（并发度、备份策略、排除列表、检测通知…）与「下一轮什么时候跑」无关，
// 保持原行为按新间隔重新计时即可。
func wakeDecision(prev, next Settings) wakeAction {
	if next.AutoApply && !prev.AutoApply {
		return wakeSooner
	}
	if next.CheckIntervalHours > 0 && next.CheckIntervalHours < prev.CheckIntervalHours {
		return wakeSooner
	}
	return wakeRestart
}

// wakeAutoLoop 通知后台循环「设置刚改过，按新间隔重新计时」。
// 非阻塞：调用方（HTTP handler）不该因为循环正忙而卡住。
func (s *Server) wakeAutoLoop(act wakeAction) {
	select {
	case s.autoWake <- act:
	default:
	}
}

// autoLoop 按 update.checkInterval 周期跑一轮「检测 + 可选自动更新」。
//
// 关掉周期检测时定时器仍会走，但不干活 —— 这样「打开开关」这个动作
// 只需要唤醒一次就能马上开始计时，而不必重建 goroutine。
//
// 启动首轮：checkOnStart 开着就 8 秒后跑一次**只读巡检**（不碰任何容器，
// 与设置页文案「延迟 8 秒执行…不会停止任何容器」一致），之后进入正常周期节奏。
// 唯一的不对称在「第二轮」：首轮要是查到了更新、而总开关是开着的，第二轮会排到
// 启动冷却期结束（15 分钟）而不是等满一个周期 —— 见 firstRoundDelay。
// ⚠ 这个开关曾经是死的：StartAutoLoop 取代了 main.go 里旧的启动巡检，
// 实现删了、设置项却留在页面上 —— 用户开了也没任何效果（2026-10-10 反馈）。
// 兜底层 runAutoCycle 里 execute && !cfg.AutoApply && trigger != "manual" 的
// 强制降级继续保护这一轮：就算将来有人把 execute 传错，也动不了容器。
func (s *Server) autoLoop(ctx context.Context) {
	bootAt := time.Now()
	timer := time.NewTimer(startupDelay(s.readSettings()))
	defer timer.Stop()
	first := true

	for {
		select {
		case <-ctx.Done():
			return
		case act := <-s.autoWake:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}
			var d time.Duration
			if first {
				// 启动巡检还没跑过：一次「保存设置」不该把它挤到一个周期之后，
				// 仍按启动延迟重排（新间隔从之后的轮次起生效）。
				d = startupDelay(s.readSettings())
			} else {
				s.autoMu.Lock()
				pending := s.autoNextAt
				s.autoMu.Unlock()
				d = wakeRoundDelay(time.Now(), bootAt, s.nextInterval(), pending, act)
			}
			s.scheduleNext(timer, d)
		case <-timer.C:
			cfg := s.readSettings()
			if first {
				first = false
				// 首轮只服务 checkOnStart；execute 恒为 false，只检查不动手。
				var available int
				if cfg.CheckOnStart {
					if sum := s.runAutoCycle(ctx, "startup", false); sum != nil {
						available = sum.Available
					}
				}
				s.scheduleNext(timer, firstRoundDelay(time.Now(), bootAt, s.nextInterval(), cfg, available))
				continue
			}
			if ok, execute := scheduledRun(cfg); ok {
				s.runAutoCycle(ctx, "schedule", execute)
			}
			s.scheduleNext(timer, s.nextInterval())
		}
	}
}

// scheduleNext 重排定时器，并把「下一轮实际排在何时」记下来。
//
// 必须记：设置页的「下次巡检」原来靠「上次巡检时间 + 间隔」推算，
// 而启动冷却期与保存设置后的提前都会让真实排期早于这个推算 ——
// 界面说 1 小时后巡检、实际 15 分钟后就跑，这种「显示比实际慢」最招人疑惑。
func (s *Server) scheduleNext(timer *time.Timer, d time.Duration) {
	timer.Reset(d)
	s.autoMu.Lock()
	s.autoNextAt = time.Now().Add(d)
	s.autoMu.Unlock()
}

// startupDelay 返回循环启动后第一轮的等待时长：checkOnStart 开着就 8 秒后
// 跑启动巡检，否则维持原来的 2 分钟静默（让面板先起来、避开启动竞态）。
func startupDelay(cfg Settings) time.Duration {
	if cfg.CheckOnStart {
		return 8 * time.Second
	}
	return 2 * time.Minute
}

// firstRoundDelay 决定「启动首轮巡检之后」到下一轮之间的等待时长。
//
// 常规情况就是等满一个周期；唯一的例外是「总开关开着、且首轮确实查到了更新」——
// 这时把下一轮提前到启动冷却期结束，让急着更新的人不必干等一整个周期
// （首轮本身是只读的，真正动手的是这一轮）。
//
// 抽成纯函数（时间与间隔都当参数传进来）是为了能对着固定时刻写断言 ——
// 「等多久」这件事一旦改错，表现是「自动更新好像不生效」，极难从日志上看出来。
func firstRoundDelay(now, bootAt time.Time, interval time.Duration, cfg Settings, available int) time.Duration {
	if !cfg.CheckOnStart || !cfg.AutoApply || available <= 0 {
		return interval
	}
	d := bootAt.Add(autoBootCooldown).Sub(now)
	// 巡检本身耗时超过冷却期（超大堆栈）时别倒着等：至少留 1 秒，
	// 免得 Reset(负值) 把定时器变成一个 0 时长的忙等触发器。
	if d < time.Second {
		d = time.Second
	}
	if d > interval {
		d = interval
	}
	return d
}

// wakeRoundDelay 决定「保存设置之后」到下一轮之间的等待时长。
//
// 大多数设置改动按新间隔重新计时（原行为）；只有 wakeSooner 那一类会提前，
// 且提前也不等于立刻就动手：至少留 30 秒缓冲，并且不早于启动冷却期结束 ——
// 免得用户刚打开总开关，就在容器启动峰值里被它重建一轮。
//
// pending 是当前已排定的下一轮时刻：普通改动按新间隔重新计时没错，
// 但不能把一个已经排得更早的轮次往后推 —— 否则「打开自动更新（15 分钟后跑）→
// 顺手改下面板地址再保存」就把那一轮悄悄顶回一个小时之后，提前等于白提。
func wakeRoundDelay(now, bootAt time.Time, interval time.Duration, pending time.Time, act wakeAction) time.Duration {
	if act != wakeSooner {
		if !pending.IsZero() {
			if until := pending.Sub(now); until > 0 && until < interval {
				return until
			}
		}
		return interval
	}
	d := 30 * time.Second
	if until := bootAt.Add(autoBootCooldown).Sub(now); until > d {
		d = until
	}
	if d > interval {
		d = interval
	}
	return d
}

// scheduledRun 决定「定时循环到点时：跑不跑 / 跑的话动不动手」。
// 返回 ok=false 表示这一轮什么都不做；ok=true 时 execute 表示允许真的更新容器。
//
// 🚨 单独抽出来是因为这里**曾经把 execute 写死成 true**：定时循环会无视
// 「检测到新版本后自动更新」这个总开关，在启动 2 分钟后真的去更新容器。
// 2026-10-09 真机事故 —— 用户从没打开过自动更新，却在总览里看到一条
// 「自动更新：更新 0 个容器，失败 1 个容器」，还连带停了一个容器。
// 判断只有两行，但它是「未经许可绝不动用户的容器」这条底线的唯一落点，
// 所以抽成纯函数并用测试钉住（见 autoupdate_test.go）。
func scheduledRun(cfg Settings) (ok, execute bool) {
	if cfg.CheckIntervalHours <= 0 {
		return false, false
	}
	return true, cfg.AutoApply
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

// runAutoCycle 跑一整轮：巡检 → 决策 → （可选）执行，并把这一轮的汇总返回。
//
// execute 为 false 时只巡检并更新候选，绝不动任何容器 —— 这就是
// 「总开关关掉时它只是看着，永远不会擅自动手」的落点。
//
// 返回值是调用方（后台循环）决定「下一轮多久之后跑」的唯一依据；
// 上一轮还在跑（撞上 autoBusy）时返回 nil，表示这一轮什么都没做。
func (s *Server) runAutoCycle(ctx context.Context, trigger string, execute bool) *AutoRunSummary {
	s.autoMu.Lock()
	if s.autoBusy {
		s.autoMu.Unlock()
		return nil
	}
	s.autoBusy = true
	s.autoMu.Unlock()
	defer func() {
		s.autoMu.Lock()
		s.autoBusy = false
		s.autoMu.Unlock()
	}()

	cfg := s.readSettings()
	// 双保险：除「用户在界面上手动点立即执行」外，任何执行都必须以总开关打开为前提。
	// 调用方（定时循环）已经判过一次，但那是一个容易在重构里被弄丢的判断 ——
	// 「未经许可动了用户的容器」是最不能容忍的故障，这里再兜一次。
	if execute && !cfg.AutoApply && trigger != "manual" {
		execute = false
	}
	started := time.Now()
	summary := &AutoRunSummary{
		StartedAt: started.UTC().Format(time.RFC3339),
		Trigger:   trigger,
		DryRun:    !execute,
	}

	cctx, cancel := context.WithTimeout(ctx, 20*time.Minute)
	defer cancel()

	// ---- 1. 巡检（只读，走 /distribution，绝不动容器）----
	results := s.up.CheckAll(cctx, s.excluded())
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
			"message":   "本轮巡检发现 " + itoa(summary.Available) + " 个容器有可用更新",
		})
	}

	// ---- 3. 总开关关着，或没有候选 ⇒ 到此为止 ----
	if !execute || summary.Available == 0 {
		s.finishCycle(summary, started, "仅巡检，未改动任何容器")
		return summary
	}

	// ---- 4. 执行 ----
	names := autoupdate.Names(cands)
	rows := s.up.UpdateMany(cctx, names, false, "auto")

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
	return summary
}

func (s *Server) finishCycle(summary *AutoRunSummary, started time.Time, note string) {
	summary.FinishedAt = time.Now().UTC().Format(time.RFC3339)
	summary.DurationMs = time.Since(started).Milliseconds()
	s.autoMu.Lock()
	s.autoLast = summary
	s.autoMu.Unlock()

	// 巡检（note != ""）与真更新是两种动作，运行记录用不同 kind 记，
	// 前端徽标才能分别显示「自动巡检 / 自动更新」，而不是一律 auto_update。
	// 动作名已经由徽标表达，消息里不再重复前缀。
	//
	// 只有**纯巡检**才写运行记录：真更新的轮次里 updater 对每个容器各写了一条
	// （「已更新」/「镜像已是最新，容器保持原样」），再写一条汇总就是
	// 同一次更新的第二行记录 —— 总览页看起来像执行了两遍。
	if note != "" {
		s.st.AddRunLog("auto_check", summary.Trigger, "done",
			itoa(summary.Available)+" 个容器有可用更新（"+note+"）", "")
	}
	s.bus.Publish("update", "auto_done", "success", map[string]any{
		"trigger": summary.Trigger,
		"updated": summary.Updated,
		"failed":  summary.Failed,
	})
	// 汇总通知由 UpdateMany 统一发一条（batch_update_done / batch_update_failed）。
	// 这里曾再发一条 auto_update_done，内容与那条汇总几乎一样 —— 一轮 4 台失败 2 台
	// 的自动更新会发出两条汇总 + 两条例失败共 4 条推送（2026-10-10 用户投诉），
	// 现在只保留 SSE 广播，通知不再重复。
	if summary.Updated > 0 || summary.Failed > 0 {
		s.bus.Publish("update", "auto_summary", "success", map[string]any{
			"trigger": summary.Trigger,
			"updated": summary.Updated,
			"failed":  summary.Failed,
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
