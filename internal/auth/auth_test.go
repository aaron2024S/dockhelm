package auth

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/aaron2024s/dockhelm/internal/store"
)

func newTestManager(t *testing.T) *Manager {
	t.Helper()
	dir := t.TempDir()
	st, err := store.Open(filepath.Join(dir, "dockhelm.db"))
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	m, err := New(filepath.Join(dir, "auth.json"), st)
	if err != nil {
		t.Fatalf("auth.New: %v", err)
	}
	if err := m.SetupPassword("correct-horse"); err != nil {
		t.Fatalf("SetupPassword: %v", err)
	}
	return m
}

// newManagerWithAttempts 造一个「登录记录里已经躺着这些尝试」的管理器。
//
// RecordLogin 只能记「现在」，测不了窗口滑动 —— 直接把受控时刻写进 db 文件再打开，
// 就等于让存储层带着一段历史启动（这正是重启面板后的真实形态）。
func newManagerWithAttempts(t *testing.T, attempts []store.LoginAttempt) *Manager {
	t.Helper()
	dir := t.TempDir()
	doc := map[string]any{"version": 1, "loginAttempts": attempts}
	b, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	dbPath := filepath.Join(dir, "dockhelm.db")
	if err := os.WriteFile(dbPath, b, 0o600); err != nil {
		t.Fatalf("写 db: %v", err)
	}
	st, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	m, err := New(filepath.Join(dir, "auth.json"), st)
	if err != nil {
		t.Fatalf("auth.New: %v", err)
	}
	return m
}

func failAt(d time.Duration, ip string) store.LoginAttempt {
	return store.LoginAttempt{TS: time.Now().Add(-d).UTC().Format(time.RFC3339), IP: ip, OK: false}
}

// TestVerifyHasNoFixedDelay 钉住「失败不再罚站 1 秒」。
//
// 旧实现每次密码错误都 time.Sleep(time.Second)：既同步占着 HTTP handler 的 goroutine，
// 又让手滑打错一个字的本人白等 1 秒。真正的防线是滑动窗口锁定（见下一条）。
// 3 次失败若还慢于 1.5 秒，说明有人把 Sleep 加回来了。
func TestVerifyHasNoFixedDelay(t *testing.T) {
	m := newTestManager(t)
	start := time.Now()
	for i := 0; i < 3; i++ {
		if _, err := m.Verify("wrong", "10.0.0.9", "ua", false); !errors.Is(err, ErrBadPassword) {
			t.Fatalf("第 %d 次错误密码：err = %v，期望 ErrBadPassword", i+1, err)
		}
	}
	if d := time.Since(start); d > 1500*time.Millisecond {
		t.Fatalf("3 次失败共耗时 %v，像是又加了固定延时（bcrypt 单纯比较约 0.1 秒）", d)
	}
}

// TestVerifyLocksOnFifthFailure 第 5 次失败必须**直接**按「已锁定」返回。
//
// 否则前端只会拿到 bad_password + remaining=0，界面上写着「密码错误，还可尝试 0 次」，
// 既不说锁了多久，也不说什么时候能重试 —— 用户只能盲目重复点。
func TestVerifyLocksOnFifthFailure(t *testing.T) {
	m := newTestManager(t)
	const ip = "10.0.0.7"

	for i := 1; i < MaxFailures; i++ {
		if _, err := m.Verify("wrong", ip, "ua", false); !errors.Is(err, ErrBadPassword) {
			t.Fatalf("第 %d 次失败：err = %v，期望 ErrBadPassword", i, err)
		}
		if left := m.LockedFor(ip); left != 0 {
			t.Fatalf("第 %d 次失败就锁了（还剩 %v），门限是 %d 次", i, left, MaxFailures)
		}
	}

	if _, err := m.Verify("wrong", ip, "ua", false); !errors.Is(err, ErrLocked) {
		t.Fatalf("第 %d 次失败：err = %v，期望 ErrLocked", MaxFailures, err)
	}

	// 锁定期内即便密码正确也进不去，且返回的仍是 ErrLocked（前端据此起倒计时）
	if _, err := m.Verify("correct-horse", ip, "ua", false); !errors.Is(err, ErrLocked) {
		t.Fatalf("锁定期内正确密码：err = %v，期望 ErrLocked", err)
	}

	// 换一个来源 IP 不受影响（锁定是按来源 IP 的滑动窗口）
	if _, err := m.Verify("wrong", "10.0.0.8", "ua", false); !errors.Is(err, ErrBadPassword) {
		t.Fatalf("另一个 IP 被连坐锁定了：%v", err)
	}
}

// TestLockedForCountsFromOldestFailure 剩余时长按**第 (n-MaxFailures+1) 早的那次失败**
// 滑出窗口来算，不是「最后一次失败 + 5 分钟」。
//
// 5 次失败里最早那次已过去 4 分钟 ⇒ 还剩约 1 分钟就解锁。若按「最后一次 + 5 分钟」算，
// 会报成近 5 分钟，用户白白多等 4 分钟 —— 这条测试就是钉死这个差异。
func TestLockedForCountsFromOldestFailure(t *testing.T) {
	const ip = "10.0.0.6"
	m := newManagerWithAttempts(t, []store.LoginAttempt{
		failAt(4*time.Minute, ip),
		failAt(3*time.Minute, ip),
		failAt(2*time.Minute, ip),
		failAt(90*time.Second, ip),
		failAt(10*time.Second, ip),
	})
	left := m.LockedFor(ip)
	if left <= 0 {
		t.Fatal("5 次窗口内失败应当处于锁定态")
	}
	if left > 2*time.Minute {
		t.Fatalf("LockedFor = %v，明显按「最后一次失败 + 窗口」算了；期望约 1 分钟", left)
	}
	if left < 30*time.Second {
		t.Fatalf("LockedFor = %v，期望约 1 分钟（最早那次 4 分钟前失败）", left)
	}
}

// TestLockedForIgnoresOutOfWindowFailures 窗口外的失败与成功记录都不算数。
func TestLockedForIgnoresOutOfWindowFailures(t *testing.T) {
	const ip = "10.0.0.5"
	// 5 条里 2 条已滑出窗口 ⇒ 窗口内只有 3 条，不该锁
	m := newManagerWithAttempts(t, []store.LoginAttempt{
		failAt(6*time.Minute, ip),
		failAt(5*time.Minute+20*time.Second, ip),
		failAt(3*time.Minute, ip),
		failAt(2*time.Minute, ip),
		failAt(time.Minute, ip),
	})
	if left := m.LockedFor(ip); left != 0 {
		t.Fatalf("窗口外失败被算进门限了：LockedFor = %v", left)
	}
	if n := m.Failures(ip); n != 3 {
		t.Fatalf("Failures = %d，期望 3", n)
	}

	// 窗口内 5 条、其中 2 条是成功 ⇒ 只有 3 条失败，不该锁
	m2 := newManagerWithAttempts(t, []store.LoginAttempt{
		failAt(3*time.Minute, ip),
		failAt(2*time.Minute, ip),
		failAt(time.Minute, ip),
		{TS: time.Now().Add(-4 * time.Minute).UTC().Format(time.RFC3339), IP: ip, OK: true},
		{TS: time.Now().Add(-30 * time.Second).UTC().Format(time.RFC3339), IP: ip, OK: true},
	})
	if left := m2.LockedFor(ip); left != 0 {
		t.Fatalf("成功记录被当成失败算进门限了：LockedFor = %v", left)
	}
}

// TestLockedHint 锁定期文案必须带上「还要等多久」。
func TestLockedHint(t *testing.T) {
	cases := []struct {
		left time.Duration
		want string
	}{
		{0, "尝试次数过多，请稍后再试"},
		{-time.Second, "尝试次数过多，请稍后再试"},
		{3200 * time.Millisecond, "尝试次数过多，请约 4 秒后再试"},
		{59 * time.Second, "尝试次数过多，请约 59 秒后再试"},
		{60 * time.Second, "尝试次数过多，请约 1 分钟后再试"},
		{90 * time.Second, "尝试次数过多，请约 2 分钟后再试"},
		{5 * time.Minute, "尝试次数过多，请约 5 分钟后再试"},
	}
	for _, c := range cases {
		if got := LockedHint(c.left); got != c.want {
			t.Errorf("LockedHint(%v) = %q，期望 %q", c.left, got, c.want)
		}
	}
}
