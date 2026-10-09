package bus

import (
	"sync"
	"testing"
	"time"
)

// 订阅者掉线与事件发布并发时必须安全。
//
// 回归点：cancel 早先会 `close(ch)`，而 Publish 是「锁内拷贝订阅者、锁外逐个发送」，
// 两者的临界区不重叠 —— 于是存在「Publish 拷到 ch → cancel 关掉 ch → Publish 发送」
// 的时序。向已关闭的 channel 发送**即使在 select+default 下也会 panic**，
// 而这个 panic 发生在 updater / watch / scheduler 这些没有 recover 的 goroutine 里，
// 会把整个进程带走（不是单个请求 500）。掉线重连的 SSE 客户端一多就必然撞上。
//
// 配合 `go test -race` 跑，同时能抓 data race。
func TestPublishCancelRace(t *testing.T) {
	b := New()
	stop := make(chan struct{})
	var subWG, pubWG sync.WaitGroup

	// 订阅者反复上线 / 掉线
	for i := 0; i < 8; i++ {
		subWG.Add(1)
		go func() {
			defer subWG.Done()
			for {
				select {
				case <-stop:
					return
				default:
				}
				_, cancel := b.Subscribe()
				cancel()
			}
		}()
	}

	// 发布方持续发事件
	for i := 0; i < 4; i++ {
		pubWG.Add(1)
		go func() {
			defer pubWG.Done()
			for j := 0; j < 3000; j++ {
				b.Publish("update", "step", "running", map[string]any{"i": j})
			}
		}()
	}

	pubWG.Wait()
	close(stop)
	subWG.Wait()
}

// 正常路径：订阅后能收到事件；取消订阅后不再收到。
func TestSubscribeThenCancel(t *testing.T) {
	b := New()
	ch, cancel := b.Subscribe()

	b.Publish("update", "step", "running", nil)
	select {
	case ev := <-ch:
		if ev.Topic != "update" || ev.Kind != "step" {
			t.Fatalf("收到的事件不对：%+v", ev)
		}
	case <-time.After(time.Second):
		t.Fatal("订阅后没收到事件")
	}

	cancel()
	cancel() // 重复取消不能崩（幂等）

	b.Publish("update", "step", "success", nil)
	select {
	case ev := <-ch:
		t.Fatalf("取消订阅后不该再收到事件，却收到 %+v", ev)
	case <-time.After(50 * time.Millisecond):
	}

	// 取消后没有任何订阅者，Publish 仍应正常完成
	if ev := b.Publish("backup", "done", "success", nil); ev.ID == 0 {
		t.Fatal("取消全部订阅后 Publish 没有正常返回")
	}
}
