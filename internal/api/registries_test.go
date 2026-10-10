package api

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/aaron2024s/dockhelm/internal/store"
)

// shouldSeedRegistries 是「预置加速源只灌一次」这条不变式唯一的守门人。
// 一旦这里判错，用户删掉的加速源就会在重启后自己长回来。
func TestShouldSeedRegistries(t *testing.T) {
	cases := []struct {
		name     string
		seedMark string
		settings string
		want     bool
	}{
		{"全新安装：既没标记也没配置", "", "", true},
		{"已经灌过一轮", "1", "", false},
		{"老用户升级：配置早就在（哪怕列表是空的）", "", `{"mirrors":[]}`, false},
		{"用户删光后重启：标记还在", "1", `{"mirrors":[]}`, false},
	}
	for _, c := range cases {
		if got := shouldSeedRegistries(c.seedMark, c.settings); got != c.want {
			t.Errorf("%s：shouldSeedRegistries(%q, %q) = %v，期望 %v",
				c.name, c.seedMark, c.settings, got, c.want)
		}
	}
}

// 预置清单是「开箱即用」的默认配置，质量要钉住。
func TestPresetMirrorsSanity(t *testing.T) {
	if len(presetMirrors) == 0 {
		t.Fatal("预置清单不能为空")
	}
	seen := map[string]bool{}
	for _, m := range presetMirrors {
		if !strings.HasPrefix(m.URL, "https://") {
			t.Errorf("预置源必须是 https：%s", m.URL)
		}
		if !m.Enabled {
			t.Errorf("预置源必须默认启用（否则「直接写进我的加速源」没意义）：%s", m.URL)
		}
		if !m.Builtin {
			t.Errorf("预置源必须带 builtin 标记（前端要显示「预置」徽标）：%s", m.URL)
		}
		key := normalizeMirror(m.URL)
		if seen[key] {
			t.Errorf("预置清单里有重复项：%s", m.URL)
		}
		seen[key] = true
	}
	// 2026-10-10 实测：dockerhub.icu 连续超时已剔除，别让它被加回来。
	if seen[normalizeMirror("https://dockerhub.icu")] {
		t.Error("dockerhub.icu 已失效，不应出现在预置清单里")
	}
}

// 端到端：全新安装灌一次、再启动不翻倍、用户删光后不复活。
func TestSeedRegistries(t *testing.T) {
	st, err := store.Open(filepath.Join(t.TempDir(), "dockhelm.db"))
	if err != nil {
		t.Fatalf("打开测试数据文件失败：%v", err)
	}
	defer st.Close()

	readMirrors := func() []MirrorConfig {
		var cfg RegistrySettings
		st.GetJSON(registrySettingsKey, &cfg)
		return cfg.Mirrors
	}

	// 首次启动：整份预置清单灌进「我的加速源」
	seedRegistries(st)
	got := readMirrors()
	if len(got) != len(presetMirrors) {
		t.Fatalf("首次灌入后应有 %d 条，实际 %d 条", len(presetMirrors), len(got))
	}
	for _, m := range got {
		if !m.Enabled || !m.Builtin {
			t.Errorf("灌进来的条目应默认启用且带 builtin 标记：%+v", m)
		}
	}

	// 再启动一次：不能翻倍
	seedRegistries(st)
	if n := len(readMirrors()); n != len(presetMirrors) {
		t.Errorf("重复启动导致了重复灌入：%d 条", n)
	}

	// 用户删光后重启：绝不能自己长回来（这是本次改动的核心不变式）
	if err := st.SetJSON(registrySettingsKey, RegistrySettings{Mirrors: []MirrorConfig{}}); err != nil {
		t.Fatalf("写回空的加速源列表失败：%v", err)
	}
	seedRegistries(st)
	if n := len(readMirrors()); n != 0 {
		t.Errorf("用户删光的加速源在重启后又出现了 %d 条", n)
	}
}
