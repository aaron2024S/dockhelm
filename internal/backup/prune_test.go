package backup

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// exists 判断路径是否存在（测试用）。
func exists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

// mkSnapshot 造一份「含卷数据」的快照：<ts>.json + <ts>.volumes/<name>.tar。
func mkSnapshot(t *testing.T, root, container, ts string, reason string, volBytes int) {
	t.Helper()
	dir := filepath.Join(root, container)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	doc := `{"_dockhelm":{"reason":"` + reason + `"},"inspect":{"Config":{"Image":"x:1"}}}`
	if err := os.WriteFile(filepath.Join(dir, ts+".json"), []byte(doc), 0o644); err != nil {
		t.Fatal(err)
	}
	if volBytes <= 0 {
		return
	}
	vdir := filepath.Join(dir, ts+".volumes")
	if err := os.MkdirAll(vdir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(vdir, container+".tar"), make([]byte, volBytes), 0o644); err != nil {
		t.Fatal(err)
	}
}

// 清理快照必须把配套的 <ts>.volumes 目录一起删掉，并把它的体积计进「释放」。
//
// 回归点：Prune 以前只删 json，卷 tar 目录原样留着 —— 列表里已经看不见它
// （json 没了），它也不会被任何一次清理碰到，几 GB 就这么永久躺着。
func TestPruneRemovesVolumeDir(t *testing.T) {
	svc := newTestService(t)
	root := svc.cfg.ContainerBackupDir()

	mkSnapshot(t, root, "web", "20260101-000000", "manual", 4096) // 旧
	mkSnapshot(t, root, "web", "20260102-000000", "manual", 4096) // 新

	res := svc.Prune(PruneOptions{KeepPerContainer: 1})
	if res.Removed != 1 {
		t.Fatalf("应清理 1 份，实际 %d 份", res.Removed)
	}
	oldJSON := filepath.Join(root, "web", "20260101-000000.json")
	oldVol := filepath.Join(root, "web", "20260101-000000.volumes")
	if exists(oldJSON) {
		t.Error("旧快照 json 没被删掉")
	}
	if exists(oldVol) {
		t.Error("旧快照的卷数据目录没被删掉（磁盘会永久泄漏）")
	}
	// 释放量必须把 4KB 卷数据算进去（json 只有几十字节）
	if res.FreedBytes < 4096 {
		t.Errorf("释放量没算卷数据：%d 字节", res.FreedBytes)
	}
	// 保留的那份连同卷数据都要在
	if !exists(filepath.Join(root, "web", "20260102-000000.json")) {
		t.Error("该保留的快照被删了")
	}
	if !exists(filepath.Join(root, "web", "20260102-000000.volumes")) {
		t.Error("该保留的卷数据被删了")
	}
}

// 保留期必须按**本地时区**解析文件名里的时间戳。
//
// 回归点：Prune 早先用 time.Parse（UTC）解析由 time.Now()（本地）生成的
// 20060102-150405。东八区下等于把每份快照都当成「晚了 8 小时」，保留期被悄悄拉长。
func TestPruneAgeUsesLocalTimezone(t *testing.T) {
	svc := newTestService(t)
	root := svc.cfg.ContainerBackupDir()

	now := time.Now()
	// 刚刚拍的（本地时间），保留期设为 1 天 —— 绝不该被删
	fresh := now.Format("20060102-150405")
	// 两天前的，该被删
	old := now.AddDate(0, 0, -2).Format("20060102-150405")
	mkSnapshot(t, root, "web", old, "manual", 0)
	mkSnapshot(t, root, "web", fresh, "manual", 0)

	res := svc.Prune(PruneOptions{MaxAgeDays: 1})
	if exists(filepath.Join(root, "web", old+".json")) {
		t.Error("两天前的快照没被清掉 —— 说明保留期算长了")
	}
	if !exists(filepath.Join(root, "web", fresh+".json")) {
		t.Error("刚刚拍的快照被误删了 —— 说明时间戳按 UTC 解析，截止点算错了")
	}
	if res.Removed != 1 {
		t.Errorf("应只清掉那一份两天前的，实际清了 %d 份", res.Removed)
	}
}

// 同一秒内的第二份快照不能覆盖第一份的卷数据目录。
//
// 回归点：json 会退让成 <ts>-1.json，卷目录却固定写 <ts>.volumes，
// 于是第二份把第一份的 tar 整个盖掉（真丢数据），而且返回的 TS 指向不存在的文件名。
func TestSameSecondSnapshotDoesNotClobberVolumeDir(t *testing.T) {
	svc := newTestService(t)
	root := svc.cfg.ContainerBackupDir()
	dir := filepath.Join(root, "web")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}

	items, _ := svc.List()
	if len(items) != 0 {
		t.Fatalf("开始前不该有快照，实际 %d 份", len(items))
	}

	// 直接验证命名规则：手动摆出两份同秒快照，走 List 看它们各自认领自己的卷目录
	ts := time.Now().Format("20060102-150405")
	mkSnapshot(t, root, "web", ts, "manual", 1024)
	mkSnapshot(t, root, "web", ts+"-1", "manual", 2048)

	got, err := svc.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("应有 2 份快照，实际 %d 份", len(got))
	}
	sizeByTS := map[string]int64{}
	for _, it := range got {
		sizeByTS[it.TS] = it.Size
	}
	// 各自算进自己的卷体积：1KB 与 2KB（json 本身几十字节）
	if sizeByTS[ts] >= sizeByTS[ts+"-1"] {
		t.Errorf("两份快照的体积应各自独立（ts=%d, ts-1=%d）—— 说明卷目录被其中一份独占了",
			sizeByTS[ts], sizeByTS[ts+"-1"])
	}
}
