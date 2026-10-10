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
	// Reason 快照的由来：manual / pre_update / scheduled。
	// 「更新前快照永不自动清理」这条保留策略靠它识别。
	Reason string `json:"reason"`
	// WithData 这份快照里是否真的打包了卷数据。
	WithData bool `json:"withData"`
}

// SnapshotOptions 备份选项。
type SnapshotOptions struct {
	Reason string
	// WithVolumes 同时把 named volume 的数据打包进来。
	// 只有卷根目录被映射进 Dockhelm 时才真的打得到，打不到会在快照里写明原因。
	WithVolumes bool
}

// Snapshot 给某个容器写一份配置快照（不含卷数据）。
func (s *Service) Snapshot(ctx context.Context, nameOrID, reason string) (*SnapshotItem, error) {
	return s.SnapshotWith(ctx, nameOrID, SnapshotOptions{Reason: reason})
}

// SnapshotWith 按选项写一份快照。
func (s *Service) SnapshotWith(ctx context.Context, nameOrID string, opt SnapshotOptions) (*SnapshotItem, error) {
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
	// 磁盘上的真实标识（可能带 -1 后缀）。json 文件名、卷目录名、返回的 TS
	// 三者必须同源 —— 以前卷目录固定用 ts，于是同一秒的第二份快照会把第一份的
	// 卷 tar 整个覆盖掉（真正丢数据），而且返回的 TS 还指向不存在的文件名。
	slot := strings.TrimSuffix(filepath.Base(p), ".json")

	// 卷数据：先打包再写 json，这样 json 里记的 Packed/Bytes 一定是真的
	var volumes VolumesMeta
	if opt.WithVolumes {
		volumes = s.PackVolumes(ctx, insp, filepath.Join(dir, slot+".volumes"))
	}

	meta := map[string]any{
		"snapshotAt": time.Now().UTC().Format(time.RFC3339),
		"reason":     opt.Reason,
		"version":    1,
	}
	if opt.WithVolumes {
		meta["volumes"] = volumes
	}
	doc := map[string]any{"_dockhelm": meta, "inspect": insp}
	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}
	// 0600：Env 里可能有数据库密码
	if err := os.WriteFile(p, b, 0o600); err != nil {
		return nil, err
	}
	total := int64(len(b))
	for _, e := range volumes.Items {
		total += e.Bytes
	}
	item := &SnapshotItem{
		Container: name,
		TS:        slot,
		Path:      p,
		Size:      total,
		Image:     imageRef(insp),
		Running:   isRunning(insp),
		Created:   time.Now().UTC().Format(time.RFC3339),
		Reason:    opt.Reason,
		WithData:  len(volumes.Items) > 0 && volumes.Items[0].Packed,
	}
	for _, e := range volumes.Items {
		if e.Packed {
			item.WithData = true
			break
		}
	}
	msg := fmt.Sprintf("%s（%.1f KB）", filepath.Base(p), float64(total)/1024)
	if opt.WithVolumes {
		packed := 0
		for _, e := range volumes.Items {
			if e.Packed {
				packed++
			}
		}
		msg += fmt.Sprintf(" · 含 %d 个卷的数据", packed)
	}
	s.notify.Emit("backup_success", map[string]string{
		"container": name, "result": "配置快照已保存", "message": msg,
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
			// 快照自带的时间比文件 mtime 可靠：数据目录一旦被复制 / 同步 / 从压缩包
			// 里解出来，mtime 就变成「复制那一刻」了，而快照里记的才是真正拍下来的时刻。
			// 优先级：_dockhelm.snapshotAt（RFC3339，权威）> 文件名里的时间戳（本地时区）> mtime。
			if t, ok := parseSnapshotTS(it.TS); ok {
				it.Created = t.UTC().Format(time.RFC3339)
			}
			if doc, err := readSnapshot(full); err == nil {
				if insp, ok := doc["inspect"].(map[string]any); ok {
					it.Image = imageRef(insp)
					it.Running = isRunning(insp)
				}
				if meta, ok := doc["_dockhelm"].(map[string]any); ok {
					if r, ok := meta["reason"].(string); ok {
						it.Reason = r
					}
					if sa, ok := meta["snapshotAt"].(string); ok {
						if t, err := time.Parse(time.RFC3339, sa); err == nil {
							it.Created = t.UTC().Format(time.RFC3339)
						}
					}
					// 卷数据块存在且有任意一项 Packed ⇒ 这份快照含数据
					if vm, ok := meta["volumes"].(map[string]any); ok {
						if items, ok := vm["items"].([]any); ok {
							for _, x := range items {
								if m, ok := x.(map[string]any); ok {
									if packed, _ := m["packed"].(bool); packed {
										it.WithData = true
										break
									}
								}
							}
						}
					}
				}
			}
			// 卷 tar 的体积要计进快照总体积，否则「占用空间」会明显偏小
			if vdir := filepath.Join(dir, it.TS+".volumes"); dirExists(vdir) {
				if entries, err := os.ReadDir(vdir); err == nil {
					for _, ve := range entries {
						if fi, err := ve.Info(); err == nil {
							it.Size += fi.Size()
						}
					}
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

// parseSnapshotTS 把快照文件名里的时间戳（20060102-150405）解析成时间。
// 写快照用的是 time.Now().Format(...)，即**本地时区**，这里必须按本地解析。
func parseSnapshotTS(ts string) (time.Time, bool) {
	t, err := time.ParseInLocation("20060102-150405", ts, time.Local)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}

// Delete 删除一份快照（连同它的卷数据一起删）。
func (s *Service) Delete(container, ts string) error {
	p := filepath.Join(s.cfg.ContainerBackupDir(), sanitize(container), sanitize(ts)+".json")
	if !strings.HasPrefix(filepath.Clean(p), filepath.Clean(s.cfg.ContainerBackupDir())) {
		return fmt.Errorf("非法路径")
	}
	if err := os.Remove(p); err != nil {
		return err
	}
	// 卷数据目录：删了 json 却留着几 GB 的 tar 是最讨厌的一种「已删除」
	_ = os.RemoveAll(filepath.Join(s.cfg.ContainerBackupDir(), sanitize(container), sanitize(ts)+".volumes"))
	return nil
}

// Import 导入一份外部快照（用户手工上传的 json）。
//
// 只做最小校验：必须是合法 JSON、且带 inspect 块 —— 没有 inspect 的文件
// 既不能展示也不能还原，收进来只会变成垃圾。容器名与时间戳由调用方给，
// 时间戳撞车时自动加后缀，绝不覆盖已有快照。
func (s *Service) Import(container, ts, content string) (*SnapshotItem, error) {
	container = strings.TrimSpace(container)
	ts = strings.TrimSpace(ts)
	if container == "" || ts == "" {
		return nil, fmt.Errorf("需要指定容器名与快照时间戳")
	}
	if strings.ContainsAny(container, `/\`) {
		return nil, fmt.Errorf("名字里不能包含路径分隔符")
	}
	// 先判时间戳格式再判 ts 里的分隔符：用户把 "2026/10/09" 粘进来时，
	// 「时间戳格式应为 20060102-150405」比「不能包含路径分隔符」有用得多。
	if _, err := time.Parse("20060102-150405", ts); err != nil {
		return nil, fmt.Errorf("时间戳格式应为 20060102-150405（例如 20261009-094200）")
	}
	if strings.ContainsAny(ts, `/\`) {
		return nil, fmt.Errorf("时间戳里不能包含路径分隔符")
	}
	var doc map[string]any
	if err := json.Unmarshal([]byte(content), &doc); err != nil {
		return nil, fmt.Errorf("不是合法的 JSON：%w", err)
	}
	insp, ok := doc["inspect"].(map[string]any)
	if !ok || len(insp) == 0 {
		return nil, fmt.Errorf("文件里没有 inspect 内容，无法作为快照使用")
	}

	dir := filepath.Join(s.cfg.ContainerBackupDir(), sanitize(container))
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	p := filepath.Join(dir, sanitize(ts)+".json")
	for i := 1; ; i++ {
		if _, err := os.Stat(p); os.IsNotExist(err) {
			break
		}
		if i > 999 {
			return nil, fmt.Errorf("同名快照太多了，换个时间戳")
		}
		p = filepath.Join(dir, fmt.Sprintf("%s-%d.json", sanitize(ts), i))
	}
	// 落盘前重新序列化：保持与其它快照一致的缩进，也顺手挡掉不可控的额外字段。
	//
	// _dockhelm 块存在就沿用文件自带的信息（我们不改写别人的元数据），
	// 缺失才补一份。下面的 effectiveReason 与返回值必须一致 ——
	// 否则「接口说你导入了，列表说这份是定时快照」这种自相矛盾迟早会被当成 bug 报上来。
	effectiveReason := "imported"
	if meta, ok := doc["_dockhelm"].(map[string]any); ok && meta != nil {
		if r, ok := meta["reason"].(string); ok && r != "" {
			effectiveReason = r
		}
	} else {
		doc["_dockhelm"] = map[string]any{
			"snapshotAt": time.Now().UTC().Format(time.RFC3339),
			"reason":     effectiveReason,
			"version":    1,
		}
	}
	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}
	if err := os.WriteFile(p, b, 0o600); err != nil {
		return nil, err
	}
	s.st.AddRunLog("backup", container, "success", "已导入外部快照 "+filepath.Base(p), "")
	return &SnapshotItem{
		Container: container,
		TS:        strings.TrimSuffix(filepath.Base(p), ".json"),
		Path:      p,
		Size:      int64(len(b)),
		Image:     imageRef(insp),
		Running:   isRunning(insp),
		Created:   time.Now().UTC().Format(time.RFC3339),
		Reason:    effectiveReason,
	}, nil
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
	// WithVolumes 连卷数据一起还原（覆盖现有文件）。
	// 默认 false —— 覆盖数据不可逆，必须是用户显式勾选的结果。
	WithVolumes bool
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

	// 创建请求体先备好并校验 —— **必须在停当前容器之前**。
	// 否则「快照里少了 Config」这种事会变成「容器停了、改了名，创建失败，再回滚」。
	var nets map[string]bool
	if list, err := s.dc.ListNetworks(ctx); err == nil {
		nets = updater.NetworkNameSet(list)
	}
	cfg, hostCfg, netCfg := updater.BuildCreateSpec(snap, nets)
	if patched, fail := updater.EnsureCreateSpec(cfg, snap); fail != "" {
		res.Message = "快照里" + fail + "，拒绝还原，容器保持原样"
		step("✗ %s", res.Message)
		return res
	} else if patched != "" {
		step("快照没有记录镜像引用，改用镜像 ID %s", shortID(patched))
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

	// 卷数据在「新容器已建好、但还没启动」的窗口里还原 —— 此刻没有任何进程
	// 在读写这份卷，覆盖它才是安全的。
	if opt.WithVolumes {
		s.UnpackVolumes(ctx, container, ts, step)
	}

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

// PruneOptions 快照清理策略。三个上限互相独立，任一为 0 表示该维度不限制。
type PruneOptions struct {
	// KeepPerContainer 每个容器保留最近 N 份。
	KeepPerContainer int
	// MaxAgeDays 超过 N 天的快照删掉。
	MaxAgeDays int
	// MaxTotalMB 全部快照的总体积上限（MB），超出时从最旧的开始删。
	MaxTotalMB int
	// KeepPreUpdate 更新前写的快照永不自动清理。
	// 它对准的是最贵的场景：自动更新把容器搞坏之后，唯一能救回来的就是
	// 更新前那一刻的配置 —— 而它恰恰是最容易被「保留最近 10 份」挤掉的。
	KeepPreUpdate bool
}

// Prune 按策略清理旧快照。
func (s *Service) Prune(opt PruneOptions) PruneResult {
	res := PruneResult{}
	items, err := s.List()
	if err != nil {
		return res
	}
	byContainer := map[string][]SnapshotItem{}
	for _, it := range items {
		byContainer[it.Container] = append(byContainer[it.Container], it)
	}
	cutoff := time.Now().AddDate(0, 0, -opt.MaxAgeDays)
	if opt.MaxAgeDays <= 0 {
		cutoff = time.Time{}
	}
	var total int64
	for _, it := range items {
		total += it.Size
	}
	limitBytes := int64(opt.MaxTotalMB) * 1024 * 1024
	// 值存整条记录，不只是「要不要删」——删除时要连带清掉 <ts>.volumes 目录，
	// 并按整份快照（含卷数据）的体积计释放量，光有路径算不出来。
	remove := map[string]SnapshotItem{}

	// protected 判定：更新前快照在开了开关时不参与任何维度的清理。
	protected := func(it SnapshotItem) bool {
		return opt.KeepPreUpdate && it.Reason == "pre_update"
	}

	for _, list := range byContainer {
		sort.Slice(list, func(i, j int) bool { return list[i].TS > list[j].TS })
		kept := 0
		for _, it := range list {
			if protected(it) {
				continue // 受保护的快照不占「保留份数」的额度，也不被天数清理
			}
			if opt.KeepPerContainer > 0 && kept >= opt.KeepPerContainer {
				remove[it.Path] = it
				continue
			}
			if !cutoff.IsZero() {
				// 快照文件名是**本地时区**的 20060102-150405（写的时候用的 time.Now()），
				// 必须用 parseSnapshotTS 按本地解析。以前这里用 time.Parse（UTC），
				// 东八区下等于把每份快照都当成「晚了 8 小时」，保留期变相被拉长。
				if t, ok := parseSnapshotTS(it.TS); ok && t.Before(cutoff) {
					remove[it.Path] = it
					continue
				}
			}
			kept++
		}
	}
	// 总量超限时从最旧的开始删，但同样跳过受保护的。
	sorted := make([]SnapshotItem, len(items))
	copy(sorted, items)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].TS < sorted[j].TS })
	for _, it := range sorted {
		if opt.MaxTotalMB > 0 && total > limitBytes && !protected(it) {
			if _, already := remove[it.Path]; !already {
				remove[it.Path] = it
				total -= it.Size
			}
		}
	}
	for p, it := range remove {
		if _, err := os.Stat(p); err != nil {
			continue
		}
		if err := os.Remove(p); err != nil {
			continue
		}
		res.Removed++
		// it.Size 已经含卷数据的体积（List 会把 <ts>.volumes 里的文件算进去），
		// 这里不要再单独累加一次，否则释放量会翻倍。
		res.FreedBytes += it.Size
		// 卷数据目录必须一起删。Delete() 一直是这么做的，Prune 以前漏了 ——
		// 结果是几 GB 的卷 tar 变成孤儿：列表里看不见（json 没了）、
		// 也不会被任何一次清理碰到，磁盘只涨不落。
		_ = os.RemoveAll(filepath.Join(filepath.Dir(p), it.TS+".volumes"))
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

// newProjectInfo 建一个空项目。
//
// 切片**必须显式初始化成空切片**：Go 的 nil 切片会序列化成 JSON `null`，
// 而前端拿到 null 再取 `.length` 会抛 TypeError、整页渲染白屏
// （0.3.0 的真实事故：有 compose 项目但一个文件都读不到时，备份页一片空白）。
// 别把这里改成 `&ProjectInfo{Project: proj}`。
func newProjectInfo(proj string) *ProjectInfo {
	return &ProjectInfo{
		Project:     proj,
		Containers:  []string{},
		ConfigFiles: []string{},
		Readable:    []ProjectFile{},
		Unreadable:  []string{},
	}
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
			p = newProjectInfo(proj)
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
