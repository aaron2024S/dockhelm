package updater

import (
	"strings"
	"testing"

	"github.com/aaron2024s/dockhelm/internal/dockerx"
)

func TestParseRef(t *testing.T) {
	cases := []struct {
		in       string
		registry string
		repo     string
		tag      string
		local    string
	}{
		{"redis:alpine", "docker.io", "library/redis", "alpine", "redis"},
		{"redis", "docker.io", "library/redis", "latest", "redis"},
		{"mtvpls/moontvplus", "docker.io", "mtvpls/moontvplus", "latest", "mtvpls/moontvplus"},
		{"ghcr.io/foo/bar:v1", "ghcr.io", "foo/bar", "v1", "ghcr.io/foo/bar"},
		{"registry-1.docker.io/library/redis:alpine", "registry-1.docker.io", "library/redis", "alpine", "registry-1.docker.io/library/redis"},
		{"localhost:5000/app:x", "localhost:5000", "app", "x", "localhost:5000/app"},
		{"nginx:1.25.3-alpine", "docker.io", "library/nginx", "1.25.3-alpine", "nginx"},
	}
	for _, c := range cases {
		got := dockerx.ParseRef(c.in)
		if got.Registry != c.registry || got.Repository != c.repo || got.Tag != c.tag || got.LocalName != c.local {
			t.Errorf("ParseRef(%q) = {reg:%q repo:%q tag:%q local:%q}, want {reg:%q repo:%q tag:%q local:%q}",
				c.in, got.Registry, got.Repository, got.Tag, got.LocalName, c.registry, c.repo, c.tag, c.local)
		}
	}
}

func TestParseRefDigest(t *testing.T) {
	r := dockerx.ParseRef("redis@sha256:abc123")
	if r.Digest != "sha256:abc123" {
		t.Fatalf("digest = %q", r.Digest)
	}
	if r.Repository != "library/redis" {
		t.Fatalf("repo = %q", r.Repository)
	}
}

// fakeInspect 造一个「用 -v ./data:/data 和一个命名卷创建」的容器 inspect 结果。
func fakeInspect() map[string]any {
	return map[string]any{
		"Id":   "9f8e7d6c5b4a39281706f5e4d3c2b1a0",
		"Name": "/moontv-redis",
		"Config": map[string]any{
			"Hostname":     "9f8e7d6c5b4a", // 等于短 ID ⇒ 应被删掉
			"Image":        "redis:alpine",
			"Env":          []any{"PATH=/usr/local/bin", "TZ=Asia/Shanghai"},
			"ExposedPorts": map[string]any{"6379/tcp": map[string]any{}},
			"Cmd":          []any{"redis-server"},
			"Labels":       map[string]any{"com.docker.compose.project": "moontv"},
		},
		"HostConfig": map[string]any{
			// 这两个必须被删掉，否则重复挂载
			"Binds":      []any{"./data:/data", "cfgvol:/etc/redis"},
			"Mounts":     []any{map[string]any{"Type": "bind", "Source": "./data", "Target": "/data"}},
			"NetworkMode": "moontv_default",
			"RestartPolicy": map[string]any{"Name": "always", "MaximumRetryCount": 0},
			"PortBindings": map[string]any{},
			"BindsExtra":   "should-stay",
		},
		"Mounts": []any{
			map[string]any{
				"Type": "bind", "Source": "/volume1/docker/moontv/data", "Destination": "/data",
				"RW": true, "Propagation": "rprivate",
			},
			map[string]any{
				"Type": "volume", "Name": "cfgvol", "Source": "/var/lib/docker/volumes/cfgvol/_data",
				"Destination": "/etc/redis", "RW": false, "Propagation": "",
			},
			map[string]any{
				"Type": "volume", "Name": "anon123", "Source": "/var/lib/docker/volumes/anon123/_data",
				"Destination": "/var/lib/anon", "RW": true,
			},
		},
		"NetworkSettings": map[string]any{
			"Networks": map[string]any{
				"moontv_default": map[string]any{
					"Aliases":   []any{"9f8e7d6c5b4a", "moontv-redis"},
					"IPAMConfig": nil,
				},
			},
		},
	}
}

func TestBuildCreateSpecClearsLegacyMountFields(t *testing.T) {
	cfg, hc, _ := BuildCreateSpec(fakeInspect(), nil)
	if _, ok := hc["Binds"]; ok {
		t.Error("HostConfig.Binds 必须被删除，否则与 Mounts 重复")
	}
	if _, ok := hc["BindsExtra"]; !ok {
		t.Error("未识别的 HostConfig 字段必须原样保留")
	}
	if hc["NetworkMode"] != "moontv_default" {
		t.Error("NetworkMode 必须保留")
	}
	if _, ok := cfg["Env"]; !ok {
		t.Error("Env 必须保留")
	}
}

func TestBuildCreateSpecMountsRebuiltFromTopLevel(t *testing.T) {
	_, hc, _ := BuildCreateSpec(fakeInspect(), nil)
	mounts, ok := hc["Mounts"].([]map[string]any)
	if !ok {
		t.Fatalf("Mounts 类型不对: %T", hc["Mounts"])
	}
	if len(mounts) != 3 {
		t.Fatalf("应重建出 3 个挂载，实际 %d", len(mounts))
	}

	// bind：保留宿主路径与只读/传播
	b := mounts[0]
	if b["Type"] != "bind" || b["Target"] != "/data" || b["Source"] != "/volume1/docker/moontv/data" {
		t.Errorf("bind 挂载不对: %+v", b)
	}
	if _, hasRO := b["ReadOnly"]; hasRO {
		t.Error("RW=true 的挂载不应该有 ReadOnly")
	}
	if _, hasOpts := b["BindOptions"]; hasOpts {
		t.Error("rprivate 是默认传播模式，不应该写回")
	}

	// 命名卷：Source 必须是卷名而不是宿主目录
	v := mounts[1]
	if v["Type"] != "volume" || v["Source"] != "cfgvol" || v["Target"] != "/etc/redis" {
		t.Errorf("命名卷不对: %+v", v)
	}
	if v["ReadOnly"] != true {
		t.Error("RW=false 的挂载应该带 ReadOnly")
	}

	// 匿名卷：必须复用同一个卷名，否则更新后数据变孤儿
	a := mounts[2]
	if a["Source"] != "anon123" {
		t.Errorf("匿名卷必须按 Name 复用，实际 Source=%v", a["Source"])
	}
}

func TestBuildCreateSpecHostname(t *testing.T) {
	insp := fakeInspect()
	cfg, _, _ := BuildCreateSpec(insp, nil)
	if _, ok := cfg["Hostname"]; ok {
		t.Error("Hostname 等于容器短 ID 时应被删除（让守护进程重新生成）")
	}

	// 用户显式设过 hostname 时要保留
	insp2 := fakeInspect()
	insp2["Config"].(map[string]any)["Hostname"] = "my-custom-host"
	cfg2, _, _ := BuildCreateSpec(insp2, nil)
	if cfg2["Hostname"] != "my-custom-host" {
		t.Error("用户显式设置的 Hostname 必须保留")
	}
}

func TestBuildCreateSpecSkipsMissingNetwork(t *testing.T) {
	// 网络已不存在 ⇒ 不应出现在 EndpointsConfig 里
	_, _, netCfg := BuildCreateSpec(fakeInspect(), map[string]bool{"other_net": true})
	if netCfg != nil {
		t.Errorf("网络不存在时应返回 nil NetworkingConfig，实际 %+v", netCfg)
	}

	// 网络存在 ⇒ 保留别名
	_, _, netCfg2 := BuildCreateSpec(fakeInspect(), map[string]bool{"moontv_default": true})
	if netCfg2 == nil {
		t.Fatal("网络存在时应有 NetworkingConfig")
	}
	eps := netCfg2["EndpointsConfig"].(map[string]any)
	if _, ok := eps["moontv_default"]; !ok {
		t.Error("应包含 moontv_default")
	}
}

func TestBuildCreateSpecHostNetworkModeSkipsNetworking(t *testing.T) {
	insp := fakeInspect()
	insp["HostConfig"].(map[string]any)["NetworkMode"] = "host"
	_, _, netCfg := BuildCreateSpec(insp, nil)
	if netCfg != nil {
		t.Error("host 网络模式不应传 EndpointsConfig")
	}
}

func TestBuildCreateSpecReadonlyBindPropagation(t *testing.T) {
	insp := fakeInspect()
	insp["Mounts"] = []any{
		map[string]any{
			"Type": "bind", "Source": "/host/sec", "Destination": "/sec",
			"RW": false, "Propagation": "rslave",
		},
	}
	_, hc, _ := BuildCreateSpec(insp, nil)
	m := hc["Mounts"].([]map[string]any)[0]
	if m["ReadOnly"] != true {
		t.Error("只读 bind 应带 ReadOnly")
	}
	opts, _ := m["BindOptions"].(map[string]any)
	if opts["Propagation"] != "rslave" {
		t.Errorf("非默认传播模式必须保留: %+v", m)
	}
}

func TestBuildCreateSpecVolumeWithoutNameIsDropped(t *testing.T) {
	// 拿不到卷名的匿名卷无法安全复用，宁可跳过也不能悄悄换一个卷
	insp := fakeInspect()
	insp["Mounts"] = []any{
		map[string]any{"Type": "volume", "Source": "/var/lib/docker/volumes/x/_data", "Destination": "/x"},
	}
	_, hc, _ := BuildCreateSpec(insp, nil)
	if _, ok := hc["Mounts"]; ok {
		t.Error("无法确定卷名的挂载应被丢弃")
	}
}

// TestEnsureCreateSpec 校验「动手之前的那道闸」。
func TestEnsureCreateSpec(t *testing.T) {
	cases := []struct {
		name     string
		cfg      map[string]any
		insp     map[string]any
		wantFail string // 非空表示期望失败，值会被当成子串匹配
		wantImg  string // 期望被补上的镜像引用
	}{
		{
			name: "配置齐全直接就过",
			cfg:  map[string]any{"Image": "nginx:1.27-alpine"},
			insp: map[string]any{"Image": "sha256:aaaa"},
		},
		{
			name:     "完全没有配置",
			cfg:      nil,
			insp:     map[string]any{"Image": "sha256:aaaa"},
			wantFail: "没有 Config",
		},
		{
			name:     "空配置",
			cfg:      map[string]any{},
			insp:     map[string]any{},
			wantFail: "没有 Config",
		},
		{
			name:    "缺镜像引用时退回顶层镜像 ID",
			cfg:     map[string]any{"Cmd": []any{"x"}},
			insp:    map[string]any{"Image": "sha256:aaaa"},
			wantImg: "sha256:aaaa",
		},
		{
			name:     "既没有引用也没有 ID 才放弃",
			cfg:      map[string]any{"Env": []any{"A=1"}},
			insp:     map[string]any{},
			wantFail: "镜像 ID",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			patched, fail := EnsureCreateSpec(tc.cfg, tc.insp)
			if tc.wantFail != "" {
				if fail == "" {
					t.Fatalf("期望失败（含 %q），却通过了", tc.wantFail)
				}
				if !strings.Contains(fail, tc.wantFail) {
					t.Fatalf("失败原因里没有 %q：%s", tc.wantFail, fail)
				}
				return
			}
			if fail != "" {
				t.Fatalf("不该失败：%s", fail)
			}
			if patched != tc.wantImg {
				t.Fatalf("补上的镜像引用应为 %q，实际 %q", tc.wantImg, patched)
			}
			if tc.wantImg != "" {
				if got, _ := tc.cfg["Image"].(string); got != tc.wantImg {
					t.Fatalf("修补结果没写回 cfg：Image=%q", got)
				}
			}
		})
	}
}

// TestBuildCreateSpecPassesEnsure 真实形状的 inspect 拼出来的请求体必须能过闸。
// 这条把「拼装」与「校验」串起来：任何一边改动导致结果不可用都会被拦住。
func TestBuildCreateSpecPassesEnsure(t *testing.T) {
	insp := fakeInspect()
	cfg, hostCfg, _ := BuildCreateSpec(insp, nil)
	if _, fail := EnsureCreateSpec(cfg, insp); fail != "" {
		t.Fatalf("从 inspect 拼出的请求体不该被判为不可用：%s", fail)
	}
	if img, _ := cfg["Image"].(string); img != "redis:alpine" {
		t.Fatalf("Configuration.Image 应被保留，实际 %q", img)
	}
	if _, ok := hostCfg["Mounts"]; !ok {
		t.Fatal("挂载应从顶层 Mounts 重建后写进 HostConfig")
	}
}
