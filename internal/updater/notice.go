// 一轮更新的结果通知（手动更新、定时自动更新共用）。
// 无论多少台、成败几何，**只发一条推送**。
//
// 文案**不区分单台与多台**：一台和十台用同一套结构 ——
//
//	共 4 个容器：更新 2 个、失败 1 个、已是最新/跳过 1 个
//	更新的容器：jellyfin、immich-server
//	失败的容器：
//	· shell-town：拉取失败
//
// 2026-10-10 之前的行为是「每台失败各一条 + 批量汇总一条 + 自动更新再汇总一条」，
// 一轮 4 台失败 2 台的自动更新能炸出 4 条推送，其中两条汇总内容几乎一样 ——
// 用户原话：「1 和 3 看着就像有点重复，这些内容不是完全可以合并为一条吗」。
// 于是合并成一条，但留下了「按台数换文案」的分支（1 台走一句话、多台走统计句）。
// 2026-10-12 用户再次要求抹平：「不用刻意强调是批量更新……就写共几个容器更新、
// 更新了几个、更新的容器是哪些、失败了几个、失败的是哪些，单个和多个完全统一」。
// 现在没有 only 分支了：任何台数都是「统计行 + 更新名单 + 失败名单」。
//
// 事件仍是 batch_update_done（普通）/ batch_update_failed（紧急，静默时段照发）。

package updater

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aaron2024s/dockhelm/internal/dockerx"
)

// batchStats 一轮更新的统计与明细，全部由 UpdateMany 从结果集里数出来。
type batchStats struct {
	// source 触发来源："auto" = 定时/启动自动更新，其余按手动表述（只影响结果短语）。
	source string
	total  int
	// updatedNames 本轮成功重建的容器名（按结果集顺序），正文点名用。
	updatedNames []string
	// updated 成功重建；skipped 已是最新或被跳过；failed 拉取/重建失败。
	updated   int
	skipped   int
	failed    int
	reclaimed int64
	// failures 每台失败容器一行：「容器名：原因」。
	failures []string
}

// BatchNotice 汇总通知的最终形态。
type BatchNotice struct {
	// Event 通知事件名：全部成功 batch_update_done，有失败 batch_update_failed。
	Event string
	// Result 结果短语，渲染成「结果：…」。
	Result string
	// Message 正文，渲染成「详情：…」，含统计与失败明细。
	Message string
	// Detail 运行记录用的完整明细（与 Message 相同的信息，换行更宽松）。
	Detail string
}

// maxFailureLines 汇总里最多逐行列出的失败条数；超出折叠成一行提示。
// NAS 上几十个容器同时拉挂的场景（比如断网后定时任务触发）不能刷出几十行推送。
const maxFailureLines = 8

// maxNameList 正文「更新的容器：」最多点名几个，超出折叠成「等 N 个」。
// 一轮更新几十台时，名单本身就能把推送撑爆。
const maxNameList = 12

// joinNameList 用「、」连接容器名，超过 limit 个则截断并补「等 N 个」。
func joinNameList(names []string, limit int) string {
	if len(names) <= limit {
		return strings.Join(names, "、")
	}
	return strings.Join(names[:limit], "、") + fmt.Sprintf(" 等 %d 个", len(names))
}

// shortFailureLine 把「容器名：拉取失败：pull failed: read tcp [240e:…]:443: …」
// 压成「容器名：拉取失败」—— 推送里刷一长串底层网络报错没人看得懂
// （2026-10-11 用户反馈：就写哪些更新失败，不用写详细信息）。
// 完整错误不丢：运行记录的 Detail 仍留全量，面板里随时可查。
// 规则：取第一个全角冒号前的动作短语（拉取失败 / 读取容器信息失败 / 自检失败…），
// 没有冒号结构的原文截到 40 字防意外超长。
func shortFailureLine(f string) string {
	i := strings.Index(f, "：")
	if i < 0 {
		return truncateRunes(f, 40)
	}
	name, msg := f[:i], f[i+len("："):]
	return name + "：" + shortPhase(msg)
}

// ShortPhase 从原始错误串里提取动作短语（「拉取失败：pull failed: …」→「拉取失败」）。
// 导出是给计划任务复用的：它也要把底层报错压成短语，不该再写第二套规则。
func ShortPhase(msg string) string { return shortPhase(msg) }

// shortPhase 从原始错误里提取动作短语。
func shortPhase(msg string) string {
	if i := strings.Index(msg, "："); i > 0 && i <= 30 {
		return msg[:i]
	}
	return truncateRunes(msg, 40)
}

func truncateRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}

// buildBatchNotice 把一轮更新的结果合成唯一的一条通知。
// 一台和多台走**同一套**文案，不再有台数分支。
func buildBatchNotice(s batchStats) BatchNotice {
	auto := s.source == "auto"

	// 结果集来自并发更新，完成顺序不确定 —— 名单先排序，同一轮重跑文案才稳定
	// （否则「jellyfin、immich-server」和「immich-server、jellyfin」会随机出现）。
	updatedNames := append([]string(nil), s.updatedNames...)
	sort.Strings(updatedNames)
	failures := append([]string(nil), s.failures...)
	sort.Strings(failures)

	label := "更新完成"
	if auto {
		label = "自动更新完成"
	}
	event := "batch_update_done"
	if s.failed > 0 {
		event = "batch_update_failed"
		label = "更新有失败"
		if auto {
			label = "自动更新有失败"
		}
	}

	summary := fmt.Sprintf("共 %d 个容器：更新 %d 个、失败 %d 个、已是最新/跳过 %d 个",
		s.total, s.updated, s.failed, s.skipped)
	if s.reclaimed > 0 {
		summary += "，清理旧镜像回收 " + dockerx.HumanBytes(s.reclaimed)
	}

	lines := []string{summary}
	if len(updatedNames) > 0 {
		lines = append(lines, "更新的容器："+joinNameList(updatedNames, maxNameList))
	}
	if len(failures) > 0 {
		shown := failures
		if len(shown) > maxFailureLines {
			shown = shown[:maxFailureLines]
		}
		lines = append(lines, "失败的容器：")
		for _, f := range shown {
			lines = append(lines, "· "+shortFailureLine(f))
		}
		if rest := len(failures) - len(shown); rest > 0 {
			lines = append(lines, fmt.Sprintf("（其余 %d 台略，详见面板更新记录）", rest))
		}
	}

	// 运行记录用全量：推送会被转发链路截断（见 truncateNotifyText），面板里不必。
	detailLines := []string{summary}
	if len(updatedNames) > 0 {
		detailLines = append(detailLines, "更新的容器："+strings.Join(updatedNames, "、"))
	}
	if len(failures) > 0 {
		detailLines = append(detailLines, "失败的容器：")
		for _, f := range failures {
			detailLines = append(detailLines, "· "+f)
		}
	}

	return BatchNotice{
		Event:   event,
		Result:  label,
		Message: strings.Join(lines, "\n"),
		Detail:  strings.Join(detailLines, "\n"),
	}
}
