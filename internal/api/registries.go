package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/aaron2024s/dockhelm/internal/store"
)

// MirrorConfig 用户配置的一个加速源。
type MirrorConfig struct {
	URL        string `json:"url"`
	Note       string `json:"note"`
	Enabled    bool   `json:"enabled"`
	LastTested string `json:"lastTested,omitempty"`
	LatencyMs  int64  `json:"latencyMs,omitempty"`
	// OK 不带 omitempty：false 是有意义的信息（测过了但不可用）。
	// 「有没有测过」由 LastTested 是否为空判断，别用 ok 是否出现来判断。
	OK  bool   `json:"ok"`
	Err string `json:"err,omitempty"`
	// Builtin 表示这条是从预置清单里灌进来的。只是个来源标记：
	// 不影响增删改，用户删掉它和删掉自建条目没有任何区别，也不会自己回来。
	Builtin bool `json:"builtin,omitempty"`
}

// RegistrySettings 加速源设置。
type RegistrySettings struct {
	Mirrors []MirrorConfig `json:"mirrors"`
	// PullMirror 指定 Dockhelm 拉取 Docker Hub 镜像时使用的加速站。
	// 留空 = 完全交给 Docker 守护进程（读 daemon.json 的 registry-mirrors）。
	PullMirror string `json:"pullMirror"`
	// Insecure 允许 http 或自签证书的私有仓库
	Insecure bool `json:"insecure"`
}

// 存键。加速源配置与「预置源已灌过」标记分开存，
// 后者是「删掉的内置源不能自己长回来」这条不变式的依据。
const (
	registrySettingsKey = "registry.settings"
	registrySeedMarkKey = "registry.seeded"
)

// presetMirrors 是随应用预置的常用加速源。
//
// 定位是「首次启动的起手配置」，**不再是一个单独的推荐栏位**：服务第一次启动时
// 整份写进「我的加速源」（默认勾选启用），之后完全归用户所有 —— 不想要哪条直接删，
// 删掉不会再回来（见 seedRegistries）。列表被删空了也能一键找回（presets + 前端空态按钮）。
//
// 清单不会自动联网刷新，所以每次改动都必须是实测过的：
// 2026-10-10 由 NAS 出网实测，剔除连续超时的 dockerhub.icu，
// 补入当时最快的候补 docker.1panel.top。按实测延迟升序排列。
var presetMirrors = []MirrorConfig{
	{URL: "https://docker.m.daocloud.io", Note: "DaoCloud", Enabled: true, Builtin: true},
	{URL: "https://docker.1ms.run", Note: "1ms.run", Enabled: true, Builtin: true},
	{URL: "https://hub.rat.dev", Note: "rat.dev", Enabled: true, Builtin: true},
	{URL: "https://docker.1panel.top", Note: "1Panel", Enabled: true, Builtin: true},
	{URL: "https://docker.1panel.live", Note: "1Panel live", Enabled: true, Builtin: true},
	{URL: "https://dockerproxy.net", Note: "dockerproxy", Enabled: true, Builtin: true},
}

// shouldSeedRegistries 判断首次启动要不要把预置加速源灌进用户的列表。
//
// 只有「从没灌过」+「从没配置过」两条同时成立才灌。这是这条功能唯一的正确性边界：
//   - 老版本升级上来的用户：registry.settings 早就存在（哪怕列表被清空过）=> 不灌；
//   - 已经删光预置源的用户：同一个键依然存在 => 不灌，删掉的东西不能自己长回来。
func shouldSeedRegistries(seedMark, settingsRaw string) bool {
	return seedMark == "" && settingsRaw == ""
}

// seedRegistries 在服务启动时调用一次：首次运行把预置加速源写进「我的加速源」。
//
// 之后再也不动用户的列表 —— 用户删一条、删光，都只是用户自己的事。
func seedRegistries(st *store.Store) {
	if st.GetSetting(registrySeedMarkKey, "") != "" {
		return
	}
	if shouldSeedRegistries("", st.GetSetting(registrySettingsKey, "")) {
		cfg := RegistrySettings{Mirrors: append([]MirrorConfig{}, presetMirrors...)}
		if err := st.SetJSON(registrySettingsKey, cfg); err != nil {
			log.Printf("写入预置加速源失败（下次启动会重试）：%v", err)
		} else {
			log.Printf("已写入 %d 个预置加速源到「我的加速源」（可在镜像加速源页逐条删除）", len(cfg.Mirrors))
		}
	}
	// 标记与灌入分开：即使这次没灌（用户早就有配置），也要置位，
	// 免得以后每次启动都重复判断一遍。
	_ = st.SetSetting(registrySeedMarkKey, "1")
}

// hGetRegistries 返回配置 + 守护进程真实生效的镜像列表。
func (s *Server) hGetRegistries(w http.ResponseWriter, r *http.Request) {
	var cfg RegistrySettings
	s.st.GetJSON(registrySettingsKey, &cfg)
	if cfg.Mirrors == nil {
		cfg.Mirrors = []MirrorConfig{}
	}

	ctx, cancel := s.ctx(r)
	defer cancel()
	daemonMirrors := []string{}
	daemonErr := ""
	if info, err := s.dc.Info(ctx); err == nil {
		daemonMirrors = info.RegistryConfig.Mirrors
	} else {
		daemonErr = err.Error()
	}

	// presets 只服务于「列表被删空后一键找回」，不再渲染成页面上独立的「推荐加速源」栏位 ——
	// 预置源在首次启动时就已经写进 settings.mirrors，用户删掉即消失。
	writeOK(w, map[string]any{
		"settings":      cfg,
		"presets":       presetMirrors,
		"daemonMirrors": daemonMirrors,
		"daemonError":   daemonErr,
		"snippet":       daemonSnippet(cfg, daemonMirrors),
		"explain": "Docker 守护进程只认 daemon.json 里的 registry-mirrors，Dockhelm 无法替它改配置。" +
			"这里做三件事：① 显示守护进程当前真正生效的加速源；② 让你维护、测速并生成 daemon.json 片段贴到 NAS 面板；" +
			"③ 如果不想动 daemon.json，可以直接选一个加速源让 Dockhelm 自己按它拉取。",
	})
}

type saveRegistriesReq struct {
	Settings RegistrySettings `json:"settings"`
}

func (s *Server) hSaveRegistries(w http.ResponseWriter, r *http.Request) {
	var in saveRegistriesReq
	if err := decodeBody(r, &in); err != nil {
		writeErr(w, http.StatusBadRequest, "请求体格式错误："+err.Error())
		return
	}
	clean := []MirrorConfig{}
	seen := map[string]bool{}
	for _, m := range in.Settings.Mirrors {
		u := strings.TrimSpace(m.URL)
		if u == "" {
			continue
		}
		key := normalizeMirror(u)
		if seen[key] {
			continue
		}
		seen[key] = true
		m.URL = u
		clean = append(clean, m)
	}
	in.Settings.Mirrors = clean
	in.Settings.PullMirror = strings.TrimSpace(in.Settings.PullMirror)

	if err := s.st.SetJSON(registrySettingsKey, in.Settings); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	// 立刻通知更新引擎换用新的加速源
	s.up.SetMirror(func() string {
		var c RegistrySettings
		s.st.GetJSON(registrySettingsKey, &c)
		return c.PullMirror
	})
	writeOK(w, map[string]any{"ok": true, "settings": in.Settings})
}

type testRegistryReq struct {
	URL string `json:"url"`
}

// hTestRegistry 测一个加速源的可达性与延迟。
//
// 从 Dockhelm 容器里探测 —— 它测到的网络路径正是拉镜像时走的路径，所以这个数字是有意义的。
func (s *Server) hTestRegistry(w http.ResponseWriter, r *http.Request) {
	var in testRegistryReq
	if err := decodeBody(r, &in); err != nil || strings.TrimSpace(in.URL) == "" {
		writeErr(w, http.StatusBadRequest, "需要提供加速源地址")
		return
	}
	res := probeMirror(r.Context(), strings.TrimSpace(in.URL))
	s.patchMirrorResult(res)
	writeOK(w, res)
}

type testAllReq struct {
	URLs []string `json:"urls"`
}

// hTestAllRegistries 批量测速，按延迟排序返回。
func (s *Server) hTestAllRegistries(w http.ResponseWriter, r *http.Request) {
	var in testAllReq
	_ = decodeBody(r, &in)
	urls := in.URLs
	if len(urls) == 0 {
		// 兜底只测「我的加速源」里的条目。预置清单不再单独补测 ——
		// 用户删掉的源就是不想再看到它，不该还出现在测速结果里。
		var cfg RegistrySettings
		s.st.GetJSON(registrySettingsKey, &cfg)
		for _, m := range cfg.Mirrors {
			urls = append(urls, m.URL)
		}
	}
	out := make([]MirrorConfig, 0, len(urls))
	sem := make(chan struct{}, 6)
	done := make(chan MirrorConfig, len(urls))
	for _, u := range urls {
		go func(u string) {
			sem <- struct{}{}
			defer func() { <-sem }()
			done <- probeMirror(r.Context(), u)
		}(u)
	}
	for range urls {
		out = append(out, <-done)
	}
	// 按「可用优先、延迟升序」
	sortMirrors(out)
	for _, m := range out {
		s.patchMirrorResult(m)
	}
	writeOK(w, map[string]any{"results": out})
}

// patchMirrorResult 把测速结果回写到已保存的镜像项里（匹配 URL 即更新）。
func (s *Server) patchMirrorResult(res MirrorConfig) {
	var cfg RegistrySettings
	s.st.GetJSON(registrySettingsKey, &cfg)
	if len(cfg.Mirrors) == 0 {
		return
	}
	key := normalizeMirror(res.URL)
	changed := false
	for i := range cfg.Mirrors {
		if normalizeMirror(cfg.Mirrors[i].URL) == key {
			cfg.Mirrors[i].LatencyMs = res.LatencyMs
			cfg.Mirrors[i].OK = res.OK
			cfg.Mirrors[i].Err = res.Err
			cfg.Mirrors[i].LastTested = res.LastTested
			changed = true
		}
	}
	if changed {
		_ = s.st.SetJSON(registrySettingsKey, cfg)
	}
}

func probeMirror(parent context.Context, raw string) MirrorConfig {
	out := MirrorConfig{
		URL: strings.TrimSpace(raw), LastTested: time.Now().UTC().Format(time.RFC3339),
	}
	base := strings.TrimSuffix(out.URL, "/")
	if !strings.HasPrefix(base, "http://") && !strings.HasPrefix(base, "https://") {
		base = "https://" + base
		out.URL = base
	}
	client := &http.Client{
		Timeout: 8 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{MinVersion: tls.VersionTLS12},
		},
		// 不跟随跳转：能跳转说明服务是活的，但我们要的是本地这一跳的延迟
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
	ctx, cancel := context.WithTimeout(parent, 8*time.Second)
	defer cancel()

	best := int64(-1)
	var lastErr string
	// /v2/ 是 Registry HTTP API 的版本探测端点：200 或 401 都说明服务可用
	for _, path := range []string{"/v2/", "/v2/library/alpine/manifests/latest"} {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+path, nil)
		if err != nil {
			lastErr = err.Error()
			continue
		}
		req.Header.Set("User-Agent", "dockhelm-probe")
		start := time.Now()
		resp, err := client.Do(req)
		elapsed := time.Since(start).Milliseconds()
		if err != nil {
			lastErr = err.Error()
			continue
		}
		_ = resp.Body.Close()
		if resp.StatusCode < 500 {
			if best < 0 || elapsed < best {
				best = elapsed
			}
			out.OK = true
			out.LatencyMs = best
			out.Err = ""
			return out
		}
		lastErr = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}
	out.OK = false
	if best >= 0 {
		out.LatencyMs = best
	}
	out.Err = lastErr
	if out.Err == "" {
		out.Err = "无法连接"
	}
	return out
}

// daemonSnippet 生成可以直接贴进 daemon.json 的片段。
func daemonSnippet(cfg RegistrySettings, daemonMirrors []string) string {
	urls := []string{}
	for _, m := range cfg.Mirrors {
		if m.Enabled {
			urls = append(urls, m.URL)
		}
	}
	if len(urls) == 0 {
		urls = daemonMirrors
	}
	if len(urls) == 0 {
		return "{\n  \"registry-mirrors\": []\n}"
	}
	var b strings.Builder
	b.WriteString("{\n  \"registry-mirrors\": [\n")
	for i, u := range urls {
		b.WriteString("    \"" + u + "\"")
		if i < len(urls)-1 {
			b.WriteString(",")
		}
		b.WriteString("\n")
	}
	b.WriteString("  ]\n}")
	return b.String()
}

func normalizeMirror(u string) string {
	u = strings.TrimSpace(strings.ToLower(u))
	u = strings.TrimSuffix(u, "/")
	return u
}

func sortMirrors(list []MirrorConfig) {
	sort.SliceStable(list, func(i, j int) bool {
		a, b := list[i], list[j]
		if a.OK != b.OK {
			return b.OK // 可用的排前面
		}
		if !a.OK {
			return a.URL < b.URL
		}
		if a.LatencyMs != b.LatencyMs {
			return a.LatencyMs < b.LatencyMs
		}
		return a.URL < b.URL
	})
}

// PullMirrorOf 读取当前配置的拉取加速源（启动时注入给更新引擎）。
func PullMirrorOf(s *Server) func() string {
	return func() string {
		var c RegistrySettings
		s.st.GetJSON(registrySettingsKey, &c)
		return c.PullMirror
	}
}
