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
	"math"
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
//
// 勾「保持登录」⇒ 持久化 Cookie，7 天；不勾 ⇒ 会话 Cookie（关掉浏览器即失效），
// 服务端仍给 24 小时上限兜底 —— 有些浏览器会「恢复上次标签页」时连会话 Cookie
// 一起恢复，没有上限的话一个忘了关的标签页可以一直活着。
// ⚠ 这两个值要和登录页文案「保持登录（7 天）」保持一致，改一处要改两处。
const (
	SessionTTLKeep = 7 * 24 * time.Hour
	SessionTTL     = 24 * time.Hour
	LockWindow     = 5 * time.Minute
	MaxFailures    = 5
	MinPasswordLen = 6
)

// ErrLocked 连续失败被临时锁定。
//
// 这条只是兜底文案：真正给用户看的话在 LockedHint 里生成（带还要等多久）。
var ErrLocked = errors.New("尝试次数过多，请稍后再试")

// ErrBadPassword 密码错误。
var ErrBadPassword = errors.New("密码错误")

// ErrNotInitialized 尚未设置密码。
var ErrNotInitialized = errors.New("尚未初始化密码")

// ErrWeakPassword 密码太短。
var ErrWeakPassword = fmt.Errorf("密码至少需要 %d 位", MinPasswordLen)

// ErrAlreadyInitialized 已经设置过密码。
var ErrAlreadyInitialized = errors.New("密码已设置")

// manager 认证数据（只存 bcrypt 哈希）。
type authFile struct {
	Version      int    `json:"version"`
	PasswordHash string `json:"passwordHash"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// Manager 认证管理器。
type Manager struct {
	path  string
	store *store.Store
	mu    sync.RWMutex
	data  authFile
	// loadWarning 启动时发现的问题（如 auth.json 损坏被隔离），由 main 负责打日志。
	loadWarning string

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
	if err := json.Unmarshal(b, &m.data); err != nil {
		// auth.json 坏了**不能**让进程起不来（NAS 上那意味着只能进 SSH 手删文件）。
		// 把它挪到一边、当作「未初始化」继续跑，并留一条 warning 给 main 打日志。
		// 这不降低安全性：能碰到这个文件的人本来就能直接删掉它，效果完全一样。
		quarantine := m.path + ".corrupt-" + time.Now().Format("20060102-150405")
		if mvErr := os.Rename(m.path, quarantine); mvErr != nil {
			return fmt.Errorf("auth.json 解析失败且无法移走：%w（原错误：%v）", mvErr, err)
		}
		m.data = authFile{}
		m.loadWarning = fmt.Sprintf(
			"auth.json 解析失败，已改名为 %s 并按「未初始化」处理，请重新设置登录密码（%v）",
			filepath.Base(quarantine), err)
		return nil
	}
	return nil
}

// LoadWarning 返回启动时的可读告警（目前只有「auth.json 损坏已被隔离」一种）。
func (m *Manager) LoadWarning() string { return m.loadWarning }

func (m *Manager) save() error {
	dir := filepath.Dir(m.path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	b, err := json.MarshalIndent(m.data, "", "  ")
	if err != nil {
		return err
	}
	// 先写临时文件、fsync、再 rename —— 与 store 的落盘方式一致。
	// 直接 os.WriteFile 是「截断 + 写」，NAS 上掉电/被杀进程会留下半截 JSON，
	// 而 auth.json 半截就等于启动时解析失败。
	tmp, err := os.CreateTemp(dir, ".auth-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() { _ = os.Remove(tmpName) }() // rename 成功后这里是空操作
	if err := tmp.Chmod(0o600); err != nil {  // 即使只是哈希，也不该让其他用户读到
		_ = tmp.Close()
		return err
	}
	if _, err := tmp.Write(b); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, m.path)
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

// SetupPassword 首次设置密码。
//
// 与 ForceSetPassword 的区别：**「还没设过密码」这个判断和写入在同一把锁里完成**。
// 首启接口原来是先 Initialized() 再 ForceSetPassword，两个并发请求可能都通过判断，
// 后一个把前一个刚设的密码覆盖掉（后者自己还拿到一个会话）。窗口很小，但一次就够。
func (m *Manager) SetupPassword(pw string) error {
	if len(pw) < MinPasswordLen {
		return ErrWeakPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.data.PasswordHash != "" {
		return ErrAlreadyInitialized
	}
	now := time.Now().UTC().Format(time.RFC3339)
	m.data.CreatedAt = now
	m.data.UpdatedAt = now
	m.data.Version = 1
	m.data.PasswordHash = string(hash)
	if err := m.save(); err != nil {
		return err
	}
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
	if m.LockedFor(ip) > 0 {
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
		m.store.RecordLogin(ip, false) // 本次失败已经记进去了
		// 这里**故意不再 Sleep 1 秒**。原来的写法有两个毛病：
		//   ① 它挡的是本进程的 HTTP handler（同步占着 goroutine），并发爆破时先被拖死的是面板自己；
		//   ② 对「本人手滑打错一个字」也要罚站 1 秒，而真正的防线是下面的滑动窗口锁定
		//      （5 次 / 5 分钟，命中即立刻 429、不比对密码），单次延时对它的边际贡献接近零。
		//
		// 触发锁定的这一次也直接按「已锁定」返回：否则前端只会看到「密码错误，还可尝试 0 次」，
		// 既不说锁了多久，也不说什么时候能重试。
		if n := m.store.CountRecentFailures(ip, LockWindow); n >= MaxFailures {
			if m.Notify != nil {
				m.Notify("login_locked", map[string]string{
					"message": fmt.Sprintf("来源 %s 连续 %d 次密码错误，已临时锁定 %s", ip, n, LockWindow),
				})
			}
			return "", ErrLocked
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

// LockedFor 返回该 IP 还需要等多久才能再次尝试；未锁定返回 0。
//
// 锁定是滑动窗口（窗口内失败 ≥ MaxFailures 就锁），所以解锁时刻不是「最后一次失败
// + 5 分钟」，而是**第 (n-MaxFailures+1) 早的那次失败滑出窗口**的瞬间 —— 只有它掉出
// 窗口，计数才会跌破 MaxFailures。按前者算会把等待时间报长，按后者才和真实放行时刻对上。
func (m *Manager) LockedFor(ip string) time.Duration {
	ts := m.store.RecentFailureTimes(ip, LockWindow)
	if len(ts) < MaxFailures {
		return 0
	}
	unlock := ts[len(ts)-MaxFailures].Add(LockWindow)
	if left := time.Until(unlock); left > 0 {
		return left
	}
	return 0
}

// LockedHint 生成给用户看的锁定提示（含还需等待多久）。
//
// 不满 1 分钟按秒报（否则「请约 1 分钟后再试」会在还剩 3 秒时骗人），
// 超过后向上取整到分钟 —— 宁可说长一点，也别让人在还剩 20 秒时反复试。
func LockedHint(left time.Duration) string {
	if left <= 0 {
		return "尝试次数过多，请稍后再试"
	}
	sec := int(math.Ceil(left.Seconds()))
	if sec < 60 {
		return fmt.Sprintf("尝试次数过多，请约 %d 秒后再试", sec)
	}
	return fmt.Sprintf("尝试次数过多，请约 %d 分钟后再试", (sec+59)/60)
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
