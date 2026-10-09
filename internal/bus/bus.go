// Package bus 提供一个极简的进程内发布订阅总线，用于把长任务（更新、拉取、备份）
// 的进度实时推给前端的 SSE 连接。
//
// 设计取舍：只保留「最近 N 条」历史，客户端带 Last-Event-ID 重连时能补发，
// 超出就把最后一条状态事件重放给新订阅者，避免刷新页面后进度条永远空着。
package bus

import (
	"sync"
	"time"
)

// Event 一条要推给前端的事件。
type Event struct {
	ID     int64          `json:"id"`
	Topic  string         `json:"topic"`  // update / schedule / backup / notify / log
	Kind   string         `json:"kind"`   // 更细的动作，如 pull_progress / container_done
	Time   string         `json:"time"`
	Data   map[string]any `json:"data,omitempty"`
	Status string         `json:"status,omitempty"` // running | success | failed | info
}

// Bus 进程内广播。
type Bus struct {
	mu       sync.RWMutex
	nextID   int64
	subs     map[int64]chan Event
	lastSeen map[string]Event // topic -> 最后一条（用于新订阅者回放）
	histSize int
	hist     []Event
}

// New 创建总线。
func New() *Bus {
	return &Bus{
		subs:     map[int64]chan Event{},
		lastSeen: map[string]Event{},
		histSize: 200,
	}
}

// Publish 发布一条事件。
func (b *Bus) Publish(topic, kind, status string, data map[string]any) Event {
	b.mu.Lock()
	b.nextID++
	ev := Event{
		ID:     b.nextID,
		Topic:  topic,
		Kind:   kind,
		Status: status,
		Time:   time.Now().UTC().Format(time.RFC3339),
		Data:   data,
	}
	b.hist = append(b.hist, ev)
	if len(b.hist) > b.histSize {
		b.hist = b.hist[len(b.hist)-b.histSize:]
	}
	b.lastSeen[topic] = ev
	subs := make([]chan Event, 0, len(b.subs))
	for _, ch := range b.subs {
		subs = append(subs, ch)
	}
	b.mu.Unlock()

	for _, ch := range subs {
		select {
		case ch <- ev:
		default:
			// 订阅者消费不过来就丢这一条，绝不阻塞发布方
		}
	}
	return ev
}

// Subscribe 订阅，返回事件通道与取消函数。缓冲区 64 条。
func (b *Bus) Subscribe() (<-chan Event, func()) {
	ch := make(chan Event, 64)
	b.mu.Lock()
	id := b.nextID + 1
	for {
		if _, exists := b.subs[id]; !exists {
			break
		}
		id++
	}
	b.subs[id] = ch
	b.mu.Unlock()

	cancel := func() {
		b.mu.Lock()
		delete(b.subs, id)
		b.mu.Unlock()
		// 这里**故意不 close(ch)**。
		//
		// Publish 是「锁内拷贝订阅者、锁外逐个发送」，而 cancel 在锁内把它摘掉。
		// 两者的临界区不重叠：一旦在 cancel 里 close，就可能出现
		// 「Publish 已拷到 ch → cancel 关掉 ch → Publish 发送」的时序，
		// 而向已关闭的 channel 发送即使在 select+default 下也会 **panic**。
		// 这个 panic 发生在发布方 goroutine（updater / watch / scheduler），
		// 没有 recover，会把整个进程带走 —— 掉线重连的 SSE 客户端多的时候很容易撞上。
		//
		// 不 close 也没有副作用：两个消费方（server.go / updates.go）都已经盯
		// r.Context().Done() 退出，`!ok` 分支只是多一层保险；摘掉订阅后没人再持有
		// 这个 channel，交给 GC 即可。
	}
	return ch, cancel
}

// History 返回最近的事件（可按 topic 过滤，topic 为空表示不过滤）。
func (b *Bus) History(topic string, limit int) []Event {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := []Event{}
	for _, ev := range b.hist {
		if topic == "" || ev.Topic == topic {
			out = append(out, ev)
		}
	}
	if limit > 0 && len(out) > limit {
		out = out[len(out)-limit:]
	}
	return out
}

// Snapshot 返回每个 topic 的最后一条事件，用于新连接的首屏渲染。
func (b *Bus) Snapshot() map[string]Event {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := map[string]Event{}
	for k, v := range b.lastSeen {
		out[k] = v
	}
	return out
}
