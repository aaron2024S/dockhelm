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
			name:        "默认状态（每 6 小时检测、自动更新关闭）⇒ 只巡检，绝不动容器",
			interval:    6,
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
// 所以非法值（含 0 与负数）必须收敛到默认的 6 小时。
func TestNormalizeInterval(t *testing.T) {
	cases := map[int]int{
		-1:  6,
		0:   6,
		2:   6,
		5:   6,
		25:  24,
		100: 24,
		1:   1,
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
}
