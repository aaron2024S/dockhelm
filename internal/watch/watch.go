// Package watch 订阅 Docker 事件流，把「容器意外退出 / 崩溃循环 / 健康检查失败」
// 以及「守护进程断开又恢复」翻译成通知事件。
package watch

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aaron2024s/dockhelm/internal/bus"
	"github.com/aaron2024s/dockhelm/internal/dockerx"
	"github.com/aaron2024s/dockhelm/internal/intent"
	"github.com/aaron2024s/dockhelm/internal/notify"
	"github.com/aaron2024s/dockhelm/internal/store"
)

// Watcher 事件观察者。
type Watcher struct {
	dc   *dockerx.Client
	nt   *notify.Manager
	st   *store.Store
	bus  *bus.Bus
	self func() string

	mu       sync.Mutex
	dieTimes map[string][]time.Time // 容器 -> 最近退出时间
	killAt   map[string]time.Time   // 容器 -> 最近一次收到停止信号（kill 事件）的时间
	oomAt    map[string]time.Time   // 容器 -> 最近一次 OOM 的时间
	healthy  bool
}

// New 创建观察者。
func New(dc *dockerx.Client, nt *notify.Manager, st *store.Store, b *bus.Bus, selfFn func() string) *Watcher {
	return &Watcher{
		dc: dc, nt: nt, st: st, bus: b, self: selfFn,
		dieTimes: map[string][]time.Time{},
		killAt:   map[string]time.Time{},
		oomAt:    map[string]time.Time{},
		healthy:  true,
	}
}

// 崩溃循环判定：10 分钟内退出 3 次以上。
const (
	crashWindow = 10 * time.Minute
	crashCount  = 3
)

// 外部主动停止的判定窗口。docker stop 先发 SIGTERM（kill 事件）再等宽限期，
// compose 的 stop_grace_period 可能配到分钟级，窗口要盖得住整段等待。
// 判定只在这个窗口内成立，过期即失效 —— 免得把随后的一次真崩溃也掩护掉。
const stopWindow = 3 * time.Minute

// Run 启动事件循环与连通性巡检，直到 ctx 取消。
func (w *Watcher) Run(ctx context.Context) {
	go w.pingLoop(ctx)
	for {
		if ctx.Err() != nil {
			return
		}
		if err := w.dc.Events(ctx, w.handle); err != nil {
			if ctx.Err() != nil {
				return
			}
			log.Printf("docker 事件流中断，2 秒后重连：%v", err)
			select {
			case <-ctx.Done():
				return
			case <-time.After(2 * time.Second):
			}
		}
	}
}

func (w *Watcher) handle(ev dockerx.Event) {
	name := ev.Actor.Attributes["name"]
	if name == "" {
		name = ev.Actor.ID
	}
	// Dockhelm 自己的容器不参与容器级告警（否则更新自己会刷屏）
	if self := w.self(); self != "" && name == self {
		return
	}

	// 状态真的变了的动作，统一在这里推一条给前端（总览页、容器页据此静默刷新）。
	// 必须放在下面那个 switch 之前：里面每条分支都可能提前 return（正常退出、dockhelm
	// 自己发起的停止…），放后面会让「停止」这类最该刷新的事件漏掉。
	// 走 container 主题 + kind=changed，和 API 里手动操作的 kind=action 区分开，
	// 容器页只对 action 做强制刷新，自动变化走防抖。
	switch ev.Action {
	case "start", "restart", "stop", "die", "kill", "pause", "unpause", "destroy":
		w.bus.Publish("container", "changed", "info", map[string]any{
			"name": name, "action": ev.Action,
		})
	}

	switch {
	case ev.Action == "die":
		code := ev.Actor.Attributes["exitCode"]
		n, _ := strconv.Atoi(code)
		// 0 = 正常结束，不告警
		if n == 0 {
			return
		}
		// dockhelm 自己发起的停止（更新 / 还原 / 计划任务 / 手动停）不是「意外退出」
		// —— stop 宽限期一到 Docker 发 SIGKILL，退出码就是 137。
		if intent.Consume(name) {
			return
		}
		// 不是 dockhelm 发起的，但可能仍是**外部主动停止**：NAS 面板 / docker CLI /
		// compose 重新部署都会走 docker stop，而它会先发 kill 事件（SIGTERM）再 die。
		// 真崩溃与 OOM 不经过这条路（后者另有 oom 事件），照常告警。
		if w.exitedByExternalStop(name, time.Now()) {
			return
		}
		image := ev.Actor.Attributes["image"]
		w.recordDie(name, image)
		w.nt.Emit("container_died", map[string]string{
			"container": name, "image": image,
			"result":  "容器意外退出",
			"message": "退出码 " + exitCodeHint(n),
		})
		w.st.AddRunLog("container", name, "failed", "容器意外退出，退出码 "+code, "")

	case strings.HasPrefix(ev.Action, "health_status"):
		if strings.Contains(ev.Action, "unhealthy") {
			// 更新 / 重启后的头几分钟新容器短暂 unhealthy 是常态，等它自己变回来
			if intent.Active(name) {
				return
			}
			image := ev.Actor.Attributes["image"]
			w.nt.Emit("health_failed", map[string]string{
				"container": name, "image": image,
				"result":  "健康检查失败",
				"message": "容器健康状态变为 unhealthy",
			})
			w.st.AddRunLog("container", name, "failed", "健康检查失败（unhealthy）", "")
		}

	case ev.Action == "kill":
		// 主动停止的第一步：docker stop 发 SIGTERM、docker kill 发 SIGKILL，
		// 守护进程都会先推一条 kill（Attributes.signal）。记下来给 die 判定用。
		w.markAt(w.killAt, name, time.Now())

	case ev.Action == "oom":
		// 内核 OOM 杀的容器：这是真告警，不能因为「有 kill 事件」被掩护掉
		// （个别守护进程版本在 OOM 时也会补一条 kill）。
		w.markAt(w.oomAt, name, time.Now())

		// start / restart 的刷新事件已经在 switch 之前的统一发布块里发过，这里不再重复。
	}
}

// markAt 记录「容器 -> 时间」，并顺手清掉过期项（容器多了也不让 map 无限长）。
func (w *Watcher) markAt(m map[string]time.Time, name string, now time.Time) {
	w.mu.Lock()
	defer w.mu.Unlock()
	m[name] = now
	for k, t := range m {
		if now.Sub(t) > stopWindow {
			delete(m, k)
		}
	}
}

// exitedByExternalStop 判断这次非零退出是否由「外部主动停止」引起：
// stopWindow 内该容器出现过 kill 事件，且同期没有 oom 事件。
// 命中后清掉 kill 记录 —— 后续窗口内若再来一次没有 kill 前置的退出，照常告警。
func (w *Watcher) exitedByExternalStop(name string, now time.Time) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	k, ok := w.killAt[name]
	if !ok || now.Sub(k) > stopWindow {
		return false
	}
	if o, ok := w.oomAt[name]; ok && now.Sub(o) <= stopWindow {
		return false
	}
	delete(w.killAt, name)
	return true
}

// recordDie 记录退出时间，并在窗口内次数超阈值时发一次崩溃循环告警。
func (w *Watcher) recordDie(name, image string) {
	w.mu.Lock()
	now := time.Now()
	list := append(w.dieTimes[name], now)
	cut := now.Add(-crashWindow)
	kept := list[:0]
	for _, t := range list {
		if t.After(cut) {
			kept = append(kept, t)
		}
	}
	w.dieTimes[name] = kept
	count := len(kept)
	w.mu.Unlock()

	if count >= crashCount {
		w.nt.Emit("container_crashloop", map[string]string{
			"container": name, "image": image,
			"result":  "容器崩溃循环",
			"message": "最近 " + w.humanWindow() + " 内已退出 " + strconv.Itoa(count) + " 次",
		})
	}
}

// exitCodeHint 把退出码翻译成人能看懂的一句话。
//
// 「退出码 137」这种在推送里等于没说：137 = 128+9 = 被 SIGKILL 强制终止，
// 而用户真正要看的是「被谁、以什么方式杀掉」。128+n 这个约定是 Docker 里最常见
// 的一类退出码，统一按信号名解释；其余自解释（0/1/2）或原样返回。
func exitCodeHint(code int) string {
	if name, ok := signalNames[code-128]; code > 128 && ok {
		return fmt.Sprintf("%d（被 %s 信号终止）", code, name)
	}
	switch code {
	case 1:
		return "1（应用自身报错退出）"
	case 2:
		return "2（命令行用法错误）"
	case 126:
		return "126（命令无法执行：没有可执行权限）"
	case 127:
		return "127（命令未找到：入口点或解释器缺失）"
	default:
		return strconv.Itoa(code)
	}
}

// signalNames 128+n 里 n 对应的信号名（覆盖常见的前 16 个）。
var signalNames = map[int]string{
	1:  "SIGHUP",
	2:  "SIGINT",
	3:  "SIGQUIT",
	4:  "SIGILL",
	6:  "SIGABRT",
	8:  "SIGFPE",
	9:  "SIGKILL",
	11: "SIGSEGV",
	13: "SIGPIPE",
	14: "SIGALRM",
	15: "SIGTERM",
}

func (w *Watcher) humanWindow() string {
	return strconv.Itoa(int(crashWindow.Minutes())) + " 分钟"
}

// pingLoop 定期探活，把「断开 / 恢复」各通知一次（而不是每 30 秒刷一条）。
func (w *Watcher) pingLoop(ctx context.Context) {
	t := time.NewTicker(30 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			err := w.dc.Ping(ctx)
			w.mu.Lock()
			wasHealthy := w.healthy
			w.healthy = err == nil
			w.mu.Unlock()
			if err != nil && wasHealthy {
				w.nt.Emit("docker_disconnected", map[string]string{
					"result":  "容器管理与事件监控暂时中断",
					"message": "Docker 守护进程当前不可达，恢复后会自动再通知一次",
				})
				w.st.AddRunLog("system", "docker", "failed", "无法连接 Docker 守护进程："+err.Error(), "")
				w.bus.Publish("system", "docker", "failed", map[string]any{"error": err.Error()})
			} else if err == nil && !wasHealthy {
				w.nt.Emit("docker_reconnected", map[string]string{
					"result":  "容器管理与事件监控已恢复",
					"message": "Docker 守护进程已恢复可用",
				})
				w.st.AddRunLog("system", "docker", "success", "Docker 守护进程已恢复", "")
				w.bus.Publish("system", "docker", "success", map[string]any{})
			}
		}
	}
}

// DetectSelf 推断 Dockhelm 自己的容器名。
//
// 容器里 hostname 默认就是容器短 ID，用它去比对即可；也支持环境变量显式指定。
func DetectSelf(ctx context.Context, dc *dockerx.Client, explicit string) (name, id string) {
	if explicit != "" {
		if insp, err := dc.Inspect(ctx, explicit); err == nil {
			return strings.TrimPrefix(str(insp["Name"]), "/"), str(insp["Id"])
		}
		return explicit, ""
	}
	host, err := os.Hostname()
	if err != nil || host == "" {
		return "", ""
	}
	list, err := dc.ListContainers(ctx)
	if err != nil {
		return "", ""
	}
	for _, c := range list {
		if strings.HasPrefix(c.ID, host) || c.Name() == host {
			return c.Name(), c.ID
		}
	}
	return "", ""
}

func str(v any) string {
	s, _ := v.(string)
	return s
}
