package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/aaron2024s/dockhelm/internal/backup"
)

func (s *Server) hListBackups(w http.ResponseWriter, r *http.Request) {
	items, err := s.bk.List()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOK(w, map[string]any{"backups": items})
}

func (s *Server) hBackupStats(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := s.ctx(r)
	defer cancel()
	writeOK(w, s.bk.GetStats(ctx))
}

type snapshotReq struct {
	Container string `json:"container"`
	Reason    string `json:"reason"`
}

func (s *Server) hSnapshot(w http.ResponseWriter, r *http.Request) {
	var in snapshotReq
	if err := decodeBody(r, &in); err != nil || strings.TrimSpace(in.Container) == "" {
		writeErr(w, http.StatusBadRequest, "需要指定容器名")
		return
	}
	ctx, cancel := s.ctx(r)
	defer cancel()
	reason := in.Reason
	if reason == "" {
		reason = "manual"
	}
	item, err := s.bk.Snapshot(ctx, in.Container, reason)
	if err != nil {
		s.nt.Emit("backup_failed", map[string]string{
			"container": in.Container, "result": "备份失败", "message": err.Error(),
		})
		writeErr(w, http.StatusBadGateway, "备份失败："+err.Error())
		return
	}
	writeOK(w, item)
}

func (s *Server) hBackupDiff(w http.ResponseWriter, r *http.Request) {
	container := r.URL.Query().Get("container")
	ts := r.URL.Query().Get("ts")
	if container == "" || ts == "" {
		writeErr(w, http.StatusBadRequest, "需要提供 container 与 ts")
		return
	}
	ctx, cancel := s.ctx(r)
	defer cancel()
	diff, err := s.bk.Diff(ctx, container, ts)
	if err != nil {
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	writeOK(w, map[string]any{"diff": diff, "count": len(diff)})
}

type restoreReq struct {
	Container           string `json:"container"`
	TS                  string `json:"ts"`
	KeepBackupContainer *bool  `json:"keepBackupContainer"`
	PreSnapshot         *bool  `json:"preSnapshot"`
}

// hRestore 用快照还原容器配置。高风险动作，前端会被要求二次确认。
func (s *Server) hRestore(w http.ResponseWriter, r *http.Request) {
	var in restoreReq
	if err := decodeBody(r, &in); err != nil || in.Container == "" || in.TS == "" {
		writeErr(w, http.StatusBadRequest, "需要提供 container 与 ts")
		return
	}
	if s.up.IsSelf(in.Container) {
		writeErr(w, http.StatusForbidden, "不允许还原 Dockhelm 自身容器")
		return
	}
	opt := backup.RestoreOptions{KeepBackupContainer: true, PreSnapshot: true}
	if in.KeepBackupContainer != nil {
		opt.KeepBackupContainer = *in.KeepBackupContainer
	}
	if in.PreSnapshot != nil {
		opt.PreSnapshot = *in.PreSnapshot
	}
	ctx, cancel := s.ctx(r)
	defer cancel()
	res := s.bk.Restore(ctx, in.Container, in.TS, opt)
	status := http.StatusOK
	if !res.OK {
		status = http.StatusBadGateway
	}
	writeJSON(w, status, res)
}

type pruneBackupsReq struct {
	KeepPerContainer int `json:"keepPerContainer"`
	MaxAgeDays       int `json:"maxAgeDays"`
	MaxTotalMB       int `json:"maxTotalMB"`
}

func (s *Server) hPruneBackups(w http.ResponseWriter, r *http.Request) {
	in := pruneBackupsReq{KeepPerContainer: 10, MaxAgeDays: 30, MaxTotalMB: 2048}
	_ = decodeBody(r, &in)
	if in.KeepPerContainer < 0 {
		in.KeepPerContainer = 10
	}
	res := s.bk.Prune(in.KeepPerContainer, in.MaxAgeDays, in.MaxTotalMB)
	writeOK(w, res)
}

func (s *Server) hDeleteBackup(w http.ResponseWriter, r *http.Request) {
	container := pathParam(r, "container")
	ts := pathParam(r, "ts")
	if err := s.bk.Delete(container, ts); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOK(w, map[string]any{"ok": true})
}

func (s *Server) hListProjects(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := s.ctx(r)
	defer cancel()
	projects, err := s.bk.ListProjects(ctx)
	if err != nil {
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	writeOK(w, map[string]any{
		"projects":     projects,
		"roots":        s.cfg.HostRoots,
		"pathMappings": s.cfg.PathMappings(),
		"note": "compose 文件是项目的唯一事实源。列表里的「看得见 / 看不见」取决于宿主目录有没有" +
			"挂进 Dockhelm —— 挂进来了（冒号右边随便叫什么，启动时会自动识别）就能读到；" +
			"没挂进来就只能记录路径，请自行备份那份 yaml。",
	})
}

// hReadProjectFile 预览一个 compose 文件的内容（只允许读备份目录与已映射的宿主路径）。
func (s *Server) hReadProjectFile(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query().Get("path")
	if strings.TrimSpace(p) == "" {
		writeErr(w, http.StatusBadRequest, "需要提供 path")
		return
	}
	local, ok := s.cfg.MapHostPath(p)
	if !ok {
		writeErr(w, http.StatusNotFound, "这个路径在 Dockhelm 容器里看不见，请把宿主目录挂进来")
		return
	}
	if !fileExistsAndSmall(local) {
		writeErr(w, http.StatusNotFound, "文件不存在或过大（上限 1MB）")
		return
	}
	b, err := os.ReadFile(local)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOK(w, map[string]any{
		"path":    local,
		"hostPath": p,
		"content": string(b),
	})
}

func fileExistsAndSmall(p string) bool {
	fi, err := os.Stat(p)
	if err != nil || fi.IsDir() {
		return false
	}
	return fi.Size() <= 1<<20
}
