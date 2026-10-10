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
