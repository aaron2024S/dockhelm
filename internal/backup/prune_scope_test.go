package backup

import "testing"

// 造两份容器快照 + 两份 compose 备份，供下面三个用例共用。
func pruneScopeFixture(t *testing.T) *Service {
	t.Helper()
	svc := newTestService(t)
	root := svc.cfg.ContainerBackupDir()
	svc.cfg.SetMounts(map[string]string{"/volume5/docker": t.TempDir()})

	for _, ts := range []string{"20260101-000000", "20260102-000000"} {
		mkSnapshot(t, root, "web", ts, "manual", 0)
		mkProjectBackup(t, svc, "immich", ts,
			map[string]string{"compose.yaml": "services: {}\n"},
			map[string]string{"compose.yaml": "/volume5/docker/immich/compose.yaml"})
	}
	return svc
}

// 容器快照页点「清理过期」：只清容器快照，compose 备份一份都不许动。
func TestPruneScopeSnapshotsSkipsProjects(t *testing.T) {
	svc := pruneScopeFixture(t)

	res := svc.Prune(PruneOptions{
		KeepPerContainer:      1,
		ProjectKeepPerProject: 1, // 就算策略也开着，这一半的清理也不该碰它
		Scope:                 ScopeSnapshots,
	})
	if res.Removed != 1 {
		t.Fatalf("应只清掉 1 份容器快照，实际 %d", res.Removed)
	}
	if items, err := svc.List(); err != nil || len(items) != 1 {
		t.Fatalf("容器快照应剩 1 份，实际 %d：%v", len(items), err)
	}
	pbs, err := svc.ListProjectBackups()
	if err != nil {
		t.Fatal(err)
	}
	if len(pbs) != 2 {
		t.Fatalf("compose 备份不该被容器这一半的清理碰到，实际剩 %d", len(pbs))
	}
}

// compose 项目页点「清理过期」：只清项目备份，容器快照一份都不许动。
func TestPruneScopeProjectsSkipsSnapshots(t *testing.T) {
	svc := pruneScopeFixture(t)

	res := svc.Prune(PruneOptions{
		KeepPerContainer:      1,
		ProjectKeepPerProject: 1,
		Scope:                 ScopeProjects,
	})
	if res.Removed != 1 {
		t.Fatalf("应只清掉 1 份项目备份，实际 %d", res.Removed)
	}
	pbs, err := svc.ListProjectBackups()
	if err != nil {
		t.Fatal(err)
	}
	if len(pbs) != 1 {
		t.Fatalf("项目备份应剩 1 份，实际 %d", len(pbs))
	}
	if items, err := svc.List(); err != nil || len(items) != 2 {
		t.Fatalf("容器快照不该被项目这一半的清理碰到，实际剩 %d：%v", len(items), err)
	}
}

// 每日自动清理不带 scope（空字符串）：两半都跑，行为与升级前一致。
func TestPruneEmptyScopeClearsBoth(t *testing.T) {
	svc := pruneScopeFixture(t)

	res := svc.Prune(PruneOptions{KeepPerContainer: 1, ProjectKeepPerProject: 1})
	if res.Removed != 2 {
		t.Fatalf("应清掉 2 份（容器 1 + 项目 1），实际 %d", res.Removed)
	}
	if items, err := svc.List(); err != nil || len(items) != 1 {
		t.Fatalf("容器快照应剩 1 份，实际 %d：%v", len(items), err)
	}
	if pbs, err := svc.ListProjectBackups(); err != nil || len(pbs) != 1 {
		t.Fatalf("项目备份应剩 1 份，实际 %d：%v", len(pbs), err)
	}
}
