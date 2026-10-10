// 项目（compose yaml）备份。
//
// 与容器配置快照的分工：
//   - 容器快照 = docker inspect = 「容器现在长什么样」，能把应用原样重建起来；
//   - 项目备份 = 项目目录里那几个 yaml / .env / override 的**原样副本** = 「源头文件」，
//     人类可读、可下载、可迁移、能进 git。
//
// 明确不做的事：**不打包项目目录里的任何数据**（数据库、缓存、上传的文件）。
// 那些属于「数据」，Dockhelm 不碰，界面上也如实说明。
package backup

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// manifestName 每份项目备份里记录「原路径 + 校验值」的文件。
// 用下划线开头，与真实文件区分开。
const manifestName = "_manifest.json"

// ProjectManifest 一份项目备份的清单。
type ProjectManifest struct {
	Project   string          `json:"project"`
	TS        string          `json:"ts"`
	Created   string          `json:"created"`
	Reason    string          `json:"reason"`
	WorkDir   string          `json:"workDir"`
	Files     []ProjectFile   `json:"files"`
	Checksums map[string]string `json:"checksums"` // 文件名 → sha256
}

// projectBackupDir 一份项目备份的落盘目录：<backupDir>/projects/<项目>/<ts>/
func (s *Service) projectBackupDir(project, ts string) string {
	return filepath.Join(s.cfg.ProjectBackupDir(), sanitize(project), sanitize(ts))
}

// ProjectSnapshot 备份一个 compose 项目当前可读的 yaml / .env / override。
//
// 只在**至少读得到一个文件**时才写：一个文件都读不到（宿主目录没挂进来）时
// 什么都不写、如实报错，绝不写一份空备份假装成功 —— 那会让用户以为 yaml 有救了。
func (s *Service) ProjectSnapshot(ctx context.Context, project, reason string) (*ProjectBackupItem, error) {
	projects, err := s.ListProjects(ctx)
	if err != nil {
		return nil, err
	}
	var target *ProjectInfo
	for i := range projects {
		if projects[i].Project == project {
			target = &projects[i]
			break
		}
	}
	if target == nil {
		return nil, fmt.Errorf("没有找到 compose 项目 %s（它可能已经不存在了）", project)
	}
	if len(target.Readable) == 0 {
		return nil, fmt.Errorf("项目 %s 的 compose 文件在 Dockhelm 容器里读不到，"+
			"请先把它的宿主目录挂进来（例如 - /volume1/docker/%s:/host/docker/%s）", project, project, project)
	}
	if reason == "" {
		reason = "manual"
	}
	ts := time.Now().Format("20060102-150405")
	dir := s.projectBackupDir(project, ts)
	for i := 1; ; i++ {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			break
		}
		ts = fmt.Sprintf("%s-%d", ts, i)
		dir = s.projectBackupDir(project, ts)
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}

	man := ProjectManifest{
		Project: project, TS: ts,
		Created: time.Now().UTC().Format(time.RFC3339),
		Reason:  reason, WorkDir: target.WorkDir,
		Files: []ProjectFile{}, Checksums: map[string]string{},
	}
	var total int64
	for _, f := range target.Readable {
		b, err := os.ReadFile(f.Path)
		if err != nil {
			continue // 读不到就跳过，别让一个坏文件毁掉整次备份
		}
		name := filepath.Base(f.Path)
		// 同名文件（不同目录的同名 yaml）加序号，避免互相覆盖
		dest := filepath.Join(dir, name)
		for i := 1; ; i++ {
			if _, err := os.Stat(dest); os.IsNotExist(err) {
				break
			}
			name = fmt.Sprintf("%s.%d%s", strings.TrimSuffix(filepath.Base(f.Path), filepath.Ext(f.Path)), i, filepath.Ext(f.Path))
			dest = filepath.Join(dir, name)
		}
		if err := os.WriteFile(dest, b, 0o600); err != nil {
			return nil, err
		}
		sum := sha256.Sum256(b)
		man.Files = append(man.Files, ProjectFile{
			Name: name, HostPath: f.HostPath, Path: dest, Size: int64(len(b)),
		})
		man.Checksums[name] = hex.EncodeToString(sum[:])
		total += int64(len(b))
	}
	if len(man.Files) == 0 {
		_ = os.RemoveAll(dir)
		return nil, fmt.Errorf("项目 %s 的文件一个都没读成功，未写入备份", project)
	}
	// manifest 里的 Path 是容器内路径，换台机器就失效；真正跨机可靠的是 HostPath。
	// 落盘时清掉 Path，避免以后有人拿它当依据。
	for i := range man.Files {
		man.Files[i].Path = ""
	}
	mb, err := json.MarshalIndent(man, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(dir, manifestName), mb, 0o600); err != nil {
		return nil, err
	}

	item := &ProjectBackupItem{
		Project: project, TS: ts, Created: man.Created,
		Size: total + int64(len(mb)), Files: fileNames(man.Files), Reason: reason,
	}
	names := strings.Join(item.Files, "、")
	if len(item.Files) > 3 {
		names = fmt.Sprintf("%s 等 %d 个文件", strings.Join(item.Files[:3], "、"), len(item.Files))
	}
	s.notify.Emit("backup_success", map[string]string{
		"container": project, "result": "项目配置已备份", "message": names,
	})
	s.st.AddRunLog("backup", project, "success", "项目配置备份 "+ts, names)
	return item, nil
}

// ProjectSnapshotAll 把所有「读得到文件」的项目各备份一份。
func (s *Service) ProjectSnapshotAll(ctx context.Context) ([]ProjectBackupItem, []BatchFailure) {
	items := []ProjectBackupItem{}
	failed := []BatchFailure{}
	projects, err := s.ListProjects(ctx)
	if err != nil {
		return items, []BatchFailure{{Container: "*", Error: err.Error()}}
	}
	for _, p := range projects {
		if len(p.Readable) == 0 {
			failed = append(failed, BatchFailure{Container: p.Project, Error: "文件在容器里看不见，已跳过"})
			continue
		}
		it, err := s.ProjectSnapshot(ctx, p.Project, "manual")
		if err != nil {
			failed = append(failed, BatchFailure{Container: p.Project, Error: err.Error()})
			continue
		}
		items = append(items, *it)
	}
	return items, failed
}

// ListProjectBackups 列出全部项目备份（新 → 旧）。
func (s *Service) ListProjectBackups() ([]ProjectBackupItem, error) {
	root := s.cfg.ProjectBackupDir()
	projDirs, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return []ProjectBackupItem{}, nil
		}
		return nil, err
	}
	out := []ProjectBackupItem{}
	for _, pd := range projDirs {
		if !pd.IsDir() {
			continue
		}
		tsDirs, err := os.ReadDir(filepath.Join(root, pd.Name()))
		if err != nil {
			continue
		}
		for _, td := range tsDirs {
			if !td.IsDir() {
				continue
			}
			dir := filepath.Join(root, pd.Name(), td.Name())
			man, err := readProjectManifest(dir)
			if err != nil {
				continue
			}
			it := ProjectBackupItem{
				Project: man.Project, TS: man.TS, Created: man.Created,
				Files: fileNames(man.Files), Reason: man.Reason,
			}
			if it.Project == "" {
				it.Project = pd.Name()
			}
			// 体积按磁盘实际占用算（含 manifest），文件被外部改动过也不会对不上
			if entries, err := os.ReadDir(dir); err == nil {
				for _, e := range entries {
					if fi, err := e.Info(); err == nil {
						it.Size += fi.Size()
					}
				}
			}
			if it.Created == "" {
				if t, ok := parseSnapshotTS(it.TS); ok {
					it.Created = t.UTC().Format(time.RFC3339)
				}
			}
			out = append(out, it)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Project != out[j].Project {
			return out[i].Project < out[j].Project
		}
		return out[i].TS > out[j].TS
	})
	return out, nil
}

// ProjectRestoreFile 还原结果里的一个文件。
type ProjectRestoreFile struct {
	Name     string `json:"name"`
	HostPath string `json:"hostPath"`
	Written  bool   `json:"written"`
	// Note 失败 / 跳过原因，或「内容与备份一致，未改动」这类说明。
	Note string `json:"note"`
}

// ProjectRestoreResult 项目还原结果。
type ProjectRestoreResult struct {
	Project string              `json:"project"`
	TS      string              `json:"ts"`
	Files   []ProjectRestoreFile `json:"files"`
	OK      bool                `json:"ok"`
	Message string              `json:"message"`
}

// RestoreProject 把一份项目备份里的 yaml / .env 写回它们的原始宿主路径。
//
// 安全护栏与容器快照一致：动手前先给**当前**状态存一份备份（preSnapshot），
// 写错了好退回来。逐文件核对校验值，与备份一致的文件不动。
func (s *Service) RestoreProject(ctx context.Context, project, ts string, preSnapshot bool) *ProjectRestoreResult {
	res := &ProjectRestoreResult{Project: project, TS: ts, Files: []ProjectRestoreFile{}}
	dir := s.projectBackupDir(project, ts)
	man, err := readProjectManifest(dir)
	if err != nil {
		res.Message = "读取备份清单失败：" + err.Error()
		return res
	}
	if preSnapshot {
		if _, err := s.ProjectSnapshot(ctx, project, "restore-pre"); err != nil {
			res.Message = "还原前先给当前状态存一份备份失败了：" + err.Error() + "。为安全起见已中止还原"
			return res
		}
	}
	written := 0
	for _, f := range man.Files {
		entry := ProjectRestoreFile{Name: f.Name, HostPath: f.HostPath}
		src := filepath.Join(dir, f.Name)
		b, err := os.ReadFile(src)
		if err != nil {
			entry.Note = "备份里的文件读不到：" + err.Error()
			res.Files = append(res.Files, entry)
			continue
		}
		local, ok := s.cfg.MapHostPath(f.HostPath)
		if !ok {
			entry.Note = "这个路径在 Dockhelm 容器里看不见，写不回去 —— 请先把宿主目录挂进来"
			res.Files = append(res.Files, entry)
			continue
		}
		if cur, err := os.ReadFile(local); err == nil {
			if sha256.Sum256(cur) == sha256.Sum256(b) {
				entry.Note = "内容与备份一致，未改动"
				entry.Written = true
				res.Files = append(res.Files, entry)
				continue
			}
		}
		if err := os.WriteFile(local, b, 0o644); err != nil {
			entry.Note = "写入失败：" + err.Error()
			res.Files = append(res.Files, entry)
			continue
		}
		entry.Written = true
		entry.Note = "已写回"
		written++
		res.Files = append(res.Files, entry)
	}
	res.OK = written > 0 || len(man.Files) == 0
	// 有文件写不回去（看不见）时不能报「成功」—— 用户以为 yaml 回来了，其实没有
	blocked := 0
	for _, f := range res.Files {
		if !f.Written {
			blocked++
		}
	}
	if blocked > 0 {
		res.OK = false
		res.Message = fmt.Sprintf("还原未完成：%d 个文件没能写回去（原因见下表）", blocked)
	} else if written == 0 {
		res.OK = true
		res.Message = "所有文件的内容都与备份一致，无需改动"
	} else {
		res.OK = true
		res.Message = fmt.Sprintf("已写回 %d 个文件", written)
	}
	s.st.AddRunLog("restore", project, "success", "已用 "+ts+" 还原项目文件", res.Message)
	return res
}

// DeleteProjectBackup 删除一份项目备份。
func (s *Service) DeleteProjectBackup(project, ts string) error {
	dir := s.projectBackupDir(project, ts)
	root := filepath.Clean(s.cfg.ProjectBackupDir())
	if !strings.HasPrefix(filepath.Clean(dir)+string(os.PathSeparator), root+string(os.PathSeparator)) {
		return fmt.Errorf("非法路径")
	}
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	// 顺手清掉空的项目目录
	pdir := filepath.Join(root, sanitize(project))
	if entries, err := os.ReadDir(pdir); err == nil && len(entries) == 0 {
		_ = os.Remove(pdir)
	}
	return nil
}

// ProjectBackupFiles 返回一份备份里的文件清单（供下载 / 还原弹窗展示）。
func (s *Service) ProjectBackupFiles(project, ts string) ([]ProjectFile, error) {
	man, err := readProjectManifest(s.projectBackupDir(project, ts))
	if err != nil {
		return nil, err
	}
	out := []ProjectFile{}
	for _, f := range man.Files {
		p := filepath.Join(s.projectBackupDir(project, ts), f.Name)
		fi, err := os.Stat(p)
		if err != nil {
			continue
		}
		out = append(out, ProjectFile{Name: f.Name, HostPath: f.HostPath, Path: p, Size: fi.Size()})
	}
	return out, nil
}

// WriteProjectBackupZip 把一份项目备份（含 manifest）写成一个 zip 流。
func (s *Service) WriteProjectBackupZip(w io.Writer, project, ts string) error {
	dir := s.projectBackupDir(project, ts)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}
	zw := zip.NewWriter(w)
	defer zw.Close()
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		fw, err := zw.Create(e.Name())
		if err != nil {
			return err
		}
		f, err := os.Open(filepath.Join(dir, e.Name()))
		if err != nil {
			return err
		}
		if _, err := io.Copy(fw, f); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}
	return nil
}

// pruneProjectBackups 按与容器快照同一套策略清理项目备份。
// 返回值：删除份数、释放字节数。opt.KeepPreUpdate 打开时 restore-pre 备份受保护。
func (s *Service) pruneProjectBackups(opt PruneOptions) (int, int64) {
	items, err := s.ListProjectBackups()
	if err != nil {
		return 0, 0
	}
	byProj := map[string][]ProjectBackupItem{}
	for _, it := range items {
		byProj[it.Project] = append(byProj[it.Project], it)
	}
	cutoff := time.Now().AddDate(0, 0, -opt.MaxAgeDays)
	if opt.MaxAgeDays <= 0 {
		cutoff = time.Time{}
	}
	removed := 0
	var freed int64
	for _, list := range byProj {
		sort.Slice(list, func(i, j int) bool { return list[i].TS > list[j].TS })
		kept := 0
		for _, it := range list {
			if opt.KeepPreUpdate && it.Reason == "restore-pre" {
				continue // 还原前备份是「退回上一步」的底牌，与更新前快照同等对待
			}
			drop := false
			if opt.KeepPerContainer > 0 && kept >= opt.KeepPerContainer {
				drop = true
			} else if !cutoff.IsZero() {
				if t, ok := parseSnapshotTS(it.TS); ok && t.Before(cutoff) {
					drop = true
				}
			}
			if !drop {
				kept++
				continue
			}
			if err := s.DeleteProjectBackup(it.Project, it.TS); err == nil {
				removed++
				freed += it.Size
			}
		}
	}
	return removed, freed
}

func readProjectManifest(dir string) (*ProjectManifest, error) {
	b, err := os.ReadFile(filepath.Join(dir, manifestName))
	if err != nil {
		return nil, err
	}
	var man ProjectManifest
	if err := json.Unmarshal(b, &man); err != nil {
		return nil, err
	}
	if man.Files == nil {
		man.Files = []ProjectFile{}
	}
	return &man, nil
}

func fileNames(files []ProjectFile) []string {
	out := []string{}
	for _, f := range files {
		out = append(out, f.Name)
	}
	return out
}
