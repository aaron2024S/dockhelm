package api

import (
	"strings"
	"testing"

	"github.com/aaron2024s/dockhelm/internal/scheduler"
)

// TestValidateScheduleNeedsTargets 钉住「不选 = 0 个容器」。
//
// 这条断言记录的是一次语义更正：早先空的目标列表等于「全部容器（排除自身与排除项）」，
// 界面上只写一行小字提示，于是「随手建个任务没勾容器」就会作用到整机所有容器。
// 现在空目标必须被拒绝，只有与容器无关的动作（清理未使用镜像）允许不带目标。
func TestValidateScheduleNeedsTargets(t *testing.T) {
	cases := []struct {
		name    string
		in      scheduleReq
		wantOK  bool
		wantMsg string
	}{
		{
			name:    "需要目标却一个都没选 ⇒ 拒绝",
			in:      scheduleReq{Name: "夜间重启", Cron: "0 3 * * *", Action: "restart"},
			wantOK:  false,
			wantMsg: "至少选择一个容器",
		},
		{
			name: "需要目标且明确选了容器 ⇒ 通过",
			in: scheduleReq{Name: "夜间重启", Cron: "0 3 * * *", Action: "restart",
				Targets: []string{"jellyfin"}},
			wantOK: true,
		},
		{
			name:   "与容器无关的动作不需要目标 ⇒ 通过",
			in:     scheduleReq{Name: "清理镜像", Cron: "0 4 * * *", Action: "prune_images"},
			wantOK: true,
		},
		{
			name:    "备份同样必须指定容器",
			in:      scheduleReq{Name: "快照", Cron: "0 4 * * *", Action: "backup"},
			wantOK:  false,
			wantMsg: "至少选择一个容器",
		},
		{
			name:    "任务名不能为空",
			in:      scheduleReq{Name: "  ", Cron: "0 3 * * *", Action: "restart", Targets: []string{"a"}},
			wantOK:  false,
			wantMsg: "任务名称",
		},
		{
			name:    "不认识的动作",
			in:      scheduleReq{Name: "x", Cron: "0 3 * * *", Action: "explode"},
			wantOK:  false,
			wantMsg: "不支持的动作",
		},
		{
			name:    "非法 cron",
			in:      scheduleReq{Name: "x", Cron: "not a cron", Action: "restart", Targets: []string{"a"}},
			wantOK:  false,
			wantMsg: "cron",
		},
	}

	s := &Server{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg, ok := s.validateSchedule(tc.in)
			if ok != tc.wantOK {
				t.Fatalf("校验结果 = %v（%q），期望 %v", ok, msg, tc.wantOK)
			}
			if !tc.wantOK && !strings.Contains(msg, tc.wantMsg) {
				t.Fatalf("拒绝理由 = %q，期望包含 %q", msg, tc.wantMsg)
			}
		})
	}
}

// TestActionNeedsTargets 确认动作表里「谁需要目标」没有再被改回默认值。
func TestActionNeedsTargets(t *testing.T) {
	for _, key := range []string{"start", "stop", "restart", "update", "backup"} {
		a, ok := scheduler.ActionMap[key]
		if !ok {
			t.Fatalf("动作表里没有 %q", key)
		}
		if !a.NeedsTargets {
			t.Fatalf("动作 %q 应当要求明确的目标", key)
		}
	}
	if a, ok := scheduler.ActionMap["prune_images"]; !ok || a.NeedsTargets {
		t.Fatalf("prune_images 不需要目标，当前 NeedsTargets=%v", a.NeedsTargets)
	}
}
