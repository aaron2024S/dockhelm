package api

import "testing"

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
