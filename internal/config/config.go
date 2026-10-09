// Package config 负责运行期配置：数据目录、监听地址、宿主路径映射等。
// 全部可用环境变量覆盖，方便在 NAS 面板的 compose 里改。
package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

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
	// HostRoots 允许 Dockhelm 直接访问的宿主目录前缀（逗号分隔），用于 compose 文件路径映射兜底。
	// 例：/volume1/docker,/volume2/apps
	HostRoots []string
	// SelfContainer 自身容器名/ID，默认为空时由容器 hostname 反查。
	SelfContainer string
	// TrustedProxies 传给 gin 的代理信任列表，默认 nil（不信任）。
	TrustedProxies []string
	// LogLevel debug|info|warn|error
	LogLevel string
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
	for _, r := range strings.Split(env("DOCKHELM_HOST_ROOTS", ""), ",") {
		r = strings.TrimSpace(r)
		if r != "" {
			roots = append(roots, filepath.Clean(r))
		}
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

// MapHostPath 把「宿主视角」的路径换算成 Dockhelm 容器内可用的路径。
//
// 推荐部署方式是「冒号两边写一样」（- /volume1/docker:/volume1/docker），此时 label 里的
// 宿主路径在容器内直接可用。但用户也可能写成 /volume1/docker:/compose 这种不一致的形式，
// 这里按 HostRoots 做兜底换算；两条都失败则返回 ok=false，调用方必须显式告知用户
// 「这个路径看不见」，绝不假装读到了。
func (c *Config) MapHostPath(hostPath string) (string, bool) {
	if hostPath == "" {
		return "", false
	}
	clean := filepath.Clean(hostPath)
	// 1) 直接可用（冒号两边一致的情形）
	if _, err := os.Stat(clean); err == nil {
		return clean, true
	}
	// 2) 在 HostRoots 下按「挂载点重映射」找：把宿主路径的每一层前缀依次试成候选
	for _, root := range c.HostRoots {
		if sub, ok := strings.CutPrefix(clean, root); ok {
			cand := filepath.Join(root, sub)
			if _, err := os.Stat(cand); err == nil {
				return cand, true
			}
		}
	}
	// 3) 同名基名搜索（用户把某个子目录单独挂进来了，如 /volume1/docker/moontv:/compose）
	base := filepath.Base(clean)
	parent := filepath.Base(filepath.Dir(clean))
	for _, root := range c.HostRoots {
		// 只试两层，避免递归扫描拖慢请求
		for _, cand := range []string{
			filepath.Join(root, base),
			filepath.Join(root, parent, base),
		} {
			if _, err := os.Stat(cand); err == nil {
				return cand, true
			}
		}
	}
	return "", false
}

// AtoiDefault 便捷解析。
func AtoiDefault(s string, def int) int {
	if n, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
		return n
	}
	return def
}
