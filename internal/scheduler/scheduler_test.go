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
	want := "共 3 个容器，成功 2 个、跳过 1 个\n" +
		"已停止：jellyfin、qbittorrent\n" +
		"跳过（本就已停止）：immich-server"
	if got != want {
		t.Fatalf("停容器回执：\n got=%q\nwant=%q", got, want)
	}

	// 有失败：统计行带失败数，末行点名失败容器。
	got = summarize("start", 2, 1, 0, []string{"a"}, nil, []string{"b"})
	want = "共 2 个容器，成功 1 个、跳过 0 个、失败 1 个\n已启动：a\n失败：b"
	if got != want {
		t.Fatalf("启动容器回执：\n got=%q\nwant=%q", got, want)
	}

	// 重启没有「跳过」语义，跳过名单恒空。
	if got := summarize("restart", 2, 2, 0, []string{"b", "a"}, nil, nil); got !=
		"共 2 个容器，成功 2 个、跳过 0 个\n已重启：a、b" {
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

// TestBrief 钉住运行记录/「上次结果」的一行简报：动作 + 统计（+ 小任务点名）。
//
// 2026-10-13 之前只有统计（「共 1 个容器：成功 1 个、跳过 0 个」），看不出这条
// 记录是启动还是停止、是备份还是更新 —— 用户要求「一看就知道哪个计划任务执行了
// 什么」。统计行的措辞与推送正文（summarize）共用 countLine，两处不许飘。
func TestBrief(t *testing.T) {
	// 小任务点名：目标 ≤ maxBriefNames 且无失败 ⇒ 连容器名一起给。
	if got := brief("restart", 2, 2, 0, 0, []string{"shell-town", "qbittorrent"}); got !=
		"重启容器：共 2 个容器，成功 2 个、跳过 0 个（qbittorrent、shell-town）" {
		t.Fatalf("重启简报 = %q", got)
	}

	// 有失败：不点名（点名会把失败原因挤掉，原因在推送正文与记录详情里）。
	got := brief("stop", 3, 1, 0, 2, []string{"a"})
	if !strings.HasPrefix(got, "停止容器：共 3 个容器，成功 1 个、跳过 0 个、失败 2 个") {
		t.Fatalf("有失败时应只报数：%q", got)
	}
	if strings.Contains(got, "（") {
		t.Fatalf("有失败时不该再点名：%q", got)
	}

	// 目标多：只报数，名字靠推送正文与记录详情。
	got = brief("update", 12, 12, 0, 0, []string{"c01", "c02"})
	if strings.Contains(got, "（") {
		t.Fatalf("目标超过 %d 个应只报数：%q", maxBriefNames, got)
	}

	// 备份动作走同一套：说明里有动作，备份记录不再只有一句「全量备份 N 个容器」。
	if got := brief("backup", 1, 1, 0, 0, []string{"nginx"}); got !=
		"备份配置快照：共 1 个容器，成功 1 个、跳过 0 个（nginx）" {
		t.Fatalf("备份简报 = %q", got)
	}

	// 未知动作原样回显，别把动作吞掉。
	if got := brief("frobnicate", 1, 1, 0, 0, nil); !strings.Contains(got, "frobnicate") {
		t.Fatalf("未知动作不该被吞掉：%q", got)
	}

	if ActionLabel("prune_images") != "清理未使用镜像" {
		t.Fatalf("清理动作的标签 = %q", ActionLabel("prune_images"))
	}
}
