package config

import (
	"os"
	"path/filepath"
	"testing"
)

// 造一个「容器内」的目录树，返回它的根。
func fakeContainer(t *testing.T, files ...string) string {
	t.Helper()
	root := t.TempDir()
	for _, f := range files {
		p := filepath.Join(root, filepath.FromSlash(f))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte("x"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

// 自动识别出的映射：宿主机路径换成容器内路径后能读到文件。
func TestMapHostPathAutoRewrite(t *testing.T) {
	root := fakeContainer(t, "moontv/docker-compose.yml")
	cfg := &Config{}
	cfg.SetMounts(map[string]string{"/volume1/docker": root})

	got, ok := cfg.MapHostPath("/volume1/docker/moontv/docker-compose.yml")
	if !ok {
		t.Fatal("应当映射成功，却返回 ok=false")
	}
	want := filepath.Join(root, "moontv", "docker-compose.yml")
	if got != want {
		t.Fatalf("映射结果 = %q，期望 %q", got, want)
	}
}

// 挂在同一个宿主路径上的映射里，最长的那个前缀应当优先命中。
func TestMapHostPathLongestPrefixWins(t *testing.T) {
	outer := fakeContainer(t, "other/keep.txt")
	inner := fakeContainer(t, "docker-compose.yml")

	cfg := &Config{}
	cfg.SetMounts(map[string]string{
		"/volume1/docker":       outer,
		"/volume1/docker/moontv": inner,
	})

	got, ok := cfg.MapHostPath("/volume1/docker/moontv/docker-compose.yml")
	if !ok {
		t.Fatal("应当映射成功，却返回 ok=false")
	}
	if got != filepath.Join(inner, "docker-compose.yml") {
		t.Fatalf("最长前缀未生效，得到 %q", got)
	}
}

// 两边一致的写法（- /volume1/docker:/volume1/docker）继续可用。
func TestMapHostPathBothSidesIdentical(t *testing.T) {
	root := fakeContainer(t, "app/compose.yaml")
	cfg := &Config{}
	cfg.SetMounts(map[string]string{root: root})

	got, ok := cfg.MapHostPath(filepath.Join(root, "app", "compose.yaml"))
	if !ok || got != filepath.Join(root, "app", "compose.yaml") {
		t.Fatalf("两边一致的情形应当直通，得到 %q ok=%v", got, ok)
	}
}

// 前缀必须按路径分段匹配：/vol 不能命中 /volume1。
func TestMapHostPathSegmentBoundary(t *testing.T) {
	root := fakeContainer(t, "docker/app/compose.yaml")
	cfg := &Config{}
	cfg.SetMounts(map[string]string{"/vol": root})

	if _, ok := cfg.MapHostPath("/volume1/docker/app/compose.yaml"); ok {
		t.Fatal("/vol 不该命中 /volume1 —— 分段边界判断失效")
	}
}

// 完全看不见的路径必须老实返回 false，不能瞎猜。
func TestMapHostPathUnreadableReturnsFalse(t *testing.T) {
	root := fakeContainer(t, "moontv/docker-compose.yml")
	cfg := &Config{}
	cfg.SetMounts(map[string]string{"/volume1/docker": root})

	if local, ok := cfg.MapHostPath("/volume9/nothing/compose.yaml"); ok {
		t.Fatalf("看不见的路径应当返回 false，却给了 %q", local)
	}
}

// DOCKHELM_HOST_ROOTS 的成对语法：宿主机=容器内。
func TestHostRootsEnvPairSyntax(t *testing.T) {
	root := fakeContainer(t, "moontv/docker-compose.yml")
	t.Setenv("DOCKHELM_HOST_ROOTS", "/volume1/docker="+root+", /volume2/apps")

	cfg := Load()

	got, ok := cfg.MapHostPath("/volume1/docker/moontv/docker-compose.yml")
	if !ok || got != filepath.Join(root, "moontv", "docker-compose.yml") {
		t.Fatalf("成对语法未生效：got=%q ok=%v", got, ok)
	}
	// 单值写法仍然表示「两边一致」
	ms := cfg.PathMappings()
	var envCount int
	for _, m := range ms {
		if m.Source == "env" {
			envCount++
		}
	}
	if envCount != 2 {
		t.Fatalf("应当解析出 2 条 env 映射，实际 %d：%+v", envCount, ms)
	}
}

// 自动映射 + env 映射同时存在时都要出现在 PathMappings 里，且带来源标记。
func TestPathMappingsMergesSources(t *testing.T) {
	root := fakeContainer(t, "a")
	t.Setenv("DOCKHELM_HOST_ROOTS", "/host="+root)

	cfg := Load()
	cfg.SetMounts(map[string]string{"/data": "/data"})

	var auto, env int
	for _, m := range cfg.PathMappings() {
		switch m.Source {
		case "auto":
			auto++
		case "env":
			env++
		}
	}
	if auto != 1 || env != 1 {
		t.Fatalf("来源标记不对：auto=%d env=%d", auto, env)
	}
}

// Summary 在没有任何映射时也要给出可读的提示（启动日志用）。
func TestSummaryWithoutMappings(t *testing.T) {
	cfg := &Config{}
	if s := cfg.Summary(); s == "" {
		t.Fatal("Summary 不该为空")
	}
}
