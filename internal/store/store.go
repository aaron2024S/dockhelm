// Package store 是 Dockhelm 的持久化层。
//
// 实现方式：单个 JSON 文档 + 原子落盘（写临时文件 → fsync → rename），
// 内存里保存全部状态并用一把互斥锁串行写入。
//
// 为什么不上 SQLite：Dockhelm 的数据量极小且有明确上限（配置几十条、历史各留 500 条、
// 会话几十条），总量不到 1 MB，全量载入内存绰绰有余。换成 JSON 之后：
//   - 依赖为零，`CGO_ENABLED=0` 交叉编译 amd64/arm64 毫无悬念；
//   - 不会出现 NAS 上 SQLite 文件被异常掉电写坏、需要手工恢复的情况；
//   - 备份数据目录 = 拷贝几个可读的 JSON，用户随时能自己看一眼。
package store

import (
	"encoding/json"
	"errors"
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
	maxRunLogs        = 500
	maxNotifyHistory  = 500
	maxLoginAttempts  = 500
)

// Store 持久化句柄。
type Store struct {
	path string
	mu   sync.Mutex
	db   data
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
			return nil, err
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

// SetSetting 写设置（不立即落盘，批量改完调用 Save）。
func (s *Store) SetSetting(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db.Settings[key] = value
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
type Session struct {
	Token     string `json:"token"`
	CreatedAt string `json:"createdAt"`
	ExpiresAt string `json:"expiresAt"`
	IP        string `json:"ip"`
	UA        string `json:"ua"`
}

// CreateSession 新建会话。
func (s *Store) CreateSession(token string, ttl time.Duration, ip, ua string) error {
	now := time.Now()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.db.Sessions = append(s.db.Sessions, Session{
		Token:     token,
		CreatedAt: now.UTC().Format(time.RFC3339),
		ExpiresAt: now.Add(ttl).UTC().Format(time.RFC3339),
		IP:        ip,
		UA:        truncate(ua, 200),
	})
	s.pruneSessionsLocked()
	return s.flushLocked()
}

// GetSession 取有效会话。
func (s *Store) GetSession(token string) *Session {
	if token == "" {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for i := range s.db.Sessions {
		if s.db.Sessions[i].Token == token {
			et, err := time.Parse(time.RFC3339, s.db.Sessions[i].ExpiresAt)
			if err != nil || time.Now().After(et) {
				return nil
			}
			cp := s.db.Sessions[i]
			return &cp
		}
	}
	return nil
}

// DeleteSession 删除会话。
func (s *Store) DeleteSession(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := s.db.Sessions[:0]
	for _, it := range s.db.Sessions {
		if it.Token != token {
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

// SessionCount 当前有效会话数。
func (s *Store) SessionCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	n := 0
	now := time.Now()
	for _, it := range s.db.Sessions {
		if et, err := time.Parse(time.RFC3339, it.ExpiresAt); err == nil && now.Before(et) {
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
		if et, err := time.Parse(time.RFC3339, it.ExpiresAt); err == nil && now.Before(et) {
			out = append(out, it)
		}
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

// CountRecentFailures 某 IP 在窗口内的失败次数。
func (s *Store) CountRecentFailures(ip string, window time.Duration) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	since := time.Now().Add(-window)
	n := 0
	for _, it := range s.db.Logins {
		if it.OK || it.IP != ip {
			continue
		}
		if ts, err := time.Parse(time.RFC3339, it.TS); err == nil && ts.After(since) {
			n++
		}
	}
	return n
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
