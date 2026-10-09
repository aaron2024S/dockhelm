// Package backup 实现容器配置快照的备份、对比与还原，以及 compose 项目的文件备份。
//
// 两个必须讲清楚的边界（UI 上也要如实呈现）：
//   - 快照 = docker inspect 的结果 = 容器**当前长什么样**，它是「结果」不是「来源」。
//     compose 项目真正的来源是那份 yaml，所以项目备份单独做。
//   - 绑定挂载的数据在**宿主机目录**上，Dockhelm 容器默认看不见它；没挂进来就
//     只能记录路径并明确告诉用户「这份数据请你自行备份」，绝不假装备过。
package backup

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/aaron2024s/dockhelm/internal/config"
	"github.com/aaron2024s/dockhelm/internal/dockerx"
	"github.com/aaron2024s/dockhelm/internal/notify"
	"github.com/aaron2024s/dockhelm/internal/store"
	"github.com/aaron2024s/dockhelm/internal/updater"
)

// Service 备份服务。
type Service struct {
	cfg    *config.Config
	dc     *dockerx.Client
	st     *store.Store
	notify *notify.Manager
}

// New 创建备份服务。
func New(cfg *config.Config, dc *dockerx.Client, st *store.Store, nt *notify.Manager) *Service {
	return &Service{cfg: cfg, dc: dc, st: st, notify: nt}
}

// SnapshotItem 一条快照记录。
type SnapshotItem struct {
	Container string `json:"container"`
	TS        string `json:"ts"`
	Path      string `json:"path"`
	Size      int64  `json:"size"`
	Image     string `json:"image"`
	Running   bool   `json:"running"`
	Created   string `json:"created"`
}

// Snapshot 给某个容器写一份配置快照。
func (s *Service) Snapshot(ctx context.Context, nameOrID, reason string) (*SnapshotItem, error) {
	insp, err := s.dc.Inspect(ctx, nameOrID)
	if err != nil {
		return nil, err
	}
	name := strings.TrimPrefix(asString(insp["Name"]), "/")
	if name == "" {
		name = nameOrID
	}
	dir := filepath.Join(s.cfg.ContainerBackupDir(), sanitize(name))
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	ts := time.Now().Format("20060102-150405")
	// 同一秒重复备份时加后缀，避免互相覆盖
	p := filepath.Join(dir, ts+".json")
	for i := 1; ; i++ {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			break
		}
		p = filepath.Join(dir, fmt.Sprintf("%s-%d.json", ts, i))
	}

	doc := map[string]any{
		"_dockhelm": map[string]any{
			"snapshotAt": time.Now().UTC().Format(time.RFC3339),
			"reason":     reason,
			"version":    1,
		},
		"inspect": insp,
	}
	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}
	// 0600：Env 里可能有数据库密码
	if err := os.WriteFile(p, b, 0o600); err != nil {
		return nil, err
	}
	item := &SnapshotItem{
		Container: name,
		TS:        ts,
		Path:      p,
		Size:      int64(len(b)),
		Image:     imageRef(insp),
		Running:   isRunning(insp),
		Created:   time.Now().UTC().Format(time.RFC3339),
	}
	s.notify.Emit("backup_success", map[string]string{
		"container": name, "result": "配置快照已保存",
		"message": fmt.Sprintf("%s（%.1f KB）", filepath.Base(p), float64(len(b))/1024),
	})
	s.st.AddRunLog("backup", name, "success", "配置快照 "+filepath.Base(p), "")
	return item, nil
}

// List 列出全部快照（按容器名、时间倒序）。
func (s *Service) List() ([]SnapshotItem, error) {
	root := s.cfg.ContainerBackupDir()
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return []SnapshotItem{}, nil
		}
		return nil, err
	}
	out := []SnapshotItem{}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dir := filepath.Join(root, e.Name())
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, f := range files {
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".json") {
				continue
			}
			full := filepath.Join(dir, f.Name())
			info, err := f.Info()
			if err != nil {
				continue
			}
			it := SnapshotItem{
				Container: e.Name(),
				TS:        strings.TrimSuffix(f.Name(), ".json"),
				Path:      full,
				Size:      info.Size(),
				Created:   info.ModTime().UTC().Format(time.RFC3339),
			}
			if doc, err := readSnapshot(full); err == nil {
				if insp, ok := doc["inspect"].(map[string]any); ok {
					it.Image = imageRef(insp)
					it.Running = isRunning(insp)
				}
			}
			out = append(out, it)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Container != out[j].Container {
			return out[i].Container < out[j].Container
		}
		return out[i].TS > out[j].TS
	})
	return out, nil
}

// Delete 删除一份快照。
func (s *Service) Delete(container, ts string) error {
	p := filepath.Join(s.cfg.ContainerBackupDir(), sanitize(container), sanitize(ts)+".json")
	if !strings.HasPrefix(filepath.Clean(p), filepath.Clean(s.cfg.ContainerBackupDir())) {
		return fmt.Errorf("非法路径")
	}
	return os.Remove(p)
}

// DiffEntry 快照与现状的一项差异。
type DiffEntry struct {
	Field    string `json:"field"`
	Snapshot string `json:"snapshot"`
	Current  string `json:"current"`
}

// Diff 对比快照与容器当前状态。
func (s *Service) Diff(ctx context.Context, container, ts string) ([]DiffEntry, error) {
	snap, err := s.loadSnapshotInspect(container, ts)
	if err != nil {
		return nil, err
	}
	cur, err := s.dc.Inspect(ctx, container)
	if err != nil {
		return nil, fmt.Errorf("容器 %s 当前不存在：%w", container, err)
	}
	out := []DiffEntry{}
	add := func(field, a, b string) {
		if strings.TrimSpace(a) != strings.TrimSpace(b) {
			out = append(out, DiffEntry{Field: field, Snapshot: a, Current: b})
		}
	}
	add("镜像", imageRef(snap), imageRef(cur))
	add("运行状态", runWord(isRunning(snap)), runWord(isRunning(cur)))
	add("命令", cmdStr(snap), cmdStr(cur))
	add("入口点", entrypointStr(snap), entrypointStr(cur))
	add("重启策略", restartPolicy(snap), restartPolicy(cur))
	add("网络模式", hostField(snap, "NetworkMode"), hostField(cur, "NetworkMode"))
	add("端口映射", portsStr(snap), portsStr(cur))
	add("挂载", mountsStr(snap), mountsStr(cur))
	add("环境变量", envStr(snap), envStr(cur))
	add("特权模式", boolField(snap, "Privileged"), boolField(cur, "Privileged"))
	add("内存限制", hostNum(snap, "Memory"), hostNum(cur, "Memory"))
	add("CPU 配额", hostNum(snap, "NanoCpus"), hostNum(cur, "NanoCpus"))
	return out, nil
}

// RestoreOptions 还原选项。
type RestoreOptions struct {
	// KeepBackupContainer 是否保留被替换掉的旧容器（默认 true，改成 __bak_ 名字停下）。
	KeepBackupContainer bool
	// PreSnapshot 还原前是否先给当前状态也存一份（安全护栏，默认 true）。
	PreSnapshot bool
}

// RestoreResult 还原结果。
type RestoreResult struct {
	Container string   `json:"container"`
	Snapshot  string   `json:"snapshot"`
	Image     string   `json:"image"`
	Steps     []string `json:"steps"`
	OK        bool     `json:"ok"`
	Message   string   `json:"message"`
}

// Restore 用快照还原一个容器。
//
// 安全护栏：还原前先给**当前**状态也存一份快照，这样改错了还能退回来。
func (s *Service) Restore(ctx context.Context, container, ts string, opt RestoreOptions) *RestoreResult {
	res := &RestoreResult{Container: container, Snapshot: ts}
	step := func(f string, a ...any) { res.Steps = append(res.Steps, fmt.Sprintf(f, a...)) }

	snap, err := s.loadSnapshotInspect(container, ts)
	if err != nil {
		res.Message = "读取快照失败：" + err.Error()
		return res
	}
	res.Image = imageRef(snap)
	step("已载入快照 %s（镜像 %s）", ts, res.Image)

	// 镜像必须已在本地，否则拒绝还原（不在还原流程里偷偷拉镜像）
	if _, err := s.dc.ImageInspect(ctx, res.Image); err != nil {
		res.Message = fmt.Sprintf("镜像 %s 不在本地，请先在「更新中心」拉取后再还原", res.Image)
		step("✗ %s", res.Message)
		return res
	}

	if opt.PreSnapshot {
		if _, err := s.Snapshot(ctx, container, "restore-pre"); err == nil {
			step("已为当前状态写入还原前快照")
		}
	}

	cur, curErr := s.dc.Inspect(ctx, container)
	wasRunning := false
	curID := ""
	if curErr == nil {
		curID = asString(cur["Id"])
		wasRunning = isRunning(cur)
		if wasRunning {
			t := 30
			if err := s.dc.ContainerAction(ctx, curID, "stop", &t); err != nil {
				res.Message = "停止当前容器失败：" + err.Error()
				step("✗ %s", res.Message)
				return res
			}
			step("当前容器已停止")
		}
		bakName := fmt.Sprintf("%s__restorebak_%s", trunc(container, 40), time.Now().Format("20060102-150405"))
		if err := s.dc.RenameContainer(ctx, curID, bakName); err != nil {
			if wasRunning {
				_ = s.dc.ContainerAction(ctx, curID, "start", nil)
			}
			res.Message = "备份当前容器失败：" + err.Error()
			step("✗ %s", res.Message)
			return res
		}
		step("当前容器已改名为 %s", bakName)
	} else {
		step("当前不存在同名容器，将直接创建")
	}

	var nets map[string]bool
	if list, err := s.dc.ListNetworks(ctx); err == nil {
		nets = updater.NetworkNameSet(list)
	}
	cfg, hostCfg, netCfg := updater.BuildCreateSpec(snap, nets)

	newID, err := s.dc.CreateContainer(ctx, container, cfg, hostCfg, netCfg)
	if err != nil {
		step("✗ 创建容器失败：%v", err)
		if curID != "" {
			if rerr := s.dc.RenameContainer(ctx, curID, container); rerr == nil {
				if wasRunning {
					_ = s.dc.ContainerAction(ctx, curID, "start", nil)
				}
				step("已回滚到原容器")
				res.Message = "还原失败（已回滚）：" + err.Error()
				return res
			}
			step("✗ 回滚也失败，请手工把 %s 改回 %s", curID[:12], container)
			res.Message = "还原失败且回滚失败：" + err.Error()
			return res
		}
		res.Message = "还原失败：" + err.Error()
		return res
	}
	step("新容器已创建（%s）", shortID(newID))

	if isRunning(snap) {
		if err := s.dc.ContainerAction(ctx, newID, "start", nil); err != nil {
			step("✗ 启动失败：%v", err)
			res.Message = "还原失败（容器已创建但启动失败）：" + err.Error()
			return res
		}
		step("新容器已启动")
	} else {
		step("快照记录的状态是「已停止」，新容器保持停止")
	}

	res.OK = true
	res.Message = "还原完成"
	step("✓ 还原完成")
	s.st.AddRunLog("restore", container, "success", "已用快照 "+ts+" 还原", strings.Join(res.Steps, "\n"))
	s.notify.Emit("restore_done", map[string]string{
		"container": container, "image": res.Image, "result": "配置已还原",
		"message": "使用快照 " + ts + " 还原（原容器已改名保留）",
	})
	return res
}

// PruneResult 清理结果。
type PruneResult struct {
	Removed    int   `json:"removed"`
	FreedBytes int64 `json:"freedBytes"`
}

// Prune 按「每容器保留 N 份 + 最大保留天数 + 总体积上限」清理旧快照。
// 手工备份不会被自动清理（文件名不含 pre-update 的视为手工）。
func (s *Service) Prune(keepPerContainer, maxAgeDays int, maxTotalMB int) PruneResult {
	res := PruneResult{}
	items, err := s.List()
	if err != nil {
		return res
	}
	byContainer := map[string][]SnapshotItem{}
	for _, it := range items {
		byContainer[it.Container] = append(byContainer[it.Container], it)
	}
	cutoff := time.Now().AddDate(0, 0, -maxAgeDays)
	if maxAgeDays <= 0 {
		cutoff = time.Time{}
	}
	var total int64
	for _, it := range items {
		total += it.Size
	}
	limitBytes := int64(maxTotalMB) * 1024 * 1024
	remove := map[string]bool{}
	for _, list := range byContainer {
		sort.Slice(list, func(i, j int) bool { return list[i].TS > list[j].TS })
		for i, it := range list {
			if keepPerContainer > 0 && i >= keepPerContainer {
				remove[it.Path] = true
			}
			if !cutoff.IsZero() {
				if t, err := time.Parse("20060102-150405", it.TS); err == nil && t.Before(cutoff) {
					remove[it.Path] = true
				}
			}
		}
	}
	for _, it := range items {
		if maxTotalMB > 0 && total > limitBytes && !remove[it.Path] {
			remove[it.Path] = true
			total -= it.Size
		}
	}
	for p := range remove {
		if err := os.Remove(p); err == nil {
			res.Removed++
			if fi, err := os.Stat(p); err == nil {
				res.FreedBytes += fi.Size()
			}
		}
	}
	// 顺手清掉空目录
	for name := range byContainer {
		dir := filepath.Join(s.cfg.ContainerBackupDir(), name)
		if entries, err := os.ReadDir(dir); err == nil && len(entries) == 0 {
			_ = os.Remove(dir)
		}
	}
	return res
}

// Stats 备份占用统计。
type Stats struct {
	Snapshots  int   `json:"snapshots"`
	SizeBytes  int64 `json:"sizeBytes"`
	Containers int   `json:"containers"`
	Dir        string `json:"dir"`
	// VolumeRootMounted 是否挂了 /var/lib/docker/volumes（决定卷数据能否备份）
	VolumeRootMounted bool `json:"volumeRootMounted"`
	// DockerRootVisible 是否能看到 Docker 数据根目录
	DockerRootVisible bool `json:"dockerRootVisible"`
	DockerRoot        string `json:"dockerRoot"`
	// PathMappings 当前生效的「宿主机路径 → 容器内路径」映射，界面上直接展示，
	// 免得用户猜 Dockhelm 到底看见了什么。
	PathMappings []config.PathMapping `json:"pathMappings"`
}

// GetStats 统计备份占用与「能看到什么」。
func (s *Service) GetStats(ctx context.Context) Stats {
	items, _ := s.List()
	st := Stats{Snapshots: len(items), Dir: s.cfg.ContainerBackupDir()}
	set := map[string]bool{}
	for _, it := range items {
		st.SizeBytes += it.Size
		set[it.Container] = true
	}
	st.Containers = len(set)
	st.PathMappings = s.cfg.PathMappings()
	// 具名卷的实体在宿主机的 /var/lib/docker/volumes，能不能读到得走映射判断
	if local, ok := s.cfg.MapHostPath("/var/lib/docker/volumes"); ok {
		st.VolumeRootMounted = dirExists(local)
	}
	if info, err := s.dc.Info(ctx); err == nil {
		st.DockerRoot = info.DockerRootDir
		if local, ok := s.cfg.MapHostPath(info.DockerRootDir); ok {
			st.DockerRootVisible = dirExists(local)
		}
	}
	return st
}

// ---------- compose 项目 ----------

// ProjectInfo 一个 compose 项目。
type ProjectInfo struct {
	Project     string   `json:"project"`
	Containers  []string `json:"containers"`
	ConfigFiles []string `json:"configFiles"` // 宿主视角的路径
	// ReadableFiles 是容器内实际读得到的文件（含解析结果）
	Readable []ProjectFile `json:"readable"`
	// Unreadable 是「知道路径但看不见」的文件，UI 必须显式标出来
	Unreadable []string `json:"unreadable"`
	WorkDir    string   `json:"workDir"`
}

// ProjectFile 一个可读的 compose 文件。
type ProjectFile struct {
	Name     string `json:"name"`
	HostPath string `json:"hostPath"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
}

// ListProjects 按 compose 项目聚合容器。
func (s *Service) ListProjects(ctx context.Context) ([]ProjectInfo, error) {
	list, err := s.dc.ListContainers(ctx)
	if err != nil {
		return nil, err
	}
	byProj := map[string]*ProjectInfo{}
	for _, c := range list {
		proj := c.ComposeProject()
		if proj == "" {
			continue
		}
		p, ok := byProj[proj]
		if !ok {
			p = &ProjectInfo{Project: proj}
			byProj[proj] = p
		}
		p.Containers = append(p.Containers, c.Name())
		if files := c.Labels["com.docker.compose.project.config_files"]; files != "" {
			for _, f := range strings.Split(files, ",") {
				if f = strings.TrimSpace(f); f != "" && !contains(p.ConfigFiles, f) {
					p.ConfigFiles = append(p.ConfigFiles, f)
				}
			}
		}
		if wd := c.Labels["com.docker.compose.project.working_dir"]; wd != "" && p.WorkDir == "" {
			p.WorkDir = wd
		}
	}
	out := make([]ProjectInfo, 0, len(byProj))
	for _, p := range byProj {
		sort.Strings(p.Containers)
		for _, hostPath := range p.ConfigFiles {
			local, ok := s.cfg.MapHostPath(hostPath)
			if !ok || !fileExists(local) {
				p.Unreadable = append(p.Unreadable, hostPath)
				continue
			}
			fi, _ := os.Stat(local)
			p.Readable = append(p.Readable, ProjectFile{
				Name: filepath.Base(local), HostPath: hostPath, Path: local, Size: fi.Size(),
			})
		}
		// .env 与 override 文件通常与 compose 文件同目录，一并探一探
		if p.WorkDir != "" {
			if local, ok := s.cfg.MapHostPath(p.WorkDir); ok {
				for _, cand := range []string{".env", "compose.override.yaml", "compose.override.yml", "docker-compose.override.yml"} {
					full := filepath.Join(local, cand)
					if fileExists(full) && !hasFile(p.Readable, full) {
						fi, _ := os.Stat(full)
						p.Readable = append(p.Readable, ProjectFile{
							Name: cand, HostPath: filepath.Join(p.WorkDir, cand), Path: full, Size: fi.Size(),
						})
					}
				}
			}
		}
		out = append(out, *p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Project < out[j].Project })
	return out, nil
}

// ---------- 内部工具 ----------

func readSnapshot(p string) (map[string]any, error) {
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	var doc map[string]any
	if err := json.Unmarshal(b, &doc); err != nil {
		return nil, err
	}
	return doc, nil
}

func (s *Service) loadSnapshotInspect(container, ts string) (map[string]any, error) {
	p := filepath.Join(s.cfg.ContainerBackupDir(), sanitize(container), sanitize(ts)+".json")
	clean := filepath.Clean(p)
	if !strings.HasPrefix(clean, filepath.Clean(s.cfg.ContainerBackupDir())) {
		return nil, fmt.Errorf("非法路径")
	}
	doc, err := readSnapshot(clean)
	if err != nil {
		return nil, err
	}
	insp, ok := doc["inspect"].(map[string]any)
	if !ok {
		// 兼容「直接存 inspect」的旧格式
		if _, hasID := doc["Id"]; hasID {
			return doc, nil
		}
		return nil, fmt.Errorf("快照格式无法识别")
	}
	return insp, nil
}

// sanitize 只保留文件名安全字符，并拒绝任何路径穿越。
func sanitize(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "\\", "_")
	name = strings.ReplaceAll(name, "..", "_")
	if name == "" || name == "." {
		return "_"
	}
	return name
}

func imageRef(insp map[string]any) string {
	if c, ok := insp["Config"].(map[string]any); ok {
		if s, ok := c["Image"].(string); ok {
			return s
		}
	}
	return ""
}

func isRunning(insp map[string]any) bool {
	if st, ok := insp["State"].(map[string]any); ok {
		b, _ := st["Running"].(bool)
		return b
	}
	return false
}

func asString(v any) string {
	s, _ := v.(string)
	return s
}

func runWord(b bool) string {
	if b {
		return "运行中"
	}
	return "已停止"
}

func hostField(insp map[string]any, key string) string {
	if hc, ok := insp["HostConfig"].(map[string]any); ok {
		return jsonish(hc[key])
	}
	return ""
}

func hostNum(insp map[string]any, key string) string {
	if hc, ok := insp["HostConfig"].(map[string]any); ok {
		return jsonish(hc[key])
	}
	return ""
}

func cmdStr(insp map[string]any) string {
	if c, ok := insp["Config"].(map[string]any); ok {
		return jsonish(c["Cmd"])
	}
	return ""
}

func entrypointStr(insp map[string]any) string {
	if c, ok := insp["Config"].(map[string]any); ok {
		return jsonish(c["Entrypoint"])
	}
	return ""
}

func boolField(insp map[string]any, key string) string {
	if hc, ok := insp["HostConfig"].(map[string]any); ok {
		if b, ok := hc[key].(bool); ok {
			if b {
				return "是"
			}
			return "否"
		}
	}
	return "否"
}

func restartPolicy(insp map[string]any) string {
	if hc, ok := insp["HostConfig"].(map[string]any); ok {
		if rp, ok := hc["RestartPolicy"].(map[string]any); ok {
			n, _ := rp["Name"].(string)
			if n == "" {
				return "no"
			}
			return n
		}
	}
	return ""
}

func portsStr(insp map[string]any) string {
	if hc, ok := insp["HostConfig"].(map[string]any); ok {
		if pb, ok := hc["PortBindings"].(map[string]any); ok && len(pb) > 0 {
			keys := make([]string, 0, len(pb))
			for k := range pb {
				keys = append(keys, k)
			}
			sort.Strings(keys)
			out := make([]string, 0, len(keys))
			for _, k := range keys {
				out = append(out, fmt.Sprintf("%s→%s", k, jsonish(pb[k])))
			}
			return strings.Join(out, ", ")
		}
	}
	return "无"
}

func mountsStr(insp map[string]any) string {
	raw, _ := insp["Mounts"].([]any)
	out := []string{}
	for _, it := range raw {
		m, ok := it.(map[string]any)
		if !ok {
			continue
		}
		src, _ := m["Source"].(string)
		dst, _ := m["Destination"].(string)
		if name, _ := m["Name"].(string); name != "" && m["Type"] == "volume" {
			src = name
		}
		rw := "rw"
		if b, ok := m["RW"].(bool); ok && !b {
			rw = "ro"
		}
		out = append(out, fmt.Sprintf("%s→%s(%s)", src, dst, rw))
	}
	sort.Strings(out)
	if len(out) == 0 {
		return "无"
	}
	return strings.Join(out, ", ")
}

func envStr(insp map[string]any) string {
	c, ok := insp["Config"].(map[string]any)
	if !ok {
		return ""
	}
	raw, _ := c["Env"].([]any)
	out := make([]string, 0, len(raw))
	for _, it := range raw {
		if s, ok := it.(string); ok {
			out = append(out, s)
		}
	}
	sort.Strings(out)
	return strings.Join(out, "\n")
}

func jsonish(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case bool:
		if t {
			return "true"
		}
		return "false"
	case float64:
		return fmt.Sprintf("%v", t)
	default:
		b, err := json.Marshal(t)
		if err != nil {
			return fmt.Sprintf("%v", t)
		}
		return string(b)
	}
}

func dirExists(p string) bool {
	if p == "" {
		return false
	}
	fi, err := os.Stat(p)
	return err == nil && fi.IsDir()
}

func fileExists(p string) bool {
	if p == "" {
		return false
	}
	fi, err := os.Stat(p)
	return err == nil && !fi.IsDir()
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func hasFile(list []ProjectFile, path string) bool {
	for _, v := range list {
		if v.Path == path {
			return true
		}
	}
	return false
}

func trunc(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func shortID(id string) string {
	id = strings.TrimPrefix(id, "sha256:")
	if len(id) > 12 {
		return id[:12]
	}
	return id
}
