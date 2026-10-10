package scheduler

import (
	"context"
	"testing"

	"github.com/aaron2024s/dockhelm/internal/store"
)

// TestExecuteSkipsEmptyTargets 钉住「空目标任务」在运行期的文案。
//
// 老版本存下来的任务可能是空目标（当时空目标 = 全部容器），现在的语义是拒绝执行。
// 文案只说结论 —— 早期版本带一串「『不选』不再等于全部容器，请编辑任务并勾选目标」，
// 在「最近执行记录」里被截断成「任务没有选择任…」，用户根本读不完。
// 这条断言还顺带保证 Execute 在空目标分支上不会去碰 Runner 里的 dc/st 等依赖
// （零值 Runner 也必须能安全走到这里）。
func TestExecuteSkipsEmptyTargets(t *testing.T) {
	r := &Runner{}
	msg, ok := r.Execute(context.Background(), store.Schedule{
		Name:   "老任务",
		Cron:   "0 3 * * *",
		Action: "restart",
	})
	if ok {
		t.Fatalf("空目标不应算执行成功：%q", msg)
	}
	if msg != "任务没有选择任何容器，已跳过" {
		t.Fatalf("运行期文案 = %q", msg)
	}
}
