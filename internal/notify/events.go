package notify

import "sort"

// EventDef 一个可订阅事件的定义。
type EventDef struct {
	Event       string `json:"event"`
	Label       string `json:"label"`
	Group       string `json:"group"`
	Level       string `json:"level"`   // urgent | normal
	Default     bool   `json:"default"` // 出厂默认是否订阅
	Description string `json:"description"`
}

// LevelUrgent / LevelNormal 事件级别。紧急事件可设置「静默时段照发」。
const (
	LevelUrgent = "urgent"
	LevelNormal = "normal"
)

// EventCatalog 全部事件。加新事件只需要在这里加一行。
var EventCatalog = []EventDef{
	// —— 更新 ——
	{Event: "update_available", Label: "检测到有更新", Group: "更新", Level: LevelNormal, Default: false,
		Description: "例行检测发现有新镜像时提醒（批量巡检可能较频繁，默认关闭）"},
	{Event: "update_success", Label: "容器更新成功", Group: "更新", Level: LevelNormal, Default: false,
		Description: "单个容器成功更新到新镜像"},
	{Event: "update_failed", Label: "容器更新失败", Group: "更新", Level: LevelUrgent, Default: true,
		Description: "更新失败，正文里会带上「是否已回滚」的结论"},
	{Event: "batch_update_done", Label: "批量更新完成", Group: "更新", Level: LevelNormal, Default: true,
		Description: "一次批量更新结束后合并成一条汇总，而不是每台一条"},

	// —— 容器 ——
	{Event: "container_died", Label: "容器意外退出", Group: "容器", Level: LevelUrgent, Default: true,
		Description: "容器以非 0 退出码结束（手工停止不算）"},
	{Event: "container_crashloop", Label: "容器崩溃循环", Group: "容器", Level: LevelUrgent, Default: true,
		Description: "短时间内反复退出（有去重窗口，不会刷屏）"},
	{Event: "health_failed", Label: "健康检查失败", Group: "容器", Level: LevelUrgent, Default: true,
		Description: "健康检查转为 unhealthy"},

	// —— 计划任务 ——
	{Event: "schedule_success", Label: "计划任务执行成功", Group: "计划任务", Level: LevelNormal, Default: false,
		Description: "定时启停/重启/更新任务执行成功"},
	{Event: "schedule_failed", Label: "计划任务执行失败", Group: "计划任务", Level: LevelUrgent, Default: true,
		Description: "定时任务执行出错"},

	// —— 备份 ——
	{Event: "backup_success", Label: "备份完成", Group: "备份", Level: LevelNormal, Default: false,
		Description: "手动或定时的配置快照备份完成"},
	{Event: "backup_failed", Label: "备份失败", Group: "备份", Level: LevelUrgent, Default: true,
		Description: "备份写入失败"},
	{Event: "restore_done", Label: "还原完成", Group: "备份", Level: LevelUrgent, Default: true,
		Description: "执行了配置还原（高风险动作，建议保留）"},

	// —— 系统 ——
	{Event: "docker_disconnected", Label: "Docker 连接断开", Group: "系统", Level: LevelUrgent, Default: true,
		Description: "无法连接 Docker 守护进程"},
	{Event: "docker_reconnected", Label: "Docker 连接恢复", Group: "系统", Level: LevelUrgent, Default: true,
		Description: "Docker 守护进程恢复可用"},

	// —— 登录与账户 ——
	{Event: "login_success", Label: "登录成功", Group: "登录与账户", Level: LevelNormal, Default: false,
		Description: "有人成功登录面板（若只想在异常时被打扰可保持关闭）"},
	{Event: "login_failed", Label: "登录失败", Group: "登录与账户", Level: LevelUrgent, Default: true,
		Description: "密码错误，正文包含来源 IP"},
	{Event: "login_locked", Label: "登录连续失败告警", Group: "登录与账户", Level: LevelUrgent, Default: true,
		Description: "同一 IP 连续失败达到阈值并触发临时锁定"},
	{Event: "password_changed", Label: "登录密码已修改", Group: "登录与账户", Level: LevelUrgent, Default: true,
		Description: "密码变更，全部既有会话已失效"},

	// —— 其它 ——
	{Event: "test", Label: "测试消息", Group: "其它", Level: LevelNormal, Default: true,
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
