// Package auth 实现单用户密码登录。
//
// 安全模型：Dockhelm 挂着 docker.sock，等价于宿主机 root，所以「不设密码就等于把
// 宿主机交给局域网」。因此首启强制设置密码，且明文永不落盘。
package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/aaron2024s/dockhelm/internal/store"
)

// CookieName 会话 Cookie 名。
const CookieName = "dockhelm_session"

// 会话时长
const (
	SessionTTL      = 7 * 24 * time.Hour
	SessionTTLKeep  = 30 * 24 * time.Hour
	LockWindow      = 5 * time.Minute
	MaxFailures     = 5
	MinPasswordLen  = 6
)

// ErrLocked 连续失败被临时锁定。
var ErrLocked = errors.New("尝试次数过多，请稍后再试")

// ErrBadPassword 密码错误。
var ErrBadPassword = errors.New("密码错误")

// ErrNotInitialized 尚未设置密码。
var ErrNotInitialized = errors.New("尚未初始化密码")

// ErrWeakPassword 密码太短。
var ErrWeakPassword = fmt.Errorf("密码至少需要 %d 位", MinPasswordLen)

// manager 认证数据（只存 bcrypt 哈希）。
type authFile struct {
	Version     int    `json:"version"`
	PasswordHash string `json:"passwordHash"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

// Manager 认证管理器。
type Manager struct {
	path  string
	store *store.Store
	mu    sync.RWMutex
	data  authFile

	// Notify 登录成功/失败的告警钩子（由 main 注入，避免循环依赖）
	Notify func(event string, vars map[string]string)
}

// New 创建认证管理器。path 为 auth.json 路径。
func New(path string, st *store.Store) (*Manager, error) {
	m := &Manager{path: path, store: st}
	if err := m.load(); err != nil {
		return nil, err
	}
	return m, nil
}

func (m *Manager) load() error {
	b, err := os.ReadFile(m.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 未初始化
		}
		return err
	}
	return json.Unmarshal(b, &m.data)
}

func (m *Manager) save() error {
	if err := os.MkdirAll(filepath.Dir(m.path), 0o700); err != nil {
		return err
	}
	b, err := json.MarshalIndent(m.data, "", "  ")
	if err != nil {
		return err
	}
	// 0600：即使只是哈希，也不该让其他用户读到
	return os.WriteFile(m.path, b, 0o600)
}

// Initialized 是否已经设置过密码。
func (m *Manager) Initialized() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data.PasswordHash != ""
}

// SetPassword 设置/修改密码。oldPassword 为空时表示首次设置（或已被重置）。
func (m *Manager) SetPassword(oldPassword, newPassword string) error {
	if len(newPassword) < MinPasswordLen {
		return ErrWeakPassword
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.data.PasswordHash != "" && oldPassword != "" {
		if bcrypt.CompareHashAndPassword([]byte(m.data.PasswordHash), []byte(oldPassword)) != nil {
			return ErrBadPassword
		}
	} else if m.data.PasswordHash != "" && oldPassword == "" {
		// 已有密码但没给旧密码 ⇒ 拒绝（防止未授权改密）
		return ErrBadPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	now := time.Now().UTC().Format(time.RFC3339)
	if m.data.CreatedAt == "" {
		m.data.CreatedAt = now
	}
	m.data.UpdatedAt = now
	m.data.Version = 1
	m.data.PasswordHash = string(hash)
	if err := m.save(); err != nil {
		return err
	}
	// 改密后其它设备的会话立即失效
	return m.store.DeleteAllSessions()
}

// ForceSetPassword 无条件覆盖密码（DOCKHELM_PASSWORD 环境变量或删除 auth.json 后的首启）。
func (m *Manager) ForceSetPassword(pw string) error {
	if len(pw) < MinPasswordLen {
		return ErrWeakPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now().UTC().Format(time.RFC3339)
	if m.data.CreatedAt == "" {
		m.data.CreatedAt = now
	}
	m.data.UpdatedAt = now
	m.data.Version = 1
	m.data.PasswordHash = string(hash)
	if err := m.save(); err != nil {
		return err
	}
	return m.store.DeleteAllSessions()
}

// Verify 校验密码，成功返回新会话 token。
func (m *Manager) Verify(password, ip, ua string, keep bool) (string, error) {
	if m.store.CountRecentFailures(ip, LockWindow) >= MaxFailures {
		return "", ErrLocked
	}
	m.mu.RLock()
	hash := m.data.PasswordHash
	m.mu.RUnlock()

	if hash == "" {
		return "", ErrNotInitialized
	}
	// 时序安全由 bcrypt 的恒定时间比较保证
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		m.store.RecordLogin(ip, false)
		// 故意延迟 1s，抬高暴力破解成本
		time.Sleep(time.Second)
		if m.store.CountRecentFailures(ip, LockWindow)+1 == MaxFailures && m.Notify != nil {
			m.Notify("login_locked", map[string]string{
				"message": fmt.Sprintf("来源 %s 连续 %d 次密码错误，已临时锁定 %s", ip, MaxFailures, LockWindow),
			})
		}
		return "", ErrBadPassword
	}

	m.store.RecordLogin(ip, true)
	ttl := SessionTTL
	if keep {
		ttl = SessionTTLKeep
	}
	token, err := randomToken(32)
	if err != nil {
		return "", err
	}
	if err := m.store.CreateSession(token, ttl, ip, ua); err != nil {
		return "", err
	}
	return token, nil
}

// Validate 校验会话 token 是否有效。
func (m *Manager) Validate(token string) bool {
	if token == "" {
		return false
	}
	return m.store.GetSession(token) != nil
}

// Logout 注销会话。
func (m *Manager) Logout(token string) error {
	return m.store.DeleteSession(token)
}

// Failures 返回某 IP 当前窗口内的失败次数（前端用来显示剩余尝试次数）。
func (m *Manager) Failures(ip string) int {
	return m.store.CountRecentFailures(ip, LockWindow)
}

func randomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// NewSecret 生成一个 base64 随机串（用于内部用途）。
func NewSecret(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// MaskSession 用于展示。
func MaskSession(t string) string {
	if len(t) <= 8 {
		return strings.Repeat("*", len(t))
	}
	return t[:4] + "…" + t[len(t)-4:]
}
