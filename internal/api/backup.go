package api

import (
	"archive/zip"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
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
	item, err := s.bk.SnapshotWith(ctx, in.Container, backup.SnapshotOptions{Reason: reason})
	if err != nil {
		s.nt.Emit("backup_failed", map[string]string{
			"container": in.Container, "result": "备份失败", "message": err.Error(),
		})
		writeErr(w, http.StatusBadGateway, "备份失败："+err.Error())
		return
	}
	writeOK(w, item)
}

type snapshotAllReq struct {
	// Containers 留空 = 备份全部容器。
	Containers []string `json:"containers"`
	Reason     string   `json:"reason"`
}

// hSnapshotAll 「立即备份全部容器」：一次动作 = 一批共用同一个时间戳的快照，
// 列表里聚合成一行。单个容器失败不中断整批。
func (s *Server) hSnapshotAll(w http.ResponseWriter, r *http.Request) {
	in := snapshotAllReq{}
	_ = decodeBody(r, &in)
	ctx, cancel := s.ctx(r)
	defer cancel()

	names := in.Containers
	if len(names) == 0 {
		list, err := s.dc.ListContainers(ctx)
		if err != nil {
			writeErr(w, http.StatusBadGateway, "读取容器列表失败："+err.Error())
			return
		}
		for _, c := range list {
			names = append(names, c.Name())
		}
	}
	sort.Strings(names)
	if len(names) == 0 {
		writeErr(w, http.StatusBadRequest, "没有可备份的容器")
		return
	}
	reason := in.Reason
	if reason == "" {
		reason = "manual"
	}
	res := s.bk.SnapshotMany(ctx, names, reason)
	if len(res.Items) == 0 {
		msg := "全部备份失败"
		if len(res.Failed) > 0 {
			msg = "备份失败：" + res.Failed[0].Error
		}
		writeErr(w, http.StatusBadGateway, msg)
		return
	}
	writeOK(w, res)
}

type importBackupReq struct {
	Container string `json:"container"`
	TS        string `json:"ts"`
	Content   string `json:"content"`
}

// hImportBackup 导入一份外部快照 json（「导入备份包」）。
func (s *Server) hImportBackup(w http.ResponseWriter, r *http.Request) {
	var in importBackupReq
	if err := decodeBody(r, &in); err != nil {
		writeErr(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	item, err := s.bk.Import(in.Container, in.TS, in.Content)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
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

// hExportBackup 把一份快照原样导出成一个 json 文件（可再导入到别的实例）。
func (s *Server) hExportBackup(w http.ResponseWriter, r *http.Request) {
	container := r.URL.Query().Get("container")
	ts := r.URL.Query().Get("ts")
	if container == "" || ts == "" {
		writeErr(w, http.StatusBadRequest, "需要提供 container 与 ts")
		return
	}
	items, err := s.bk.List()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, it := range items {
		if it.Container != container || it.TS != ts {
			continue
		}
		b, err := os.ReadFile(it.Path)
		if err != nil {
			writeErr(w, http.StatusInternalServerError, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition",
			fmt.Sprintf("attachment; filename=%q", container+"-"+ts+".json"))
		_, _ = w.Write(b)
		return
	}
	writeErr(w, http.StatusNotFound, "没有这份快照")
}

// hExportBatch 把「同一批（同一个时间戳）的全部快照」打包成一个 zip。
//
// 整批导出让浏览器下 24 次文件是不礼貌的（会被拦、也难整理），
// 一次动作 = 一个压缩包，与列表里的「一批」恰好对应。
func (s *Server) hExportBatch(w http.ResponseWriter, r *http.Request) {
	ts := r.URL.Query().Get("ts")
	if ts == "" {
		writeErr(w, http.StatusBadRequest, "需要提供 ts")
		return
	}
	items, err := s.bk.List()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	picked := []backup.SnapshotItem{}
	for _, it := range items {
		if it.TS == ts {
			picked = append(picked, it)
		}
	}
	if len(picked) == 0 {
		writeErr(w, http.StatusNotFound, "没有这一批快照")
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", "snapshots-"+ts+".zip"))
	zw := zip.NewWriter(w)
	defer zw.Close()
	for _, it := range picked {
		b, err := os.ReadFile(it.Path)
		if err != nil {
			continue
		}
		fw, err := zw.Create(fmt.Sprintf("%s-%s.json", it.Container, it.TS))
		if err != nil {
			return
		}
		if _, err := fw.Write(b); err != nil {
			return
		}
	}
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
	KeepPerContainer *int  `json:"keepPerContainer"`
	MaxAgeDays       *int  `json:"maxAgeDays"`
	MaxTotalMB       *int  `json:"maxTotalMB"`
	KeepPreUpdate    *bool `json:"keepPreUpdate"`
}

// hPruneBackups 手动触发一次清理。不带参数时用设置页里的保留策略，
// 这样「点一下立即清理」与「每天自动清理」的行为完全一致。
func (s *Server) hPruneBackups(w http.ResponseWriter, r *http.Request) {
	cfg := s.readSettings()
	in := pruneBackupsReq{}
	_ = decodeBody(r, &in)

	opt := backup.PruneOptions{
		KeepPerContainer: cfg.BackupKeepPerContainer,
		MaxAgeDays:       cfg.BackupMaxAgeDays,
		MaxTotalMB:       cfg.BackupMaxTotalMB,
		KeepPreUpdate:    cfg.BackupKeepPreUpdate,
	}
	if in.KeepPerContainer != nil {
		opt.KeepPerContainer = *in.KeepPerContainer
	}
	if in.MaxAgeDays != nil {
		opt.MaxAgeDays = *in.MaxAgeDays
	}
	if in.MaxTotalMB != nil {
		opt.MaxTotalMB = *in.MaxTotalMB
	}
	if in.KeepPreUpdate != nil {
		opt.KeepPreUpdate = *in.KeepPreUpdate
	}
	writeOK(w, s.bk.Prune(opt))
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
		"note": "compose 文件是项目的唯一事实源。备份 = 把它（含 .env / override）原样存一份副本，" +
			"可下载、可还原。标「不可见」的文件是因为宿主目录没挂进 Dockhelm —— " +
			"挂进来了（冒号右边随便叫什么，启动时会自动识别）就能读到并备份。",
	})
}

type projectSnapshotReq struct {
	Project string `json:"project"`
}

// hProjectSnapshot 备份单个 compose 项目。
func (s *Server) hProjectSnapshot(w http.ResponseWriter, r *http.Request) {
	var in projectSnapshotReq
	if err := decodeBody(r, &in); err != nil || strings.TrimSpace(in.Project) == "" {
		writeErr(w, http.StatusBadRequest, "需要指定项目名")
		return
	}
	ctx, cancel := s.ctx(r)
	defer cancel()
	item, err := s.bk.ProjectSnapshot(ctx, in.Project, "manual")
	if err != nil {
		s.nt.Emit("backup_failed", map[string]string{
			"container": in.Project, "result": "项目备份失败", "message": err.Error(),
		})
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	writeOK(w, item)
}

// hProjectSnapshotAll 备份全部「看得见文件」的项目；看不见的如实列在 failed 里。
func (s *Server) hProjectSnapshotAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := s.ctx(r)
	defer cancel()
	items, failed := s.bk.ProjectSnapshotAll(ctx)
	if len(items) == 0 && len(failed) > 0 {
		writeErr(w, http.StatusBadGateway, "没有可备份的项目："+failed[0].Error)
		return
	}
	writeOK(w, map[string]any{"items": items, "failed": failed, "total": len(items) + len(failed)})
}

// hProjectBackupDownload 下载一份项目备份：带 file 参数下载单个文件，否则打包成 zip。
func (s *Server) hProjectBackupDownload(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	project := q.Get("project")
	ts := q.Get("ts")
	if project == "" || ts == "" {
		writeErr(w, http.StatusBadRequest, "需要提供 project 与 ts")
		return
	}
	if name := strings.TrimSpace(q.Get("file")); name != "" {
		// 只允许下载这份备份清单里的文件 —— 绝不能让 ?file=../../dockhelm.db 读出去
		files, err := s.bk.ProjectBackupFiles(project, ts)
		if err != nil {
			writeErr(w, http.StatusNotFound, "没有这份备份")
			return
		}
		for _, f := range files {
			if f.Name != name {
				continue
			}
			b, err := os.ReadFile(f.Path)
			if err != nil {
				writeErr(w, http.StatusInternalServerError, err.Error())
				return
			}
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			w.Header().Set("Content-Disposition",
				fmt.Sprintf("attachment; filename=%q", f.Name))
			_, _ = w.Write(b)
			return
		}
		writeErr(w, http.StatusNotFound, "这份备份里没有该文件")
		return
	}
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition",
		fmt.Sprintf("attachment; filename=%q", project+"-"+ts+".zip"))
	if err := s.bk.WriteProjectBackupZip(w, project, ts); err != nil {
		// 响应头已经发出去了，这里只能把错误写进日志
		log.Printf("下载项目备份失败：%v", err)
	}
}

type projectRestoreReq struct {
	Project     string `json:"project"`
	TS          string `json:"ts"`
	PreSnapshot *bool  `json:"preSnapshot"`
}

// hProjectRestore 把项目备份里的 yaml / .env 写回原始宿主路径。
func (s *Server) hProjectRestore(w http.ResponseWriter, r *http.Request) {
	var in projectRestoreReq
	if err := decodeBody(r, &in); err != nil || in.Project == "" || in.TS == "" {
		writeErr(w, http.StatusBadRequest, "需要提供 project 与 ts")
		return
	}
	pre := true
	if in.PreSnapshot != nil {
		pre = *in.PreSnapshot
	}
	ctx, cancel := s.ctx(r)
	defer cancel()
	res := s.bk.RestoreProject(ctx, in.Project, in.TS, pre)
	status := http.StatusOK
	if !res.OK {
		status = http.StatusBadGateway
	}
	writeJSON(w, status, res)
}

func (s *Server) hDeleteProjectBackup(w http.ResponseWriter, r *http.Request) {
	project := pathParam(r, "project")
	ts := pathParam(r, "ts")
	if err := s.bk.DeleteProjectBackup(project, ts); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeOK(w, map[string]any{"ok": true})
}

// hReadProjectFile 预览一个 compose 文件的内容。
//
// 白名单只有一份：**刚刚列出来的那份 compose 文件清单**。
//
// 绝不能直接拿 MapHostPath 的结果去读 —— 它为了让「两边路径一致」这种挂法能用，
// 留了两条兜底（「原样可用」和「基名启发」），对容器内**任意**已存在的路径都会
// 返回 ok。于是 ?path=/data/dockhelm.db 就能把通知渠道的 Bot Token、SMTP 密码
// 整份读出来，?path=/proc/self/environ 还能读到 DOCKHELM_PASSWORD。
// 这个接口的用途只是「看一眼 compose 文件」，不需要那么宽的权限。
func (s *Server) hReadProjectFile(w http.ResponseWriter, r *http.Request) {
	p := strings.TrimSpace(r.URL.Query().Get("path"))
	if p == "" {
		writeErr(w, http.StatusBadRequest, "需要提供 path")
		return
	}
	ctx, cancel := s.ctx(r)
	defer cancel()

	projects, err := s.bk.ListProjects(ctx)
	if err != nil {
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	clean := filepath.Clean(p)
	local := ""
	for _, pr := range projects {
		for _, f := range pr.Readable {
			if filepath.Clean(f.HostPath) == clean {
				local = f.Path
				break
			}
		}
		if local != "" {
			break
		}
	}
	if local == "" {
		writeErr(w, http.StatusNotFound,
			"只能查看 compose 项目清单里的文件；这个路径不在清单里（或它在 Dockhelm 容器里看不见）")
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
		"path":     local,
		"hostPath": p,
		"content":  string(b),
	})
}

func fileExistsAndSmall(p string) bool {
	fi, err := os.Stat(p)
	if err != nil || fi.IsDir() {
		return false
	}
	return fi.Size() <= 1<<20
}
