package intent

import (
	"testing"
	"time"
)

func TestConsumeWithoutMark(t *testing.T) {
	if Consume("c-nobody") {
		t.Fatal("没有登记就 Consume 必须返回 false，否则任何外部崩溃都会被吞掉")
	}
}

func TestMarkConsumeOnce(t *testing.T) {
	Mark("c-once")
	if !Consume("c-once") {
		t.Fatal("Mark 后第一次 Consume 应豁免")
	}
	if Consume("c-once") {
		t.Fatal("豁免是一次性的：第二次 Consume 必须返回 false，之后的真崩溃要照常告警")
	}
}

func TestMarkCount(t *testing.T) {
	// 更新流程会在同一容器上 Mark 两次（停旧容器 + 失败时清理新容器）
	Mark("c-count")
	Mark("c-count")
	if !Consume("c-count") || !Consume("c-count") {
		t.Fatal("Mark 两次应核销两次")
	}
	if Consume("c-count") {
		t.Fatal("核销完不应再有豁免")
	}
}

func TestConsumeIsPerContainer(t *testing.T) {
	Mark("c-a")
	if Consume("c-b") {
		t.Fatal("豁免必须按容器名隔离")
	}
	Consume("c-a")
}

func TestActiveWindow(t *testing.T) {
	if Active("c-active") {
		t.Fatal("未 Mark 时不应处于活跃窗口")
	}
	Mark("c-active")
	if !Active("c-active") {
		t.Fatal("Mark 后应处于活跃窗口（更新期间健康抖动不算失败）")
	}
	Release("c-active")
	if Active("c-active") {
		t.Fatal("Release 应提前结束健康事件豁免窗口")
	}
	// Release 只清健康窗口，die 豁免的计数保留给事件流异步核销
	if !Consume("c-active") {
		t.Fatal("Release 不应清掉还没核销的 die 豁免")
	}
}

func TestTTLExpiry(t *testing.T) {
	Mark("c-ttl")
	// 把登记时间拨回 TTL 之前，模拟「die 事件一直没来」的兜底过期
	mu.Lock()
	pending["c-ttl"] = []time.Time{time.Now().Add(-2 * busyTTL)}
	busyUntil["c-ttl"] = time.Now().Add(-time.Minute)
	mu.Unlock()
	if Consume("c-ttl") {
		t.Fatal("过期的登记必须作废，否则真崩溃会被陈年豁免吞掉")
	}
	if Active("c-ttl") {
		t.Fatal("过期的健康窗口必须失效")
	}
}
