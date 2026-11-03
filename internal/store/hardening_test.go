package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestCorruptDBIsQuarantined 钉住「数据文件坏了也要能开机」。
//
// 反例（改之前）：json.Unmarshal 失败直接 return err ⇒ main 里 log.Fatalf ⇒
// 容器进入重启循环，用户只能在 NAS 上 SSH 手删文件。而 auth.json 一直有这层兜底。
func TestCorruptDBIsQuarantined(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "dockhelm.db")
	if err := os.WriteFile(path, []byte(`{"version":1,"settings":{`), 0o600); err != nil {
		t.Fatal(err)
	}

	st, err := Open(path)
	if err != nil {
		t.Fatalf("数据文件损坏时不该返回错误（会导致进程起不来），实际：%v", err)
	}
	if st.LoadWarning() == "" {
		t.Fatal("损坏时必须有可读告警，否则用户只会觉得「设置全没了」")
	}
	if v := st.GetSetting("update.concurrency", "2"); v != "2" {
		t.Fatalf("应按空数据启动，实际读到 %q", v)
	}

	// 坏文件必须留档，不能直接丢掉 —— 里面可能还有手工恢复的余地
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, e := range entries {
		if strings.Contains(e.Name(), ".corrupt-") {
			found = true
		}
	}
	if !found {
		t.Fatalf("坏文件没有被改名留档，目录内容：%v", entries)
	}
	// 新文件必须是可解析的
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var probe map[string]any
	if err := json.Unmarshal(b, &probe); err != nil {
		t.Fatalf("启动后应立刻写出一份可解析的新文件，实际：%v", err)
	}
}

// TestSessionTokenIsHashedAtRest 钉住「数据文件里没有可直接使用的 token」。
//
// 备份数据目录、把整个目录传网盘，都不该等于把登录态一起送出去。
func TestSessionTokenIsHashedAtRest(t *testing.T) {
	path := filepath.Join(t.TempDir(), "dockhelm.db")
	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	const token = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
	if err := st.CreateSession(token, time.Hour, "192.168.3.64", "ua"); err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(raw), token) {
		t.Fatal("数据文件里出现了明文会话 token")
	}
	if !strings.Contains(string(raw), hashToken(token)) {
		t.Fatal("数据文件里应有 token 的哈希")
	}
	// 但带原文回来必须还能认出来
	if st.GetSession(token) == nil {
		t.Fatal("哈希之后，带原 token 的请求必须仍然有效")
	}
}

// TestSessionIdleExpiry 连续 48 小时没访问 ⇒ 失效（Cookie 还在也不认）。
func TestSessionIdleExpiry(t *testing.T) {
	st, err := Open(filepath.Join(t.TempDir(), "dockhelm.db"))
	if err != nil {
		t.Fatal(err)
	}
	if err := st.CreateSession("tok", 7*24*time.Hour, "ip", "ua"); err != nil {
		t.Fatal(err)
	}
	if st.GetSession("tok") == nil {
		t.Fatal("刚建的会话应有效")
	}

	// 把最后访问时间推到 49 小时前
	st.mu.Lock()
	st.db.Sessions[0].LastSeen = time.Now().Add(-49 * time.Hour).UTC().Format(time.RFC3339)
	st.mu.Unlock()

	if st.GetSession("tok") != nil {
		t.Fatal("空闲超过 48 小时的会话必须失效")
	}
	if st.SessionCount() != 0 {
		t.Fatalf("空闲会话不该计入会话数，实际 %d", st.SessionCount())
	}
}

// TestLegacyPlainTokenStillWorks 老版本落盘的是明文 token，第一次带旧 Cookie 进来
// 要能原地升级成哈希 —— 升级后把所有用户踢下线是不能接受的。
func TestLegacyPlainTokenStillWorks(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "dockhelm.db")
	legacy := map[string]any{
		"version": 1,
		"sessions": []map[string]any{{
			"token":     "legacy-plain-token",
			"createdAt": time.Now().Add(-time.Hour).UTC().Format(time.RFC3339),
			"expiresAt": time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
			"ip":        "192.168.3.64",
			"ua":        "ua",
		}},
	}
	b, _ := json.Marshal(legacy)
	if err := os.WriteFile(path, b, 0o600); err != nil {
		t.Fatal(err)
	}

	st, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if st.GetSession("legacy-plain-token") == nil {
		t.Fatal("旧版本的明文 token 在升级后必须仍然有效（否则所有人被强制登出）")
	}
	after, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(after), "legacy-plain-token") {
		t.Fatal("旧 token 应被就地迁移成哈希")
	}
}
