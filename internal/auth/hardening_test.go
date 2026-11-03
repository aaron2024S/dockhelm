package auth

import (
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aaron2024s/dockhelm/internal/store"
)

// TestCheckPasswordRules 密码强度规则。
//
// 这个面板等价于宿主机 root，6 位密码在离线爆破下只是几分钟的事；
// 纯数字 8 位同理（数字恰好是用户最爱用的那种）。
func TestCheckPasswordRules(t *testing.T) {
	cases := []struct {
		pw   string
		want error
	}{
		{"short7c", ErrWeakPassword},    // 7 位
		{"12345678", ErrDigitsOnly},     // 纯数字
		{"000000000000", ErrDigitsOnly}, // 更长的纯数字
		{"1234567a", nil},               // 8 位、含字母
		{"correct horse", nil},          // 长口令
		{strings.Repeat("x", 72), nil},  // bcrypt 上限内的长密码
	}
	for _, c := range cases {
		err := checkPassword(c.pw)
		if !errors.Is(err, c.want) && !(err == nil && c.want == nil) {
			t.Errorf("checkPassword(%q) = %v，期望 %v", c.pw, err, c.want)
		}
	}
}

// TestWeakPasswordRejectedOnEveryEntryPoint 三个入口（首次设置 / 改密 / 环境变量）
// 必须用同一套规则，否则「某个入口能设弱密码」就是绕过的口子。
func TestWeakPasswordRejectedOnEveryEntryPoint(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(filepath.Join(dir, "dockhelm.db"))
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	m, err := New(filepath.Join(dir, "auth.json"), st)
	if err != nil {
		t.Fatalf("auth.New: %v", err)
	}

	if err := m.SetupPassword("123456"); !errors.Is(err, ErrWeakPassword) {
		t.Fatalf("首次设置 6 位密码应被拒（ErrWeakPassword），实际 %v", err)
	}
	if err := m.ForceSetPassword("12345678"); !errors.Is(err, ErrDigitsOnly) {
		t.Fatalf("环境变量强制重置为纯数字应被拒（ErrDigitsOnly），实际 %v", err)
	}
	if err := m.SetPassword("", "12345678"); !errors.Is(err, ErrDigitsOnly) {
		t.Fatalf("改密入口的纯数字密码应被拒（ErrDigitsOnly），实际 %v", err)
	}
	if err := m.SetupPassword("correct-horse"); err != nil {
		t.Fatalf("强密码应被接受，实际 %v", err)
	}
}
