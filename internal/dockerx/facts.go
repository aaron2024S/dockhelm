package dockerx

import (
	"context"
	"strings"
)

// ContainerFacts 一次 inspect 里「顺序启动」真正需要的几个事实。
//
// 为什么不直接让调用方啃 Inspect 返回的 map[string]any：
// 「这个容器归不归本程序管」取决于 HostConfig.RestartPolicy.Name，
// 「现在要不要启动它」取决于 State，两处都要逐层类型断言 ——
// 放在业务代码里会写成一串开关。抽成结构体，顺便可以脱离 Docker 单测。
type ContainerFacts struct {
	ID       string
	Name     string // 去掉前导斜线的容器名
	Image    string
	State    string // running / exited / paused / created / dead …
	Policy   string // no / on-failure / always / unless-stopped
	MaxRetry int    // 仅 on-failure 有意义，其余为 0
	Project  string // compose 项目名；非 compose 容器为空
}

// Facts 读取容器当前状态、重启策略与 compose 项目（一次 inspect 拿全）。
func (c *Client) Facts(ctx context.Context, idOrName string) (ContainerFacts, error) {
	insp, err := c.Inspect(ctx, idOrName)
	if err != nil {
		return ContainerFacts{}, err
	}
	return FactsFromInspect(insp), nil
}

// RestartPolicy 只取重启策略（语义清晰的窄接口，页面上「策略」一列用它）。
func (c *Client) RestartPolicy(ctx context.Context, idOrName string) (name string, maxRetry int, err error) {
	f, err := c.Facts(ctx, idOrName)
	if err != nil {
		return "", 0, err
	}
	return f.Policy, f.MaxRetry, nil
}

// FactsFromInspect 从 inspect 的原始结果里抽事实（纯函数，便于单测）。
func FactsFromInspect(insp map[string]any) ContainerFacts {
	f := ContainerFacts{}
	f.ID, _ = insp["Id"].(string)
	f.Name = strings.TrimPrefix(jsonStr(insp["Name"]), "/")

	cfg, _ := insp["Config"].(map[string]any)
	f.Image, _ = cfg["Image"].(string)
	if labels, ok := cfg["Labels"].(map[string]any); ok {
		f.Project, _ = labels["com.docker.compose.project"].(string)
	}

	if st, ok := insp["State"].(map[string]any); ok {
		f.State, _ = st["Status"].(string)
		// 某些精简实现只给 Running 布尔值
		if f.State == "" {
			if running, _ := st["Running"].(bool); running {
				f.State = "running"
			} else {
				f.State = "exited"
			}
		}
	}

	if hc, ok := insp["HostConfig"].(map[string]any); ok {
		if rp, ok := hc["RestartPolicy"].(map[string]any); ok {
			f.Policy, _ = rp["Name"].(string)
			if n, ok := rp["MaximumRetryCount"].(float64); ok {
				f.MaxRetry = int(n)
			}
		}
	}
	return f
}

func jsonStr(v any) string {
	s, _ := v.(string)
	return s
}
