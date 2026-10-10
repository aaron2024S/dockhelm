package api

import (
	"testing"
	"time"
)

// TestScheduledRun 钉住「定时循环在没有许可时绝不动手」。
//
// 这条断言的存在本身就是一次事故的记录：定时循环原先写死了 execute=true，
// 于是「从没打开过自动更新」的用户也会被它在启动 2 分钟后更新容器。
// 只要这个测试在，那种改法就过不去。
func TestScheduledRun(t *testing.T) {
	cases := []struct {
		name        string
		interval    int
		autoApply   bool
		wantOK      bool
		wantExecute bool
	}{
		{
			name:        "关掉周期检测 ⇒ 什么都不做",
			interval:    0,
			autoApply:   false,
			wantOK:      false,
			wantExecute: false,
		},
		{
			name:        "关掉周期检测、但自动更新开着 ⇒ 仍然什么都不做",
			interval:    0,
			autoApply:   true,
			wantOK:      false,
			wantExecute: false,
		},
		{
			name:        "默认状态（每 1 小时检测、自动更新关闭）⇒ 只巡检，绝不动容器",
			interval:    1,
			autoApply:   false,
			wantOK:      true,
			wantExecute: false,
		},
		{
			name:        "两个都开 ⇒ 允许更新",
			interval:    6,
			autoApply:   true,
			wantOK:      true,
			wantExecute: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotOK, gotExecute := scheduledRun(Settings{
				CheckIntervalHours: tc.interval,
				AutoApply:          tc.autoApply,
			})
			if gotOK != tc.wantOK || gotExecute != tc.wantExecute {
				t.Fatalf("scheduledRun = (ok=%v, execute=%v)，期望 (ok=%v, execute=%v)",
					gotOK, gotExecute, tc.wantOK, tc.wantExecute)
			}
		})
	}
}

// TestStartupDelay 钉住启动首轮的等待时长：checkOnStart 开了才提前巡检。
//
// 背景：这个开关曾经是死的 —— StartAutoLoop 取代了 main.go 里旧的启动巡检，
// 实现删了、设置项还留在页面上，用户开了也毫无效果（2026-10-10 用户反馈
// 「重建容器后等了很久都没检测到更新」）。
func TestStartupDelay(t *testing.T) {
	cases := []struct {
		name         string
		checkOnStart bool
		want         time.Duration
	}{
		{"开关开 ⇒ 8 秒后启动巡检", true, 8 * time.Second},
		{"开关关 ⇒ 维持 2 分钟静默", false, 2 * time.Minute},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := startupDelay(Settings{CheckOnStart: tc.checkOnStart})
			if got != tc.want {
				t.Fatalf("startupDelay = %v，期望 %v", got, tc.want)
			}
		})
	}
}

// TestNormalizeInterval 钉住「检测间隔绝不落成 0」。
//
// 背景：旧版本设置页的下拉框有「关闭」档，会在库里存下 0。周期检测是自动更新的
// 唯一触发源 —— 一旦读到 0，自动更新就永远不会发生，而设置页还显示「已开启」。
// 所以非法值（含 0 与负数）必须收敛到默认的 1 小时。
func TestNormalizeInterval(t *testing.T) {
	cases := map[int]int{
		-1:  1,
		0:   1,
		5:   1,
		25:  24,
		100: 24,
		1:   1,
		2:   2,
		3:   3,
		6:   6,
		12:  12,
		24:  24,
	}
	for in, want := range cases {
		if got := normalizeInterval(in); got != want {
			t.Errorf("normalizeInterval(%d) = %d，期望 %d", in, got, want)
		}
	}
	// 可选值里绝不能出现 0，且每个可选值都必须被原样接受。
	for _, c := range checkIntervalChoices {
		if c <= 0 {
			t.Fatalf("检测频率可选值出现了 <= 0 的档：%d（0 会让自动更新永久失效）", c)
		}
		if got := normalizeInterval(c); got != c {
			t.Errorf("normalizeInterval(%d) = %d，应当原样返回", c, got)
		}
	}
	// 2 小时档是按用户要求新增的：删掉它前端下拉会少一档，等于功能回退。
	var has2 bool
	for _, c := range checkIntervalChoices {
		if c == 2 {
			has2 = true
		}
	}
	if !has2 {
		t.Fatalf("检测频率可选值里必须有 2 小时档，实际是 %v", checkIntervalChoices)
	}
}

// TestWakeDecision 钉住「哪两类设置改动会把下一轮提前」。
//
// 只有「刚打开自动更新」与「把检测间隔调小」才提前：前者用户显然在等着看它动，
// 后者改完还按旧的长间隔等，等于这次改动没生效。其余改动一律按新间隔重新计时 ——
// 否则每改一次并发度就顺手多打一轮镜像仓库。
func TestWakeDecision(t *testing.T) {
	base := Settings{CheckIntervalHours: 6}
	cases := []struct {
		name string
		prev Settings
		next Settings
		want wakeAction
	}{
		{"自动更新由关变开 ⇒ 提前", base, Settings{CheckIntervalHours: 6, AutoApply: true}, wakeSooner},
		{"自动更新由开变关 ⇒ 不提前", Settings{CheckIntervalHours: 6, AutoApply: true}, base, wakeRestart},
		{"间隔调小（6→1）⇒ 提前", base, Settings{CheckIntervalHours: 1}, wakeSooner},
		{"间隔调大（1→6）⇒ 不提前", Settings{CheckIntervalHours: 1}, Settings{CheckIntervalHours: 6}, wakeRestart},
		{"间隔没变 ⇒ 不提前", base, Settings{CheckIntervalHours: 6, Concurrency: 8}, wakeRestart},
		{"只改排除列表 ⇒ 不提前", base, Settings{CheckIntervalHours: 6, Exclude: []string{"web"}}, wakeRestart},
		{"间隔是非法值 0 ⇒ 不参与比较、不提前", base, Settings{CheckIntervalHours: 0}, wakeRestart},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := wakeDecision(tc.prev, tc.next); got != tc.want {
				t.Fatalf("wakeDecision = %v，期望 %v", got, tc.want)
			}
		})
	}
}

// TestFirstRoundDelay 钉住「启动首轮巡检之后，第二轮等多久」。
//
// 这是「重启后自动更新要干等一个完整周期」这个体验问题的落点：
// 首轮只读巡检查到了更新、且总开关开着时，第二轮提前到启动冷却期结束；
// 其余情况一律维持「等满一个周期」的老行为。
func TestFirstRoundDelay(t *testing.T) {
	boot := time.Date(2026, 10, 10, 12, 0, 0, 0, time.UTC)
	// 启动 8 秒后：正是首轮只读巡检跑完、要排第二轮的时刻。
	after := boot.Add(8 * time.Second)
	hour := time.Hour
	on := Settings{CheckOnStart: true, AutoApply: true}

	cases := []struct {
		name      string
		now       time.Time
		interval  time.Duration
		cfg       Settings
		available int
		want      time.Duration
	}{
		{"查到更新 + 总开关开 ⇒ 提前到启动冷却期结束", after, hour, on, 3, autoBootCooldown - 8*time.Second},
		{"查到更新但总开关关着 ⇒ 等满一个周期", after, hour, Settings{CheckOnStart: true}, 3, hour},
		{"总开关开着但没查到更新 ⇒ 等满一个周期", after, hour, on, 0, hour},
		{"没开启动巡检 ⇒ 等满一个周期", after, hour, Settings{AutoApply: true}, 3, hour},
		{"首轮巡检耗时超过冷却期 ⇒ 至少留 1 秒", boot.Add(20 * time.Minute), hour, on, 3, time.Second},
		{"冷却期比间隔还长 ⇒ 不超出间隔", after, 5 * time.Second, on, 3, 5 * time.Second},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := firstRoundDelay(tc.now, boot, tc.interval, tc.cfg, tc.available); got != tc.want {
				t.Fatalf("firstRoundDelay = %v，期望 %v", got, tc.want)
			}
		})
	}
}

// TestWakeRoundDelay 钉住「保存设置之后，下一轮等多久」。
func TestWakeRoundDelay(t *testing.T) {
	boot := time.Date(2026, 10, 10, 12, 0, 0, 0, time.UTC)
	hour := time.Hour
	// 早就过了启动冷却期（用户运行中改设置就是这种情形）。
	settled := boot.Add(3 * time.Hour)
	// 刚启动 8 秒 —— 此时改设置仍要尊重启动冷却期。
	fresh := boot.Add(8 * time.Second)
	// 已排定的更早一轮：boot+15min（正是启动冷却期提前安排出来的那种排期）。
	soonerRound := boot.Add(autoBootCooldown)

	cases := []struct {
		name     string
		now      time.Time
		interval time.Duration
		pending  time.Time
		act      wakeAction
		want     time.Duration
	}{
		{"普通改动 ⇒ 按新间隔重新计时", settled, hour, time.Time{}, wakeRestart, hour},
		{"值得提前的改动 ⇒ 30 秒后跑", settled, hour, time.Time{}, wakeSooner, 30 * time.Second},
		{"刚启动就改 ⇒ 仍等冷却期结束", fresh, hour, time.Time{}, wakeSooner, autoBootCooldown - 8*time.Second},
		{"提前也不超过间隔", settled, 10 * time.Second, time.Time{}, wakeSooner, 10 * time.Second},
		{"普通改动不把已排定的更早一轮往后推", fresh, hour, soonerRound, wakeRestart, autoBootCooldown - 8*time.Second},
		{"已排定的轮次比新间隔还晚 ⇒ 照新间隔", fresh, 5 * time.Minute, soonerRound, wakeRestart, 5 * time.Minute},
		{"已排定的轮次已过期 ⇒ 照新间隔", boot.Add(20 * time.Minute), hour, soonerRound, wakeRestart, hour},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := wakeRoundDelay(tc.now, boot, tc.interval, tc.pending, tc.act); got != tc.want {
				t.Fatalf("wakeRoundDelay = %v，期望 %v", got, tc.want)
			}
		})
	}
}
