package backup

import (
	"errors"
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

// 2026-10-12 用户反馈：backup_failed 的推送正文直塞 err.Error()，用户看到的
// 是一句英文 docker 报错。这里钉住「翻成人话」的映射，以及「中文错误原样放行」
// （compose 项目那几条错误本来就是我们自己写的中文，不能被模板覆盖掉）。
func TestSnapshotFailReason(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Error response from daemon: No such container: web", "容器不存在或已被删除"},
		{"mkdir /data/backup: permission denied", "没有写入备份目录的权限"},
		{"write /data/backup/x.json: no space left on device", "磁盘空间不足，无法写入快照"},
		{"Cannot connect to the Docker daemon at unix:///var/run/docker.sock", "无法连接 Docker 守护进程"},
		{"unexpected EOF", "读取容器配置失败"},
		// 中文错误原样保留
		{"项目 blog 的 compose 文件在 Dockhelm 容器里读不到，请先把它的宿主目录挂进来",
			"项目 blog 的 compose 文件在 Dockhelm 容器里读不到，请先把它的宿主目录挂进来"},
	}
	for _, c := range cases {
		if got := SnapshotFailReason(errors.New(c.in)); got != c.want {
			t.Errorf("SnapshotFailReason(%q) = %q，期望 %q", c.in, got, c.want)
		}
	}
	if got := SnapshotFailReason(nil); got != "备份失败" {
		t.Errorf("nil 错误应返回「备份失败」，实际 %q", got)
	}
}
