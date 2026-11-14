package api

import (
	"testing"

	"github.com/aaron2024s/dockhelm/internal/auth"
)

// TestSessionCookie 钉住「保持登录」的语义 —— 这是登录页那句文案背后的真实行为。
//
// 勾选 ⇒ 持久化 Cookie（MaxAge = 7 天，浏览器关掉也留着）；
// 不勾 ⇒ **绝不下发 MaxAge**，浏览器按会话 Cookie 处理，关掉即失效。
// 这条反了的话，「一次性登录」会悄悄变成长期登录 —— 属于安全语义，不能漂。
func TestSessionCookie(t *testing.T) {
	keep := sessionCookie("tok", true, false)
	if keep.MaxAge != 7*24*3600 {
		t.Fatalf("勾选保持登录：MaxAge = %d，期望 7 天（604800）", keep.MaxAge)
	}
	if !keep.HttpOnly {
		t.Fatal("会话 Cookie 必须 HttpOnly")
	}
	if keep.Path != "/" {
		t.Fatalf("Path = %q，期望 /", keep.Path)
	}

	oneShot := sessionCookie("tok", false, false)
	if oneShot.MaxAge != 0 {
		t.Fatalf("不勾保持登录：MaxAge = %d，期望 0（会话 Cookie，关浏览器即失效）", oneShot.MaxAge)
	}
	if oneShot.Expires.Unix() != -62135596800 && !oneShot.Expires.IsZero() {
		t.Fatalf("不勾时不该设 Expires，实际 %v", oneShot.Expires)
	}

	// Secure 只在反代声明 https 时才置位（局域网 HTTP 下置位会导致登录不了）
	if !sessionCookie("tok", true, true).Secure {
		t.Fatal("X-Forwarded-Proto=https 时应置 Secure")
	}
	if sessionCookie("tok", true, false).Secure {
		t.Fatal("HTTP 部署下不应置 Secure")
	}

	// 服务端 TTL 与 Cookie 语义配套
	if sessionTTL(true).Hours() != 7*24 {
		t.Fatalf("勾选：服务端 TTL = %v，期望 168h", sessionTTL(true))
	}
	if sessionTTL(false).Hours() != 24 {
		t.Fatalf("不勾：服务端 TTL = %v，期望 24h 兜底", sessionTTL(false))
	}
}

// TestNormalizeLockPolicy 登录保护这两项读写都要收进合法区间。
//
// 直接改 JSON 把值写坏（0 / 负数 / 超大）不该让面板变成「永不锁定」或
// 「一次就锁死」—— 读的时候过一遍比只在保存时校验更可靠。
func TestNormalizeLockPolicy(t *testing.T) {
	def := auth.MaxFailures
	defMin := int(auth.LockWindow.Minutes())
	cases := []struct {
		in      int
		wantMax int
		wantMin int
	}{
		{5, 5, 5}, // 正常值原样保留
		{0, def, defMin},   // 没配过 ⇒ 默认
		{-3, def, defMin},  // 负值 ⇒ 默认
		{1, auth.MinFailureLimit, auth.MinLockMinutes}, // 太小 ⇒ 下限
		{999, auth.MaxFailureLimit, auth.MaxLockMinutes}, // 太大 ⇒ 上限
	}
	for _, c := range cases {
		gotMax, gotMin := normalizeLockPolicy(c.in, c.in)
		if gotMax != c.wantMax || gotMin != c.wantMin {
			t.Fatalf("normalizeLockPolicy(%d, %d) = (%d, %d)，期望 (%d, %d)",
				c.in, c.in, gotMax, gotMin, c.wantMax, c.wantMin)
		}
	}
}
