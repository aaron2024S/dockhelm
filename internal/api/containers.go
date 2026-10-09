package api

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aaron2024s/dockhelm/internal/dockerx"
	"github.com/aaron2024s/dockhelm/internal/updater"
)

// ContainerView 是给前端的容器视图模型。
type ContainerView struct {
	Name        string            `json:"name"`
	ID          string            `json:"id"`
	ShortID     string            `json:"shortId"`
	Image       string            `json:"image"`
	ImageID     string            `json:"imageId"`
	State       string            `json:"state"`
	Status      string            `json:"status"`
	Created     int64             `json:"created"`
	Ports       []string          `json:"ports"`
	Project     string            `json:"project"`
	Labels      map[string]string `json:"labels"`
	HasUpdate   bool              `json:"hasUpdate"`
	UpdateKnown bool              `json:"updateKnown"`
	Self        bool              `json:"self"`
	Excluded    bool              `json:"excluded"`
	Health      string            `json:"health"`
}

func (s *Server) hListContainers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := s.ctx(r)
	defer cancel()
	list, err := s.dc.ListContainers(ctx)
	if err != nil {
		writeErr(w, http.StatusBadGateway, "读取容器列表失败："+err.Error())
		return
	}

	// 把最新的检测结果合并进来（缓存，不在这里触发网络请求）
	s.checkMu.RLock()
	cache := map[string]updater.CheckResult{}
	for _, c := range s.checkCache {
		cache[c.Container] = c
	}
	checkedAt := s.checkAt
	s.checkMu.RUnlock()

	ex := s.excluded()
	out := make([]ContainerView, 0, len(list))
	for _, c := range list {
		name := c.Name()
		v := ContainerView{
			Name: name, ID: c.ID, ShortID: shortID(c.ID),
			Image: c.Image, ImageID: shortID(c.ImageID),
			State: c.State, Status: c.Status, Created: c.Created,
			Project: c.ComposeProject(), Labels: c.Labels,
			Ports: formatPorts(c.Ports),
		}
		v.Self = s.up.IsSelf(name)
		v.Excluded = ex[name]
		if cr, ok := cache[name]; ok {
			v.UpdateKnown = true
			v.HasUpdate = cr.Status == updater.StatusUpdateAvailable
		}
		out = append(out, v)
	}
	sort.Slice(out, func(i, j int) bool {
		// 有更新的排前面，然后是运行中的
		if out[i].HasUpdate != out[j].HasUpdate {
			return out[i].HasUpdate
		}
		if (out[i].State == "running") != (out[j].State == "running") {
			return out[i].State == "running"
		}
		return out[i].Name < out[j].Name
	})
	writeOK(w, map[string]any{
		"containers": out,
		"checkedAt":  timeOrEmpty(checkedAt),
		"total":      len(out),
	})
}

func (s *Server) hInspectContainer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := s.ctx(r)
	defer cancel()
	name := pathParam(r, "name")
	insp, err := s.dc.Inspect(ctx, name)
	if err != nil {
		writeErr(w, dockerStatus(err), err.Error())
		return
	}
	writeOK(w, map[string]any{"inspect": insp, "summary": summarizeInspect(insp)})
}

// summarizeInspect 抽出一份更适合界面展示的摘要，避免前端去啃原始 inspect。
func summarizeInspect(insp map[string]any) map[string]any {
	cfg, _ := insp["Config"].(map[string]any)
	hc, _ := insp["HostConfig"].(map[string]any)
	st, _ := insp["State"].(map[string]any)
	ns, _ := insp["NetworkSettings"].(map[string]any)

	out := map[string]any{
		"id":         str(insp["Id"]),
		"name":       strings.TrimPrefix(str(insp["Name"]), "/"),
		"created":    str(insp["Created"]),
		"image":      str(cfg["Image"]),
		"imageId":    shortID(str(insp["Image"])),
		"cmd":        strSlice(cfg["Cmd"]),
		"entrypoint": strSlice(cfg["Entrypoint"]),
		"workingDir": str(cfg["WorkingDir"]),
		"user":       str(cfg["User"]),
		"env":        maskEnv(strSlice(cfg["Env"])),
		"labels":     cfg["Labels"],
		"restart":    restartName(hc),
		"privileged": hc["Privileged"],
		"portBindings": hc["PortBindings"],
		"networkMode": str(hc["NetworkMode"]),
		"health":     "",
		"startedAt":  str(st["StartedAt"]),
		"finishedAt": str(st["FinishedAt"]),
		"exitCode":   st["ExitCode"],
		"restarts":   st["RestartCount"],
	}
	if h, ok := st["Health"].(map[string]any); ok {
		out["health"] = str(h["Status"])
	}
	if m := maps[any](hc["Memory"]); m != nil {
		out["memory"] = m
	}
	if c := hc["NanoCpus"]; c != nil && c != float64(0) {
		out["cpus"] = c
	}
	// 挂载：区分「看得见」与「看不见」的数据
	mounts := []map[string]any{}
	if raw, ok := insp["Mounts"].([]any); ok {
		for _, it := range raw {
			m, ok := it.(map[string]any)
			if !ok {
				continue
			}
			entry := map[string]any{
				"type":        str(m["Type"]),
				"source":      str(m["Source"]),
				"destination": str(m["Destination"]),
				"rw":          m["RW"],
			}
			if name := str(m["Name"]); name != "" {
				entry["name"] = name
			}
			if str(m["Type"]) == "bind" {
				entry["note"] = "绑定挂载：数据在宿主机目录上，Dockhelm 只有把它挂进来才看得到"
			} else if str(m["Type"]) == "volume" {
				entry["note"] = "命名卷：数据在 /var/lib/docker/volumes 下"
			}
			mounts = append(mounts, entry)
		}
	}
	out["mounts"] = mounts

	// 网络
	nets := []map[string]any{}
	if ns != nil {
		if nm, ok := ns["Networks"].(map[string]any); ok {
			for k, v := range nm {
				entry := map[string]any{"name": k}
				if m, ok := v.(map[string]any); ok {
					entry["ip"] = str(m["IPAddress"])
					entry["aliases"] = strSlice(m["Aliases"])
					entry["gateway"] = str(m["Gateway"])
					entry["mac"] = str(m["MacAddress"])
				}
				nets = append(nets, entry)
			}
		}
	}
	sort.Slice(nets, func(i, j int) bool { return str(nets[i]["name"]) < str(nets[j]["name"]) })
	out["networks"] = nets
	return out
}

// maskEnv 对疑似敏感的键做打码（界面上仍然可见，但默认不铺开）。
func maskEnv(env []string) []map[string]string {
	out := make([]map[string]string, 0, len(env))
	for _, e := range env {
		k, v, _ := strings.Cut(e, "=")
		sensitive := isSensitiveKey(k)
		item := map[string]string{"key": k, "value": v}
		if sensitive {
			item["sensitive"] = "true"
		}
		out = append(out, item)
	}
	sort.Slice(out, func(i, j int) bool { return out[i]["key"] < out[j]["key"] })
	return out
}

var sensitiveWords = []string{"PASSWORD", "PASSWD", "SECRET", "TOKEN", "APIKEY", "API_KEY", "ACCESS_KEY", "PRIVATE", "CREDENTIAL", "AUTH"}

func isSensitiveKey(k string) bool {
	up := strings.ToUpper(k)
	for _, w := range sensitiveWords {
		if strings.Contains(up, w) {
			return true
		}
	}
	return false
}

func (s *Server) hContainerLogs(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := s.ctx(r)
	defer cancel()
	name := pathParam(r, "name")
	tail := queryInt(r, "tail", 200)
	follow := queryBool(r, "follow", false)

	rc, err := s.dc.Logs(ctx, name, tail, follow, queryBool(r, "timestamps", false))
	if err != nil {
		writeErr(w, dockerStatus(err), err.Error())
		return
	}
	defer rc.Close()

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("X-Accel-Buffering", "no")
	if follow {
		flusher, _ := w.(http.Flusher)
		buf := make([]byte, 32*1024)
		demux := dockerx.DemuxLogs(rc)
		for {
			n, err := demux.Read(buf)
			if n > 0 {
				if _, werr := w.Write(buf[:n]); werr != nil {
					return
				}
				if flusher != nil {
					flusher.Flush()
				}
			}
			if err != nil {
				return
			}
		}
	}
	// 非 follow：整体解码后一次性返回（Docker 的日志是带帧的多路复用流）
	_, _ = io.Copy(w, dockerx.DemuxLogs(rc))
}

func (s *Server) hContainerStats(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := s.ctx(r)
	defer cancel()
	st, err := s.dc.StatsOneShot(ctx, pathParam(r, "name"))
	if err != nil {
		writeErr(w, dockerStatus(err), err.Error())
		return
	}
	writeOK(w, map[string]any{"stats": st, "summary": summarizeStats(st)})
}

// summarizeStats 从 docker stats 的原始结构里算出人看得懂的百分比。
func summarizeStats(st map[string]any) map[string]any {
	out := map[string]any{}
	if cpu := maps[any](st["cpu_stats"]); cpu != nil {
		total := num(cpu["system_cpu_usage"])
		pre := maps[any](st["precpu_stats"])
		preTotal := float64(0)
		if pre != nil {
			preTotal = num(pre["system_cpu_usage"])
		}
		usage := float64(0)
		if u := maps[any](cpu["cpu_usage"]); u != nil {
			usage = num(u["total_usage"])
		}
		preUsage := float64(0)
		if pre != nil {
			if u := maps[any](pre["cpu_usage"]); u != nil {
				preUsage = num(u["total_usage"])
			}
		}
		online := num(cpu["online_cpus"])
		if online == 0 {
			online = 1
		}
		cpuDelta := usage - preUsage
		sysDelta := total - preTotal
		if sysDelta > 0 && cpuDelta > 0 {
			out["cpuPercent"] = cpuDelta / sysDelta * online * 100
		} else {
			out["cpuPercent"] = 0.0
		}
		out["onlineCpus"] = online
	}
	if mem := maps[any](st["memory_stats"]); mem != nil {
		usage := num(mem["usage"])
		limit := num(mem["limit"])
		if cache := num(maps[any](mem["stats"])["inactive_file"]); cache > 0 && usage > cache {
			usage -= cache
		}
		out["memUsage"] = usage
		out["memLimit"] = limit
		if limit > 0 {
			out["memPercent"] = usage / limit * 100
		}
	}
	if nets := maps[any](st["networks"]); nets != nil {
		rx, tx := float64(0), float64(0)
		for _, v := range nets {
			if m, ok := v.(map[string]any); ok {
				rx += num(m["rx_bytes"])
				tx += num(m["tx_bytes"])
			}
		}
		out["netRx"] = rx
		out["netTx"] = tx
	}
	return out
}

type actionReq struct {
	Action  string `json:"action"`
	Timeout int    `json:"timeout"`
}

// hContainerAction 启动 / 停止 / 重启 / 暂停 / 恢复。
func (s *Server) hContainerAction(w http.ResponseWriter, r *http.Request) {
	var in actionReq
	if err := decodeBody(r, &in); err != nil {
		writeErr(w, http.StatusBadRequest, "请求体格式错误")
		return
	}
	name := pathParam(r, "name")
	allowed := map[string]bool{"start": true, "stop": true, "restart": true, "pause": true, "unpause": true, "kill": true}
	if !allowed[in.Action] {
		writeErr(w, http.StatusBadRequest, "不支持的动作："+in.Action)
		return
	}
	ctx, cancel := s.ctx(r)
	defer cancel()
	var timeout *int
	if in.Action != "start" && in.Action != "pause" && in.Action != "unpause" {
		t := in.Timeout
		if t <= 0 {
			t = 30
		}
		timeout = &t
	}
	if in.Action == "kill" {
		timeout = nil
	}
	if err := s.dc.ContainerAction(ctx, name, in.Action, timeout); err != nil {
		writeErr(w, dockerStatus(err), err.Error())
		return
	}
	s.st.AddRunLog("container", name, "success", "执行 "+in.Action, "")
	s.bus.Publish("container", "action", "success", map[string]any{"name": name, "action": in.Action})
	writeOK(w, map[string]any{"ok": true})
}

type renameReq struct {
	NewName string `json:"newName"`
}

// hRenameContainer 重命名容器。
func (s *Server) hRenameContainer(w http.ResponseWriter, r *http.Request) {
	var in renameReq
	if err := decodeBody(r, &in); err != nil || strings.TrimSpace(in.NewName) == "" {
		writeErr(w, http.StatusBadRequest, "新名称不能为空")
		return
	}
	old := pathParam(r, "name")
	if s.up.IsSelf(old) {
		writeErr(w, http.StatusForbidden, "不允许重命名 Dockhelm 自身容器")
		return
	}
	ctx, cancel := s.ctx(r)
	defer cancel()
	if err := s.dc.RenameContainer(ctx, old, in.NewName); err != nil {
		writeErr(w, dockerStatus(err), err.Error())
		return
	}
	s.st.AddRunLog("container", old, "success", "改名为 "+in.NewName, "")
	writeOK(w, map[string]any{"ok": true})
}

// hRemoveContainer 删除容器。volumes=true 会连匿名卷一起删 —— 这是不可逆的。
func (s *Server) hRemoveContainer(w http.ResponseWriter, r *http.Request) {
	name := pathParam(r, "name")
	if s.up.IsSelf(name) {
		writeErr(w, http.StatusForbidden, "不允许删除 Dockhelm 自身容器")
		return
	}
	force := queryBool(r, "force", false)
	vols := queryBool(r, "volumes", false)

	ctx, cancel := s.ctx(r)
	defer cancel()

	// 删之前先看清楚「会连带删掉哪些卷」，并把它们回给前端展示
	warn := []map[string]any{}
	if insp, err := s.dc.Inspect(ctx, name); err == nil {
		if raw, ok := insp["Mounts"].([]any); ok {
			for _, it := range raw {
				m, ok := it.(map[string]any)
				if !ok {
					continue
				}
				if str(m["Type"]) != "volume" {
					continue
				}
				nm := str(m["Name"])
				warn = append(warn, map[string]any{
					"name": nm,
					"dest": str(m["Destination"]),
					"anonymousDetected": strings.HasPrefix(str(m["Source"]), "/var/lib/docker/volumes/"),
				})
			}
		}
	}

	if err := s.dc.RemoveContainer(ctx, name, force, vols); err != nil {
		writeErr(w, dockerStatus(err), err.Error())
		return
	}
	msg := "已删除容器"
	if vols {
		msg += "（同时删除了匿名卷）"
	}
	s.st.AddRunLog("container", name, "success", msg, "")
	writeOK(w, map[string]any{"ok": true, "message": msg, "volumes": warn})
}

// hExportContainer 导出容器配置（就是 docker inspect 的 JSON），供用户自己存档。
func (s *Server) hExportContainer(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := s.ctx(r)
	defer cancel()
	name := pathParam(r, "name")
	insp, err := s.dc.Inspect(ctx, name)
	if err != nil {
		writeErr(w, dockerStatus(err), err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.json"`, safeFileName(name)))
	writeOK(w, insp)
}

// ---------- 镜像 / 网络 / 卷 ----------

func (s *Server) hListImages(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := s.ctx(r)
	defer cancel()
	list, err := s.dc.ListImages(ctx, queryBool(r, "all", false))
	if err != nil {
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	type imageView struct {
		ID         string   `json:"id"`
		ShortID    string   `json:"shortId"`
		Tags       []string `json:"tags"`
		Digests    []string `json:"digests"`
		Size       int64    `json:"size"`
		Created    int64    `json:"created"`
		Containers int64    `json:"containers"`
		Dangling   bool     `json:"dangling"`
	}
	out := make([]imageView, 0, len(list))
	var total int64
	for _, im := range list {
		tags := []string{}
		for _, t := range im.RepoTags {
			if t != "<none>:<none>" {
				tags = append(tags, t)
			}
		}
		short := func(ds []string) []string {
			o := []string{}
			for _, d := range ds {
				if i := strings.LastIndex(d, "@"); i >= 0 {
					d = d[i+1:]
				}
				o = append(o, shortID(d))
			}
			return o
		}
		out = append(out, imageView{
			ID: im.ID, ShortID: shortID(im.ID), Tags: tags, Digests: short(im.RepoDigests),
			Size: im.Size, Created: im.Created, Containers: im.Containers,
			Dangling: len(tags) == 0,
		})
		total += im.Size
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Created > out[j].Created })
	writeOK(w, map[string]any{"images": out, "totalSize": total, "count": len(out)})
}

func (s *Server) hRemoveImage(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := s.ctx(r)
	defer cancel()
	id := pathParam(r, "id")
	if err := s.dc.RemoveImage(ctx, id, queryBool(r, "force", false), queryBool(r, "noprune", false)); err != nil {
		writeErr(w, dockerStatus(err), err.Error())
		return
	}
	s.st.AddRunLog("image", shortID(id), "success", "已删除镜像", "")
	writeOK(w, map[string]any{"ok": true})
}

func (s *Server) hPruneImages(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := s.ctx(r)
	defer cancel()
	freed, err := s.dc.PruneImages(ctx)
	if err != nil {
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	s.st.AddRunLog("image", "-", "success", fmt.Sprintf("清理悬空镜像释放 %.1f MB", float64(freed)/1024/1024), "")
	writeOK(w, map[string]any{"ok": true, "freedBytes": freed})
}

func (s *Server) hListNetworks(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := s.ctx(r)
	defer cancel()
	list, err := s.dc.ListNetworks(ctx)
	if err != nil {
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	writeOK(w, map[string]any{"networks": list})
}

func (s *Server) hListVolumes(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := s.ctx(r)
	defer cancel()
	list, err := s.dc.ListVolumes(ctx)
	if err != nil {
		writeErr(w, http.StatusBadGateway, err.Error())
		return
	}
	writeOK(w, map[string]any{"volumes": list})
}

func (s *Server) hRemoveVolume(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := s.ctx(r)
	defer cancel()
	name := pathParam(r, "name")
	if err := s.dc.RemoveVolume(ctx, name, queryBool(r, "force", false)); err != nil {
		writeErr(w, dockerStatus(err), err.Error())
		return
	}
	s.st.AddRunLog("volume", name, "success", "已删除卷", "")
	writeOK(w, map[string]any{"ok": true})
}

// ---------- 助手 ----------

func formatPorts(ports []struct {
	IP          string `json:"IP"`
	PrivatePort uint16 `json:"PrivatePort"`
	PublicPort  uint16 `json:"PublicPort"`
	Type        string `json:"Type"`
}) []string {
	out := []string{}
	for _, p := range ports {
		if p.PublicPort > 0 {
			host := ""
			if p.IP != "" && p.IP != "0.0.0.0" && p.IP != "::" {
				host = p.IP + ":"
			}
			out = append(out, fmt.Sprintf("%s%d→%d/%s", host, p.PublicPort, p.PrivatePort, p.Type))
		} else {
			out = append(out, fmt.Sprintf("%d/%s", p.PrivatePort, p.Type))
		}
	}
	sort.Strings(out)
	return out
}

func restartName(hc map[string]any) string {
	if hc == nil {
		return ""
	}
	if rp, ok := hc["RestartPolicy"].(map[string]any); ok {
		if n := str(rp["Name"]); n != "" {
			return n
		}
	}
	return "no"
}

func str(v any) string {
	s, _ := v.(string)
	return s
}

func strSlice(v any) []string {
	raw, ok := v.([]any)
	if !ok {
		if ss, ok := v.([]string); ok {
			return ss
		}
		return []string{}
	}
	out := make([]string, 0, len(raw))
	for _, it := range raw {
		out = append(out, str(it))
	}
	return out
}

func num(v any) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case float32:
		return float64(t)
	case int64:
		return float64(t)
	case int:
		return float64(t)
	case uint64:
		return float64(t)
	}
	return 0
}

// maps 是一个泛型小工具，把 any 断言成 map[string]T。
func maps[T any](v any) map[string]T {
	m, _ := v.(map[string]T)
	return m
}

func shortID(id string) string {
	id = strings.TrimPrefix(id, "sha256:")
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

func safeFileName(s string) string {
	s = strings.ReplaceAll(s, "/", "_")
	s = strings.ReplaceAll(s, "\\", "_")
	s = strings.ReplaceAll(s, "..", "_")
	if s == "" {
		return "container"
	}
	return s
}

func timeOrEmpty(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}

func dockerStatus(err error) int {
	var se *dockerx.StatusError
	if ok := asStatus(err, &se); ok {
		switch se.Code {
		case http.StatusNotFound:
			return http.StatusNotFound
		case http.StatusConflict:
			return http.StatusConflict
		}
	}
	return http.StatusBadGateway
}

func asStatus(err error, target **dockerx.StatusError) bool {
	se, ok := err.(*dockerx.StatusError)
	if ok {
		*target = se
	}
	return ok
}

// parseID 把路径里的字符串转成 int64。
func parseID(s string) (int64, error) { return strconv.ParseInt(s, 10, 64) }
