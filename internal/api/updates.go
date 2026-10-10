package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/aaron2024s/dockhelm/internal/updater"
)

func (s *Server) hListUpdates(w http.ResponseWriter, r *http.Request) {
	s.checkMu.RLock()
	cache := make([]updater.CheckResult, len(s.checkCache))
	copy(cache, s.checkCache)
	at := s.checkAt
	s.checkMu.RUnlock()

	summary := map[string]int{
		"updateAvailable": 0, "upToDate": 0, "unknown": 0, "noUpstream": 0,
	}
	for _, c := range cache {
		switch c.Status {
		case updater.StatusUpdateAvailable:
			summary["updateAvailable"]++
		case updater.StatusUpToDate:
			summary["upToDate"]++
		case updater.StatusNoUpstream:
			summary["noUpstream"]++
		default:
			summary["unknown"]++
		}
	}
	writeOK(w, map[string]any{
		"results":   cache,
		"checkedAt": timeOrEmpty(at),
		"summary":   summary,
		"excluded":  keys(s.excluded()),
		"selfName":  s.up.SelfName(),
	})
}

type checkReq struct {
	Names []string `json:"names"`
}

// hCheckUpdates 巡检（只读）：走守护进程的 /distribution 比对远端摘要与本地摘要，
// 不拉层、不停容器。拉取后再比对镜像 ID 的那套「权威判定」属于更新流程，不属于检测。
func (s *Server) hCheckUpdates(w http.ResponseWriter, r *http.Request) {
	var in checkReq
	_ = decodeBody(r, &in)

	s.bus.Publish("update", "check_start", "running", map[string]any{"count": len(in.Names)})
	started := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
	defer cancel()

	ex := s.excluded()
	var results []updater.CheckResult
	if len(in.Names) == 0 {
		results = s.up.CheckAll(ctx, ex)
	} else {
		for _, n := range in.Names {
			if ex[n] {
				continue
			}
			res, err := s.up.Check(ctx, n)
			if err != nil {
				results = append(results, updater.CheckResult{
					Container: n, Status: updater.StatusUnknown, Reason: err.Error(),
					CheckedAt: time.Now().UTC().Format(time.RFC3339),
				})
				continue
			}
			results = append(results, *res)
		}
	}

	// 只有「全量巡检」才覆盖缓存，单容器检查不冲掉别的结果
	if len(in.Names) == 0 {
		s.checkMu.Lock()
		s.checkCache = results
		s.checkAt = time.Now()
		s.checkMu.Unlock()
	} else {
		s.mergeCheckResults(results)
	}

	avail := 0
	for _, c := range results {
		if c.Status == updater.StatusUpdateAvailable {
			avail++
		}
	}
	s.bus.Publish("update", "check_done", "success", map[string]any{
		"checked": len(results), "updateAvailable": avail,
		"durationMs": time.Since(started).Milliseconds(),
	})
	if avail > 0 {
		names := []string{}
		for _, c := range results {
			if c.Status == updater.StatusUpdateAvailable {
				names = append(names, c.Container)
			}
		}
		s.nt.Emit("update_available", map[string]string{
			"container": strings.Join(names, ", "),
			"result":    "检测到新版本",
			"message":   "共 " + itoa(avail) + " 个容器有可用更新",
		})
	}
	writeOK(w, map[string]any{"results": results, "checkedAt": time.Now().UTC().Format(time.RFC3339)})
}

func (s *Server) mergeCheckResults(rows []updater.CheckResult) {
	if len(rows) == 0 {
		return
	}
	s.checkMu.Lock()
	defer s.checkMu.Unlock()
	byName := map[string]int{}
	for i, c := range s.checkCache {
		byName[c.Container] = i
	}
	for _, row := range rows {
		if i, ok := byName[row.Container]; ok {
			s.checkCache[i] = row
		} else {
			s.checkCache = append(s.checkCache, row)
			byName[row.Container] = len(s.checkCache) - 1
		}
	}
	s.checkAt = time.Now()
}

type applyReq struct {
	Names []string `json:"names"`
	Force bool     `json:"force"`
}

// hApplyUpdates 执行更新。整个批次在后台跑，进度通过 SSE 推送。
func (s *Server) hApplyUpdates(w http.ResponseWriter, r *http.Request) {
	var in applyReq
	if err := decodeBody(r, &in); err != nil {
		writeErr(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	ex := s.excluded()
	names := []string{}
	for _, n := range in.Names {
		n = strings.TrimSpace(n)
		if n == "" || ex[n] || s.up.IsSelf(n) {
			continue
		}
		names = append(names, n)
	}
	if len(names) == 0 {
		writeErr(w, http.StatusBadRequest, "没有可更新的容器（可能都被排除了）")
		return
	}

	// 后台执行，立刻返回，前端靠 SSE 看进度
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
		defer cancel()
		results := s.up.UpdateMany(ctx, names, in.Force, "manual")

		// 更新完把缓存里的状态刷成最新
		rows := []updater.CheckResult{}
		for _, res := range results {
			switch res.Status {
			case updater.ResultUpToDate:
				rows = append(rows, updater.CheckResult{
					Container: res.Container, Image: res.Image,
					Status: updater.StatusUpToDate, Reason: "刚拉取过，镜像未变化",
					CheckedAt: time.Now().UTC().Format(time.RFC3339),
				})
			case updater.ResultUpdated:
				rows = append(rows, updater.CheckResult{
					Container: res.Container, Image: res.Image,
					Status: updater.StatusUpToDate, Reason: "刚刚更新到最新",
					CheckedAt: time.Now().UTC().Format(time.RFC3339),
				})
			}
		}
		s.mergeCheckResults(rows)
	}()

	writeOK(w, map[string]any{"accepted": len(names), "names": names})
}

// hUpdateStream 更新进度的 SSE 流（与全局事件流分开，方便前端只看更新）。
func (s *Server) hUpdateStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeErr(w, http.StatusInternalServerError, "当前连接不支持流式推送")
		return
	}
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-transform")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)

	for _, ev := range s.bus.History("update", 50) {
		writeSSE(w, ev)
	}
	flusher.Flush()

	ch, cancel := s.bus.Subscribe()
	defer cancel()
	ping := time.NewTicker(20 * time.Second)
	defer ping.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case ev, open := <-ch:
			if !open {
				return
			}
			if ev.Topic != "update" {
				continue
			}
			writeSSE(w, ev)
			flusher.Flush()
		case <-ping.C:
			_, _ = w.Write([]byte(": ping\n\n"))
			flusher.Flush()
		}
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	buf := [20]byte{}
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		buf[i] = '-'
	}
	return string(buf[i:])
}

func keys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
