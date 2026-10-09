package api

import (
	"net/http"
	"strings"

	"github.com/aaron2024s/dockhelm/internal/notify"
	"github.com/aaron2024s/dockhelm/internal/store"
)

func (s *Server) hNotifyPresets(w http.ResponseWriter, r *http.Request) {
	writeOK(w, map[string]any{
		"presets": notify.Presets,
		"vars":    notify.TemplateVars,
	})
}

func (s *Server) hNotifyEvents(w http.ResponseWriter, r *http.Request) {
	writeOK(w, map[string]any{
		"catalog": notify.SortedEvents(),
		"groups":  notify.Groups(),
		"current": s.st.ListNotifyEvents(),
	})
}

type saveEventsReq struct {
	Events []store.NotifyEvent `json:"events"`
}

func (s *Server) hSaveNotifyEvents(w http.ResponseWriter, r *http.Request) {
	var in saveEventsReq
	if err := decodeBody(r, &in); err != nil {
		writeErr(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	for i := range in.Events {
		e := &in.Events[i]
		if e.Level != notify.LevelUrgent && e.Level != notify.LevelNormal {
			if d, ok := notify.EventMap[e.Event]; ok {
				e.Level = d.Level
			} else {
				e.Level = notify.LevelNormal
			}
		}
	}
	if err := s.st.BulkUpsertNotifyEvents(in.Events); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOK(w, map[string]any{"ok": true, "current": s.st.ListNotifyEvents()})
}

func (s *Server) hNotifySettings(w http.ResponseWriter, r *http.Request) {
	writeOK(w, map[string]any{
		"settings": s.nt.GetSettings(),
		"channels": channelCount(s),
	})
}

func channelCount(s *Server) map[string]int {
	channels, _ := s.st.ListChannels()
	enabled := 0
	for _, c := range channels {
		if c.Enabled {
			enabled++
		}
	}
	return map[string]int{"total": len(channels), "enabled": enabled}
}

func (s *Server) hSaveNotifySettings(w http.ResponseWriter, r *http.Request) {
	var in notify.Settings
	if err := decodeBody(r, &in); err != nil {
		writeErr(w, http.StatusBadRequest, "请求体格式错误："+err.Error())
		return
	}
	if in.DedupeMinutes < 0 {
		in.DedupeMinutes = 10
	}
	if in.DailyLimit < 0 {
		in.DailyLimit = 200
	}
	if in.QuietMode != "digest" && in.QuietMode != "drop" {
		in.QuietMode = "digest"
	}
	if err := s.nt.SaveSettings(in); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	// 面板地址同步到通用设置里，供模板变量 {{url}} 使用
	writeOK(w, map[string]any{"ok": true, "settings": s.nt.GetSettings()})
}

func (s *Server) hListChannels(w http.ResponseWriter, r *http.Request) {
	channels, err := s.st.ListChannels()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if channels == nil {
		channels = []store.Channel{}
	}
	writeOK(w, map[string]any{"channels": channels})
}

type channelReq struct {
	Name    string         `json:"name"`
	Type    string         `json:"type"`
	Enabled bool           `json:"enabled"`
	Config  map[string]any `json:"config"`
}

func (s *Server) validateChannel(in channelReq) (string, bool) {
	if strings.TrimSpace(in.Name) == "" {
		return "渠道名称不能为空", false
	}
	p, ok := notify.PresetMap[in.Type]
	if !ok {
		return "不支持的渠道类型：" + in.Type, false
	}
	for _, f := range p.Fields {
		if !f.Required {
			continue
		}
		if v, _ := in.Config[f.Key].(string); strings.TrimSpace(v) == "" {
			return "缺少必填项：" + f.Label, false
		}
	}
	return "", true
}

func (s *Server) hCreateChannel(w http.ResponseWriter, r *http.Request) {
	var in channelReq
	if err := decodeBody(r, &in); err != nil {
		writeErr(w, http.StatusBadRequest, "请求体格式错误："+err.Error())
		return
	}
	if msg, ok := s.validateChannel(in); !ok {
		writeErr(w, http.StatusBadRequest, msg)
		return
	}
	if in.Config == nil {
		in.Config = map[string]any{}
	}
	// 补上预设默认值，避免用户没填的字段渲染出空串
	for k, v := range notify.DefaultConfig(in.Type) {
		if _, exists := in.Config[k]; !exists {
			in.Config[k] = v
		}
	}
	ch, err := s.st.CreateChannel(store.Channel{
		Name: in.Name, Type: in.Type, Enabled: in.Enabled, Config: in.Config,
	})
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOK(w, ch)
}

func (s *Server) hUpdateChannel(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(pathParam(r, "id"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "渠道 ID 无效")
		return
	}
	var in channelReq
	if err := decodeBody(r, &in); err != nil {
		writeErr(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	if msg, ok := s.validateChannel(in); !ok {
		writeErr(w, http.StatusBadRequest, msg)
		return
	}
	if err := s.st.UpdateChannel(store.Channel{
		ID: id, Name: in.Name, Type: in.Type, Enabled: in.Enabled, Config: in.Config,
	}); err != nil {
		writeErr(w, http.StatusNotFound, err.Error())
		return
	}
	ch, _ := s.st.GetChannel(id)
	writeOK(w, ch)
}

func (s *Server) hDeleteChannel(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(pathParam(r, "id"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "渠道 ID 无效")
		return
	}
	if err := s.st.DeleteChannel(id); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOK(w, map[string]any{"ok": true})
}

// hTestChannel 发一条测试消息。同步返回结果，用户马上能看到成功还是失败。
func (s *Server) hTestChannel(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(pathParam(r, "id"))
	if err != nil {
		writeErr(w, http.StatusBadRequest, "渠道 ID 无效")
		return
	}
	sent, msg := s.nt.EmitSync("test", map[string]string{
		"result":  "渠道连通性测试",
		"message": "如果你看到这条消息，说明这个渠道配置是可用的。",
	}, id)
	if sent == 0 {
		writeJSON(w, http.StatusBadGateway, map[string]any{"ok": false, "error": msg})
		return
	}
	writeOK(w, map[string]any{"ok": true, "message": msg})
}

func (s *Server) hNotifyHistory(w http.ResponseWriter, r *http.Request) {
	records, err := s.st.ListNotifyRecords(queryInt(r, "limit", 100))
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	if records == nil {
		records = []store.NotifyRecord{}
	}
	writeOK(w, map[string]any{"history": records})
}

func (s *Server) hClearNotifyHistory(w http.ResponseWriter, r *http.Request) {
	s.st.ClearNotifyHistory()
	writeOK(w, map[string]any{"ok": true})
}
