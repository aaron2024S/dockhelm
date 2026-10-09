package backup

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// 快照列表里的「时间」必须来自快照自己，而不是文件 mtime。
//
// 理由：数据目录一旦被复制 / rsync / 从压缩包里解出来，mtime 就变成「复制那一刻」，
// 界面上会出现「几十份快照全是刚刚」的假象，而快照里记的才是真正拍下来的时刻。
// 优先级：_dockhelm.snapshotAt（权威，RFC3339）> 文件名时间戳（本地时区）> mtime。
func TestListUsesSnapshotTimeNotMtime(t *testing.T) {
	svc := newTestService(t)
	root := svc.cfg.ContainerBackupDir()

	write := func(container, ts string, doc map[string]any, mtime time.Time) {
		t.Helper()
		dir := filepath.Join(root, container)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		raw, err := json.Marshal(doc)
		if err != nil {
			t.Fatal(err)
		}
		p := filepath.Join(dir, ts+".json")
		if err := os.WriteFile(p, raw, 0o644); err != nil {
			t.Fatal(err)
		}
		if err := os.Chtimes(p, mtime, mtime); err != nil {
			t.Fatal(err)
		}
	}

	// 一眼能认出来的假 mtime：它绝对不该出现在结果里
	fakeMtime := time.Date(2020, 1, 2, 3, 4, 5, 0, time.UTC)

	write("redis", "20260101-120000", map[string]any{
		"_dockhelm": map[string]any{"snapshotAt": "2026-01-01T04:00:00Z", "reason": "pre_update"},
		"inspect": map[string]any{
			"Config": map[string]any{"Image": "redis:7.4-alpine"},
			"State":  map[string]any{"Running": true},
		},
	}, fakeMtime)

	// 老快照 / 手工导入的快照可能没有 snapshotAt：退回到文件名里的时间戳
	write("redis", "20261009-094200", map[string]any{
		"_dockhelm": map[string]any{"reason": "manual"},
		"inspect":   map[string]any{"Config": map[string]any{"Image": "redis:7-alpine"}},
	}, fakeMtime)

	items, err := svc.List()
	if err != nil {
		t.Fatalf("List 失败：%v", err)
	}
	if len(items) != 2 {
		t.Fatalf("应有 2 份快照，实际 %d 份", len(items))
	}
	// 同一容器内按时间戳倒序
	if items[0].TS != "20261009-094200" {
		t.Errorf("应按时间倒序，拿到 %q", items[0].TS)
	}

	byTS := map[string]SnapshotItem{}
	for _, it := range items {
		byTS[it.TS] = it
	}

	a := byTS["20260101-120000"]
	if a.Created != "2026-01-01T04:00:00Z" {
		t.Errorf("有 snapshotAt 时应以它为准，得到 %q", a.Created)
	}
	if a.Image != "redis:7.4-alpine" || !a.Running {
		t.Errorf("镜像 / 运行状态应从嵌套 inspect 里取，得到 image=%q running=%v", a.Image, a.Running)
	}
	if a.Reason != "pre_update" {
		t.Errorf("来源应为 pre_update，得到 %q", a.Reason)
	}

	b := byTS["20261009-094200"]
	want := time.Date(2026, 10, 9, 9, 42, 0, 0, time.Local).UTC().Format(time.RFC3339)
	if b.Created != want {
		t.Errorf("没有 snapshotAt 时按本地时区解析文件名时间戳：want %q got %q", want, b.Created)
	}
	if b.Running {
		t.Error("快照里没有 State.Running 时应视为已停止")
	}
	if b.Created == fakeMtime.Format(time.RFC3339) {
		t.Error("结果里出现了文件 mtime —— 说明退回到了最不可靠的那一档")
	}
}
