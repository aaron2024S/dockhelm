package updater

import (
	"fmt"
	"strings"
	"testing"
)

// 汇总通知是用户唯一能收到的一条批量结果，格式错了就是信息丢了。
// 这组用例把 buildBatchNotice 的输出钉住：事件名、级别语义、失败明细与折叠。

func TestNoticeAllSuccess(t *testing.T) {
	n := buildBatchNotice(batchStats{
		source: "manual", total: 4, updated: 3, skipped: 1,
		reclaimed: 1610612736,
	})
	if n.Event != "batch_update_done" {
		t.Fatalf("全部成功应发 batch_update_done，实际 %s", n.Event)
	}
	if n.Result != "批量更新完成" {
		t.Fatalf("手动批量文案不应带「自动」，实际 %q", n.Result)
	}
	if !strings.Contains(n.Message, "共 4 个容器：更新 3、已是最新/跳过 1、失败 0") {
		t.Fatalf("统计句不符：%q", n.Message)
	}
	if !strings.Contains(n.Message, "回收 1.5 GB") {
		t.Fatalf("应包含回收空间：%q", n.Message)
	}
	if strings.Contains(n.Message, "失败明细") {
		t.Fatalf("没有失败就不该有明细段：%q", n.Message)
	}
}

func TestNoticeAutoSourceWording(t *testing.T) {
	n := buildBatchNotice(batchStats{source: "auto", total: 4, updated: 2, skipped: 0, failed: 2,
		failures: []string{"msedge：拉取失败", "music-assistant-server：拉取失败"}})
	if n.Event != "batch_update_failed" {
		t.Fatalf("有失败应发 batch_update_failed，实际 %s", n.Event)
	}
	if n.Result != "自动更新有失败" {
		t.Fatalf("自动触发文案应为「自动更新有失败」，实际 %q", n.Result)
	}
	// 用户投诉的场景：两条失败 + 两条重复汇总。合并后必须一条里能看全。
	for _, want := range []string{"更新 2", "失败 2", "msedge", "music-assistant-server"} {
		if !strings.Contains(n.Message, want) {
			t.Fatalf("汇总里缺 %q：%q", want, n.Message)
		}
	}
	if strings.Count(n.Message, "· ") != 2 {
		t.Fatalf("失败明细应恰好两行：%q", n.Message)
	}
}

func TestNoticeManualFailureWording(t *testing.T) {
	n := buildBatchNotice(batchStats{source: "manual", total: 2, updated: 1, skipped: 0, failed: 1,
		failures: []string{"web：启动失败"}})
	if n.Result != "批量更新有失败" {
		t.Fatalf("手动批量失败文案不符：%q", n.Result)
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

func TestNoticeReusedAndNoReclaim(t *testing.T) {
	n := buildBatchNotice(batchStats{source: "manual", total: 3, updated: 3, reused: 2})
	if !strings.Contains(n.Message, "2 个容器复用了本轮已拉取的镜像") {
		t.Fatalf("复用句缺失：%q", n.Message)
	}
	if strings.Contains(n.Message, "回收") {
		t.Fatalf("没有回收就不该提回收：%q", n.Message)
	}
}

// 2026-10-11 用户反馈：推送里只写「哪些更新失败」，底层网络报错（read tcp […]）
// 一概不要 —— 那串东西没人看得懂，还把消息撑到被转发链路腰斩。
func TestNoticeShortensFailureDetail(t *testing.T) {
	n := buildBatchNotice(batchStats{source: "auto", total: 3, updated: 0, skipped: 0, failed: 3,
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
