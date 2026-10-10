package backup

import (
	"os"
	"path/filepath"
	"testing"
)

// 项目备份的还原回路：备份里记的是**宿主路径**，还原时按当前挂载映射换算回容器内路径再写。
//
// 这是整套「compose yaml 备份」的地基 —— 映射算错就会写到别的地方去，
// 而用户的 yaml 被写坏是没有后悔药的。
func TestRestoreProjectWritesBackToHostPath(t *testing.T) {
	svc := newTestService(t)
	// 宿主 /volume5/docker 映射到容器内的临时目录
	root := t.TempDir()
	svc.cfg.SetMounts(map[string]string{"/volume5/docker": root})

	hostYAML := "/volume5/docker/immich/compose.yaml"
	if err := os.MkdirAll(filepath.Join(root, "immich"), 0o755); err != nil {
		t.Fatal(err)
	}
	// 现网上的 yaml 是「改坏了」的版本
	if err := os.WriteFile(filepath.Join(root, "immich", "compose.yaml"), []byte("broken: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	mkProjectBackup(t, svc, "immich", "20261010-120000",
		map[string]string{"compose.yaml": "services:\n  immich:\n    image: x:1\n"},
		map[string]string{"compose.yaml": hostYAML})

	res := svc.RestoreProject(nil, "immich", "20261010-120000", false)
	if !res.OK {
		t.Fatalf("还原应当成功，实际：%s · %+v", res.Message, res.Files)
	}
	got, err := os.ReadFile(filepath.Join(root, "immich", "compose.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "services:\n  immich:\n    image: x:1\n" {
		t.Fatalf("写回去的内容不对：%q", got)
	}
}

// 路径看不见（宿主目录没挂进来）时必须**拒绝并如实说明**，绝不能报成功。
//
// 这一条是「不可见」这个状态唯一真正的危险所在：报成功而没写回去，
// 用户会以为 yaml 已经救回来了。
func TestRestoreProjectRefusesInvisiblePath(t *testing.T) {
	svc := newTestService(t)
	svc.cfg.SetMounts(map[string]string{"/volume5/docker": t.TempDir()})

	mkProjectBackup(t, svc, "moontv", "20261010-120000",
		map[string]string{"compose.yaml": "a: 1\n"},
		map[string]string{"compose.yaml": "/volume9/elsewhere/compose.yaml"})

	res := svc.RestoreProject(nil, "moontv", "20261010-120000", false)
	if res.OK {
		t.Fatal("路径看不见却报成功 —— 用户会以为文件已经写回去了")
	}
	if len(res.Files) != 1 || res.Files[0].Written {
		t.Fatalf("该文件应当标记为未写入：%+v", res.Files)
	}
}

// 内容与备份一致的文件不该被重写（保持 mtime，也避免无意义的写盘）。
func TestRestoreProjectSkipsIdenticalFile(t *testing.T) {
	svc := newTestService(t)
	root := t.TempDir()
	svc.cfg.SetMounts(map[string]string{"/volume5/docker": root})
	if err := os.MkdirAll(filepath.Join(root, "npm"), 0o755); err != nil {
		t.Fatal(err)
	}
	body := "services: {}\n"
	full := filepath.Join(root, "npm", "docker-compose.yml")
	if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	before, _ := os.Stat(full)

	mkProjectBackup(t, svc, "npm", "20261010-120000",
		map[string]string{"docker-compose.yml": body},
		map[string]string{"docker-compose.yml": "/volume5/docker/npm/docker-compose.yml"})

	res := svc.RestoreProject(nil, "npm", "20261010-120000", false)
	if !res.OK {
		t.Fatalf("应当成功：%s", res.Message)
	}
	if res.Files[0].Note != "内容与备份一致，未改动" {
		t.Fatalf("应当识别出一致并跳过，实际：%q", res.Files[0].Note)
	}
	after, _ := os.Stat(full)
	if !after.ModTime().Equal(before.ModTime()) {
		t.Error("内容一致的文件不该被重写")
	}
}

// 项目备份也要跟着保留策略一起清 —— 用户只配一次策略，
// 不能出现「快照被清了、项目备份无限涨」。
func TestPruneProjectBackups(t *testing.T) {
	svc := newTestService(t)
	svc.cfg.SetMounts(map[string]string{"/volume5/docker": t.TempDir()})

	for _, ts := range []string{"20260101-000000", "20260102-000000", "20260103-000000"} {
		mkProjectBackup(t, svc, "immich", ts,
			map[string]string{"compose.yaml": "services: {}\n"},
			map[string]string{"compose.yaml": "/volume5/docker/immich/compose.yaml"})
	}
	items, err := svc.ListProjectBackups()
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 3 {
		t.Fatalf("应有 3 份项目备份，实际 %d", len(items))
	}

	res := svc.Prune(PruneOptions{KeepPerContainer: 1})
	if res.Removed != 2 {
		t.Fatalf("应清掉 2 份旧的项目备份，实际 %d", res.Removed)
	}
	left, _ := svc.ListProjectBackups()
	if len(left) != 1 || left[0].TS != "20260103-000000" {
		t.Fatalf("留下的应当是最近那份，实际 %+v", left)
	}
	if res.FreedBytes <= 0 {
		t.Error("释放字节数应当被统计")
	}
}

// 删除项目备份不能越界：无论传什么参数，备份目录之外的文件都不能被碰到。
//
// 这里不假设「一定会报错」—— sanitize 先把 `..` 化掉也是一种正确防线；
// 真正要钉住的是**结果**：外面的文件还在。
func TestDeleteProjectBackupNeverEscapesRoot(t *testing.T) {
	svc := newTestService(t)
	outside := filepath.Join(filepath.Dir(svc.cfg.ProjectBackupDir()), "keepme.txt")
	if err := os.WriteFile(outside, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, bad := range []string{"../..", "..", ".", "/", "a/../../.."} {
		_ = svc.DeleteProjectBackup(bad, "../..")
		if _, err := os.Stat(outside); err != nil {
			t.Fatalf("参数 %q 把备份目录之外的 %s 删掉了", bad, outside)
		}
	}
}

// 列出项目备份时，坏掉的 manifest（半个 json / 没有 manifest）不能让整个列表失败。
func TestListProjectBackupsSkipsBrokenEntries(t *testing.T) {
	svc := newTestService(t)
	mkProjectBackup(t, svc, "good", "20260101-000000",
		map[string]string{"compose.yaml": "a: 1\n"},
		map[string]string{"compose.yaml": "/volume5/docker/good/compose.yaml"})
	// 一个没有 manifest 的目录
	broken := filepath.Join(svc.cfg.ProjectBackupDir(), "bad", "20260101-000000")
	if err := os.MkdirAll(broken, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(broken, "compose.yaml"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	items, err := svc.ListProjectBackups()
	if err != nil {
		t.Fatalf("列表不该整体失败：%v", err)
	}
	if len(items) != 1 || items[0].Project != "good" {
		t.Fatalf("应当只列出完好的那份，实际 %+v", items)
	}
}
