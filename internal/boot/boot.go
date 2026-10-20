// Package boot 实现「顺序启动」：dockhelm 起来之后，按用户排好的清单，
// 串行地把容器一个一个拉起来。
//
// 边界（2026-10-19 定稿，别擅自扩展）：
//
//   - 只管**启动**；关机的逆序停止不做；
//   - **不改**容器的重启策略（用户自己在 NAS 面板改，程序只读出来给页面显示）；
//   - 没有快照、没有后台采样 —— 唯一依据是清单里每行的 auto 开关；
//   - 只碰「受管」容器（restartPolicy 为 no / on-failure）。always / unless-stopped
//     的容器 Docker 会在开机时自己拉起，比本功能更早，改了开关和延时对它们也没用，
//     所以一概不碰，页面上单独一组只读展示。
//
// ⚠ 运行时行为（方案 §八，别当成 bug 改）：**每次进程启动都会执行一次**，
// 包括镜像更新导致的自重启。那种情况下绝大多数容器在跑、会被跳过，
// 只有「开关开着但当前停着」的容器会被拉起 —— 这是预期行为。
package boot

import (
	"fmt"
	"sort"
	"strings"
)

// 存储键。只用现有的 GetJSON / SetJSON，不改 DB schema。
const (
	planKey = "boot.plan"
	// lastKey 上一次执行的结构化结果。
	//
	// 为什么不用运行记录现算：页面顶部那一行要写成「上次执行：今天 08:12 · 成功 47、跳过 3」，
	// 而运行记录里只有一句完整的话（「顺序启动 50 个容器，成功 47、跳过 3」）。
	// 从那句话里正则抠数字既脆弱又蠢，所以顺手存一份结构化的。
	// 它只是个展示用的回执，不是快照、不参与任何执行决策。
	lastKey = "boot.last"
)

// 默认值。首次打开页面、或参数缺省时用它。
const (
	DefaultInitialDelaySec = 30 // 等 Docker 稳定下来的初始等待
	DefaultWaitTimeoutSec  = 60 // 单个容器等它进入 running 的上限
	DefaultDelaySec        = 3  // 默认延时：新加入的容器用它，也是「应用到全部行」的来源

	maxInitialDelaySec = 1800
	maxWaitTimeoutSec  = 600
	minWaitTimeoutSec  = 5
	maxDelaySec        = 300
)

// Item 清单里的一行。数组顺序就是启动次序，拖拽改的就是它。
type Item struct {
	Name     string `json:"name"`
	Auto     bool   `json:"auto"`     // 开机启动开关；关掉 = 永远不参与
	DelaySec int    `json:"delaySec"` // 这一行**启动成功后**等 N 秒再起下一个
}

// Plan 整份清单 + 执行参数。
type Plan struct {
	Enabled         bool   `json:"enabled"`
	InitialDelaySec int    `json:"initialDelaySec"`
	WaitTimeoutSec  int    `json:"waitTimeoutSec"`
	DefaultDelaySec int    `json:"defaultDelaySec"`
	Items           []Item `json:"items"`
}

// DefaultPlan 出厂清单：开关为开、延时取默认。
func DefaultPlan() Plan {
	return Plan{
		Enabled:         true,
		InitialDelaySec: DefaultInitialDelaySec,
		WaitTimeoutSec:  DefaultWaitTimeoutSec,
		DefaultDelaySec: DefaultDelaySec,
		Items:           []Item{},
	}
}

// Normalize 把参数与条目收敛到合法区间（负数/越界一律纠正，不报错）。
func (p Plan) Normalize() Plan {
	if p.InitialDelaySec < 0 {
		p.InitialDelaySec = DefaultInitialDelaySec
	}
	if p.InitialDelaySec > maxInitialDelaySec {
		p.InitialDelaySec = maxInitialDelaySec
	}
	// 0 是**合法值**（不等待），只有「没填」才会走到默认值 —— 但 JSON 里两者都是 0，
	// 所以这里只拦越界：缺省由 LoadPlan 的 DefaultPlan 兜底，不在这里猜。
	if p.WaitTimeoutSec <= 0 {
		p.WaitTimeoutSec = DefaultWaitTimeoutSec
	}
	if p.WaitTimeoutSec < minWaitTimeoutSec {
		p.WaitTimeoutSec = minWaitTimeoutSec
	}
	if p.WaitTimeoutSec > maxWaitTimeoutSec {
		p.WaitTimeoutSec = maxWaitTimeoutSec
	}
	if p.DefaultDelaySec < 0 || p.DefaultDelaySec > maxDelaySec {
		p.DefaultDelaySec = DefaultDelaySec
	}

	out := make([]Item, 0, len(p.Items))
	seen := map[string]bool{}
	for _, it := range p.Items {
		it.Name = strings.TrimSpace(it.Name)
		if it.Name == "" || seen[it.Name] {
			continue
		}
		seen[it.Name] = true
		if it.DelaySec < 0 {
			it.DelaySec = p.DefaultDelaySec
		}
		if it.DelaySec > maxDelaySec {
			it.DelaySec = maxDelaySec
		}
		out = append(out, it)
	}
	p.Items = out
	return p
}

// ContainerInfo 一次容器快照里本包关心的字段。
type ContainerInfo struct {
	Name    string `json:"name"`
	ID      string `json:"id"`
	State   string `json:"state"`
	Policy  string `json:"restartPolicy"`
	Project string `json:"project"`
}

// Managed 这个容器归不归本功能管 —— 判定只看重启策略：
//
//	no / on-failure → 归（开机时 Docker 不会自己拉起它）
//	always / unless-stopped → 不归（Docker 会先于本功能把它拉起来）
//
// 策略读不到（空串）时按「不归」处理：宁可不动，也不要在一个状态未知的容器上
// 替用户做决定。
func Managed(policy string) bool {
	return policy == "no" || policy == "on-failure"
}

// Merge 把清单与当前容器对齐，返回新的清单：
//
//   - 已删除的容器：剔除（它的条目留着也没意义）；
//   - 新出现的容器：追加到末尾，auto=true、delaySec=DefaultDelaySec；
//   - 已有条目：**顺序与开关原样保留**（哪怕它此刻归 Docker 管，位置也不丢）。
//
// 追加顺序 = compose 项目分组 + 项目内按名称；没有 compose 项目的排在最后。
// 这就是「第一次打开页面时」的初始顺序（方案 §八 第 3 条）。
func (p Plan) Merge(infos []ContainerInfo) Plan {
	exists := make(map[string]ContainerInfo, len(infos))
	for _, c := range infos {
		exists[c.Name] = c
	}

	kept := make([]Item, 0, len(p.Items))
	seen := map[string]bool{}
	for _, it := range p.Items {
		if _, ok := exists[it.Name]; !ok {
			continue
		}
		if seen[it.Name] {
			continue
		}
		seen[it.Name] = true
		kept = append(kept, it)
	}

	added := make([]ContainerInfo, 0, len(infos))
	for _, c := range infos {
		if !seen[c.Name] {
			added = append(added, c)
		}
	}
	sort.Slice(added, func(i, j int) bool {
		a, b := added[i], added[j]
		if a.Project != b.Project {
			// 没有 compose 项目的散装容器排在有项目名的后面
			if a.Project == "" {
				return false
			}
			if b.Project == "" {
				return true
			}
			return a.Project < b.Project
		}
		return a.Name < b.Name
	})
	for _, c := range added {
		kept = append(kept, Item{Name: c.Name, Auto: true, DelaySec: p.DefaultDelaySec})
	}

	p.Items = kept
	return p
}

// StepKind 一次扫描里某一行的去向。
type StepKind string

const (
	StepStart     StepKind = "start"     // 停着 ⇒ 启动
	StepRunning   StepKind = "running"   // 已在跑 ⇒ 跳过（计入统计）
	StepDisabled  StepKind = "disabled"  // 开关关着 ⇒ 跳过（**不**计入统计）
	StepUnmanaged StepKind = "unmanaged" // 归 Docker 管 ⇒ 不碰（**不**计入统计）
	StepMissing   StepKind = "missing"   // 清单里有、Docker 里没了 ⇒ 记失败
)

// Step 扫描结果的一行。
type Step struct {
	Name     string
	Kind     StepKind
	DelaySec int
	Policy   string
	State    string
}

// Counted 这一行是否进入统计（运行记录里的「N 个容器」）。
//
// 统计口径（方案 §2.4）：只有「受管 + 开关开着」的才算。
// 归 Docker 管的那组既不算成功也不算跳过，完全不进数字。
func (s Step) Counted() bool {
	switch s.Kind {
	case StepStart, StepRunning, StepMissing:
		return true
	}
	return false
}

// Resolve 按清单顺序逐行决定这次要做什么（纯函数，不碰 Docker）。
func Resolve(p Plan, infos []ContainerInfo) []Step {
	byName := make(map[string]ContainerInfo, len(infos))
	for _, c := range infos {
		byName[c.Name] = c
	}
	steps := make([]Step, 0, len(p.Items))
	for _, it := range p.Items {
		c, ok := byName[it.Name]
		if !ok {
			steps = append(steps, Step{Name: it.Name, Kind: StepMissing, DelaySec: it.DelaySec})
			continue
		}
		st := Step{Name: it.Name, DelaySec: it.DelaySec, Policy: c.Policy, State: c.State}
		switch {
		case !Managed(c.Policy):
			st.Kind = StepUnmanaged
		case !it.Auto:
			st.Kind = StepDisabled
		case c.State == "running":
			st.Kind = StepRunning
		default:
			st.Kind = StepStart
		}
		steps = append(steps, st)
	}
	return steps
}

// startCostSec 估算「启动一个容器并等它 running」的平均耗时。
// 只用于页面上的「预计 X 分 Y 秒」，不参与任何执行决策。
const startCostSec = 8

// EstimateSeconds 估算一次顺序启动要多久：初始等待 + 每个要启动的容器的
// （平均启动耗时 + 该行延时）。已在运行/被跳过的不耗时间。
func EstimateSeconds(p Plan, steps []Step) int {
	total := p.InitialDelaySec
	for _, st := range steps {
		if st.Kind == StepStart {
			total += startCostSec + st.DelaySec
		}
	}
	return total
}

// maxSummaryNames 运行记录里最多点几个失败容器，超出折叠。
const maxSummaryNames = 8

// Summary 运行记录 / 「上次执行」的一行简报（方案 §六）。
//
//	全部正常：顺序启动 50 个容器，成功 47、跳过 3
//	有失败：  顺序启动 50 个容器，成功 45、跳过 3、失败 2：nginx、foo
//	全跳过：  顺序启动：50 个容器都已在运行，本次未启动任何容器
func Summary(total, success, skipped int, failed []string) string {
	switch {
	case total == 0:
		return "顺序启动：没有需要启动的容器"
	case len(failed) > 0:
		return fmt.Sprintf("顺序启动 %d 个容器，成功 %d、跳过 %d、失败 %d：%s",
			total, success, skipped, len(failed), foldNames(failed))
	case success == 0:
		return fmt.Sprintf("顺序启动：%d 个容器都已在运行，本次未启动任何容器", total)
	default:
		return fmt.Sprintf("顺序启动 %d 个容器，成功 %d、跳过 %d", total, success, skipped)
	}
}

// foldNames 排序 + 截断，让简报里的名单稳定可读。
func foldNames(names []string) string {
	sorted := append([]string(nil), names...)
	sort.Strings(sorted)
	if len(sorted) > maxSummaryNames {
		return strings.Join(sorted[:maxSummaryNames], "、") + fmt.Sprintf("…等 %d 个", len(sorted))
	}
	return strings.Join(sorted, "、")
}
