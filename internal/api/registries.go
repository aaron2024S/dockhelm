package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

// MirrorConfig 用户配置的一个加速源。
type MirrorConfig struct {
	URL        string `json:"url"`
	Note       string `json:"note"`
	Enabled    bool   `json:"enabled"`
	LastTested string `json:"lastTested,omitempty"`
	LatencyMs  int64  `json:"latencyMs,omitempty"`
	OK         bool   `json:"ok,omitempty"`
	Err        string `json:"err,omitempty"`
	// Builtin 表示这是内置建议列表里带过来的
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

// builtinMirrors 是内置的建议清单。**不再硬编码成"必用"**，
// 只是给用户一个起点，是否启用、顺序如何完全由用户决定。
var builtinMirrors = []MirrorConfig{
	{URL: "https://docker.m.daocloud.io", Note: "DaoCloud", Enabled: false, Builtin: true},
	{URL: "https://docker.1ms.run", Note: "1ms.run", Enabled: false, Builtin: true},
	{URL: "https://dockerproxy.net", Note: "dockerproxy", Enabled: false, Builtin: true},
	{URL: "https://docker.1panel.live", Note: "1Panel", Enabled: false, Builtin: true},
	{URL: "https://hub.rat.dev", Note: "rat.dev", Enabled: false, Builtin: true},
	{URL: "https://dockerhub.icu", Note: "dockerhub.icu", Enabled: false, Builtin: true},
}

// hGetRegistries 返回配置 + 守护进程真实生效的镜像列表。
func (s *Server) hGetRegistries(w http.ResponseWriter, r *http.Request) {
	var cfg RegistrySettings
	s.st.GetJSON("registry.settings", &cfg)
	if cfg.Mirrors == nil {
		cfg.Mirrors = []MirrorConfig{}
	}
	// 把内置但尚未添加的项补到列表尾部（标记 builtin，默认不启用）
	have := map[string]bool{}
	for _, m := range cfg.Mirrors {
		have[normalizeMirror(m.URL)] = true
	}
	suggestions := []MirrorConfig{}
	for _, b := range builtinMirrors {
		if !have[normalizeMirror(b.URL)] {
			suggestions = append(suggestions, b)
		}
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

	writeOK(w, map[string]any{
		"settings":      cfg,
		"suggestions":   suggestions,
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

	if err := s.st.SetJSON("registry.settings", in.Settings); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	// 立刻通知更新引擎换用新的加速源
	s.up.SetMirror(func() string {
		var c RegistrySettings
		s.st.GetJSON("registry.settings", &c)
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
		var cfg RegistrySettings
		s.st.GetJSON("registry.settings", &cfg)
		for _, m := range cfg.Mirrors {
			urls = append(urls, m.URL)
		}
		for _, b := range builtinMirrors {
			urls = append(urls, b.URL)
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
	s.st.GetJSON("registry.settings", &cfg)
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
		_ = s.st.SetJSON("registry.settings", cfg)
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
		s.st.GetJSON("registry.settings", &c)
		return c.PullMirror
	}
}
