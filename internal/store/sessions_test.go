package store

import (
	"path/filepath"
	"testing"
	"time"
)

// TestDeleteOtherSessions 钉住「强制登出其他会话」的核心不变式：
// **当前会话必须原样保留**，只作废其余的。这条反了的话，用户点一下按钮
// 会把自己也踢下线 —— 那是改密码才该有的效果。
func TestDeleteOtherSessions(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "dockhelm.db"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	mk := func(token string) {
		if err := st.CreateSession(token, time.Hour, "192.168.3.64", "ua"); err != nil {
			t.Fatalf("CreateSession(%s): %v", token, err)
		}
	}
	mk("current")
	mk("stale-1")
	mk("stale-2")

	removed, err := st.DeleteOtherSessions("current")
	if err != nil {
		t.Fatalf("DeleteOtherSessions: %v", err)
	}
	if removed != 2 {
		t.Fatalf("removed = %d，期望 2", removed)
	}
	if st.SessionCount() != 1 {
		t.Fatalf("SessionCount = %d，期望 1", st.SessionCount())
	}
	if st.GetSession("current") == nil {
		t.Fatal("当前会话被误删 —— 这是绝不允许的")
	}
	if st.GetSession("stale-1") != nil || st.GetSession("stale-2") != nil {
		t.Fatal("残留会话没有被作废")
	}

	// 没有其他会话时再点一次：无事发生，0 个
	removed, err = st.DeleteOtherSessions("current")
	if err != nil || removed != 0 {
		t.Fatalf("重复调用 = (%d, %v)，期望 (0, nil)", removed, err)
	}
}
