// Package autoupdate 决定「检测到有更新之后，要不要自动动容器」。
//
// 这里只放决策，不碰 Docker —— 纯函数才写得动单测，也才能让前端拿到一份
// 「哪些会被更新、哪些被跳过、为什么」的完整预览。
//
// 语义上刻意与 dockerCopilot 的 autoupdateplan 对齐，但有一处刻意的差别：
// dockerCopilot 先判「未检测到更新」再判排除，于是被排除的容器在预览里
// 一律显示成「未检测到更新」，看不出是用户主动排除的。这里把保护 / 排除
// 提到前面判，预览里能直接看出「这个容器是你自己排除掉的」。
package autoupdate

import (
	"path"
	"strings"
)

// selfKeyword 兜底识别自身容器。除了配置里的自身名，容器名里带这个关键字
// 的一律当作自身 —— 否则一旦改名，自动更新会把自己重启掉，
// 而「更新到一半把自己停了」是最难排查的故障。
const selfKeyword = "dockhelm"

// Container 是筛选阶段需要的容器最小信息集。
type Container struct {
	Name      string
	Image     string
	ImageID   string
	Running   bool
	HasUpdate bool
}

// Candidate 单个容器的筛选结论，直接序列化给前端做候选预览。
type Candidate struct {
	Name       string `json:"name"`
	Image      string `json:"image"`
	Running    bool   `json:"running"`
	HasUpdate  bool   `json:"hasUpdate"`
	WillUpdate bool   `json:"willUpdate"`
	// Protected 命中自身保护（优先级高于排除列表）。
	Protected bool `json:"protected"`
	// Excluded 命中用户配置的排除列表。
	Excluded bool `json:"excluded"`
	// Reason 被跳过的原因；WillUpdate 为真时为空串。
	Reason string `json:"reason"`
}

// Options 筛选参数。
type Options struct {
	// SelfName Dockhelm 自身容器名（与关键字匹配取或）。
	SelfName string
	// Exclude 排除列表，支持 redis* 这类通配写法。
	Exclude []string
}

// Plan 依据配置筛选出本轮应当自动更新的容器。
func Plan(containers []Container, opt Options) []Candidate {
	out := make([]Candidate, 0, len(containers))
	for _, c := range containers {
		item := Candidate{
			Name:      c.Name,
			Image:     c.Image,
			Running:   c.Running,
			HasUpdate: c.HasUpdate,
		}
		switch {
		case IsSelf(c.Name, opt.SelfName):
			item.Protected = true
			item.Reason = "Dockhelm 自身，永不自动更新"
		case MatchesExclude(opt.Exclude, c.Name, c.Image):
			item.Excluded = true
			item.Reason = "命中排除列表"
		case !c.HasUpdate:
			item.Reason = "未检测到更新"
		default:
			item.WillUpdate = true
		}
		out = append(out, item)
	}
	return out
}

// Names 取出会被更新的容器名，保持传入顺序（排序稳定的预览更好读）。
func Names(cands []Candidate) []string {
	out := make([]string, 0, len(cands))
	for _, c := range cands {
		if c.WillUpdate {
			out = append(out, c.Name)
		}
	}
	return out
}

// Count 统计会被更新的容器数。
func Count(cands []Candidate) int {
	n := 0
	for _, c := range cands {
		if c.WillUpdate {
			n++
		}
	}
	return n
}

// IsSelf 判断容器名是否指向 Dockhelm 自身。
func IsSelf(name, selfName string) bool {
	if name == "" {
		return false
	}
	if selfName != "" && strings.EqualFold(name, selfName) {
		return true
	}
	return strings.Contains(strings.ToLower(name), selfKeyword)
}

// MatchesExclude 判断容器是否命中排除列表。
// 容器名、镜像全名、镜像仓库名三者都参与匹配，支持 redis* 这类通配写法。
func MatchesExclude(exclude []string, containerName, image string) bool {
	if len(exclude) == 0 {
		return false
	}
	repo := StripTag(image)
	for _, pattern := range exclude {
		if MatchPattern(pattern, containerName) ||
			MatchPattern(pattern, image) ||
			MatchPattern(pattern, repo) {
			return true
		}
	}
	return false
}

// MatchPattern 大小写不敏感地做精确匹配或通配匹配。
// 通配表达式本身非法时只退化成精确匹配，不中断整个匹配流程。
func MatchPattern(pattern, value string) bool {
	if pattern == "" || value == "" {
		return false
	}
	lowerPattern := strings.ToLower(pattern)
	lowerValue := strings.ToLower(value)
	if lowerPattern == lowerValue {
		return true
	}
	matched, err := path.Match(lowerPattern, lowerValue)
	if err != nil {
		return false
	}
	return matched
}

// StripTag 去掉镜像的 tag，只看仓库名。
// 需要同时考虑 registry 端口号里的冒号，例如 registry.local:5000/nginx:latest。
func StripTag(image string) string {
	lastSlash := strings.LastIndex(image, "/")
	lastColon := strings.LastIndex(image, ":")
	if lastColon > lastSlash {
		return image[:lastColon]
	}
	return image
}
