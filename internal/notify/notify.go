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
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aaron2024s/dockhelm/internal/bus"
	"github.com/aaron2024s/dockhelm/internal/dockerx"
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
	m := &Manager{
		st:   st,
		bus:  b,
		host: HostLabel(),
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

// maxDigestLinesPerEvent 每组事件最多逐条列几个；超出折叠成一行。
// 静默一整夜遇上崩溃循环，缓冲里能攒几十条 —— 汇总本身不能变成刷屏。
const maxDigestLinesPerEvent = 4

// maxDigestFailures 聚合块里最多逐条列几台失败容器。
const maxDigestFailures = 6

// 更新类事件的聚合变量名 —— 必须与 updater.Agg.FillVars 写入的键一致。
// 不直接 import updater：notify 是被它依赖的下游（靠 Notifier 接口解耦），
// 反向引用会把这条线绕回来。键名改动两边要一起改。
const (
	varAggUpdated   = "aggUpdated"
	varAggFailed    = "aggFailed"
	varAggSkipped   = "aggSkipped"
	varAggReclaimed = "aggReclaimed"
	varAggNames     = "aggNames"
	varAggFailures  = "aggFailures"
)

// isUpdateEvent 这类事件带结构化统计，可以跨多条求和。
func isUpdateEvent(ev string) bool {
	return ev == "batch_update_done" || ev == "batch_update_failed"
}

// updateGroupLines 把同一组更新事件压成「组头 + 合计统计 + 名单」。
//
// 2026-10-27 用户反馈：夜里 5 次自动更新，汇总里是 5 行一模一样的长句
// （「· rsshub — 共 1 个容器：更新 1 个、失败 0 个、已是最新/跳过 0 个，清理旧镜像回收 1.3 GB」），
// 「太长了吧，其实可以像单条消息那样总结起来」。求和后就是一条统计句 + 一份名单，
// 和单次更新的文案同构。
//
// 只有一条时也走这里：至少能把失败容器与原因列出来（逐条模式只取 message 首行）。
//
// ok=false 表示这批事件没带结构化统计（老版本攒下的、或别的生产者），
// 调用方回落到逐条列出。
func updateGroupLines(ev string, list []buffered) ([]string, bool) {
	var updated, failed, skipped int
	var reclaimed int64
	names, failures := []string{}, []string{}

	for _, it := range list {
		u, err := strconv.Atoi(it.Vars[varAggUpdated])
		if err != nil {
			return nil, false
		}
		f, err := strconv.Atoi(it.Vars[varAggFailed])
		if err != nil {
			return nil, false
		}
		sk, err := strconv.Atoi(it.Vars[varAggSkipped])
		if err != nil {
			return nil, false
		}
		rc, _ := strconv.ParseInt(it.Vars[varAggReclaimed], 10, 64)
		updated += u
		failed += f
		skipped += sk
		reclaimed += rc
		names = append(names, splitJoin(it.Vars[varAggNames], "、")...)
		failures = append(failures, splitJoin(it.Vars[varAggFailures], "\n")...)
	}

	lines := []string{fmt.Sprintf("%s（%d）", Label(ev), len(list))}
	stats := fmt.Sprintf("共 %d 次更新：更新 %d 个、失败 %d 个、已是最新/跳过 %d 个",
		len(list), updated, failed, skipped)
	if reclaimed > 0 {
		stats += "，清理旧镜像回收 " + dockerx.HumanBytes(reclaimed)
	}
	// 只有一次更新时，组头看不出是哪一轮 —— 把时刻补在统计句尾部。
	if len(list) == 1 {
		stats += "（" + list[0].At.Format("15:04") + "）"
	}
	lines = append(lines, stats)

	if names = uniqueSorted(names); len(names) > 0 {
		shown := names
		if len(shown) > 12 {
			shown = shown[:12]
		}
		line := "更新的容器：" + strings.Join(shown, "、")
		if rest := len(names) - len(shown); rest > 0 {
			line += fmt.Sprintf(" 等 %d 个", len(names))
		}
		lines = append(lines, line)
	}
	if failures = uniqueSorted(failures); len(failures) > 0 {
		lines = append(lines, "失败的容器：")
		shown := failures
		if len(shown) > maxDigestFailures {
			shown = shown[:maxDigestFailures]
		}
		for _, f := range shown {
			lines = append(lines, "· "+f)
		}
		if rest := len(failures) - len(shown); rest > 0 {
			lines = append(lines, fmt.Sprintf("· 另有 %d 条，详见面板「推送历史」", rest))
		}
	}
	return lines, true
}

// splitJoin 按分隔符切开并丢掉空串。
func splitJoin(s, sep string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	out := []string{}
	for _, p := range strings.Split(s, sep) {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// uniqueSorted 去重 + 排序（汇总里的名单必须与遍历顺序无关）。
func uniqueSorted(in []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	sort.Strings(out)
	return out
}

// digestMessage 把静默期攒下的事件合成一条汇总。返回 (标题, 正文)。
//
// 2026-10-26 用户反馈：原来只有「· 容器意外退出 ×2」这种计数，晚上挂了两台容器，
// 早上完全看不出是哪两台、什么原因；而且同一批事件每次顺序都不一样
// （`for ev, n := range counts` 走 Go map，遍历顺序随机），
// 标题也没有普通消息那个语气图标，混在通知列表里不像同一个来源。
// 现在：按事件分组 → 组序固定（事件目录顺序）→ 组内逐条点名容器与原因，
// 并补齐普通消息那套「主机 / 时段 / 结果 / 面板」头部。
func (m *Manager) digestMessage(items []buffered, now time.Time) (string, string) {
	if len(items) == 0 {
		return "", ""
	}

	groups := map[string][]buffered{}
	order := make([]string, 0, len(items))
	for _, it := range items {
		if _, ok := groups[it.Event]; !ok {
			order = append(order, it.Event)
		}
		groups[it.Event] = append(groups[it.Event], it)
	}
	// 组序：按事件目录顺序（与「事件订阅」页的排列一致），目录外的按 key 兜底。
	rank := make(map[string]int, len(EventCatalog))
	for i, d := range EventCatalog {
		rank[d.Event] = i
	}
	rankOf := func(ev string) int {
		if r, ok := rank[ev]; ok {
			return r
		}
		return len(EventCatalog) + 1
	}
	sort.SliceStable(order, func(i, j int) bool {
		ri, rj := rankOf(order[i]), rankOf(order[j])
		if ri != rj {
			return ri < rj
		}
		return order[i] < order[j]
	})

	// 标题图标取「最坏语气」：汇总里只要有一条失败，就不该挂个提示图标。
	toneRank := map[string]int{ToneSuccess: 0, ToneInfo: 1, ToneWarning: 2, ToneFailed: 3}
	worst := ToneSuccess
	toneCount := map[string]int{}
	for _, it := range items {
		t := ToneOf(it.Event)
		toneCount[t]++
		if toneRank[t] > toneRank[worst] {
			worst = t
		}
	}

	title := fmt.Sprintf("%s Dockhelm 静默时段汇总 · %d 条", toneIcon[worst], len(items))

	// 时段：首条与末条的实际时间（比写死「23:00–07:00」更准 —— 静默期可能是自定义的，
	// 而且最后一条之后到汇总之间还隔着一段轮询间隔）。
	first, last := items[0].At, items[0].At
	for _, it := range items[1:] {
		if it.At.Before(first) {
			first = it.At
		}
		if it.At.After(last) {
			last = it.At
		}
	}

	lines := []string{title, "主机：" + m.host}
	if first.Format("2006-01-02") == last.Format("2006-01-02") {
		lines = append(lines, fmt.Sprintf("时段：%s → %s（%s）",
			first.Format("01-02 15:04"), last.Format("15:04"), humanSpan(last.Sub(first))))
	} else {
		lines = append(lines, fmt.Sprintf("时段：%s → %s（%s）",
			first.Format("01-02 15:04"), last.Format("01-02 15:04"), humanSpan(last.Sub(first))))
	}
	if parts := toneCountParts(toneCount); parts != "" {
		lines = append(lines, "结果："+parts)
	}

	for _, ev := range order {
		list := groups[ev]
		// 更新事件：求和成一条统计句，别刷同样的长句（1 条也走，这样失败原因能列出来）
		if isUpdateEvent(ev) {
			if block, ok := updateGroupLines(ev, list); ok {
				lines = append(lines, block...)
				continue
			}
		}
		lines = append(lines, fmt.Sprintf("%s（%d）", Label(ev), len(list)))
		shown := list
		if len(shown) > maxDigestLinesPerEvent {
			shown = shown[:maxDigestLinesPerEvent]
		}
		for _, it := range shown {
			lines = append(lines, digestItemLine(it))
		}
		if rest := len(list) - len(shown); rest > 0 {
			lines = append(lines, fmt.Sprintf("· 另有 %d 条，详见面板「推送历史」", rest))
		}
	}

	if u := m.st.GetSetting(KeyPanelURL, ""); u != "" {
		lines = append(lines, "面板："+u)
	}
	return title, truncateNotifyText(strings.Join(lines, "\n"))
}

// digestItemLine 汇总里的一行明细：「· 容器名 — 原因（03:12）」。
// 缺项就少一段，绝不输出「· — （03:12）」这种占位。
func digestItemLine(it buffered) string {
	container := strings.TrimSpace(it.Vars["container"])
	detail := strings.TrimSpace(it.Vars["message"])
	if i := strings.IndexByte(detail, '\n'); i >= 0 {
		detail = strings.TrimSpace(detail[:i])
	}
	if detail == "" {
		detail = strings.TrimSpace(it.Vars["result"])
	}
	hm := it.At.Format("15:04")
	switch {
	case container != "" && detail != "":
		return fmt.Sprintf("· %s — %s（%s）", container, detail, hm)
	case container != "":
		return fmt.Sprintf("· %s（%s）", container, hm)
	case detail != "":
		return fmt.Sprintf("· %s（%s）", detail, hm)
	default:
		return "· " + hm
	}
}

// toneCountParts 把语气计数写成「成功 3 · 失败 1」，零项不出现。
func toneCountParts(counts map[string]int) string {
	names := []struct {
		tone string
		word string
	}{
		{ToneFailed, "失败"},
		{ToneWarning, "警告"},
		{ToneInfo, "提示"},
		{ToneSuccess, "成功"},
	}
	out := []string{}
	for _, n := range names {
		if counts[n.tone] > 0 {
			out = append(out, fmt.Sprintf("%s %d", n.word, counts[n.tone]))
		}
	}
	return strings.Join(out, " · ")
}

// humanSpan 把时长写成「7 小时 46 分」，不足 1 分钟按 1 分钟算。
func humanSpan(d time.Duration) string {
	if d < time.Minute {
		return "不足 1 分钟"
	}
	d = d.Round(time.Minute)
	h, m := int(d.Hours()), int(d.Minutes())%60
	switch {
	case h == 0:
		return fmt.Sprintf("%d 分钟", m)
	case m == 0:
		return fmt.Sprintf("%d 小时", h)
	default:
		return fmt.Sprintf("%d 小时 %d 分", h, m)
	}
}

func (m *Manager) sendDigest(items []buffered) {
	title, text := m.digestMessage(items, time.Now())
	if title == "" {
		return
	}
	m.deliver("digest", LevelNormal, title, text, nil, true)
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
	// bus 可能为 nil（单测里 New(st, nil)）——Emit 跑在独立 goroutine 里，
	// 这里 panic 会直接把整个进程带走，所以宁可少推一条事件。
	if m.bus != nil {
		m.bus.Publish("notify", "delivered", "info", map[string]any{
			"event": event, "level": level, "title": title, "text": text,
			"channels": len(channels), "sent": sent,
		})
	}
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
//
// 行序有意把「主机 / 时间」钉在头部（2026-10-11 用户反馈：主机原来挂在最后一行，
// 消息一长被转发链路腰斩时就永远看不到是哪台机器）。详情/明细这类大块头放后面，
// 超长截断砍的只会是它们的尾部。
func (m *Manager) compose(event, titleOverride, textOverride string, vars map[string]string) (string, string) {
	if titleOverride != "" && textOverride != "" {
		return titleOverride, textOverride
	}
	label := Label(event)
	// 图标按「语气」（成功/失败/警告/提示）走，不再按打扰强度分铃铛/警灯：
	// 「更新完成」配铃铛没有语义，「容器退出」配警灯又过重（2026-10-12 用户反馈）。
	icon := Icon(event)
	head := label
	if c := vars["container"]; c != "" {
		head = fmt.Sprintf("%s · %s", label, c)
	}
	title := fmt.Sprintf("%s Dockhelm %s", icon, head)

	lines := []string{title, "主机：" + m.host, "时间：" + time.Now().Format("2006-01-02 15:04:05")}
	// 「结果」与事件名相同就不另起一行 —— 标题里已经写着「容器意外退出 / Docker 连接
	// 断开」，再来一句「结果：容器意外退出」纯属重复（2026-10-12 用户反馈）。
	// 只有结果**补充了新信息**（如「全量备份有失败」「配置已还原」）才输出。
	if r := vars["result"]; r != "" && r != label {
		lines = append(lines, "结果："+r)
	}
	if vars["message"] != "" {
		lines = append(lines, "详情："+vars["message"])
	}
	if vars["image"] != "" {
		lines = append(lines, "镜像："+vars["image"])
	}
	if u := m.st.GetSetting(KeyPanelURL, ""); u != "" {
		lines = append(lines, "面板："+u)
	}
	return title, truncateNotifyText(strings.Join(lines, "\n"))
}

// maxNotifyTextBytes 推送正文的安全上限（UTF-8 字节）。
// 各渠道上限不一（企业微信 4096、Telegram 4096，一些转发链路更小），超长消息
// 会被渠道或转发链静默截断 —— 用户实测批量失败明细一长，推送被腰斩在半句话上。
// 在源头按行截断：主机/时间/结果在头部不会被切掉，砍的只会是明细尾部。
const maxNotifyTextBytes = 1800

func truncateNotifyText(text string) string {
	if len(text) <= maxNotifyTextBytes {
		return text
	}
	used := 0
	for i, ln := range strings.Split(text, "\n") {
		if used+len(ln) > maxNotifyTextBytes {
			kept := strings.TrimSpace(strings.Join(strings.Split(text, "\n")[:i], "\n"))
			if kept == "" {
				// 第一行就装不下：硬切字节也要保住头部信息。
				return string([]rune(text)[:maxNotifyTextBytes/4]) + "……（内容过长已截断）"
			}
			return kept + "\n……（内容过长已截断）"
		}
		used += len(ln) + 1
	}
	return text
}

// ---------- 主机标识 ----------

// HostLabel 通知里显示的「主机」。
//
// 容器里 os.Hostname() 是随机的容器 ID（如 eba791308b83），用户根本认不出
// 这是哪台机器 —— 用户原话：「主机名改为机器的具体型号，比如 dx4600」。
// 读取顺序：
//  1. DOCKHELM_HOST_LABEL 环境变量（用户显式指定，最高优先）
//  2. DMI product_name / board_name（x86 NAS 与主板；容器里 /sys 就是宿主的，
//     所以容器内通常照样读得到）
//  3. 设备树 model（ARM 设备，挂载了才能读到）
//  4. 回落 os.Hostname()
func HostLabel() string {
	if v := strings.TrimSpace(os.Getenv("DOCKHELM_HOST_LABEL")); v != "" {
		return v
	}
	for _, p := range []string{
		"/sys/class/dmi/id/product_name",
		"/sys/class/dmi/id/board_name",
		"/proc/device-tree/model",
		"/sys/firmware/devicetree/base/model",
	} {
		b, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		s := strings.TrimSpace(strings.TrimRight(string(b), "\x00\n\r\t "))
		if s == "" || genericBoardName(s) {
			continue
		}
		return s
	}
	if h, err := os.Hostname(); err == nil && h != "" {
		return h
	}
	return "unknown"
}

// genericBoardName 过滤 DMI 里的占位值 —— 一堆 x86 主板出厂就写着
// 「System Product Name」，拿它当主机标识等于没换。
func genericBoardName(s string) bool {
	switch strings.ToLower(s) {
	case "system product name", "to be filled by o.e.m.", "default string",
		"none", "empty", "not specified", "not defined", "undefined",
		"unknown", "n/a", "type1", "product name", "generic", "oem":
		return true
	}
	return false
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
	values := buildValues(event, vars, m.host, m.st.GetSetting(KeyPanelURL, ""), title, text)
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
