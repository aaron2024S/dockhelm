// Package notify 实现事件通知：渠道管理、事件订阅、静默时段、防轰炸与投递。
//
// 三条硬约束（都是从「告警系统最终被关掉」的常见死因倒推出来的）：
//  1. 投递必须异步，且失败只记日志 —— 绝不能拖慢或阻断更新主流程；
//  2. 必须去重 —— 崩溃循环的容器一分钟能产生几十条事件；
//  3. 静默时段要能「攒着，时段结束汇总发一条」，而不是简单地把通知全吞掉。
package notify

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aaron2024s/dockhelm/internal/bus"
	"github.com/aaron2024s/dockhelm/internal/store"
	"github.com/aaron2024s/dockhelm/internal/version"
)

// 设置键
const (
	KeyEnabled       = "notify.enabled"
	KeyQuietEnabled  = "notify.quietEnabled"
	KeyQuietStart    = "notify.quietStart"
	KeyQuietEnd      = "notify.quietEnd"
	KeyQuietMode     = "notify.quietNormalMode" // digest | drop
	KeyQuietUrgent   = "notify.quietUrgentSend"
	KeyDedupeMinutes = "notify.dedupeWindow"
	KeyDailyLimit    = "notify.dailyLimit"
	KeyPanelURL      = "notify.panelURL"
)

// Manager 通知管理器。
type Manager struct {
	st     *store.Store
	bus    *bus.Bus
	client *http.Client
	host   string

	mu      sync.Mutex
	buffer  []buffered // 静默期攒下的普通事件
	stopped chan struct{}
	seq     int64
}

type buffered struct {
	Event string
	Vars  map[string]string
	At    time.Time
}

// New 创建通知管理器。
func New(st *store.Store, b *bus.Bus) *Manager {
	host, _ := os.Hostname()
	m := &Manager{
		st:   st,
		bus:  b,
		host: host,
		client: &http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig:     &tls.Config{MinVersion: tls.VersionTLS12},
				MaxIdleConnsPerHost: 4,
			},
		},
		stopped: make(chan struct{}),
	}
	return m
}

// Start 启动后台循环（静默期汇总 + 每日上限重置）。
func (m *Manager) Start() {
	go func() {
		t := time.NewTicker(30 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-m.stopped:
				return
			case <-t.C:
				m.maybeFlushDigest()
			}
		}
	}()
}

// Stop 停止后台循环。
func (m *Manager) Stop() { close(m.stopped) }

// ---------- 设置读取 ----------

func (m *Manager) boolSetting(key string, def bool) bool {
	v := m.st.GetSetting(key, "")
	if v == "" {
		return def
	}
	return v == "1" || v == "true"
}

func (m *Manager) intSetting(key string, def int) int {
	v := m.st.GetSetting(key, "")
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

// Settings 汇总通知相关设置（给前端）。
type Settings struct {
	Enabled       bool   `json:"enabled"`
	QuietEnabled  bool   `json:"quietEnabled"`
	QuietStart    string `json:"quietStart"`
	QuietEnd      string `json:"quietEnd"`
	QuietMode     string `json:"quietNormalMode"`
	QuietUrgent   bool   `json:"quietUrgentSend"`
	DedupeMinutes int    `json:"dedupeWindow"`
	DailyLimit    int    `json:"dailyLimit"`
	PanelURL      string `json:"panelURL"`
	SentToday     int    `json:"sentToday"`
	InQuiet       bool   `json:"inQuietHours"`
}

// GetSettings 读取设置。
func (m *Manager) GetSettings() Settings {
	return Settings{
		Enabled:       m.boolSetting(KeyEnabled, true),
		QuietEnabled:  m.boolSetting(KeyQuietEnabled, false),
		QuietStart:    m.st.GetSetting(KeyQuietStart, "23:00"),
		QuietEnd:      m.st.GetSetting(KeyQuietEnd, "07:00"),
		QuietMode:     m.st.GetSetting(KeyQuietMode, "digest"),
		QuietUrgent:   m.boolSetting(KeyQuietUrgent, true),
		DedupeMinutes: m.intSetting(KeyDedupeMinutes, 10),
		DailyLimit:    m.intSetting(KeyDailyLimit, 200),
		PanelURL:      m.st.GetSetting(KeyPanelURL, ""),
		SentToday:     m.st.CountNotifySince(startOfToday()),
		InQuiet:       m.inQuietHours(time.Now()),
	}
}

// SaveSettings 保存设置。
func (m *Manager) SaveSettings(s Settings) error {
	set := func(k, v string) error { return m.st.SetSetting(k, v) }
	b := func(v bool) string {
		if v {
			return "1"
		}
		return "0"
	}
	if err := set(KeyEnabled, b(s.Enabled)); err != nil {
		return err
	}
	if err := set(KeyQuietEnabled, b(s.QuietEnabled)); err != nil {
		return err
	}
	if err := set(KeyQuietStart, s.QuietStart); err != nil {
		return err
	}
	if err := set(KeyQuietEnd, s.QuietEnd); err != nil {
		return err
	}
	if err := set(KeyQuietMode, s.QuietMode); err != nil {
		return err
	}
	if err := set(KeyQuietUrgent, b(s.QuietUrgent)); err != nil {
		return err
	}
	if err := set(KeyDedupeMinutes, strconv.Itoa(s.DedupeMinutes)); err != nil {
		return err
	}
	if err := set(KeyDailyLimit, strconv.Itoa(s.DailyLimit)); err != nil {
		return err
	}
	return set(KeyPanelURL, s.PanelURL)
}

// ---------- 静默时段 ----------

func parseHHMM(s string) (int, int, bool) {
	parts := strings.Split(strings.TrimSpace(s), ":")
	if len(parts) != 2 {
		return 0, 0, false
	}
	h, err1 := strconv.Atoi(parts[0])
	mn, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || h < 0 || h > 23 || mn < 0 || mn > 59 {
		return 0, 0, false
	}
	return h, mn, true
}

// inQuietHours 判断当前是否处于静默时段（支持跨天，如 23:00–07:00）。
func (m *Manager) inQuietHours(now time.Time) bool {
	if !m.boolSetting(KeyQuietEnabled, false) {
		return false
	}
	sh, sm, ok1 := parseHHMM(m.st.GetSetting(KeyQuietStart, "23:00"))
	eh, em, ok2 := parseHHMM(m.st.GetSetting(KeyQuietEnd, "07:00"))
	if !ok1 || !ok2 {
		return false
	}
	cur := now.Hour()*60 + now.Minute()
	start := sh*60 + sm
	end := eh*60 + em
	if start == end {
		return false
	}
	if start < end {
		return cur >= start && cur < end
	}
	// 跨天
	return cur >= start || cur < end
}

// maybeFlushDigest 静默时段结束后，把攒下的普通事件汇总成一条发出。
func (m *Manager) maybeFlushDigest() {
	if m.inQuietHours(time.Now()) {
		return
	}
	m.mu.Lock()
	if len(m.buffer) == 0 {
		m.mu.Unlock()
		return
	}
	items := m.buffer
	m.buffer = nil
	m.mu.Unlock()
	m.sendDigest(items)
}

func (m *Manager) sendDigest(items []buffered) {
	if len(items) == 0 {
		return
	}
	counts := map[string]int{}
	for _, it := range items {
		counts[it.Event]++
	}
	lines := []string{}
	for ev, n := range counts {
		lines = append(lines, fmt.Sprintf("· %s ×%d", Label(ev), n))
	}
	body := fmt.Sprintf("静默时段内累计 %d 条通知：\n%s", len(items), strings.Join(lines, "\n"))
	m.deliver("digest", LevelNormal, "静默时段汇总（"+strconv.Itoa(len(items))+" 条）", body, nil, true)
}

// ---------- 投递 ----------

// Emit 投递一个事件。实现 Notifier 接口：**必须非阻塞**。
func (m *Manager) Emit(event string, vars map[string]string) {
	if !m.boolSetting(KeyEnabled, true) {
		return
	}
	if vars == nil {
		vars = map[string]string{}
	}
	go m.emit(event, vars)
}

// EmitSync 同步投递（仅用于「测试渠道」按钮，需要把结果返给用户）。
func (m *Manager) EmitSync(event string, vars map[string]string, channelID int64) (int, string) {
	if vars == nil {
		vars = map[string]string{}
	}
	return m.sendToChannel(channelID, event, vars)
}

func (m *Manager) emit(event string, vars map[string]string) {
	// 事件是否被订阅
	if !m.eventEnabled(event) {
		return
	}
	level := LevelOf(event)

	// 防轰炸：同容器同事件去重窗口
	dedupeKey := event + ":" + vars["container"]
	window := time.Duration(m.intSetting(KeyDedupeMinutes, 10)) * time.Minute
	if window > 0 {
		if last := m.st.LastNotifyTime(dedupeKey); !last.IsZero() && time.Since(last) < window {
			return
		}
	}

	// 静默时段
	if m.inQuietHours(time.Now()) {
		if level == LevelUrgent && m.boolSetting(KeyQuietUrgent, true) {
			m.deliver(event, level, "", "", vars, false)
			m.st.MarkNotifyTime(dedupeKey)
			return
		}
		if m.st.GetSetting(KeyQuietMode, "digest") == "drop" {
			return
		}
		m.mu.Lock()
		m.buffer = append(m.buffer, buffered{Event: event, Vars: vars, At: time.Now()})
		if len(m.buffer) > 500 {
			m.buffer = m.buffer[len(m.buffer)-500:]
		}
		m.mu.Unlock()
		m.st.MarkNotifyTime(dedupeKey)
		return
	}

	// 每日上限
	limit := m.intSetting(KeyDailyLimit, 200)
	if limit > 0 && m.st.CountNotifySince(startOfToday()) >= limit {
		m.st.AddNotifyRecord(event, level, "已超出每日推送上限", "本条通知被丢弃（可在设置里调整上限）", "", false, "daily limit")
		return
	}

	m.deliver(event, level, "", "", vars, false)
	m.st.MarkNotifyTime(dedupeKey)
}

func (m *Manager) eventEnabled(event string) bool {
	for _, e := range m.st.ListNotifyEvents() {
		if e.Event == event {
			return e.Enabled
		}
	}
	// 没有配置过 ⇒ 用出厂默认
	d, ok := EventMap[event]
	if !ok {
		return true
	}
	return d.Default
}

// deliver 把事件发往全部已启用的渠道。
func (m *Manager) deliver(event, level, titleOverride, textOverride string, vars map[string]string, skipQuiet bool) {
	title, text := m.compose(event, titleOverride, textOverride, vars)
	channels, err := m.st.ListChannels()
	if err != nil {
		return
	}
	sent := 0
	for _, ch := range channels {
		if !ch.Enabled {
			continue
		}
		if ok, msg := m.sendOne(ch, event, title, text, vars); ok {
			sent++
			m.st.AddNotifyRecord(event, level, title, text, ch.Name, true, "")
		} else {
			m.st.AddNotifyRecord(event, level, title, text, ch.Name, false, msg)
		}
	}
	m.bus.Publish("notify", "delivered", "info", map[string]any{
		"event": event, "level": level, "title": title, "text": text,
		"channels": len(channels), "sent": sent,
	})
}

// sendToChannel 只发某个渠道（测试按钮用）。
func (m *Manager) sendToChannel(id int64, event string, vars map[string]string) (int, string) {
	ch, err := m.st.GetChannel(id)
	if err != nil {
		return 0, "渠道不存在"
	}
	title, text := m.compose(event, "", "", vars)
	ok, msg := m.sendOne(*ch, event, title, text, vars)
	if ok {
		m.st.AddNotifyRecord(event, LevelOf(event), title, text, ch.Name, true, "")
		return 1, "发送成功"
	}
	m.st.AddNotifyRecord(event, LevelOf(event), title, text, ch.Name, false, msg)
	return 0, msg
}

// compose 组装标题与正文。
func (m *Manager) compose(event, titleOverride, textOverride string, vars map[string]string) (string, string) {
	if titleOverride != "" && textOverride != "" {
		return titleOverride, textOverride
	}
	label := Label(event)
	icon := "🔔"
	if LevelOf(event) == LevelUrgent {
		icon = "🚨"
	}
	head := label
	if c := vars["container"]; c != "" {
		head = fmt.Sprintf("%s · %s", label, c)
	}
	title := fmt.Sprintf("%s Dockhelm %s", icon, head)

	lines := []string{title}
	if vars["result"] != "" {
		lines = append(lines, "结果："+vars["result"])
	}
	if vars["message"] != "" {
		lines = append(lines, "详情："+vars["message"])
	}
	if vars["image"] != "" {
		lines = append(lines, "镜像："+vars["image"])
	}
	lines = append(lines, "时间："+time.Now().Format("2006-01-02 15:04:05"), "主机："+m.host)
	if u := m.st.GetSetting(KeyPanelURL, ""); u != "" {
		lines = append(lines, "面板："+u)
	}
	return title, strings.Join(lines, "\n")
}

// sendOne 投递到单个渠道，失败重试 3 次。返回 (是否成功, 错误信息)。
func (m *Manager) sendOne(ch store.Channel, event, title, text string, vars map[string]string) (bool, string) {
	if ch.Type == "smtp" {
		return m.sendSMTPWithRetry(ch, title, text)
	}
	var lastErr string
	for attempt := 0; attempt < 3; attempt++ {
		ok, err := m.sendHTTP(ch, event, title, text, vars)
		if ok {
			return true, ""
		}
		lastErr = err
		time.Sleep(time.Duration(attempt+1) * 700 * time.Millisecond)
	}
	return false, lastErr
}

func (m *Manager) sendHTTP(ch store.Channel, event, title, text string, vars map[string]string) (bool, string) {
	host, _ := os.Hostname()
	values := buildValues(event, vars, host, m.st.GetSetting(KeyPanelURL, ""), title, text)
	values["timestamp"] = strconv.FormatInt(time.Now().UnixMilli(), 10)
	values["sign"] = dingtalkSign(ch, values["timestamp"])
	// 飞书签名用的是另一套算法，单独放一个变量
	if ch.Type == "feishu" {
		values["sign"] = feishuSign(ch, values["timestamp"])
	}
	if _, ok := values["group"]; !ok {
		values["group"] = "dockhelm"
	}

	method, rawURL, headers, body, ct, err := RenderRequest(ch.Type, ch.Config, values)
	if err != nil {
		return false, err.Error()
	}
	if rawURL == "" {
		return false, "渠道未配置请求地址"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var rdr io.Reader
	if body != "" {
		rdr = strings.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, rawURL, rdr)
	if err != nil {
		return false, err.Error()
	}
	if ct != "" {
		req.Header.Set("Content-Type", ct)
	}
	req.Header.Set("User-Agent", version.AppName+"/"+version.Version)
	for k, v := range headers {
		if v == "" {
			continue
		}
		req.Header.Set(k, v)
	}
	resp, err := m.client.Do(req)
	if err != nil {
		return false, err.Error()
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
	if resp.StatusCode >= 300 {
		return false, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}
	// 有些服务即便 200 也会在 body 里报错（Telegram / 企业微信 / 钉钉都这样）
	if msg := extractProviderError(raw); msg != "" {
		return false, msg
	}
	return true, ""
}

// extractProviderError 尽量从响应体里识别「200 但失败」。
func extractProviderError(raw []byte) string {
	s := strings.TrimSpace(string(raw))
	if s == "" || s[0] != '{' {
		return ""
	}
	var m map[string]any
	if json.Unmarshal(raw, &m) != nil {
		return ""
	}
	if ok, exists := m["ok"]; exists {
		if b, isBool := ok.(bool); isBool && !b {
			if d, _ := m["description"].(string); d != "" {
				return d
			}
			return "服务返回 ok=false"
		}
	}
	if code, exists := m["errcode"]; exists {
		if n, isNum := code.(float64); isNum && n != 0 {
			if d, _ := m["errmsg"].(string); d != "" {
				return fmt.Sprintf("errcode=%d %s", int(n), d)
			}
			return fmt.Sprintf("errcode=%d", int(n))
		}
	}
	if code, exists := m["code"]; exists {
		if n, isNum := code.(float64); isNum && n != 0 {
			if d, _ := m["msg"].(string); d != "" {
				return fmt.Sprintf("code=%d %s", int(n), d)
			}
		}
	}
	return ""
}

// ---------- 签名 ----------

func dingtalkSign(ch store.Channel, timestamp string) string {
	secret, _ := ch.Config["secret"].(string)
	if secret == "" {
		return ""
	}
	s := hmac.New(sha256.New, []byte(secret))
	s.Write([]byte(timestamp + "\n" + secret))
	return base64.StdEncoding.EncodeToString(s.Sum(nil))
}

func feishuSign(ch store.Channel, timestamp string) string {
	secret, _ := ch.Config["secret"].(string)
	if secret == "" {
		return ""
	}
	s := hmac.New(sha256.New, []byte(timestamp+"\n"+secret))
	s.Write([]byte{})
	return base64.StdEncoding.EncodeToString(s.Sum(nil))
}

// ---------- SMTP ----------

func (m *Manager) sendSMTPWithRetry(ch store.Channel, title, text string) (bool, string) {
	var lastErr string
	for attempt := 0; attempt < 3; attempt++ {
		if err := m.sendSMTP(ch, title, text); err == nil {
			return true, ""
		} else {
			lastErr = err.Error()
		}
		time.Sleep(time.Duration(attempt+1) * 700 * time.Millisecond)
	}
	return false, lastErr
}

func (m *Manager) sendSMTP(ch store.Channel, title, text string) error {
	get := func(k string) string { s, _ := ch.Config[k].(string); return strings.TrimSpace(s) }
	host := get("host")
	port := get("port")
	if host == "" || port == "" {
		return fmt.Errorf("SMTP 服务器或端口未配置")
	}
	username := get("username")
	password := get("password")
	from := get("from")
	if from == "" {
		from = username
	}
	toRaw := get("to")
	if toRaw == "" {
		return fmt.Errorf("收件人未配置")
	}
	to := []string{}
	for _, p := range strings.Split(toRaw, ",") {
		if p = strings.TrimSpace(p); p != "" {
			to = append(to, p)
		}
	}
	tlsMode := strings.ToLower(get("tls"))
	if tlsMode == "" {
		if port == "465" {
			tlsMode = "ssl"
		} else {
			tlsMode = "starttls"
		}
	}

	subject := "Dockhelm 通知 · " + sanitizeHeader(title)
	body := "From: " + from + "\r\nTo: " + strings.Join(to, ",") + "\r\n" +
		"Subject: =?UTF-8?B?" + base64.StdEncoding.EncodeToString([]byte(subject)) + "?=\r\n" +
		"MIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n" +
		"Date: " + time.Now().Format(time.RFC1123Z) + "\r\n\r\n" + text + "\r\n"

	addr := host + ":" + port
	var auth smtp.Auth
	if username != "" {
		auth = smtp.PlainAuth("", username, password, host)
	}

	if tlsMode == "ssl" {
		conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: host, MinVersion: tls.VersionTLS12})
		if err != nil {
			return err
		}
		defer conn.Close()
		c, err := smtp.NewClient(conn, host)
		if err != nil {
			return err
		}
		defer c.Close()
		if auth != nil {
			if err := c.Auth(auth); err != nil {
				return err
			}
		}
		if err := c.Mail(from); err != nil {
			return err
		}
		for _, rcpt := range to {
			if err := c.Rcpt(rcpt); err != nil {
				return err
			}
		}
		w, err := c.Data()
		if err != nil {
			return err
		}
		if _, err := w.Write([]byte(body)); err != nil {
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}
		return c.Quit()
	}

	if tlsMode == "none" {
		c, err := smtp.Dial(addr)
		if err != nil {
			return err
		}
		defer c.Close()
		if auth != nil {
			if err := c.Auth(auth); err != nil {
				return err
			}
		}
		if err := c.Mail(from); err != nil {
			return err
		}
		for _, rcpt := range to {
			if err := c.Rcpt(rcpt); err != nil {
				return err
			}
		}
		w, err := c.Data()
		if err != nil {
			return err
		}
		if _, err := w.Write([]byte(body)); err != nil {
			return err
		}
		if err := w.Close(); err != nil {
			return err
		}
		return c.Quit()
	}

	// STARTTLS
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	if err := c.StartTLS(&tls.Config{ServerName: host, MinVersion: tls.VersionTLS12}); err != nil {
		return err
	}
	if auth != nil {
		if err := c.Auth(auth); err != nil {
			return err
		}
	}
	if err := c.Mail(from); err != nil {
		return err
	}
	for _, rcpt := range to {
		if err := c.Rcpt(rcpt); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write([]byte(body)); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return c.Quit()
}

func sanitizeHeader(s string) string {
	s = strings.ReplaceAll(s, "\r", " ")
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}

func startOfToday() time.Time {
	n := time.Now()
	return time.Date(n.Year(), n.Month(), n.Day(), 0, 0, 0, 0, n.Location())
}

// Logf 供内部记录（失败绝不 fatal）。
var Logf = log.Printf
