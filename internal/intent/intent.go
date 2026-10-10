// Package intent 记录「这次容器退出是 dockhelm 自己发起的」。
//
// 背景：watch 模块订阅 Docker 事件流，把所有非 0 退出的 die 事件记成
// 「容器意外退出」并告警。但更新 / 还原 / 计划任务 / 手动停止也会停容器 ——
// stop 宽限期一到 Docker 发 SIGKILL，退出码就是 137，于是一次完全成功的
// 更新旁边总躺着一条「容器意外退出，退出码 137 · 失败」，纯属虚惊。
//
// 解法是发起方在动手停之前 Mark 一下，watch 收到 die 时 Consume 豁免：
// 两边各自独立（updater 不 import watch，watch 也不 import updater），
// 都只依赖这个无依赖的小包。容器外的 `docker stop` / OOM 等真·意外退出
// 没有登记，照常告警。
//
// 两条生命周期规则：
//
//   - die 事件用「计数核销」：Mark 一次豁免一次。事件流是异步的，
//     多步流程（更新里的停旧容器 + 失败清理新容器）各 Mark 各的，
//     谁的 die 先到核销谁的；事件不来（容器本来就是停的）由 TTL 过期兜底。
//   - 健康事件用「时间窗」：更新/重启后的头几分钟新容器短暂 unhealthy
//     是常态，等它自己变回来，不算失败。
package intent

import (
	"sync"
	"time"
)

const (
	// busyTTL 兜底窗口：Mark 之后 die 迟迟不来（容器本来就是停的）或
	// 流程忘了 Release，豁免最多保留这么久，免得之后的真崩溃被永久吞掉。
	busyTTL = 15 * time.Minute
)

var (
	mu        sync.Mutex
	pending   = map[string][]time.Time{} // 待核销的 die 豁免（每次 Mark 记一个时间戳）
	busyUntil = map[string]time.Time{}   // 健康事件豁免窗口的截止时刻
)

// Mark 声明「接下来对 name 的停止/重启是 dockhelm 发起的，其退出不算意外」。
// 若随后 ContainerAction 返回错误（没停成、不会有 die 事件），应立即 Release。
func Mark(name string) {
	now := time.Now()
	mu.Lock()
	defer mu.Unlock()
	pending[name] = append(pending[name], now)
	busyUntil[name] = now.Add(busyTTL)
}

// Release 提前结束 name 的健康事件豁免窗口（多步流程结束时调用）。
// die 豁免的计数**故意不清**：事件流异步，可能比流程本身晚到，
// 留给 Consume 核销或 TTL 过期。
func Release(name string) {
	mu.Lock()
	defer mu.Unlock()
	delete(busyUntil, name)
}

// Consume 报告并核销一次 dockhelm 发起的退出。watch 的 die 事件用它：
// 返回 true 表示这次退出是预期的，不该告警、不该写失败日志。
func Consume(name string) bool {
	now := time.Now()
	mu.Lock()
	defer mu.Unlock()
	list := pending[name]
	kept := list[:0]
	consumed := false
	for _, t := range list {
		if !consumed && now.Sub(t) < busyTTL {
			consumed = true // 核销最早的一次登记
			continue
		}
		if now.Sub(t) < busyTTL {
			kept = append(kept, t)
		}
	}
	if len(kept) == 0 {
		delete(pending, name)
	} else {
		pending[name] = kept
	}
	return consumed
}

// Active 报告 name 是否还处在「正在被 dockhelm 操作」的窗口内（非破坏性查询）。
// watch 的健康事件用它，避免更新/重启期间的健康抖动被记成失败。
func Active(name string) bool {
	mu.Lock()
	defer mu.Unlock()
	t, ok := busyUntil[name]
	return ok && time.Now().Before(t)
}
