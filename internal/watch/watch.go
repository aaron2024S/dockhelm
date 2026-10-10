// Package watch 订阅 Docker 事件流，把「容器意外退出 / 崩溃循环 / 健康检查失败」
// 以及「守护进程断开又恢复」翻译成通知事件。
package watch

import (
	"context"
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
	healthy  bool
}

// New 创建观察者。
func New(dc *dockerx.Client, nt *notify.Manager, st *store.Store, b *bus.Bus, selfFn func() string) *Watcher {
	return &Watcher{
		dc: dc, nt: nt, st: st, bus: b, self: selfFn,
		dieTimes: map[string][]time.Time{},
		healthy:  true,
	}
}

// 崩溃循环判定：10 分钟内退出 3 次以上。
const (
	crashWindow = 10 * time.Minute
	crashCount  = 3
)

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

	switch {
	case ev.Action == "die":
		code := ev.Actor.Attributes["exitCode"]
		n, _ := strconv.Atoi(code)
		// 0 = 正常结束，不告警
		if n == 0 {
			return
		}
		// dockhelm 自己发起的停止（更新 / 还原 / 计划任务 / 手动停）不是「意外退出」
		// —— stop 宽限期一到 Docker 发 SIGKILL，退出码就是 137。没登记的 die
		// （docker CLI、OOM、程序崩溃）照常告警。
		if intent.Consume(name) {
			return
		}
		w.recordDie(name)
		image := ev.Actor.Attributes["image"]
		w.nt.Emit("container_died", map[string]string{
			"container": name, "image": image,
			"result":  "容器意外退出",
			"message": "退出码 " + code,
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

	case ev.Action == "start", ev.Action == "restart":
		w.bus.Publish("container", "event", "info", map[string]any{
			"name": name, "action": ev.Action,
		})
	}
}

// recordDie 记录退出时间，并在窗口内次数超阈值时发一次崩溃循环告警。
func (w *Watcher) recordDie(name string) {
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
			"container": name,
			"result":    "容器崩溃循环",
			"message":   "最近 " + w.humanWindow() + " 内已退出 " + strconv.Itoa(count) + " 次",
		})
	}
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
					"result": "Docker 连接断开", "message": err.Error(),
				})
				w.st.AddRunLog("system", "docker", "failed", "无法连接 Docker 守护进程："+err.Error(), "")
				w.bus.Publish("system", "docker", "failed", map[string]any{"error": err.Error()})
			} else if err == nil && !wasHealthy {
				w.nt.Emit("docker_reconnected", map[string]string{
					"result": "Docker 连接恢复", "message": "守护进程已恢复可用",
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
