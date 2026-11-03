// Package store 是 Dockhelm 的持久化层。
//
// 实现方式：单个 JSON 文档 + 原子落盘（写临时文件 → fsync → rename），
// 内存里保存全部状态并用一把互斥锁串行写入。
//
// 为什么不上 SQLite：Dockhelm 的数据量极小且有明确上限（配置几十条、运行记录默认 100 条、
// 推送与登录历史各留 500 条、会话几十条），总量不到 1 MB，全量载入内存绰绰有余。换成 JSON 之后：
//   - 依赖为零，`CGO_ENABLED=0` 交叉编译 amd64/arm64 毫无悬念；
//   - 不会出现 NAS 上 SQLite 文件被异常掉电写坏、需要手工恢复的情况；
//   - 备份数据目录 = 拷贝几个可读的 JSON，用户随时能自己看一眼。
package store

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 各类历史的保留上限（超出即截断最旧的）。
const (
	// maxRunLogs 是「设置里没存过 log.retention」时的运行记录保留条数，
	// 必须与 api 侧 readSettings 的默认值保持一致 —— 两处不一致会出现
	// 「设置页显示 100、实际留了 500」这种自相矛盾的状态。
	maxRunLogs       = 100
	maxNotifyHistory = 500
	maxLoginAttempts = 500
)

// SessionIdle 会话连续多久没被使用就失效。
//
// 「保持登录 7 天」管的是 Cookie 寿命，不是「这串 token 永远有效」：
// 一台被忘在角落的旧手机、一个导出的数据目录，都会把这串 token 一起带走。
// 48 小时没访问就作废，正常每天都开面板的人完全感觉不到。
const SessionIdle = 48 * time.Hour

// sessionTouchInterval 两次「刷新最后访问时间」之间的最小间隔。
//
// GetSession 在每一次 API 请求上都会跑，每次都落盘等于每个请求一次 fsync。
// 5 分钟刷一次足够精确（失效判定的粒度本来就是小时级）。
const sessionTouchInterval = 5 * time.Minute

// Store 持久化句柄。
type Store struct {
	path string
	mu   sync.Mutex
	db   data

	// loadWarning 启动时发现的问题（如数据文件损坏被隔离），由 main 负责打日志。
	loadWarning string
	// lastTouchFlush 上次把「会话最后访问时间」落盘的时刻。
	lastTouchFlush time.Time
}

// data 是落盘的完整文档。
type data struct {
	Version   int                    `json:"version"`
	Settings  map[string]string      `json:"settings"`
	Sessions  []Session              `json:"sessions"`
	Schedules []Schedule             `json:"schedules"`
	Channels  []Channel              `json:"channels"`
	Events    map[string]NotifyEvent `json:"notifyEvents"`
	RunLogs   []RunLog               `json:"runLogs"`
	NotifyLog []NotifyRecord         `json:"notifyHistory"`
	Logins    []LoginAttempt         `json:"loginAttempts"`
	Seq       int64                  `json:"seq"`
}

func emptyData() data {
	return data{
		Version:  1,
		Settings: map[string]string{},
		Events:   map[string]NotifyEvent{},
	}
}

// Open 载入（必要时创建）存储文件。
func Open(path string) (*Store, error) {
	s := &Store{path: path, db: emptyData()}
	b, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return nil, err
			}
			return s, s.flushLocked()
		}
		return nil, err
	}
	if len(b) > 0 {
		if err := json.Unmarshal(b, &s.db); err != nil {
			// 数据文件坏了**不能**让进程起不来：NAS 上那意味着面板进不去、
			// 容器在重启循环里打转，用户只能 SSH 手删文件（而很多人根本没有 SSH）。
			// 这里把坏文件改名留档、按空数据启动，并在日志里说清楚发生了什么。
			//
			// 不降低安全性：能改这个文件的人本来就能直接删掉它。
			// 与 auth.json 的处理方式对齐（那边一直有这层兜底，store 之前没有）。
			quarantine := path + ".corrupt-" + time.Now().Format("20060102-150405")
			if mvErr := os.Rename(path, quarantine); mvErr != nil {
				return nil, fmt.Errorf("数据文件解析失败且无法移走：%w（原错误：%v）", mvErr, err)
			}
			s.db = emptyData()
			s.loadWarning = fmt.Sprintf(
				"数据文件解析失败，已改名为 %s 并按空数据启动（计划任务、通知渠道、设置需要重新配置；原错误：%v）",
				filepath.Base(quarantine), err)
			return s, s.flushLocked()
		}
	}
	if s.db.Settings == nil {
		s.db.Settings = map[string]string{}
	}
	if s.db.Events == nil {
		s.db.Events = map[string]NotifyEvent{}
	}
	return s, nil
}

// LoadWarning 返回启动时的可读告警（目前只有「数据文件损坏已被隔离」一种）。
func (s *Store) LoadWarning() string { return s.loadWarning }

// Close 立即落盘（进程退出时调用）。
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.flushLocked()
}

// flushLocked 原子写盘。调用方必须持有锁。
func (s *Store) flushLocked() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(&s.db, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	if _, err := f.Write(b); err != nil {
		f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// ---------- settings ----------

// GetSetting 读设置。
func (s *Store) GetSetting(key, def string) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, ok := s.db.Settings[key]; ok && v != "" {
		return v
	}
	return def
}

// SetSetting 写单个设置项并立刻落盘。
//
// ⚠ 每次调用都是一次「整份 JSON 写盘 + fsync」。一次要改多项时用 SetMany，
// 否则保存一次设置会把整个数据文件写十几遍（机械盘/SD 卡上很明显）。
func (s *Store) SetSetting(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db.Settings[key] = value
	return s.flushLocked()
}

// SetMany 批量写设置项，只落盘一次。
func (s *Store) SetMany(kv map[string]string) error {
	if len(kv) == 0 {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for k, v := range kv {
		s.db.Settings[k] = v
	}
	return s.flushLocked()
}

// GetJSON 反序列化设置项。
func (s *Store) GetJSON(key string, v any) {
	raw := s.GetSetting(key, "")
	if raw == "" {
		return
	}
	_ = json.Unmarshal([]byte(raw), v)
}

// SetJSON 序列化写入设置项。
func (s *Store) SetJSON(key string, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.SetSetting(key, string(b))
}

// AllSettings 返回全部设置（浅拷贝）。
func (s *Store) AllSettings() map[string]string {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[string]string, len(s.db.Settings))
	for k, v := range s.db.Settings {
		out[k] = v
	}
	return out
}

// ---------- sessions ----------

// Session 登录会话。
//
// Token 字段存的是 **SHA-256(token)**，不是原文 —— 数据文件一旦被拷走
// （备份、误传到网盘），拿到的人也没法直接拿去当 Cookie 用。
type Session struct {
	Token     string `json:"token"`
	CreatedAt string `json:"createdAt"`
	ExpiresAt string `json:"expiresAt"`
	LastSeen  string `json:"lastSeen,omitempty"`
	IP        string `json:"ip"`
	UA        string `json:"ua"`
}

// hashToken 会话 token 的落盘形式。
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// CreateSession 新建会话。
func (s *Store) CreateSession(token string, ttl time.Duration, ip, ua string) error {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db.Sessions = append(s.db.Sessions, Session{
		Token:     hashToken(token),
		CreatedAt: now.UTC().Format(time.RFC3339),
		ExpiresAt: now.Add(ttl).UTC().Format(time.RFC3339),
		LastSeen:  now.UTC().Format(time.RFC3339),
		IP:        ip,
		UA:        truncate(ua, 200),
	})
	s.pruneSessionsLocked()
	return s.flushLocked()
}

// GetSession 取有效会话（token 为 Cookie 里的原文）。
//
// 同时刷新最后访问时间，用于「空闲超时」。刷新按 5 分钟节流，
// 否则每个 API 请求都要落盘一次。
func (s *Store) GetSession(token string) *Session {
	if token == "" {
		return nil
	}
	h := hashToken(token)
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.db.Sessions {
		it := &s.db.Sessions[i]
		// 旧版本落盘的是明文 token。第一次带着旧 Cookie 进来时原地升级成哈希，
		// 用户不会被强制登出。
		legacy := it.Token == token
		if it.Token != h && !legacy {
			continue
		}
		et, err := time.Parse(time.RFC3339, it.ExpiresAt)
		if err != nil || now.After(et) {
			return nil
		}
		if t, err := time.Parse(time.RFC3339, it.LastSeen); err == nil && now.Sub(t) > SessionIdle {
			return nil
		}
		it.LastSeen = now.UTC().Format(time.RFC3339)
		if legacy {
			it.Token = h
			_ = s.flushLocked()
			s.lastTouchFlush = now
		} else if now.Sub(s.lastTouchFlush) > sessionTouchInterval {
			if err := s.flushLocked(); err == nil {
				s.lastTouchFlush = now
			}
		}
		cp := *it
		return &cp
	}
	return nil
}

// DeleteSession 删除会话。
func (s *Store) DeleteSession(token string) error {
	h := hashToken(token)
	s.mu.Lock()
	defer s.mu.Unlock()
	out := s.db.Sessions[:0]
	for _, it := range s.db.Sessions {
		if it.Token != h && it.Token != token {
			out = append(out, it)
		}
	}
	s.db.Sessions = out
	return s.flushLocked()
}

// DeleteAllSessions 清空会话（改密后调用）。
func (s *Store) DeleteAllSessions() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db.Sessions = nil
	return s.flushLocked()
}

// DeleteOtherSessions 强制登出 current 之外的所有会话（其他浏览器/设备上的残留登录），
// 当前会话原样保留。返回被登出的会话数。
func (s *Store) DeleteOtherSessions(current string) (int, error) {
	h := hashToken(current)
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Session, 0, len(s.db.Sessions))
	removed := 0
	for _, it := range s.db.Sessions {
		if it.Token == h || it.Token == current {
			out = append(out, it)
		} else {
			removed++
		}
	}
	s.db.Sessions = out
	return removed, s.flushLocked()
}

// SessionCount 当前有效会话数。
func (s *Store) SessionCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	now := time.Now()
	for _, it := range s.db.Sessions {
		if et, err := time.Parse(time.RFC3339, it.ExpiresAt); err == nil && now.Before(et) {
			if t, err := time.Parse(time.RFC3339, it.LastSeen); err == nil && now.Sub(t) > SessionIdle {
				continue
			}
			n++
		}
	}
	return n
}

// PurgeExpiredSessions 清理过期项。
func (s *Store) PurgeExpiredSessions() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.pruneSessionsLocked() {
		_ = s.flushLocked()
	}
}

func (s *Store) pruneSessionsLocked() bool {
	now := time.Now()
	out := s.db.Sessions[:0]
	for _, it := range s.db.Sessions {
		et, err := time.Parse(time.RFC3339, it.ExpiresAt)
		if err != nil || !now.Before(et) {
			continue
		}
		if t, err := time.Parse(time.RFC3339, it.LastSeen); err == nil && now.Sub(t) > SessionIdle {
			continue
		}
		out = append(out, it)
	}
	changed := len(out) != len(s.db.Sessions)
	s.db.Sessions = out
	return changed
}

// ---------- login attempts ----------

// LoginAttempt 一条登录尝试。
type LoginAttempt struct {
	TS string `json:"ts"`
	IP string `json:"ip"`
	OK bool   `json:"ok"`
}

// RecordLogin 记录登录尝试。
func (s *Store) RecordLogin(ip string, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db.Logins = append(s.db.Logins, LoginAttempt{
		TS: time.Now().UTC().Format(time.RFC3339), IP: ip, OK: ok,
	})
	if len(s.db.Logins) > maxLoginAttempts {
		s.db.Logins = s.db.Logins[len(s.db.Logins)-maxLoginAttempts:]
	}
	_ = s.flushLocked()
}

// ClearLoginFailures 清掉某 IP 的失败记录（登录成功后调用）。
//
// 锁定的语义是「**连续**错 N 次」，所以成功一次就该把窗口清零。
// 不清的话「错 4 次 → 成功 → 再错 1 次」也会被锁，和界面上那句
// 「连续错 5 次会临时锁定」对不上 —— 用户会觉得莫名其妙。
// 首启设置密码的失败也走这套计数，成功设置后一并清零，免得刚设完就被锁。
func (s *Store) ClearLoginFailures(ip string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := s.db.Logins[:0]
	changed := false
	for _, it := range s.db.Logins {
		if !it.OK && it.IP == ip {
			changed = true
			continue
		}
		out = append(out, it)
	}
	s.db.Logins = out
	if changed {
		_ = s.flushLocked()
	}
}

// CountRecentFailures 某 IP 在窗口内的失败次数。
func (s *Store) CountRecentFailures(ip string, window time.Duration) int {
	return len(s.RecentFailureTimes(ip, window))
}

// RecentFailureTimes 某 IP 在窗口内的失败时刻，按时间升序返回。
//
// 给「还要锁多久」用：锁定是滑动窗口，解锁时刻 = 第 (n-MaxFailures+1) 早的那次
// 失败滑出窗口的瞬间。只想知道次数的话用 CountRecentFailures。
func (s *Store) RecentFailureTimes(ip string, window time.Duration) []time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	since := time.Now().Add(-window)
	out := []time.Time{}
	for _, it := range s.db.Logins {
		if it.OK || it.IP != ip {
			continue
		}
		if ts, err := time.Parse(time.RFC3339, it.TS); err == nil && ts.After(since) {
			out = append(out, ts)
		}
	}
	// 记录是按调用顺序追加的，正常已经有序；时间回拨等极端情况下排一次更稳。
	sort.Slice(out, func(i, j int) bool { return out[i].Before(out[j]) })
	return out
}

// CountRecentGlobalFailures 全局限定窗口内的失败次数。
func (s *Store) CountRecentGlobalFailures(window time.Duration) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	since := time.Now().Add(-window)
	n := 0
	for _, it := range s.db.Logins {
		if it.OK {
			continue
		}
		if ts, err := time.Parse(time.RFC3339, it.TS); err == nil && ts.After(since) {
			n++
		}
	}
	return n
}

// PurgeOldLoginAttempts 清理 7 天前的登录记录。
func (s *Store) PurgeOldLoginAttempts() {
	s.mu.Lock()
	defer s.mu.Unlock()
	cut := time.Now().Add(-7 * 24 * time.Hour)
	out := s.db.Logins[:0]
	for _, it := range s.db.Logins {
		if ts, err := time.Parse(time.RFC3339, it.TS); err == nil && ts.After(cut) {
			out = append(out, it)
		}
	}
	s.db.Logins = out
	_ = s.flushLocked()
}

// RecentLogins 最近的登录尝试（新的在前）。
func (s *Store) RecentLogins(limit int) []LoginAttempt {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]LoginAttempt, 0, limit)
	for i := len(s.db.Logins) - 1; i >= 0 && len(out) < limit; i-- {
		out = append(out, s.db.Logins[i])
	}
	return out
}

// ---------- schedules ----------

// Schedule 计划任务。
type Schedule struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name"`
	Cron        string   `json:"cron"`
	Action      string   `json:"action"`  // start | stop | restart | update | backup | prune_images
	Targets     []string `json:"targets"` // 容器名列表，空 = 全部
	Enabled     bool     `json:"enabled"`
	LastRun     string   `json:"lastRun"`
	LastStatus  string   `json:"lastStatus"`
	LastMessage string   `json:"lastMessage"`
	CreatedAt   string   `json:"createdAt"`
	NextRun     string   `json:"nextRun,omitempty"` // 由调度器填充
}

// ListSchedules 列出全部计划任务。
func (s *Store) ListSchedules() ([]Schedule, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Schedule, len(s.db.Schedules))
	copy(out, s.db.Schedules)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

// GetSchedule 取单条。
func (s *Store) GetSchedule(id int64) (*Schedule, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.db.Schedules {
		if s.db.Schedules[i].ID == id {
			cp := s.db.Schedules[i]
			return &cp, nil
		}
	}
	return nil, errors.New("schedule not found")
}

// CreateSchedule 新建。
func (s *Store) CreateSchedule(in Schedule) (*Schedule, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db.Seq++
	in.ID = s.db.Seq
	in.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	if in.Targets == nil {
		in.Targets = []string{}
	}
	s.db.Schedules = append(s.db.Schedules, in)
	if err := s.flushLocked(); err != nil {
		return nil, err
	}
	cp := in
	return &cp, nil
}

// UpdateSchedule 更新。
func (s *Store) UpdateSchedule(in Schedule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.db.Schedules {
		if s.db.Schedules[i].ID == in.ID {
			in.CreatedAt = s.db.Schedules[i].CreatedAt
			in.LastRun = s.db.Schedules[i].LastRun
			in.LastStatus = s.db.Schedules[i].LastStatus
			in.LastMessage = s.db.Schedules[i].LastMessage
			if in.Targets == nil {
				in.Targets = []string{}
			}
			s.db.Schedules[i] = in
			return s.flushLocked()
		}
	}
	return errors.New("schedule not found")
}

// SetScheduleResult 记录执行结果。
func (s *Store) SetScheduleResult(id int64, status, message string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.db.Schedules {
		if s.db.Schedules[i].ID == id {
			s.db.Schedules[i].LastRun = time.Now().UTC().Format(time.RFC3339)
			s.db.Schedules[i].LastStatus = status
			s.db.Schedules[i].LastMessage = message
			return s.flushLocked()
		}
	}
	return errors.New("schedule not found")
}

// DeleteSchedule 删除。
func (s *Store) DeleteSchedule(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := s.db.Schedules[:0]
	for _, it := range s.db.Schedules {
		if it.ID != id {
			out = append(out, it)
		}
	}
	s.db.Schedules = out
	return s.flushLocked()
}

// ---------- run logs ----------

// RunLog 一条运行记录。
type RunLog struct {
	ID      int64  `json:"id"`
	TS      string `json:"ts"`
	Kind    string `json:"kind"`
	Ref     string `json:"ref"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

// AddRunLog 追加运行记录。
func (s *Store) AddRunLog(kind, ref, status, message, detail string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db.Seq++
	s.db.RunLogs = append(s.db.RunLogs, RunLog{
		ID: s.db.Seq, TS: time.Now().UTC().Format(time.RFC3339),
		Kind: kind, Ref: ref, Status: status, Message: message, Detail: detail,
	})
	limit := s.runLogLimitLocked()
	if len(s.db.RunLogs) > limit {
		s.db.RunLogs = s.db.RunLogs[len(s.db.RunLogs)-limit:]
	}
	_ = s.flushLocked()
}

// runLogLimitLocked 取运行记录的保留条数（调用方必须已持锁）。
//
// 以前这里是硬编码的 500，而设置页里那个「日志保留条数」只存不读 ——
// 用户改了完全没效果。现在按设置生效，并夹在 [50, 100000]：
// 下限防止手改文件写成 0（那会让日志一条都不剩），上限防止把落盘文件撑爆。
func (s *Store) runLogLimitLocked() int {
	limit := maxRunLogs
	if v, ok := s.db.Settings["log.retention"]; ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			limit = n
		}
	}
	if limit < 50 {
		limit = 50
	}
	if limit > 100000 {
		limit = 100000
	}
	return limit
}

// ListRunLogs 按时间倒序取运行记录。
func (s *Store) ListRunLogs(limit int) ([]RunLog, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if limit <= 0 || limit > len(s.db.RunLogs) {
		limit = len(s.db.RunLogs)
	}
	out := make([]RunLog, 0, limit)
	for i := len(s.db.RunLogs) - 1; i >= 0 && len(out) < limit; i-- {
		out = append(out, s.db.RunLogs[i])
	}
	return out, nil
}

// LastRunLog 取某一类（kind）最近的一条运行记录。
//
// 「顺序启动」页顶部要显示「上次执行：今天 08:12 · 成功 47、跳过 3」。
// 不能用 ListRunLogs(N) 再自己筛：开机一轮之后可能已经积累了几百条别的记录，
// 那条 boot 早就掉出窗口，页面就会永远显示「从未执行」。
func (s *Store) LastRunLog(kind string) (RunLog, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := len(s.db.RunLogs) - 1; i >= 0; i-- {
		if s.db.RunLogs[i].Kind == kind {
			return s.db.RunLogs[i], true
		}
	}
	return RunLog{}, false
}

// ClearRunLogs 清空运行记录。
func (s *Store) ClearRunLogs() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db.RunLogs = nil
	_ = s.flushLocked()
}

// ---------- notify ----------

// Channel 通知渠道。
type Channel struct {
	ID        int64          `json:"id"`
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Enabled   bool           `json:"enabled"`
	Config    map[string]any `json:"config"`
	CreatedAt string         `json:"createdAt"`
}

// ListChannels 列出渠道。
func (s *Store) ListChannels() ([]Channel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Channel, len(s.db.Channels))
	copy(out, s.db.Channels)
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

// GetChannel 取单条渠道。
func (s *Store) GetChannel(id int64) (*Channel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.db.Channels {
		if s.db.Channels[i].ID == id {
			cp := s.db.Channels[i]
			return &cp, nil
		}
	}
	return nil, errors.New("channel not found")
}

// CreateChannel 新建渠道。
func (s *Store) CreateChannel(in Channel) (*Channel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db.Seq++
	in.ID = s.db.Seq
	in.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	if in.Config == nil {
		in.Config = map[string]any{}
	}
	s.db.Channels = append(s.db.Channels, in)
	if err := s.flushLocked(); err != nil {
		return nil, err
	}
	cp := in
	return &cp, nil
}

// UpdateChannel 更新渠道。
func (s *Store) UpdateChannel(in Channel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.db.Channels {
		if s.db.Channels[i].ID == in.ID {
			in.CreatedAt = s.db.Channels[i].CreatedAt
			if in.Config == nil {
				in.Config = map[string]any{}
			}
			s.db.Channels[i] = in
			return s.flushLocked()
		}
	}
	return errors.New("channel not found")
}

// DeleteChannel 删除渠道。
func (s *Store) DeleteChannel(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := s.db.Channels[:0]
	for _, it := range s.db.Channels {
		if it.ID != id {
			out = append(out, it)
		}
	}
	s.db.Channels = out
	return s.flushLocked()
}

// NotifyEvent 事件订阅配置。
type NotifyEvent struct {
	Event   string `json:"event"`
	Enabled bool   `json:"enabled"`
	Level   string `json:"level"` // urgent | normal
}

// ListNotifyEvents 列出事件订阅。
func (s *Store) ListNotifyEvents() []NotifyEvent {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]NotifyEvent, 0, len(s.db.Events))
	for _, v := range s.db.Events {
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Event < out[j].Event })
	return out
}

// UpsertNotifyEvent 写事件订阅。
func (s *Store) UpsertNotifyEvent(event string, enabled bool, level string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db.Events[event] = NotifyEvent{Event: event, Enabled: enabled, Level: level}
	return s.flushLocked()
}

// BulkUpsertNotifyEvents 批量写事件订阅。
func (s *Store) BulkUpsertNotifyEvents(rows []NotifyEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, r := range rows {
		s.db.Events[r.Event] = r
	}
	return s.flushLocked()
}

// NotifyRecord 推送历史。
type NotifyRecord struct {
	ID      int64  `json:"id"`
	TS      string `json:"ts"`
	Event   string `json:"event"`
	Level   string `json:"level"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	OK      bool   `json:"ok"`
	ErrMsg  string `json:"errmsg"`
	Channel string `json:"channel"`
}

// AddNotifyRecord 写推送历史。
func (s *Store) AddNotifyRecord(event, level, title, body, channel string, ok bool, errmsg string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db.Seq++
	s.db.NotifyLog = append(s.db.NotifyLog, NotifyRecord{
		ID: s.db.Seq, TS: time.Now().UTC().Format(time.RFC3339),
		Event: event, Level: level, Title: title, Body: body,
		Channel: channel, OK: ok, ErrMsg: errmsg,
	})
	if len(s.db.NotifyLog) > maxNotifyHistory {
		s.db.NotifyLog = s.db.NotifyLog[len(s.db.NotifyLog)-maxNotifyHistory:]
	}
	_ = s.flushLocked()
}

// ListNotifyRecords 推送历史。
func (s *Store) ListNotifyRecords(limit int) ([]NotifyRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if limit <= 0 || limit > len(s.db.NotifyLog) {
		limit = len(s.db.NotifyLog)
	}
	out := make([]NotifyRecord, 0, limit)
	for i := len(s.db.NotifyLog) - 1; i >= 0 && len(out) < limit; i-- {
		out = append(out, s.db.NotifyLog[i])
	}
	return out, nil
}

// ClearNotifyHistory 清空推送历史。
func (s *Store) ClearNotifyHistory() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db.NotifyLog = nil
	_ = s.flushLocked()
}

// CountNotifySince 统计某时刻之后成功推送的条数（用于每日上限）。
func (s *Store) CountNotifySince(since time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	for _, it := range s.db.NotifyLog {
		if !it.OK {
			continue
		}
		if ts, err := time.Parse(time.RFC3339, it.TS); err == nil && ts.After(since) {
			n++
		}
	}
	return n
}

// LastNotifyTime 去重键最近一次发送时间。
func (s *Store) LastNotifyTime(key string) time.Time {
	raw := s.GetSetting("dedupe:"+key, "")
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}
	}
	return t
}

// MarkNotifyTime 记录去重键时间。
func (s *Store) MarkNotifyTime(key string) {
	_ = s.SetSetting("dedupe:"+key, time.Now().UTC().Format(time.RFC3339))
}

// DedupeKeys 返回全部去重键（用于清理）。
func (s *Store) DedupeKeys() []string {
	all := s.AllSettings()
	out := []string{}
	for k := range all {
		if strings.HasPrefix(k, "dedupe:") {
			out = append(out, k)
		}
	}
	return out
}

// DeleteSetting 删除设置项。
func (s *Store) DeleteSetting(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.db.Settings, key)
	return s.flushLocked()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
