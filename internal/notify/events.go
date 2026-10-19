package notify

import "sort"

// EventDef 一个可订阅事件的定义。
type EventDef struct {
	Event       string `json:"event"`
	Label       string `json:"label"`
	Group       string `json:"group"`
	Level       string `json:"level"`   // urgent | normal
	Tone        string `json:"tone"`    // success | failed | warning | info
	Default     bool   `json:"default"` // 出厂默认是否订阅
	Description string `json:"description"`
}

// LevelUrgent / LevelNormal 事件级别。紧急事件可设置「静默时段照发」。
const (
	LevelUrgent = "urgent"
	LevelNormal = "normal"
)

// Tone 决定推送标题的图标（以及 ntfy 通知栏图标），与 Level 是两个正交的轴：
//
//	Tone  = 语气（成功 / 失败 / 警告 / 提示）—— 给眼睛看的，一眼判断结果好坏
//	Level = 打扰强度（urgent / normal）—— 决定静默时段是否放行、优先级、提示音
//
// 两者不能互相推导：health_failed 是「警告」语气却属于 urgent（要立刻打扰），
// docker_reconnected 是「成功」语气也属于 urgent（链路恢复同样重要）。
// 也别把两者混用 —— 想让某条消息「更响」改 Level，想让标题换图标改 Tone。
const (
	ToneSuccess = "success"
	ToneFailed  = "failed"
	ToneWarning = "warning"
	ToneInfo    = "info"
)

// toneIcon 每种语气对应的标题图标。
var toneIcon = map[string]string{
	ToneSuccess: "✅",
	ToneFailed:  "❌",
	ToneWarning: "⚠️",
	ToneInfo:    "ℹ️",
}

// toneTag ntfy 的 tags 参数（emoji shortcode，决定通知栏左侧的小图标）。
// 必须与 toneIcon 一一对应，否则会出现「标题是 ✅、通知栏是警灯」的割裂。
var toneTag = map[string]string{
	ToneSuccess: "white_check_mark",
	ToneFailed:  "x",
	ToneWarning: "warning",
	ToneInfo:    "information_source",
}

// ToneOf 返回事件的语气（未知事件 / 未标注按提示处理）。
func ToneOf(event string) string {
	if d, ok := EventMap[event]; ok && d.Tone != "" {
		return d.Tone
	}
	return ToneInfo
}

// Icon 返回事件在推送标题前的图标。
func Icon(event string) string {
	if ic, ok := toneIcon[ToneOf(event)]; ok {
		return ic
	}
	return toneIcon[ToneInfo]
}

// EventCatalog 全部事件。加新事件只需要在这里加一行。
//
// Tone 逐个标注（不能靠 Level 推）：标的是「这条消息读起来是好消息还是坏消息」。
var EventCatalog = []EventDef{
	// —— 更新 ——
	{Event: "update_available", Label: "检测到有更新", Group: "更新", Level: LevelNormal, Tone: ToneInfo, Default: false,
		Description: "例行检测发现有新镜像时提醒（批量巡检可能较频繁，默认关闭）"},
	// update_success / update_failed 已于 2026-10-12 删除：它们的唯一生产者是
	// 计划任务里的单容器更新，而那批信息现在由 schedule_success / schedule_failed
	// 一条承载（正文点名到容器），两个事件并存只会让一次计划任务更新失败发两条
	// 紧急推送。历史记录里若还有这两个 key，前端 legacyEventLabels 会兜底成中文。
	//
	// ⚠ 下列两个事件承载**所有**更新轮次（手动 / 自动、1 台 / 多台，见 updater.UpdateMany）：
	// Label 不能带「批量」二字，否则 1 台时标题会出现「批量更新完成 · 1 个容器」这种组合。
	{Event: "batch_update_done", Label: "更新完成", Group: "更新", Level: LevelNormal, Tone: ToneSuccess, Default: true,
		Description: "一次更新的结果汇总：共几个容器、更新了几个（点名）、失败几个（点名）——一次更新只有这一条，1 台也一样"},
	{Event: "batch_update_failed", Label: "更新有失败", Group: "更新", Level: LevelUrgent, Tone: ToneFailed, Default: true,
		Description: "一次更新里有容器失败时的汇总（整次只有这一条），正文列出失败的容器与原因"},

	// —— 容器 ——
	{Event: "container_died", Label: "容器意外退出", Group: "容器", Level: LevelUrgent, Tone: ToneFailed, Default: true,
		Description: "容器以非 0 退出码结束（手工停止不算）"},
	{Event: "container_crashloop", Label: "容器崩溃循环", Group: "容器", Level: LevelUrgent, Tone: ToneFailed, Default: true,
		Description: "短时间内反复退出（有去重窗口，不会刷屏）"},
	{Event: "health_failed", Label: "健康检查失败", Group: "容器", Level: LevelUrgent, Tone: ToneWarning, Default: true,
		Description: "健康检查转为 unhealthy"},

	// —— 计划任务 ——
	{Event: "schedule_success", Label: "计划任务执行成功", Group: "计划任务", Level: LevelNormal, Tone: ToneSuccess, Default: false,
		Description: "定时启停/重启/更新/备份任务成功时提醒，正文点名到具体容器"},
	{Event: "schedule_failed", Label: "计划任务执行失败", Group: "计划任务", Level: LevelUrgent, Tone: ToneFailed, Default: true,
		Description: "定时任务执行出错（含计划任务里的容器更新失败），正文点名失败容器；这是计划任务唯一的回执"},

	// —— 顺序启动 ——
	// 一次开机/一次手动执行只有这两条之一，正文点名失败容器与原因。
	// 成功那条默认关闭：每次开机都推一条「启动完成」属于噪音，
	// 但想让手机确认「机器已经起好了」的人可以打开。
	{Event: "boot_done", Label: "顺序启动完成", Group: "顺序启动", Level: LevelNormal, Tone: ToneSuccess, Default: false,
		Description: "按顺序启动跑完且没有失败时提醒：启动了几个、跳过几个"},
	{Event: "boot_failed", Label: "顺序启动失败", Group: "顺序启动", Level: LevelUrgent, Tone: ToneFailed, Default: true,
		Description: "有容器没能启动（含超时、启动后立即退出、容器已不存在），正文点名原因"},

	// —— 备份 ——
	{Event: "backup_success", Label: "备份完成", Group: "备份", Level: LevelNormal, Tone: ToneSuccess, Default: false,
		Description: "手动或定时的配置快照备份完成"},
	{Event: "backup_failed", Label: "备份失败", Group: "备份", Level: LevelUrgent, Tone: ToneFailed, Default: true,
		Description: "备份写入失败"},
	// 还原成功但对系统是高风险的既有用户，语气用「警告」而非「成功」：
	// 图标要提醒人「刚刚变了东西」，不是庆祝。
	{Event: "restore_done", Label: "还原完成", Group: "备份", Level: LevelUrgent, Tone: ToneWarning, Default: true,
		Description: "执行了配置还原（高风险动作，建议保留）"},

	// —— 系统 ——
	{Event: "docker_disconnected", Label: "Docker 连接断开", Group: "系统", Level: LevelUrgent, Tone: ToneFailed, Default: true,
		Description: "无法连接 Docker 守护进程"},
	{Event: "docker_reconnected", Label: "Docker 连接恢复", Group: "系统", Level: LevelUrgent, Tone: ToneSuccess, Default: true,
		Description: "Docker 守护进程恢复可用"},

	// —— 登录与账户 ——
	{Event: "login_success", Label: "登录成功", Group: "登录与账户", Level: LevelNormal, Tone: ToneInfo, Default: false,
		Description: "有人成功登录面板（若只想在异常时被打扰可保持关闭）"},
	// 单次密码错误是「有人在试密码」，还不算系统故障 ⇒ 警告；
	// 触发锁定则升级成 failed，与登录失败形成强弱两档。
	{Event: "login_failed", Label: "登录失败", Group: "登录与账户", Level: LevelUrgent, Tone: ToneWarning, Default: true,
		Description: "密码错误，正文包含来源 IP"},
	{Event: "login_locked", Label: "登录连续失败告警", Group: "登录与账户", Level: LevelUrgent, Tone: ToneFailed, Default: true,
		Description: "同一 IP 连续失败达到阈值并触发临时锁定"},
	{Event: "password_changed", Label: "登录密码已修改", Group: "登录与账户", Level: LevelUrgent, Tone: ToneWarning, Default: true,
		Description: "密码变更，全部既有会话已失效"},

	// —— 其它 ——
	{Event: "test", Label: "测试消息", Group: "其它", Level: LevelNormal, Tone: ToneInfo, Default: true,
		Description: "点击渠道「测试」按钮时发出的消息"},
}

// EventMap 便于按 key 查找。
var EventMap = func() map[string]EventDef {
	m := make(map[string]EventDef, len(EventCatalog))
	for _, e := range EventCatalog {
		m[e.Event] = e
	}
	return m
}()

// Label 返回事件的中文名（未知事件回落成 key 本身）。
func Label(event string) string {
	if d, ok := EventMap[event]; ok {
		return d.Label
	}
	return event
}

// LevelOf 返回事件级别（未知事件按普通处理）。
func LevelOf(event string) string {
	if d, ok := EventMap[event]; ok {
		return d.Level
	}
	return LevelNormal
}

// Groups 返回分组顺序（保持目录中的出现顺序）。
func Groups() []string {
	seen := map[string]bool{}
	out := []string{}
	for _, e := range EventCatalog {
		if !seen[e.Group] {
			seen[e.Group] = true
			out = append(out, e.Group)
		}
	}
	return out
}

// SortedEvents 按分组、再按目录顺序返回。
func SortedEvents() []EventDef {
	out := make([]EventDef, len(EventCatalog))
	copy(out, EventCatalog)
	order := map[string]int{}
	for i, g := range Groups() {
		order[g] = i
	}
	sort.SliceStable(out, func(i, j int) bool { return order[out[i].Group] < order[out[j].Group] })
	return out
}
