// 批量更新的汇总通知：一轮多容器批量（手动批量、定时自动更新共用）
// 无论多少台、成败几何，**只发一条推送**。
//
// 2026-10-10 之前的行为是「每台失败各一条 + 批量汇总一条 + 自动更新再汇总一条」，
// 一轮 4 台失败 2 台的自动更新能炸出 4 条推送，其中两条汇总内容几乎一样 ——
// 用户原话：「1 和 3 看着就像有点重复，这些内容不是完全可以合并为一条吗」。
// 现在统一为：全部成功 → batch_update_done（普通）；有失败 → batch_update_failed
// （紧急，静默时段照发），失败容器与原因逐行列在正文里。

package updater

import (
	"fmt"
	"strings"
)

// batchStats 一轮批量更新的统计与失败明细，全部由 UpdateMany 从结果集里数出来。
type batchStats struct {
	// source 触发来源："auto" = 定时/启动自动更新，其余按手动批量表述。
	source string
	total  int
	// updated 成功重建；skipped 已是最新或被跳过；failed 拉取/重建失败。
	updated  int
	skipped  int
	failed   int
	reused   int // 本轮复用了同批已拉取镜像的容器数
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

// buildBatchNotice 把一轮批量结果合成唯一的一条通知。
func buildBatchNotice(s batchStats) BatchNotice {
	label := "批量更新完成"
	if s.source == "auto" {
		label = "自动更新完成"
	}
	event := "batch_update_done"
	if s.failed > 0 {
		event = "batch_update_failed"
		label = "批量更新有失败"
		if s.source == "auto" {
			label = "自动更新有失败"
		}
	}

	summary := fmt.Sprintf("共 %d 个容器：更新 %d、已是最新/跳过 %d、失败 %d",
		s.total, s.updated, s.skipped, s.failed)
	if s.reclaimed > 0 {
		summary += "，清理旧镜像回收 " + humanBytes(s.reclaimed)
	}
	if s.reused > 0 {
		summary += fmt.Sprintf("，%d 个容器复用了本轮已拉取的镜像", s.reused)
	}

	lines := []string{summary}
	detail := summary
	if len(s.failures) > 0 {
		shown := s.failures
		if len(shown) > maxFailureLines {
			shown = shown[:maxFailureLines]
		}
		lines = append(lines, "失败明细：")
		for _, f := range shown {
			lines = append(lines, "· "+f)
		}
		if rest := len(s.failures) - len(shown); rest > 0 {
			lines = append(lines, fmt.Sprintf("（其余 %d 台略，详见面板更新记录）", rest))
		}
		detail = strings.Join(append([]string{summary}, s.failures...), "\n")
	}

	return BatchNotice{
		Event:   event,
		Result:  label,
		Message: strings.Join(lines, "\n"),
		Detail:  detail,
	}
}
