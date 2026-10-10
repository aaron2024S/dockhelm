// Package version 保存应用版本与作者信息。
// 版本号是唯一事实源，前端「关于」页与 /api/about 都从这里取。
package version

import (
	"runtime"
	"runtime/debug"
	"strings"
	"time"
)

const (
	// AppName 应用显示名。
	AppName = "Dockhelm"
	// AppNameCN 中文名（用于页面副标题）。
	AppNameCN = "Dockhelm · 容器舵手"
	// Tagline 一句话介绍。
	Tagline = "看得清、控得准的 Docker 容器管理面板"
	// Description 详细介绍（关于页展示）。
	Description = "Dockhelm 是一个自托管的 Docker 容器管理面板：容器与镜像管理、同源更新的更新中心、" +
		"定时启停计划任务、可自由配置的镜像加速源、配置快照备份与还原、以及多渠道事件通知。" +
		"它由 Docker 守护进程本身来解析镜像仓库，因此「检测到的更新」与「实际拉取的镜像」永远是同一个来源。"
	// Author 作者。
	Author = "aaron2024S"
	// AuthorURL 作者主页。
	AuthorURL = "https://github.com/aaron2024S"
	// RepoURL 项目仓库。
	RepoURL = "https://github.com/aaron2024S/dockhelm"
	// License 许可协议。
	License = "MIT"
)

// Version 语义化版本号。发新版时改这里。
// 声明为 var 而不是 const，是为了让构建时可以用
// `-ldflags "-X .../internal/version.Version=x.y.z"` 注入。
var Version = "0.3.1"

// BuildInfo 由 ldflags 注入（见 Dockerfile / 构建脚本），默认值用于本地 `go run`。
var (
	Commit    = ""
	BuildTime = ""
	GoVersion = ""
)

func init() {
	// 本地直接构建时从 Go 自身读取（-ldflags 未注入时的兜底）。
	if Commit == "" {
		if bi, ok := debug.ReadBuildInfo(); ok {
			for _, s := range bi.Settings {
				if s.Key == "vcs.revision" {
					Commit = s.Value
				}
				if s.Key == "vcs.time" && BuildTime == "" {
					BuildTime = s.Value
				}
			}
		}
	}
	if GoVersion == "" {
		GoVersion = runtime.Version()
	}
}

// Info 返回给 API 的结构。
type Info struct {
	Name        string `json:"name"`
	NameCN      string `json:"nameCN"`
	Version     string `json:"version"`
	Tagline     string `json:"tagline"`
	Description string `json:"description"`
	Author      string `json:"author"`
	AuthorURL   string `json:"authorURL"`
	RepoURL     string `json:"repoURL"`
	License     string `json:"license"`
	Commit      string `json:"commit"`
	CommitShort string `json:"commitShort"`
	BuildTime   string `json:"buildTime"`
	GoVersion   string `json:"goVersion"`
	Platform    string `json:"platform"`
	StartedAt   string `json:"startedAt"`
}

// Runtime 进程启动时间（About 页显示"已运行"）。
var Runtime = time.Now()

// Get 组装当前构建信息。
func Get() Info {
	commit := strings.TrimSpace(Commit)
	short := commit
	if len(short) > 7 {
		short = short[:7]
	}
	return Info{
		Name:        AppName,
		NameCN:      AppNameCN,
		Version:     Version,
		Tagline:     Tagline,
		Description: Description,
		Author:      Author,
		AuthorURL:   AuthorURL,
		RepoURL:     RepoURL,
		License:     License,
		Commit:      commit,
		CommitShort: short,
		BuildTime:   BuildTime,
		GoVersion:   GoVersion,
		Platform:    runtime.GOOS + "/" + runtime.GOARCH,
		StartedAt:   Runtime.UTC().Format(time.RFC3339),
	}
}
