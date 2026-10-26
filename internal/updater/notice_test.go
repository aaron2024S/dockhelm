package updater

import (
	"fmt"
	"strings"
	"testing"
)

// 一轮更新的结果通知是用户唯一能收到的一条更新回执，格式错了就是信息丢了。
// 这组用例把 buildBatchNotice 的输出钉住：一台与多台文案完全同构、更新/失败名单、
// 事件名与级别语义、明细折叠、底层报错不外泄。

// 2026-10-12 用户要求：不要按台数换文案，就统一写「共几个容器、更新了几个、
// 更新的容器是哪些、失败了几个、失败的是哪些」。下面两条分别钉住 1 台与多台
// 的输出 —— 除了数字与名单，结构必须逐行一致。
func TestNoticeSingleContainerUsesSameShape(t *testing.T) {
	one := buildBatchNotice(batchStats{
		source: "auto", total: 1, updated: 1, updatedNames: []string{"jellyfin"},
	})
	if one.Event != "batch_update_done" {
		t.Fatalf("事件应为 batch_update_done，实际 %s", one.Event)
	}
	if one.Result != "自动更新完成" {
		t.Fatalf("结果短语应为「自动更新完成」，实际 %q", one.Result)
	}
	lines := strings.Split(one.Message, "\n")
	if len(lines) != 2 {
		t.Fatalf("一台的正文应为「统计 + 更新名单」两行，实际 %d 行：%q", len(lines), one.Message)
	}
	if lines[0] != "共 1 个容器：更新 1 个、失败 0 个、已是最新/跳过 0 个" {
		t.Fatalf("统计行不符：%q", lines[0])
	}
	if lines[1] != "更新的容器：jellyfin" {
		t.Fatalf("更新名单行不符：%q", lines[1])
	}

	// 多台：同样的两行结构，只是数字与名单不同。
	many := buildBatchNotice(batchStats{
		source: "auto", total: 3, updated: 2, skipped: 1,
		updatedNames: []string{"jellyfin", "immich-server"},
	})
	mlines := strings.Split(many.Message, "\n")
	if len(mlines) != 2 {
		t.Fatalf("多台正文行数应与单台一致（2 行），实际 %d：%q", len(mlines), many.Message)
	}
	if mlines[0] != "共 3 个容器：更新 2 个、失败 0 个、已是最新/跳过 1 个" {
		t.Fatalf("多台统计行不符：%q", mlines[0])
	}
	if mlines[1] != "更新的容器：immich-server、jellyfin" {
		t.Fatalf("多台更新名单不符（应按名字排序）：%q", mlines[1])
	}
	if strings.Contains(many.Message, "批量") {
		t.Fatalf("文案不该再出现「批量」：%q", many.Message)
	}
}

func TestNoticeSingleContainerFailure(t *testing.T) {
	n := buildBatchNotice(batchStats{
		source: "auto", total: 1, failed: 1,
		failures: []string{"shell-town：拉取失败：pull failed: read tcp [240e::1]:443: connection reset"},
	})
	if n.Event != "batch_update_failed" {
		t.Fatalf("事件应为 batch_update_failed，实际 %s", n.Event)
	}
	if n.Result != "自动更新有失败" {
		t.Fatalf("结果短语应为「自动更新有失败」，实际 %q", n.Result)
	}
	if !strings.Contains(n.Message, "共 1 个容器：更新 0 个、失败 1 个、已是最新/跳过 0 个") {
		t.Fatalf("统计行不符：%q", n.Message)
	}
	if !strings.Contains(n.Message, "失败的容器：") || !strings.Contains(n.Message, "· shell-town：拉取失败") {
		t.Fatalf("失败名单不符：%q", n.Message)
	}
	if strings.Contains(n.Message, "更新的容器") {
		t.Fatalf("没有成功的容器就不该有更新名单行：%q", n.Message)
	}
	for _, banned := range []string{"read tcp", "pull failed", "240e", "批量"} {
		if strings.Contains(n.Message, banned) {
			t.Fatalf("推送正文不该出现 %q：%q", banned, n.Message)
		}
	}
}

func TestNoticeAllSuccess(t *testing.T) {
	n := buildBatchNotice(batchStats{
		source: "manual", total: 4, updated: 3, skipped: 1,
		updatedNames: []string{"a", "b", "c"},
		reclaimed:    1610612736,
	})
	if n.Event != "batch_update_done" {
		t.Fatalf("全部成功应发 batch_update_done，实际 %s", n.Event)
	}
	if n.Result != "更新完成" {
		t.Fatalf("手动轮次文案不应带「自动」或「批量」，实际 %q", n.Result)
	}
	if !strings.Contains(n.Message, "共 4 个容器：更新 3 个、失败 0 个、已是最新/跳过 1 个") {
		t.Fatalf("统计句不符：%q", n.Message)
	}
	if !strings.Contains(n.Message, "更新的容器：a、b、c") {
		t.Fatalf("更新名单缺失：%q", n.Message)
	}
	if !strings.Contains(n.Message, "回收 1.5 GB") {
		t.Fatalf("应包含回收空间：%q", n.Message)
	}
	if strings.Contains(n.Message, "失败的容器") {
		t.Fatalf("没有失败就不该有失败段：%q", n.Message)
	}
}

func TestNoticeAutoSourceWording(t *testing.T) {
	n := buildBatchNotice(batchStats{source: "auto", total: 4, updated: 2, failed: 2,
		updatedNames: []string{"a", "b"},
		failures:     []string{"msedge：拉取失败", "music-assistant-server：拉取失败"}})
	if n.Event != "batch_update_failed" {
		t.Fatalf("有失败应发 batch_update_failed，实际 %s", n.Event)
	}
	if n.Result != "自动更新有失败" {
		t.Fatalf("自动触发文案应为「自动更新有失败」，实际 %q", n.Result)
	}
	// 用户投诉的场景：两条失败 + 两条重复汇总。合并后必须一条里能看全。
	for _, want := range []string{"更新 2 个", "失败 2 个", "更新的容器：a、b",
		"· msedge：拉取失败", "· music-assistant-server：拉取失败"} {
		if !strings.Contains(n.Message, want) {
			t.Fatalf("汇总里缺 %q：%q", want, n.Message)
		}
	}
	if strings.Count(n.Message, "· ") != 2 {
		t.Fatalf("失败名单应恰好两行：%q", n.Message)
	}
}

func TestNoticeManualFailureWording(t *testing.T) {
	n := buildBatchNotice(batchStats{source: "manual", total: 2, updated: 1, failed: 1,
		updatedNames: []string{"web"}, failures: []string{"db：启动失败"}})
	if n.Result != "更新有失败" {
		t.Fatalf("手动失败文案不符：%q", n.Result)
	}
}

func TestNoticeFailureFold(t *testing.T) {
	failures := make([]string, 0, maxFailureLines+3)
	for i := 0; i < maxFailureLines+3; i++ {
		failures = append(failures, fmt.Sprintf("c%d：拉取超时", i))
	}
	n := buildBatchNotice(batchStats{source: "auto", total: len(failures), failed: len(failures),
		failures: failures})
	if got := strings.Count(n.Message, "· "); got != maxFailureLines {
		t.Fatalf("明细应折叠到 %d 行，实际 %d", maxFailureLines, got)
	}
	if !strings.Contains(n.Message, "其余 3 台略") {
		t.Fatalf("折叠提示缺失：%q", n.Message)
	}
	// 运行记录 detail 不折叠 —— 面板里空间管够，全量留底。
	if got := strings.Count(n.Detail, "拉取超时"); got != len(failures) {
		t.Fatalf("Detail 应保留全部失败（%d），实际 %d", len(failures), got)
	}
}

// 更新名单也要有上限：一轮几十台时名单本身就能把推送撑爆。
func TestNoticeNameListFold(t *testing.T) {
	names := make([]string, 0, maxNameList+5)
	for i := 0; i < maxNameList+5; i++ {
		names = append(names, fmt.Sprintf("c%d", i))
	}
	n := buildBatchNotice(batchStats{source: "manual", total: len(names), updated: len(names), updatedNames: names})
	if !strings.Contains(n.Message, fmt.Sprintf("等 %d 个", len(names))) {
		t.Fatalf("超长名单应折叠成「等 N 个」：%q", n.Message)
	}
	// 名单用一个「、」多一个名字：折叠后恰好点名 maxNameList 个。
	listLine := ""
	for _, l := range strings.Split(n.Message, "\n") {
		if strings.HasPrefix(l, "更新的容器：") {
			listLine = l
		}
	}
	if listLine == "" {
		t.Fatalf("缺更新名单行：%q", n.Message)
	}
	if got := strings.Count(listLine, "、"); got != maxNameList-1 {
		t.Fatalf("折叠后应恰好点名 %d 个，实际 %d 个：%q", maxNameList, got+1, listLine)
	}
	// 运行记录里不折叠，方便面板点名。
	if !strings.Contains(n.Detail, fmt.Sprintf("c%d", len(names)-1)) {
		t.Fatalf("Detail 应保留全部名单：%q", n.Detail)
	}
}

// 2026-10-11 用户反馈：推送里只写「哪些更新失败」，底层网络报错（read tcp […]）
// 一概不要 —— 那串东西没人看得懂，还把消息撑到被转发链路腰斩。
func TestNoticeShortensFailureDetail(t *testing.T) {
	n := buildBatchNotice(batchStats{source: "auto", total: 3, failed: 3,
		failures: []string{
			"music-assistant-server：拉取失败：pull failed: read tcp [240e:355:482f:5680:4006:3465:8422:e]:42218->[2606:4700:c800:2::154]:443: read: connection reset by peer",
			"msedge：拉取失败：context deadline exceeded",
			"rsshub：读取容器信息失败：Get \"http://%2Fvar%2Frun%2Fdocker.sock/v1.43/containers/rsshub/json\": dial unix /var/run/docker.sock: connect: connection refused",
		}})
	// 正文里只留动作短语，绝不出现底层报错片段。
	for _, banned := range []string{"read tcp", "pull failed", "Get \"", "dial unix", "connection reset"} {
		if strings.Contains(n.Message, banned) {
			t.Fatalf("推送正文不该出现底层报错 %q：%q", banned, n.Message)
		}
	}
	for _, want := range []string{
		"· music-assistant-server：拉取失败",
		"· msedge：拉取失败",
		"· rsshub：读取容器信息失败",
	} {
		if !strings.Contains(n.Message, want) {
			t.Fatalf("明细行缺失 %q：%q", want, n.Message)
		}
	}
	// 运行记录 Detail 仍留全量 —— 面板里要能查到根因。
	if !strings.Contains(n.Detail, "connection reset by peer") || !strings.Contains(n.Detail, "dial unix") {
		t.Fatalf("Detail 应保留完整报错：%q", n.Detail)
	}
}

// 没有冒号结构的错误原文也要压短，不能整段照搬。
func TestShortFailureLineNoColon(t *testing.T) {
	long := strings.Repeat("网络异常", 50) // 200 字
	got := shortFailureLine("web：" + long)
	if got != "web："+truncateRunes(long, 40) {
		t.Fatalf("无结构错误应截到 40 字：%q", got)
	}
}

// 标题里点名谁：2026-10-26 用户反馈「四条更新失败长得一模一样，都写 1 个容器」。
// 标题的「· xxx」来自 notice.Container，必须是容器名，不能是数量。
func TestNoticeContainerNamesTheSubject(t *testing.T) {
	cases := []struct {
		name string
		in   batchStats
		want string
	}{
		{"单台失败点名该容器", batchStats{total: 1, failed: 1, failures: []string{"jellyfin：拉取失败"}}, "jellyfin"},
		{"多台失败首台加等N个", batchStats{total: 4, failed: 3, failures: []string{"a：拉取失败", "b：拉取失败", "c：拉取失败"}}, "a 等 3 个"},
		{"单台成功点名该容器", batchStats{total: 1, updated: 1, updatedNames: []string{"nginx"}}, "nginx"},
		{"多台成功首台加等N个", batchStats{total: 3, updated: 3, updatedNames: []string{"a", "b", "c"}}, "a 等 3 个"},
		{"全部跳过才是数量", batchStats{total: 3, skipped: 3}, "3 个容器"},
	}
	for _, c := range cases {
		if got := buildBatchNotice(c.in).Container; got != c.want {
			t.Errorf("%s：Container 应为 %q，实际 %q", c.name, c.want, got)
		}
	}
	// 名字里带全角冒号以外的冒号不能把名字截错
	n := buildBatchNotice(batchStats{
		total: 1, failed: 1,
		failures: []string{"my-container：拉取失败：pull failed: timeout"},
	})
	if n.Container != "my-container" {
		t.Fatalf("应从第一个全角冒号处切名字，实际 %q", n.Container)
	}
}
