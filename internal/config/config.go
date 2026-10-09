// Package config 负责运行期配置：数据目录、监听地址、宿主路径映射等。
// 全部可用环境变量覆盖，方便在 NAS 面板的 compose 里改。
package config

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

// mountPair 一条「宿主机路径 → 容器内路径」的映射。
type mountPair struct {
	host      string
	container string
}

// PathMapping 对外的映射描述，供界面展示与排障。
type PathMapping struct {
	Host      string `json:"host"`
	Container string `json:"container"`
	// Source：auto = 从自身容器的挂载自动识别；env = DOCKHELM_HOST_ROOTS 显式声明。
	Source string `json:"source"`
	// Visible 表示容器内这个路径当前真的存在。
	Visible bool `json:"visible"`
}

// Config 运行期配置。
type Config struct {
	// DataDir 持久化目录（数据库、备份、日志）。默认 /data。
	DataDir string
	// Listen 监听地址，默认 :8080。
	Listen string
	// DockerHost Docker 守护进程地址，默认读 DOCKER_HOST，否则 unix:///var/run/docker.sock。
	DockerHost string
	// ForcePassword 若设置，则启动时把登录密码强制覆盖为该值（忘记密码的兜底手段）。
	ForcePassword string
	// HostRoots 手动声明的宿主目录映射（逗号分隔）。支持两种写法：
	//
	//	/volume1/docker                       两边一致的直通写法
	//	/volume1/docker=/host/docker          宿主机路径=容器内路径
	//
	// 正常情况下不需要配置 —— 启动时会自动读自身容器的 Mounts 得到映射；
	// 这个变量是「自动识别失败」或「想在界面上看得更明确」时的兜底。
	HostRoots []string
	// SelfContainer 自身容器名/ID，默认为空时由容器 hostname 反查。
	SelfContainer string
	// TrustedProxies 传给 gin 的代理信任列表，默认 nil（不信任）。
	TrustedProxies []string
	// LogLevel debug|info|warn|error
	LogLevel string

	mu     sync.RWMutex
	auto   []mountPair // 自动识别（自身容器挂载）
	manual []mountPair // HostRoots 解析结果
}

func env(key, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return def
}

// Load 读取环境变量并返回配置。
func Load() *Config {
	dataDir := env("DOCKHELM_DATA", "/data")
	if !filepath.IsAbs(dataDir) {
		if abs, err := filepath.Abs(dataDir); err == nil {
			dataDir = abs
		}
	}

	roots := []string{}
	manual := []mountPair{}
	for _, r := range strings.Split(env("DOCKHELM_HOST_ROOTS", ""), ",") {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}
		host, container := r, r
		if i := strings.Index(r, "="); i > 0 {
			host = strings.TrimSpace(r[:i])
			container = strings.TrimSpace(r[i+1:])
		}
		host, container = filepath.Clean(host), filepath.Clean(container)
		if host == "" || container == "" {
			continue
		}
		if host != container {
			roots = append(roots, host+"="+container)
		} else {
			roots = append(roots, host)
		}
		manual = append(manual, mountPair{host: host, container: container})
	}

	return &Config{
		DataDir:        dataDir,
		Listen:         env("DOCKHELM_LISTEN", ":8080"),
		DockerHost:     env("DOCKER_HOST", "unix:///var/run/docker.sock"),
		ForcePassword:  strings.TrimSpace(os.Getenv("DOCKHELM_PASSWORD")),
		HostRoots:      roots,
		SelfContainer:  strings.TrimSpace(os.Getenv("DOCKHELM_SELF")),
		TrustedProxies: splitCSV(os.Getenv("DOCKHELM_TRUSTED_PROXIES")),
		LogLevel:       env("DOCKHELM_LOG_LEVEL", "info"),
		manual:         manual,
	}
}

func splitCSV(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// ---- 数据目录下的固定路径 ----

func (c *Config) DBPath() string            { return filepath.Join(c.DataDir, "dockhelm.db") }
func (c *Config) AuthPath() string          { return filepath.Join(c.DataDir, "auth.json") }
func (c *Config) BackupDir() string         { return filepath.Join(c.DataDir, "backups") }
func (c *Config) ContainerBackupDir() string { return filepath.Join(c.BackupDir(), "containers") }
func (c *Config) ProjectBackupDir() string  { return filepath.Join(c.BackupDir(), "projects") }
func (c *Config) LogDir() string            { return filepath.Join(c.DataDir, "logs") }

// EnsureDirs 创建所有需要的目录。
func (c *Config) EnsureDirs() error {
	for _, d := range []string{c.DataDir, c.BackupDir(), c.ContainerBackupDir(), c.ProjectBackupDir(), c.LogDir()} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return err
		}
	}
	return nil
}

// MapHostPath 把「宿主视角」的路径换算成 Dockhelm 容器内真正读得到的路径。
//
// 换算优先级：
//  1. 挂载映射 —— 启动时从自身容器的 Mounts 自动识别的，加上 DOCKHELM_HOST_ROOTS 显式声明的；
//     按宿主机路径「最长前缀」匹配，命中后把这段前缀换成容器内路径；
//  2. 原样可用 —— 也就是 - /volume1/docker:/volume1/docker 这种两边一致的写法；
//  3. 基名启发式 —— 用户只把某个子目录单独挂了进来（如 /volume1/docker/moontv:/compose）。
//
// 全都失败就返回 ok=false：调用方必须显式告诉用户「这个路径看不见」，绝不假装读到了。
func (c *Config) MapHostPath(hostPath string) (string, bool) {
	if strings.TrimSpace(hostPath) == "" {
		return "", false
	}
	clean := filepath.Clean(hostPath)
	pairs := c.effective()

	// 1) 前缀替换（最长前缀优先，见 effective 的排序）
	for _, m := range pairs {
		if cand, ok := rewritePrefix(clean, m.host, m.container); ok && pathExists(cand) {
			return cand, true
		}
	}
	// 2) 原样可用
	if pathExists(clean) {
		return clean, true
	}
	// 3) 基名启发式：只试两层，避免递归扫描拖慢请求
	base := filepath.Base(clean)
	parent := filepath.Base(filepath.Dir(clean))
	for _, m := range pairs {
		for _, cand := range []string{
			filepath.Join(m.container, base),
			filepath.Join(m.container, parent, base),
		} {
			if pathExists(cand) {
				return cand, true
			}
		}
	}
	return "", false
}

// SetMounts 用自身容器的挂载信息（宿主机路径 → 容器内路径）填充自动映射。
// 只在启动时调用一次，覆盖语义。
func (c *Config) SetMounts(m map[string]string) {
	pairs := make([]mountPair, 0, len(m))
	for host, container := range m {
		host, container = strings.TrimSpace(host), strings.TrimSpace(container)
		if host == "" || container == "" {
			continue
		}
		pairs = append(pairs, mountPair{host: filepath.Clean(host), container: filepath.Clean(container)})
	}
	c.mu.Lock()
	c.auto = pairs
	c.mu.Unlock()
}

// PathMappings 返回当前生效的全部映射（供界面展示与排障），顺序即实际匹配顺序。
func (c *Config) PathMappings() []PathMapping {
	auto, manual := c.raw()
	out := make([]PathMapping, 0, len(auto)+len(manual))
	for _, p := range auto {
		out = append(out, PathMapping{Host: p.host, Container: p.container, Source: "auto", Visible: pathExists(p.container)})
	}
	for _, p := range manual {
		out = append(out, PathMapping{Host: p.host, Container: p.container, Source: "env", Visible: pathExists(p.container)})
	}
	sort.Slice(out, func(i, j int) bool {
		if len(out[i].Host) != len(out[j].Host) {
			return len(out[i].Host) > len(out[j].Host)
		}
		return out[i].Host < out[j].Host
	})
	return out
}

// Summary 一行文字描述当前映射，给启动日志用。
func (c *Config) Summary() string {
	ms := c.PathMappings()
	if len(ms) == 0 {
		return "没有识别到任何挂载映射（compose 文件路径将按原样尝试）"
	}
	parts := make([]string, 0, len(ms))
	for _, m := range ms {
		s := m.Host + " → " + m.Container
		if m.Host == m.Container {
			s = m.Host + "（两边一致）"
		}
		if !m.Visible {
			s += " [容器内不存在]"
		}
		if m.Source == "env" {
			s += " [env]"
		}
		parts = append(parts, s)
	}
	return strings.Join(parts, ", ")
}

func (c *Config) raw() (auto, manual []mountPair) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return append([]mountPair(nil), c.auto...), append([]mountPair(nil), c.manual...)
}

// effective 返回合并去重、按宿主机路径长度降序的映射 —— 顺序就是「最长前缀优先」的匹配顺序。
func (c *Config) effective() []mountPair {
	auto, manual := c.raw()
	all := make([]mountPair, 0, len(auto)+len(manual))
	all = append(all, auto...)
	all = append(all, manual...)
	sort.SliceStable(all, func(i, j int) bool { return len(all[i].host) > len(all[j].host) })

	out := make([]mountPair, 0, len(all))
	seen := map[string]bool{}
	for _, p := range all {
		if p.host == "" || seen[p.host] {
			continue
		}
		seen[p.host] = true
		out = append(out, p)
	}
	return out
}

// rewritePrefix 把 p 里的 from 前缀换成 to。按路径分段匹配，避免 /vol 误命中 /volume1。
func rewritePrefix(p, from, to string) (string, bool) {
	if from == "" || from == "/" {
		return "", false
	}
	if p == from {
		return to, true
	}
	if strings.HasPrefix(p, from+"/") {
		return filepath.Join(to, strings.TrimPrefix(p, from+"/")), true
	}
	return "", false
}

func pathExists(p string) bool {
	if p == "" {
		return false
	}
	_, err := os.Stat(p)
	return err == nil
}

// AtoiDefault 便捷解析。
func AtoiDefault(s string, def int) int {
	if n, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
		return n
	}
	return def
}
