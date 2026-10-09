package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/aaron2024s/dockhelm/internal/scheduler"
	"github.com/aaron2024s/dockhelm/internal/store"
)

func (s *Server) hScheduleActions(w http.ResponseWriter, r *http.Request) {
	writeOK(w, map[string]any{"actions": scheduler.Actions, "cronPresets": cronPresets()})
}

type cronPreset struct {
	Label string `json:"label"`
	Cron  string `json:"cron"`
}

func cronPresets() []cronPreset {
	return []cronPreset{
		{"每天凌晨 3 点", "0 3 * * *"},
		{"每天凌晨 4 点", "0 4 * * *"},
		{"每小时", "0 * * * *"},
		{"每 6 小时", "0 */6 * * *"},
		{"每 30 分钟", "*/30 * * * *"},
		{"每周一凌晨 3 点", "0 3 * * 1"},
		{"每月 1 号凌晨 3 点", "0 3 1 * *"},
	}
}

func (s *Server) hListSchedules(w http.ResponseWriter, r *http.Request) {
	list, err := s.st.ListSchedules()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	next := s.sch.NextRuns()
	for i := range list {
		if v, ok := next[list[i].ID]; ok {
			list[i].NextRun = v
		} else if list[i].Enabled {
			list[i].NextRun = ""
		}
	}
	writeOK(w, map[string]any{"schedules": list, "serverTime": time.Now().UTC().Format(time.RFC3339)})
}

type scheduleReq struct {
	Name    string   `json:"name"`
	Cron    string   `json:"cron"`
	Action  string   `json:"action"`
	Targets []string `json:"targets"`
	Enabled bool     `json:"enabled"`
}

func (s *Server) validateSchedule(in scheduleReq) (string, bool) {
	if strings.TrimSpace(in.Name) == "" {
		return "任务名称不能为空", false
	}
	if _, ok := scheduler.ActionMap[in.Action]; !ok {
		return "不支持的动作：" + in.Action, false
	}
	if err := scheduler.ValidateCron(in.Cron); err != nil {
		return err.Error(), false
	}
	if a := scheduler.ActionMap[in.Action]; a.NeedsTargets && len(in.Targets) == 0 {
		// 空 = 全部，这是允许的；但备份任务作用于全部容器通常不是本意，给个软提示在 UI 上做
		_ = a
	}
	return "", true
}

func (s *Server) hCreateSchedule(w http.ResponseWriter, r *http.Request) {
	var in scheduleReq
	if err := decodeBody(r, &in); err != nil {
		writeErr(w, http.StatusBadRequest, "请求体格式错误："+err.Error())
		return
	}
	if msg, ok := s.validateSchedule(in); !ok {
		writeErr(w, http.StatusBadRequest, msg)
		return
	}
	created, err := s.st.CreateSchedule(store.Schedule{
		Name: in.Name, Cron: strings.TrimSpace(in.Cron), Action: in.Action,
		Targets: in.Targets, Enabled: in.Enabled,
	})
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := s.sch.Reload(); err != nil {
		writeErr(w, http.StatusInternalServerError, "任务已保存但调度器重载失败："+err.Error())
		return
	}
	s.st.AddRunLog("schedule", created.Name, "success", "已创建计划任务", created.Cron)
	writeOK(w, created)
}

func (s *Server) hUpdateSchedule(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(pathParam(r, "id"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "任务 ID 无效")
		return
	}
	var in scheduleReq
	if err := decodeBody(r, &in); err != nil {
		writeErr(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	if msg, ok := s.validateSchedule(in); !ok {
		writeErr(w, http.StatusBadRequest, msg)
		return
	}
	cur, err := s.st.GetSchedule(id)
	if err != nil {
		writeErr(w, http.StatusNotFound, "任务不存在")
		return
	}
	cur.Name = in.Name
	cur.Cron = strings.TrimSpace(in.Cron)
	cur.Action = in.Action
	cur.Targets = in.Targets
	cur.Enabled = in.Enabled
	if err := s.st.UpdateSchedule(*cur); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := s.sch.Reload(); err != nil {
		writeErr(w, http.StatusInternalServerError, "任务已保存但调度器重载失败："+err.Error())
		return
	}
	writeOK(w, cur)
}

func (s *Server) hDeleteSchedule(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(pathParam(r, "id"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "任务 ID 无效")
		return
	}
	if err := s.st.DeleteSchedule(id); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	_ = s.sch.Reload()
	writeOK(w, map[string]any{"ok": true})
}

// hRunSchedule 立即运行一次。
func (s *Server) hRunSchedule(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(pathParam(r, "id"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "任务 ID 无效")
		return
	}
	msg, ok := s.sch.RunNow(id)
	writeOK(w, map[string]any{"ok": ok, "message": msg})
}
