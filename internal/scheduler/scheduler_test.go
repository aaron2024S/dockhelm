package scheduler

import (
	"context"
	"fmt"
	"strings"
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
	out := r.Execute(context.Background(), store.Schedule{
		Name:   "老任务",
		Cron:   "0 3 * * *",
		Action: "restart",
	})
	if out.OK {
		t.Fatalf("空目标不应算执行成功：%q", out.Brief)
	}
	if out.Brief != "任务没有选择任何容器，请检查任务配置" {
		t.Fatalf("运行期文案 = %q", out.Brief)
	}
	// 简报与正文相同，且不该带原始报错详情。
	if out.Full != out.Brief || out.Detail != "" {
		t.Fatalf("空目标分支的产出异常：full=%q detail=%q", out.Full, out.Detail)
	}
}

// TestSummarize 钉住任务回执的正文结构：统计行 + 已动作名单 + 跳过原因 + 失败名单。
// 以前只有一句「执行完成：成功 N 个，跳过 M 个」，看不出是哪几台。
func TestSummarize(t *testing.T) {
	got := summarize("stop", 3, 2, 1,
		[]string{"qbittorrent", "jellyfin"}, []string{"immich-server"}, nil)
	want := "共 3 个容器：成功 2 个、跳过 1 个\n" +
		"已停止：jellyfin、qbittorrent\n" +
		"跳过（本就已停止）：immich-server"
	if got != want {
		t.Fatalf("停容器回执：\n got=%q\nwant=%q", got, want)
	}

	// 有失败：统计行带失败数，末行点名失败容器。
	got = summarize("start", 2, 1, 0, []string{"a"}, nil, []string{"b"})
	want = "共 2 个容器：成功 1 个、跳过 0 个、失败 1 个\n已启动：a\n失败：b"
	if got != want {
		t.Fatalf("启动容器回执：\n got=%q\nwant=%q", got, want)
	}

	// 重启没有「跳过」语义，跳过名单恒空。
	if got := summarize("restart", 2, 2, 0, []string{"b", "a"}, nil, nil); got !=
		"共 2 个容器：成功 2 个、跳过 0 个\n已重启：a、b" {
		t.Fatalf("重启回执 = %q", got)
	}
}

// TestFoldNames 名单排序 + 超限折叠。
func TestFoldNames(t *testing.T) {
	if got := foldNames([]string{"c", "a", "b"}); got != "a、b、c" {
		t.Fatalf("名单应排序后输出：%q", got)
	}
	names := make([]string, 0, maxScheduleNames+3)
	for i := 0; i < maxScheduleNames+3; i++ {
		names = append(names, fmt.Sprintf("c%02d", i))
	}
	got := foldNames(names)
	if !strings.Contains(got, fmt.Sprintf("等 %d 个", len(names))) {
		t.Fatalf("超限名单应折叠：%q", got)
	}
}
