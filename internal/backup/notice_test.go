package backup

import (
	"fmt"
	"strings"
	"testing"
)

// TestBatchNotice 钉住全量备份的事件升级规则与文案。
//
// 回归点：以前这里无论失败几台都发 backup_success（普通级、出厂默认**不订阅**），
// 「24 台只成功 1 台」也只是一条默认收不到的「备份完成」——整体失败却悄无声息。
func TestBatchNotice(t *testing.T) {
	// 全部成功 → 普通级 backup_success，状态 success。
	ev, result, status, msg := batchNotice(3, 3, nil)
	if ev != "backup_success" || result != "全量备份完成" || status != "success" {
		t.Fatalf("全成功应走 backup_success：ev=%s result=%s status=%s", ev, result, status)
	}
	if msg != "共 3 个容器：成功 3、失败 0" {
		t.Fatalf("全成功正文 = %q", msg)
	}

	// 哪怕只有一台失败，也必须升级成 backup_failed（紧急、默认订阅）。
	ev, result, status, msg = batchNotice(3, 2, []string{"qbittorrent"})
	if ev != "backup_failed" || status != "failed" {
		t.Fatalf("有失败必须升级成 backup_failed：ev=%s status=%s", ev, status)
	}
	if result != "全量备份有失败" {
		t.Fatalf("结果短语 = %q", result)
	}
	if !strings.Contains(msg, "共 3 个容器：成功 2、失败 1") ||
		!strings.Contains(msg, "失败：qbittorrent") {
		t.Fatalf("有失败正文应含统计与失败名单：%q", msg)
	}

	// 24 台只成 1 台 —— 最坏的场景也必须报失败。
	ev, _, _, msg = batchNotice(24, 1, []string{"a", "b", "c"})
	if ev != "backup_failed" {
		t.Fatalf("大面积失败必须报 backup_failed，实际 %s", ev)
	}
	if !strings.Contains(msg, "成功 1、失败 3") {
		t.Fatalf("正文统计不符：%q", msg)
	}

	// 失败名单排序 + 超限折叠。
	many := make([]string, 0, maxBackupNames+4)
	for i := maxBackupNames + 3; i >= 0; i-- {
		many = append(many, fmt.Sprintf("c%02d", i))
	}
	_, _, _, msg = batchNotice(len(many), 0, many)
	if !strings.Contains(msg, fmt.Sprintf("等 %d 个", len(many))) {
		t.Fatalf("超限名单应折叠：%q", msg)
	}
	if strings.Index(msg, "c00") > strings.Index(msg, "c01") {
		t.Fatalf("失败名单应先排序：%q", msg)
	}
}
