package backup

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/aaron2024s/dockhelm/internal/config"
	"github.com/aaron2024s/dockhelm/internal/store"
)

// newTestService 造一个只依赖本地文件系统的备份服务（不连 Docker）。
func newTestService(t *testing.T) *Service {
	t.Helper()
	root := t.TempDir()
	t.Setenv("DOCKHELM_DATA", root)
	cfg := config.Load()
	if err := cfg.EnsureDirs(); err != nil {
		t.Fatalf("创建数据目录失败：%v", err)
	}
	st, err := store.Open(cfg.DBPath())
	if err != nil {
		t.Fatalf("打开数据文件失败：%v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return New(cfg, nil, st, nil)
}

// 造一份「项目备份」的磁盘现场：<projects>/<project>/<ts>/<文件> + _manifest.json。
func mkProjectBackup(t *testing.T, svc *Service, project, ts string, files map[string]string, hostPaths map[string]string) {
	t.Helper()
	dir := svc.projectBackupDir(project, ts)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	man := ProjectManifest{
		Project: project, TS: ts, Created: "2026-10-10T04:00:00Z",
		Reason: "manual", Files: []ProjectFile{}, Checksums: map[string]string{},
	}
	for name, body := range files {
		if err := os.WriteFile(dir+string(os.PathSeparator)+name, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
		man.Files = append(man.Files, ProjectFile{Name: name, HostPath: hostPaths[name]})
	}
	writeManifest(t, dir, man)
}

func writeManifest(t *testing.T, dir string, man ProjectManifest) {
	t.Helper()
	b, err := json.MarshalIndent(man, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dir+string(os.PathSeparator)+manifestName, b, 0o600); err != nil {
		t.Fatal(err)
	}
}
