package backup

import (
	"archive/tar"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
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

// tarDir/untarDir 必须严格互逆：这是一切卷备份的地基。
func TestTarRoundTrip(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	if err := os.MkdirAll(filepath.Join(src, "sub", "deep"), 0o755); err != nil {
		t.Fatal(err)
	}
	files := map[string]string{
		"a.txt":           "hello",
		"sub/b.bin":       strings.Repeat("x", 4096),
		"sub/deep/c.json": `{"k":1}`,
		"空 格 名.txt":       "unicode",
	}
	for name, body := range files {
		if err := os.WriteFile(filepath.Join(src, filepath.FromSlash(name)), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	// 符号链接：只记链接本身，不跟随
	symlinkChecked := false
	linkPath := filepath.Join(src, "link")
	if err := os.Symlink("a.txt", linkPath); err != nil {
		t.Logf("跳过符号链接检查（当前环境不支持）：%v", err)
	} else if fi, err := os.Lstat(linkPath); err != nil || fi.Mode()&os.ModeSymlink == 0 {
		// Windows 上没开开发者模式时，os.Symlink 会退化成普通文件
		t.Logf("跳过符号链接检查（创建出来的不是真正的链接，mode=%v）", fi.Mode())
	} else {
		symlinkChecked = true
	}

	tarPath := filepath.Join(dir, "out.tar")
	size, err := tarDir(context.Background(), src, tarPath)
	if err != nil {
		t.Fatalf("打包失败：%v", err)
	}
	if size <= 0 {
		t.Fatalf("打包结果大小为 %d，应大于 0", size)
	}

	dst := filepath.Join(dir, "dst")
	if err := os.MkdirAll(dst, 0o755); err != nil {
		t.Fatal(err)
	}
	n, err := untarDir(context.Background(), tarPath, dst)
	if err != nil {
		t.Fatalf("解包失败：%v", err)
	}
	if n < len(files) {
		t.Fatalf("解出 %d 个条目，少于写入的 %d 个", n, len(files))
	}

	for name, want := range files {
		got, err := os.ReadFile(filepath.Join(dst, filepath.FromSlash(name)))
		if err != nil {
			t.Fatalf("读回 %s 失败：%v", name, err)
		}
		if string(got) != want {
			t.Fatalf("%s 内容不符\n got=%q\nwant=%q", name, got, want)
		}
	}
	if symlinkChecked {
		if target, err := os.Readlink(filepath.Join(dst, "link")); err != nil || target != "a.txt" {
			t.Fatalf("符号链接未保留：target=%q err=%v", target, err)
		}
	}
}

// 解包必须挡住路径穿越：快照可能来自用户上传，不能让它写到卷目录之外。
func TestUntarRejectsTraversal(t *testing.T) {
	dir := t.TempDir()
	// 手写一个含 ../ 的 tar
	tarPath := filepath.Join(dir, "evil.tar")
	f, err := os.Create(tarPath)
	if err != nil {
		t.Fatal(err)
	}
	tw := tar.NewWriter(f)
	if err := tw.WriteHeader(&tar.Header{
		Typeflag: tar.TypeReg, Name: "../escaped.txt", Mode: 0o644, Size: 4,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write([]byte("evil")); err != nil {
		t.Fatal(err)
	}
	_ = tw.Close()
	_ = f.Close()

	dst := filepath.Join(dir, "dst")
	_ = os.MkdirAll(dst, 0o755)
	if _, err := untarDir(context.Background(), tarPath, dst); err == nil {
		t.Fatal("含越界路径的 tar 必须被拒绝")
	}
	if _, err := os.Stat(filepath.Join(dir, "escaped.txt")); err == nil {
		t.Fatal("越界文件被写出来了，路径穿越防护失效")
	}
}

func TestImportValidatesInput(t *testing.T) {
	s := newTestService(t)

	cases := []struct {
		name      string
		container string
		ts        string
		content   string
		wantErr   string
	}{
		{"空容器名", "", "20261009-094200", "{}", "容器名"},
		{"非法时间戳", "redis", "2026-10-09", "{}", "时间戳格式"},
		{"不是 JSON", "redis", "20261009-094200", "not json", "合法的 JSON"},
		{"缺 inspect", "redis", "20261009-094200", `{"foo":1}`, "没有 inspect"},
		{"路径分隔符", "a/b", "20261009-094200", `{"inspect":{"Id":"x"}}`, "路径分隔符"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := s.Import(c.container, c.ts, c.content)
			if err == nil {
				t.Fatalf("应当报错（含 %q），实际通过了", c.wantErr)
			}
			if !strings.Contains(err.Error(), c.wantErr) {
				t.Fatalf("错误信息 %q 未包含 %q", err.Error(), c.wantErr)
			}
		})
	}
}

func TestImportWritesAndLists(t *testing.T) {
	s := newTestService(t)

	insp := map[string]any{
		"Id":   "abc123",
		"Name": "/redis",
		"Config": map[string]any{
			"Image": "redis:alpine",
		},
		"State": map[string]any{"Running": true},
	}
	doc := map[string]any{"inspect": insp}
	b, err := json.Marshal(doc)
	if err != nil {
		t.Fatal(err)
	}

	item, err := s.Import("redis", "20261009-094200", string(b))
	if err != nil {
		t.Fatalf("导入失败：%v", err)
	}
	if item.Container != "redis" || item.TS != "20261009-094200" {
		t.Fatalf("导入结果不符：%+v", item)
	}
	if item.Image != "redis:alpine" {
		t.Fatalf("镜像解析错误：%q", item.Image)
	}

	list, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].Container != "redis" {
		t.Fatalf("列表应当出现刚导入的快照，实际 %+v", list)
	}
	if list[0].Reason != "imported" {
		t.Fatalf("文件里没有 _dockhelm 块时，来源应补成 imported，实际 %q", list[0].Reason)
	}

	// 带 _dockhelm 块的文件：沿用文件自带的来源，且接口返回值与列表必须一致
	withMeta := `{"_dockhelm":{"reason":"scheduled","version":1},"inspect":{"Id":"y","Name":"/nginx","Config":{"Image":"nginx:alpine"}}}`
	m, err := s.Import("nginx", "20261001-080000", withMeta)
	if err != nil {
		t.Fatalf("导入带元数据的快照失败：%v", err)
	}
	if m.Reason != "scheduled" {
		t.Fatalf("接口应回显文件自带的来源 scheduled，实际 %q", m.Reason)
	}
	nl, _ := s.List()
	for _, it := range nl {
		if it.Container == "nginx" {
			if it.Reason != m.Reason {
				t.Fatalf("列表来源 %q 与接口返回 %q 不一致", it.Reason, m.Reason)
			}
		}
	}

	// 同一时间戳再导一次：不能覆盖，也不能报错
	again, err := s.Import("redis", "20261009-094200", string(b))
	if err != nil {
		t.Fatalf("重复导入应当自动加后缀而不是失败：%v", err)
	}
	if again.TS == item.TS {
		t.Fatalf("重复导入覆盖了已有快照（TS 都是 %s）", again.TS)
	}
	list, _ = s.List()
	// redis 两份（原 + 重名加后缀），nginx 一份
	if len(list) != 3 {
		t.Fatalf("应当有三份快照，实际 %d 份", len(list))
	}
}

// Delete 必须连卷数据一起删掉 —— 只删 json 会剩下一堆几十 MB 的 tar。
func TestDeleteRemovesVolumeData(t *testing.T) {
	s := newTestService(t)

	dir := filepath.Join(s.cfg.ContainerBackupDir(), "redis")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "20261009-094200.json"), []byte(`{"inspect":{"Id":"x"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	vdir := filepath.Join(dir, "20261009-094200.volumes")
	if err := os.MkdirAll(vdir, 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(vdir, "data.tar"), []byte("xxxx"), 0o600); err != nil {
		t.Fatal(err)
	}

	if err := s.Delete("redis", "20261009-094200"); err != nil {
		t.Fatalf("删除失败：%v", err)
	}
	if _, err := os.Stat(vdir); !os.IsNotExist(err) {
		t.Fatal("卷数据目录没有被一起删除")
	}
}
