package backup

import (
	"os"
	"path/filepath"
	"testing"
)

// 「清空」必须只清自己那一棵树：清容器快照不能顺手把 compose 备份也干掉，
// 反过来也一样。两边的数据形态完全不同（json + 卷 tar / 目录 + 清单），
// 一个 RemoveAll 写错根目录就能把用户两样底牌一起抹掉。
func TestClearOnlyTouchesContainerSnapshots(t *testing.T) {
	svc := newTestService(t)
	root := svc.cfg.ContainerBackupDir()
	svc.cfg.SetMounts(map[string]string{"/volume5/docker": t.TempDir()})

	mkSnapshot(t, root, "web", "20260101-000000", "manual", 2048)
	mkSnapshot(t, root, "web", "20260102-000000", "manual", 0)
	mkSnapshot(t, root, "db", "20260103-000000", "pre_update", 1024)
	mkProjectBackup(t, svc, "immich", "20260101-000000",
		map[string]string{"compose.yaml": "services: {}\n"},
		map[string]string{"compose.yaml": "/volume5/docker/immich/compose.yaml"})

	removed, freed := svc.Clear()
	if removed != 3 {
		t.Fatalf("应删掉 3 份容器快照，实际 %d", removed)
	}
	if freed <= 0 {
		t.Error("释放字节数应当被统计（含卷数据）")
	}
	if left, err := svc.List(); err != nil || len(left) != 0 {
		t.Fatalf("容器快照应当清空，剩余 %d：%v", len(left), err)
	}
	// 卷数据目录不能变成孤儿
	if exists(filepath.Join(root, "web", "20260101-000000.volumes")) {
		t.Error("卷数据目录没被删掉")
	}
	// compose 备份必须原样还在
	pbs, err := svc.ListProjectBackups()
	if err != nil || len(pbs) != 1 {
		t.Fatalf("compose 备份不该被清，实际剩 %d：%v", len(pbs), err)
	}
	// 目录本身要留着，否则下一次备份还得靠 MkdirAll 兜底
	if !exists(root) {
		t.Error("容器快照根目录不该被删掉")
	}
}

func TestClearProjectsOnlyTouchesProjectBackups(t *testing.T) {
	svc := newTestService(t)
	root := svc.cfg.ContainerBackupDir()
	svc.cfg.SetMounts(map[string]string{"/volume5/docker": t.TempDir()})

	mkSnapshot(t, root, "web", "20260101-000000", "manual", 0)
	for _, ts := range []string{"20260101-000000", "20260102-000000"} {
		mkProjectBackup(t, svc, "immich", ts,
			map[string]string{"compose.yaml": "services: {}\n"},
			map[string]string{"compose.yaml": "/volume5/docker/immich/compose.yaml"})
	}
	mkProjectBackup(t, svc, "gitea", "20260103-000000",
		map[string]string{"compose.yaml": "services: {}\n"},
		map[string]string{"compose.yaml": "/volume5/docker/gitea/compose.yaml"})

	removed, _ := svc.ClearProjects()
	if removed != 3 {
		t.Fatalf("应删掉 3 份项目备份，实际 %d", removed)
	}
	if pbs, err := svc.ListProjectBackups(); err != nil || len(pbs) != 0 {
		t.Fatalf("项目备份应当清空，剩余 %d：%v", len(pbs), err)
	}
	// 空的项目目录要一起收掉，不然项目列表里会留一堆 0 份的空壳
	if exists(filepath.Join(svc.cfg.ProjectBackupDir(), "gitea")) {
		t.Error("空的项目目录应当被顺手删掉")
	}
	// 容器快照一份都不能少
	if items, err := svc.List(); err != nil || len(items) != 1 {
		t.Fatalf("容器快照不该被清，实际剩 %d：%v", len(items), err)
	}
}

// 空目录下点「清空」不能报错、不能把根目录删掉。
// 前端拿到的 removed 是 0，提示「没有可清空的备份」，不该出现 500。
func TestClearOnEmptyIsNoop(t *testing.T) {
	svc := newTestService(t)
	if r, _ := svc.Clear(); r != 0 {
		t.Fatalf("空目录应删 0 份，实际 %d", r)
	}
	if r, _ := svc.ClearProjects(); r != 0 {
		t.Fatalf("空目录应删 0 份，实际 %d", r)
	}
	if !exists(svc.cfg.ContainerBackupDir()) || !exists(svc.cfg.ProjectBackupDir()) {
		t.Fatal("根目录不该被删掉")
	}
	if _, err := os.Stat(svc.cfg.ContainerBackupDir()); err != nil {
		t.Fatal(err)
	}
}
