package api

import (
	"errors"
	"net/http"

	"github.com/aaron2024s/dockhelm/internal/boot"
)

// bootRow 清单里的一行 + 该容器的实时状态。
//
// managed 由**重启策略**实时算出（no / on-failure 为 true），前端据此分成两组：
// 归本功能管的可编辑，归 Docker 管的只读。分组不落库。
type bootRow struct {
	Name     string `json:"name"`
	State    string `json:"state"`
	Policy   string `json:"restartPolicy"`
	Managed  bool   `json:"managed"`
	Project  string `json:"project"`
	Auto     bool   `json:"auto"`
	DelaySec int    `json:"delaySec"`
}

// hGetBoot 清单 + 容器实时状态 + 上次执行 + 当前进度。
func (s *Server) hGetBoot(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := s.ctx(r)
	defer cancel()

	infos, err := s.boot.Snapshot(ctx)
	if err != nil {
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	plan := s.boot.PlanMerged(infos)

	byName := make(map[string]boot.ContainerInfo, len(infos))
	for _, c := range infos {
		byName[c.Name] = c
	}
	rows := make([]bootRow, 0, len(plan.Items))
	for _, it := range plan.Items {
		c := byName[it.Name]
		rows = append(rows, bootRow{
			Name: it.Name, State: c.State, Policy: c.Policy, Managed: boot.Managed(c.Policy),
			Project: c.Project, Auto: it.Auto, DelaySec: it.DelaySec,
		})
	}

	selfName := s.boot.SelfName()
	selfPolicy := ""
	if c, ok := byName[selfName]; ok {
		selfPolicy = c.Policy
	}

	var lastRun any
	if lr, ok := s.boot.LastResult(); ok {
		lastRun = lr
	} else if l, ok := s.st.LastRunLog("boot"); ok {
		// 没有结构化回执（老数据 / 手清过库）时退回运行记录本身，
		// 前端拿不到分项数字就只显示那句话。
		lastRun = map[string]any{"ts": l.TS, "status": l.Status, "message": l.Message}
	}

	writeOK(w, map[string]any{
		"plan":       plan,
		"containers": rows,
		"self":       map[string]any{"name": selfName, "restartPolicy": selfPolicy},
		"lastRun":    lastRun,
		"status":     s.boot.Status(),
	})
}

// hSaveBoot 保存整份清单（拖拽、改延时、开关、参数都走它）。
func (s *Server) hSaveBoot(w http.ResponseWriter, r *http.Request) {
	var in boot.Plan
	if err := decodeBody(r, &in); err != nil {
		writeErr(w, http.StatusBadRequest, "请求体格式错误："+err.Error())
		return
	}
	if err := s.boot.SavePlan(in); err != nil {
		writeErr(w, http.StatusInternalServerError, "保存失败："+err.Error())
		return
	}
	writeOK(w, map[string]any{"plan": s.boot.LoadPlan()})
}

type bootRunReq struct {
	// SkipDelay 跳过开头的初始等待，用于「改完立刻验证」。
	SkipDelay bool `json:"skipDelay"`
}

// hRunBoot 立即执行一次（异步，立刻返回；进度看 /api/boot/status）。
func (s *Server) hRunBoot(w http.ResponseWriter, r *http.Request) {
	var in bootRunReq
	if err := decodeBody(r, &in); err != nil {
		writeErr(w, http.StatusBadRequest, "请求体格式错误："+err.Error())
		return
	}
	if err := s.boot.Start(boot.TriggerManual, in.SkipDelay); err != nil {
		if errors.Is(err, boot.ErrBusy) {
			writeErr(w, http.StatusConflict, err.Error())
			return
		}
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOK(w, map[string]any{"ok": true})
}

// hBootStatus 当前进度（页面 2s 轮询）。
func (s *Server) hBootStatus(w http.ResponseWriter, r *http.Request) {
	writeOK(w, map[string]any{"status": s.boot.Status()})
}
