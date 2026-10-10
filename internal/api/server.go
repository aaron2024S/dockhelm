// Package api 提供 HTTP 接口（含 SSE 实时进度）。
//
// 用标准库 net/http + Go 1.22 起的 ServeMux 路由模式（"GET /api/x/{id}"），
// 不引任何 Web 框架：路由、中间件、JSON 编解码都能用标准库直接写完，
// 少一个依赖就少一份构建和供应链风险。
package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aaron2024s/dockhelm/internal/auth"
	"github.com/aaron2024s/dockhelm/internal/backup"
	"github.com/aaron2024s/dockhelm/internal/bus"
	"github.com/aaron2024s/dockhelm/internal/config"
	"github.com/aaron2024s/dockhelm/internal/dockerx"
	"github.com/aaron2024s/dockhelm/internal/notify"
	"github.com/aaron2024s/dockhelm/internal/scheduler"
	"github.com/aaron2024s/dockhelm/internal/store"
	"github.com/aaron2024s/dockhelm/internal/updater"
	"github.com/aaron2024s/dockhelm/internal/version"
)

// Server 持有全部依赖。
type Server struct {
	cfg    *config.Config
	st     *store.Store
	dc     *dockerx.Client
	auth   *auth.Manager
	up     *updater.Updater
	sch    *scheduler.Scheduler
	bk     *backup.Service
	nt     *notify.Manager
	bus    *bus.Bus
	mux    *http.ServeMux
	static http.Handler

	// 更新检测结果的缓存
	checkMu    sync.RWMutex
	checkCache []updater.CheckResult
	checkAt    time.Time

	// 后台「周期检测 + 自动更新」循环的状态
	autoWake      chan wakeAction
	autoMu        sync.Mutex
	autoBusy      bool
	autoLast      *AutoRunSummary
	autoLastCheck time.Time
	// autoNextAt 下一轮巡检实际排定的时刻（定时器是唯一真相；
	// 冷却期/提前逻辑会让「上次巡检 + 间隔」这种推算失真）。
	autoNextAt time.Time
}

// Deps 构造参数。
type Deps struct {
	Cfg    *config.Config
	Store  *store.Store
	Docker *dockerx.Client
	Auth   *auth.Manager
	Update *updater.Updater
	Sched  *scheduler.Scheduler
	Backup *backup.Service
	Notify *notify.Manager
	Bus    *bus.Bus
	Static http.Handler
}

// New 创建 API 服务器。
func New(d Deps) *Server {
	s := &Server{
		cfg: d.Cfg, st: d.Store, dc: d.Docker, auth: d.Auth,
		up: d.Update, sch: d.Sched, bk: d.Backup, nt: d.Notify, bus: d.Bus,
		mux: http.NewServeMux(), static: d.Static,
		autoWake: make(chan wakeAction, 1),
	}
	// 设置是常驻状态的真相，进程一起来就先把它推给更新引擎，
	// 免得「重启之后策略回到默认值、直到用户手动存一次设置才生效」。
	s.applyPolicy(s.readSettings())
	// 首次启动把预置的常用加速源直接写进「我的加速源」（只灌一次，之后用户说了算）。
	seedRegistries(d.Store)
	s.routes()
	return s
}

// Handler 返回根 handler。
func (s *Server) Handler() http.Handler {
	return s.recoverMW(s.logMW(s.mux))
}

// ---------- 路由 ----------

func (s *Server) routes() {
	m := s.mux

	// —— 公开接口 ——
	m.HandleFunc("GET /api/health", s.hHealth)
	m.HandleFunc("GET /api/about", s.hAbout)
	m.HandleFunc("GET /api/session", s.hSession)
	m.HandleFunc("POST /api/login", s.hLogin)
	m.HandleFunc("POST /api/setup", s.hSetup)
	m.HandleFunc("POST /api/logout", s.hLogout)

	// —— 需要登录 ——
	a := s.requireAuth

	m.HandleFunc("POST /api/account/password", a(s.hChangePassword))
	m.HandleFunc("POST /api/account/sessions/logout-others", a(s.hLogoutOthers))
	m.HandleFunc("GET /api/account", a(s.hAccount))

	m.HandleFunc("GET /api/overview", a(s.hOverview))
	m.HandleFunc("GET /api/usage", a(s.hUsage))
	m.HandleFunc("GET /api/system", a(s.hSystem))
	m.HandleFunc("GET /api/events/stream", a(s.hEventStream))

	m.HandleFunc("GET /api/containers", a(s.hListContainers))
	m.HandleFunc("GET /api/containers/{name}", a(s.hInspectContainer))
	m.HandleFunc("GET /api/containers/{name}/logs", a(s.hContainerLogs))
	m.HandleFunc("GET /api/containers/{name}/stats", a(s.hContainerStats))
	m.HandleFunc("GET /api/containers/{name}/export", a(s.hExportContainer))
	m.HandleFunc("POST /api/containers/{name}/action", a(s.hContainerAction))
	m.HandleFunc("POST /api/containers/{name}/rename", a(s.hRenameContainer))
	m.HandleFunc("DELETE /api/containers/{name}", a(s.hRemoveContainer))

	m.HandleFunc("GET /api/images", a(s.hListImages))
	m.HandleFunc("DELETE /api/images/{id}", a(s.hRemoveImage))
	m.HandleFunc("POST /api/images/prune", a(s.hPruneImages))

	m.HandleFunc("GET /api/networks", a(s.hListNetworks))
	m.HandleFunc("GET /api/volumes", a(s.hListVolumes))
	m.HandleFunc("DELETE /api/volumes/{name}", a(s.hRemoveVolume))

	m.HandleFunc("GET /api/updates", a(s.hListUpdates))
	m.HandleFunc("POST /api/updates/check", a(s.hCheckUpdates))
	m.HandleFunc("POST /api/updates/apply", a(s.hApplyUpdates))
	m.HandleFunc("GET /api/updates/stream", a(s.hUpdateStream))
	m.HandleFunc("GET /api/updates/auto", a(s.hGetAutoUpdate))
	m.HandleFunc("POST /api/updates/auto-run", a(s.hRunAutoUpdate))

	m.HandleFunc("GET /api/schedules", a(s.hListSchedules))
	m.HandleFunc("POST /api/schedules", a(s.hCreateSchedule))
	m.HandleFunc("PUT /api/schedules/{id}", a(s.hUpdateSchedule))
	m.HandleFunc("DELETE /api/schedules/{id}", a(s.hDeleteSchedule))
	m.HandleFunc("POST /api/schedules/{id}/run", a(s.hRunSchedule))
	m.HandleFunc("POST /api/schedules/{id}/enabled", a(s.hSetScheduleEnabled))
	m.HandleFunc("GET /api/schedules/actions", a(s.hScheduleActions))

	m.HandleFunc("GET /api/registries", a(s.hGetRegistries))
	m.HandleFunc("PUT /api/registries", a(s.hSaveRegistries))
	m.HandleFunc("POST /api/registries/test", a(s.hTestRegistry))
	m.HandleFunc("POST /api/registries/test-all", a(s.hTestAllRegistries))

	m.HandleFunc("GET /api/backups", a(s.hListBackups))
	m.HandleFunc("GET /api/backups/stats", a(s.hBackupStats))
	m.HandleFunc("POST /api/backups/snapshot", a(s.hSnapshot))
	m.HandleFunc("POST /api/backups/snapshot-all", a(s.hSnapshotAll))
	m.HandleFunc("GET /api/backups/export", a(s.hExportBackup))
	m.HandleFunc("GET /api/backups/export-batch", a(s.hExportBatch))
	m.HandleFunc("GET /api/backups/diff", a(s.hBackupDiff))
	m.HandleFunc("POST /api/backups/restore", a(s.hRestore))
	m.HandleFunc("POST /api/backups/prune", a(s.hPruneBackups))
	m.HandleFunc("POST /api/backups/import", a(s.hImportBackup))
	m.HandleFunc("DELETE /api/backups/{container}/{ts}", a(s.hDeleteBackup))
	m.HandleFunc("GET /api/backups/projects", a(s.hListProjects))
	m.HandleFunc("POST /api/backups/projects/snapshot", a(s.hProjectSnapshot))
	m.HandleFunc("POST /api/backups/projects/snapshot-all", a(s.hProjectSnapshotAll))
	m.HandleFunc("POST /api/backups/projects/restore", a(s.hProjectRestore))
	m.HandleFunc("GET /api/backups/projects/download", a(s.hProjectBackupDownload))
	m.HandleFunc("DELETE /api/backups/projects/{project}/{ts}", a(s.hDeleteProjectBackup))
	m.HandleFunc("GET /api/backups/file", a(s.hReadProjectFile))

	m.HandleFunc("GET /api/notify/presets", a(s.hNotifyPresets))
	m.HandleFunc("GET /api/notify/events", a(s.hNotifyEvents))
	m.HandleFunc("PUT /api/notify/events", a(s.hSaveNotifyEvents))
	m.HandleFunc("GET /api/notify/settings", a(s.hNotifySettings))
	m.HandleFunc("PUT /api/notify/settings", a(s.hSaveNotifySettings))
	m.HandleFunc("GET /api/notify/channels", a(s.hListChannels))
	m.HandleFunc("POST /api/notify/channels", a(s.hCreateChannel))
	m.HandleFunc("PUT /api/notify/channels/{id}", a(s.hUpdateChannel))
	m.HandleFunc("DELETE /api/notify/channels/{id}", a(s.hDeleteChannel))
	m.HandleFunc("POST /api/notify/channels/{id}/test", a(s.hTestChannel))
	m.HandleFunc("GET /api/notify/history", a(s.hNotifyHistory))
	m.HandleFunc("DELETE /api/notify/history", a(s.hClearNotifyHistory))

	m.HandleFunc("GET /api/settings", a(s.hGetSettings))
	m.HandleFunc("PUT /api/settings", a(s.hSaveSettings))
	// PATCH 只改请求体里出现过的键 —— 设置项分散在三个页面里，整份替换会互相覆盖
	m.HandleFunc("PATCH /api/settings", a(s.hPatchSettings))
	m.HandleFunc("GET /api/logs", a(s.hLogs))
	m.HandleFunc("DELETE /api/logs", a(s.hClearLogs))

	// —— 前端静态资源（SPA 回退）——
	m.HandleFunc("/", s.hStatic)
}

// ---------- 中间件 ----------

func (s *Server) logMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") && r.URL.Path != "/api/health" {
			log.Printf("%s %s", r.Method, r.URL.Path)
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) recoverMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic on %s %s: %v", r.Method, r.URL.Path, rec)
				writeErr(w, http.StatusInternalServerError, "服务器内部错误")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// requireAuth 校验会话 Cookie。
func (s *Server) requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(auth.CookieName)
		if err != nil || !s.auth.Validate(c.Value) {
			writeErr(w, http.StatusUnauthorized, "未登录或会话已过期")
			return
		}
		next(w, r)
	}
}

// ---------- 响应助手 ----------

type apiError struct {
	Error string `json:"error"`
	Code  string `json:"code,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	if v != nil {
		_ = json.NewEncoder(w).Encode(v)
	}
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeErrCode(w, status, msg, "")
}

func writeErrCode(w http.ResponseWriter, status int, msg, code string) {
	writeJSON(w, status, apiError{Error: msg, Code: code})
}

func writeOK(w http.ResponseWriter, v any) { writeJSON(w, http.StatusOK, v) }

func decodeBody(r *http.Request, v any) error {
	defer r.Body.Close()
	dec := json.NewDecoder(io.LimitReader(r.Body, 4<<20))
	dec.DisallowUnknownFields()
	err := dec.Decode(v)
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

// pathParam 取路径参数（ServeMux 的 {name}）。
func pathParam(r *http.Request, key string) string {
	return r.PathValue(key)
}

func queryBool(r *http.Request, key string, def bool) bool {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	return v == "1" || v == "true" || v == "yes"
}

func queryInt(r *http.Request, key string, def int) int {
	v := r.URL.Query().Get(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

// clientIP 取客户端 IP。只有配了 TrustedProxies 才信任 X-Forwarded-For。
func (s *Server) clientIP(r *http.Request) string {
	if len(s.cfg.TrustedProxies) > 0 {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			if i := strings.IndexByte(xff, ','); i > 0 {
				return strings.TrimSpace(xff[:i])
			}
			return strings.TrimSpace(xff)
		}
		if xr := r.Header.Get("X-Real-IP"); xr != "" {
			return strings.TrimSpace(xr)
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// ---------- 静态资源（SPA） ----------

func (s *Server) hStatic(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/api/") {
		writeErr(w, http.StatusNotFound, "接口不存在")
		return
	}
	if s.static == nil {
		writeErr(w, http.StatusNotFound, "前端资源未内嵌")
		return
	}
	// index.html 绝不缓存，避免升级后仍加载旧页面
	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		w.Header().Set("Cache-Control", "no-cache, must-revalidate")
	}
	s.static.ServeHTTP(w, r)
}

// ---------- 实时事件（SSE） ----------

func (s *Server) hEventStream(w http.ResponseWriter, r *http.Request) {
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

	// 首屏补发：每个 topic 的最后一条状态，避免刷新后进度条空白
	for _, ev := range s.bus.Snapshot() {
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
			writeSSE(w, ev)
			flusher.Flush()
		case <-ping.C:
			_, _ = fmt.Fprint(w, ": ping\n\n")
			flusher.Flush()
		}
	}
}

func writeSSE(w io.Writer, ev bus.Event) {
	b, err := json.Marshal(ev)
	if err != nil {
		return
	}
	_, _ = fmt.Fprintf(w, "id: %d\nevent: %s\ndata: %s\n\n", ev.ID, ev.Topic, b)
}

// ---------- 通用：运行上下文 ----------

func (s *Server) ctx(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), 10*time.Minute)
}

// ---------- 排除列表 ----------

// excluded 返回用户配置的「永不自动更新」容器集合（含 Dockhelm 自身）。
func (s *Server) excluded() map[string]bool {
	out := map[string]bool{}
	if self := s.up.SelfName(); self != "" {
		out[self] = true
	}
	var list []string
	s.st.GetJSON("exclude.containers", &list)
	for _, n := range list {
		if n = strings.TrimSpace(n); n != "" {
			out[n] = true
		}
	}
	return out
}

// ---------- 版本信息 ----------

func (s *Server) hAbout(w http.ResponseWriter, r *http.Request) {
	info := version.Get()
	// 顺带把运行时长与 Docker 版本带给「关于」页，省一次请求
	extra := map[string]any{
		"uptimeSeconds": int(time.Since(version.Runtime).Seconds()),
		"dataDir":       s.cfg.DataDir,
		"dockerHost":    s.dc.Host(),
		// 端口可被环境变量覆盖，把「实际生效的值 + 它来自哪个变量」都露出来，
		// 「关于」页直接展示，省得去翻启动日志。
		"listen":       s.cfg.Listen,
		"listenSource": s.cfg.ListenSource,
	}
	if v := s.dc.APIVersion(); v != "" {
		extra["dockerApiVersion"] = v
	}
	if insp, err := s.dc.Info(r.Context()); err == nil {
		extra["dockerVersion"] = insp.ServerVersion
		extra["dockerOs"] = insp.OperatingSystem
		extra["dockerArch"] = insp.Architecture
		extra["dockerKernel"] = insp.KernelVersion
		extra["dockerRoot"] = insp.DockerRootDir
		extra["containersTotal"] = insp.Containers
		extra["imagesTotal"] = insp.Images
	}
	writeOK(w, map[string]any{"about": info, "runtime": extra})
}

func (s *Server) hHealth(w http.ResponseWriter, r *http.Request) {
	status := "ok"
	dockerOK := true
	if err := s.dc.Ping(r.Context()); err != nil {
		dockerOK = false
		status = "degraded"
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  status,
		"version": version.Version,
		"docker":  dockerOK,
		"time":    time.Now().UTC().Format(time.RFC3339),
	})
}

// ---------- 设置 ----------

// Settings 通用设置。
//
// 这里的字段名就是「设置页上看得见的那几个旋钮」—— 不搞一套抽象的配置模型，
// 因为一个自托管的单用户面板最终只需要「这些开关各自是什么值」这一个真相。
type Settings struct {
	Exclude      []string `json:"exclude"`
	PanelURL     string   `json:"panelURL"`
	Concurrency  int      `json:"concurrency"`
	LogRetention int      `json:"logRetention"`
	CheckOnStart bool     `json:"checkOnStart"`

	// —— 检测 ——
	// CheckIntervalHours 周期性自动巡检的间隔（小时）。取值必为 checkIntervalChoices
	// 之一（不允许 0 —— 见 normalizeInterval：0 会让自动更新永久失效）。
	CheckIntervalHours int `json:"checkIntervalHours"`
	// NotifyOnCheck 巡检发现新版本时推一条通知。
	NotifyOnCheck bool `json:"notifyOnCheck"`

	// —— 更新策略 ——
	// AutoApply 自动更新总开关。关闭时定时任务只检测、绝不动容器。
	AutoApply bool `json:"autoApply"`
	// PullOnce 同一轮批量更新里，同一个镜像只下载一次。
	PullOnce bool `json:"pullOnce"`
	// BackupBefore 重建容器前先写一份配置快照。
	BackupBefore bool `json:"backupBefore"`
	// CleanupAfter 更新成功后清理没有任何容器引用的旧镜像。
	CleanupAfter bool `json:"cleanupAfter"`
	// DirectFirst 显式写了域名的镜像（ghcr.io 等）直连，不套加速源。
	DirectFirst bool `json:"directFirst"`

	// —— 备份保留策略 ——
	BackupKeepPerContainer int  `json:"backupKeepPerContainer"`
	BackupMaxAgeDays       int  `json:"backupMaxAgeDays"`
	BackupMaxTotalMB       int  `json:"backupMaxTotalMB"`
	BackupKeepPreUpdate    bool `json:"backupKeepPreUpdate"`
}

// 检测间隔的可选值（小时）。前端下拉框与后端校验共用这一份，避免两边说法不一致。
//
// 🚨 这里绝不能有 0：周期检测是自动更新的唯一触发源，一旦允许「关闭」，
// 用户会在「自动更新：已开启」的状态下永远等不到动作 —— 没有检测轮次，
// 就永远发现不了更新，也就永远不会更新。
var checkIntervalChoices = []int{1, 2, 3, 6, 12, 24}

func (s *Server) readSettings() Settings {
	var out Settings
	out.Concurrency = 2
	out.LogRetention = 100
	out.CheckOnStart = true
	out.CheckIntervalHours = 1
	out.NotifyOnCheck = false
	out.AutoApply = false
	out.PullOnce = true
	out.BackupBefore = true
	out.CleanupAfter = true
	out.DirectFirst = true
	out.BackupKeepPerContainer = 10
	out.BackupMaxAgeDays = 30
	out.BackupMaxTotalMB = 2048
	out.BackupKeepPreUpdate = true

	readInt := func(key string, dst *int) {
		if v := s.st.GetSetting(key, ""); v != "" {
			if n, err := strconv.Atoi(v); err == nil {
				*dst = n
			}
		}
	}
	readBool := func(key string, dst *bool) {
		if v := s.st.GetSetting(key, ""); v != "" {
			*dst = v == "1"
		}
	}

	readInt("update.concurrency", &out.Concurrency)
	readBool("update.checkOnStart", &out.CheckOnStart)
	readInt("log.retention", &out.LogRetention)
	readInt("update.checkInterval", &out.CheckIntervalHours)
	// 老库里可能存过 0（旧版本下拉框有「关闭」档）⇒ 收敛到默认值，
	// 否则这个用户升级后周期检测与自动更新会一直是死的。
	out.CheckIntervalHours = normalizeInterval(out.CheckIntervalHours)
	readBool("update.notifyOnCheck", &out.NotifyOnCheck)
	readBool("update.autoApply", &out.AutoApply)
	readBool("update.pullOnce", &out.PullOnce)
	readBool("update.backupBefore", &out.BackupBefore)
	readBool("update.cleanupAfter", &out.CleanupAfter)
	readBool("update.directFirst", &out.DirectFirst)
	readInt("backup.keepPerContainer", &out.BackupKeepPerContainer)
	readInt("backup.maxAgeDays", &out.BackupMaxAgeDays)
	readInt("backup.maxTotalMB", &out.BackupMaxTotalMB)
	readBool("backup.keepPreUpdate", &out.BackupKeepPreUpdate)

	out.PanelURL = s.st.GetSetting(notify.KeyPanelURL, "")
	s.st.GetJSON("exclude.containers", &out.Exclude)
	if out.Exclude == nil {
		out.Exclude = []string{}
	}
	return out
}

// applyPolicy 把设置里的更新策略推给更新引擎。
// 引擎是常驻对象，设置改了必须立刻生效，否则「改了并发度要重启才生效」这种事
// 一定会被当成 bug 报上来。
func (s *Server) applyPolicy(in Settings) {
	opt := s.up.Options()
	opt.Apply(updater.Policy{
		Concurrency:  in.Concurrency,
		PullOnce:     in.PullOnce,
		BackupBefore: in.BackupBefore,
		CleanupAfter: in.CleanupAfter,
		DirectFirst:  in.DirectFirst,
	})
	s.up.SetOptions(opt)
}

func (s *Server) hGetSettings(w http.ResponseWriter, r *http.Request) {
	writeOK(w, s.readSettings())
}

func (s *Server) hSaveSettings(w http.ResponseWriter, r *http.Request) {
	var in Settings
	if err := decodeBody(r, &in); err != nil {
		writeErr(w, http.StatusBadRequest, "请求体格式错误："+err.Error())
		return
	}
	s.persistSettings(in)
	writeOK(w, s.readSettings())
}

// hPatchSettings 局部更新设置：只改请求体里**真正出现过**的那些键，其余保持原值。
//
// 为什么需要它：设置项被拆在三个页面里（设置页 / 镜像加速源页 / 备份与恢复页），
// 而 PUT /api/settings 是「整份替换」—— 每页手里那份都是进页面时拉的快照，
// 谁后保存谁就把别人期间的改动悄悄抹掉（在设置页把并发度改成 4，再去备份页
// 保存保留策略，并发度就回到了旧值）。
//
// 用 map[string]json.RawMessage 才能区分「字段没出现」和「字段是零值」——
// 结构体解码做不到这一点，这是 PUT 覆盖问题绕不开的一步。
func (s *Server) hPatchSettings(w http.ResponseWriter, r *http.Request) {
	var raw map[string]json.RawMessage
	if err := decodeBody(r, &raw); err != nil {
		writeErr(w, http.StatusBadRequest, "请求体格式错误："+err.Error())
		return
	}
	if len(raw) == 0 {
		writeErr(w, http.StatusBadRequest, "请求体是空的，没有要修改的设置项")
		return
	}
	// 以服务端当前值为底，把请求里出现的键盖上去
	cur, err := json.Marshal(s.readSettings())
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	merged := map[string]json.RawMessage{}
	if err := json.Unmarshal(cur, &merged); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	for k, v := range raw {
		merged[k] = v
	}
	body, err := json.Marshal(merged)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	var next Settings
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.DisallowUnknownFields() // 键名写错要报错，不能悄悄忽略
	if err := dec.Decode(&next); err != nil {
		writeErr(w, http.StatusBadRequest, "设置项无法识别："+err.Error())
		return
	}
	s.persistSettings(next)
	writeOK(w, s.readSettings())
}

// persistSettings 校验 + 落盘 + 推给常驻对象。PUT 与 PATCH 共用这一份写入口。
func (s *Server) persistSettings(in Settings) {
	// 先留一份「改之前」的快照：等会儿唤醒后台循环时，要靠它判断这次改动
	// 是不是「用户明显在等着看动作」的那两类（见 wakeDecision）。
	prev := s.readSettings()
	if in.Concurrency < 1 || in.Concurrency > 8 {
		in.Concurrency = 2
	}
	in.CheckIntervalHours = normalizeInterval(in.CheckIntervalHours)
	if in.BackupKeepPerContainer < 0 {
		in.BackupKeepPerContainer = 10
	}
	if in.BackupMaxAgeDays < 0 {
		in.BackupMaxAgeDays = 30
	}
	if in.BackupMaxTotalMB < 0 {
		in.BackupMaxTotalMB = 2048
	}

	_ = s.st.SetSetting("update.concurrency", strconv.Itoa(in.Concurrency))
	_ = s.st.SetSetting("update.checkOnStart", boolStr(in.CheckOnStart))
	_ = s.st.SetSetting("log.retention", strconv.Itoa(in.LogRetention))
	_ = s.st.SetSetting(notify.KeyPanelURL, in.PanelURL)
	_ = s.st.SetSetting("update.checkInterval", strconv.Itoa(in.CheckIntervalHours))
	_ = s.st.SetSetting("update.notifyOnCheck", boolStr(in.NotifyOnCheck))
	_ = s.st.SetSetting("update.autoApply", boolStr(in.AutoApply))
	_ = s.st.SetSetting("update.pullOnce", boolStr(in.PullOnce))
	_ = s.st.SetSetting("update.backupBefore", boolStr(in.BackupBefore))
	_ = s.st.SetSetting("update.cleanupAfter", boolStr(in.CleanupAfter))
	_ = s.st.SetSetting("update.directFirst", boolStr(in.DirectFirst))
	_ = s.st.SetSetting("backup.keepPerContainer", strconv.Itoa(in.BackupKeepPerContainer))
	_ = s.st.SetSetting("backup.maxAgeDays", strconv.Itoa(in.BackupMaxAgeDays))
	_ = s.st.SetSetting("backup.maxTotalMB", strconv.Itoa(in.BackupMaxTotalMB))
	_ = s.st.SetSetting("backup.keepPreUpdate", boolStr(in.BackupKeepPreUpdate))
	_ = s.st.SetJSON("exclude.containers", in.Exclude)

	s.applyPolicy(in)
	// 检测频率/自动更新总开关可能刚被改过，唤醒后台循环重新计时；
	// 若属于「刚打开自动更新」或「把间隔调小了」，则把下一轮提前（见 wakeDecision）。
	s.wakeAutoLoop(wakeDecision(prev, in))
}

// normalizeInterval 把检测间隔收敛到允许的取值上，防止存进一个
// 「每 0.3 小时跑一次」这种会把镜像仓库打限流的值。
func normalizeInterval(h int) int {
	for _, c := range checkIntervalChoices {
		if c == h {
			return h
		}
	}
	if h > 24 {
		return 24
	}
	// 非合法值（含老库里的 0 / 负数）一律收敛到默认 1 小时 —— 绝不返回 0，
	// 0 会让周期循环彻底停摆，自动更新也跟着永久失效。
	return 1
}

func boolStr(b bool) string {
	if b {
		return "1"
	}
	return "0"
}

// ---------- 运行日志 ----------

func (s *Server) hLogs(w http.ResponseWriter, r *http.Request) {
	logs, err := s.st.ListRunLogs(queryInt(r, "limit", 200))
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOK(w, map[string]any{"logs": logs})
}

func (s *Server) hClearLogs(w http.ResponseWriter, r *http.Request) {
	s.st.ClearRunLogs()
	writeOK(w, map[string]any{"ok": true})
}

// ---------- 系统信息 ----------

func (s *Server) hSystem(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := s.ctx(r)
	defer cancel()
	insp, err := s.dc.Info(ctx)
	if err != nil {
		writeErr(w, http.StatusBadGateway, "无法连接 Docker 守护进程："+err.Error())
		return
	}
	diskFree, diskTotal, diskPath := diskUsage(s.cfg.DataDir)
	writeOK(w, map[string]any{
		"docker": insp,
		"host": map[string]any{
			"name":      hostname(),
			"dataDir":   s.cfg.DataDir,
			"diskPath":  diskPath,
			"diskFree":  diskFree,
			"diskTotal": diskTotal,
		},
	})
}

func diskUsage(path string) (free, total uint64, real string) {
	real = path
	for p := filepath.Clean(path); p != "/" && p != "."; p = filepath.Dir(p) {
		if fi, err := os.Stat(p); err == nil && fi.IsDir() {
			real = p
			break
		}
	}
	free, total = statfs(real)
	return
}

func hostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return h
}
