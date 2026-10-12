package watch

import (
	"testing"
	"time"
)

// 「外部主动停止」判定：docker stop 会先发 kill 事件（SIGTERM）再发 die，
// 这类非零退出不该报「容器意外退出」；真崩溃/OOM 没有 kill 前置，必须照常报。
func TestExitedByExternalStop(t *testing.T) {
	now := time.Now()
	cases := []struct {
		name  string
		kill  time.Duration // <0 表示没记录过
		oom   time.Duration
		want  bool
		after string
	}{
		{name: "无 kill 记录（程序崩溃）", kill: -1, oom: -1, want: false},
		{name: "kill 刚发生（docker stop 强杀）", kill: 0, oom: -1, want: true},
		{name: "kill 在宽限期内（慢退出被强杀）", kill: 50 * time.Second, oom: -1, want: true},
		{name: "kill 刚好在窗口边界内", kill: stopWindow - time.Second, oom: -1, want: true},
		{name: "kill 已过期（随后的一次真崩溃）", kill: stopWindow + time.Second, oom: -1, want: false},
		{name: "kill + 同期 oom（OOM 不能被掩护）", kill: 0, oom: 0, want: false},
		{name: "kill 新、oom 已过期", kill: 0, oom: stopWindow + time.Second, want: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			w := New(nil, nil, nil, nil, nil)
			if c.kill >= 0 {
				w.killAt["app"] = now.Add(-c.kill)
			}
			if c.oom >= 0 {
				w.oomAt["app"] = now.Add(-c.oom)
			}
			if got := w.exitedByExternalStop("app", now); got != c.want {
				t.Fatalf("exitedByExternalStop = %v，期望 %v", got, c.want)
			}
		})
	}
}

// 命中豁免后必须清掉 kill 记录：否则同一次 stop 的标记会掩护掉窗口内
// 紧随其后的一次真崩溃。
func TestExternalStopConsumesKillMark(t *testing.T) {
	now := time.Now()
	w := New(nil, nil, nil, nil, nil)
	w.killAt["app"] = now
	if !w.exitedByExternalStop("app", now) {
		t.Fatal("第一次应判定为外部停止")
	}
	if w.exitedByExternalStop("app", now) {
		t.Fatal("kill 记录应已被消费，第二次不该再豁免")
	}
}

// 不同容器互不干扰：A 的 stop 不能掩护 B 的崩溃。
func TestExternalStopIsPerContainer(t *testing.T) {
	now := time.Now()
	w := New(nil, nil, nil, nil, nil)
	w.killAt["a"] = now
	if !w.exitedByExternalStop("a", now) {
		t.Fatal("a 应豁免")
	}
	if w.exitedByExternalStop("b", now) {
		t.Fatal("b 没有 kill 记录，不该豁免")
	}
}

// markAt 记录时间并清理过期项（容器多了 map 不能无限增长）。
func TestMarkAtPrunesStale(t *testing.T) {
	now := time.Now()
	w := New(nil, nil, nil, nil, nil)
	w.markAt(w.killAt, "old", now.Add(-stopWindow-time.Second))
	w.markAt(w.killAt, "new", now)
	if _, ok := w.killAt["old"]; ok {
		t.Fatal("过期项应被清掉")
	}
	if _, ok := w.killAt["new"]; !ok {
		t.Fatal("未过期项应保留")
	}
}
