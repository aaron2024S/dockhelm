// 卷数据打包 / 解包。
//
// 边界（必须在 UI 上如实呈现，不能假装备过）：
//   - named volume 的实体在宿主机的 /var/lib/docker/volumes/<名字>/_data，
//     只有当这个目录被映射进 Dockhelm 时才读得到，读得到才谈得上打包。
//   - bind mount 的数据在宿主机任意目录上，Dockhelm 默认看不见；
//     我们只记录路径，绝不谎称已经备份。
//
// 这里刻意不做「边读边压缩」的花活：卷数据经常是几百 MB 到几 GB，
// 先写 tar 再落盘让失败点清晰（写了一半就是写了一半，不会留下半个 gzip）。
package backup

import (
	"archive/tar"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// VolumeEntry 一个挂载点在快照里的记录。
type VolumeEntry struct {
	Name        string `json:"name"`
	Type        string `json:"type"` // volume | bind | tmpfs | npipe
	Source      string `json:"source"`
	Destination string `json:"destination"`
	// Packed 是否真的把数据打包进了快照。
	Packed bool `json:"packed"`
	// Bytes 打包后的字节数（Packed 为 false 时是 0）。
	Bytes int64 `json:"bytes"`
	// Note 为什么没打包 / 或打包了什么，直接给用户看。
	Note string `json:"note"`
}

// VolumesMeta 快照里 "volumes" 块的内容。
type VolumesMeta struct {
	Included bool          `json:"included"`
	Items    []VolumeEntry `json:"items"`
	// Warning 整体性的警告（例如卷根不可见），没有则为空。
	Warning string `json:"warning,omitempty"`
}

// volumeEntriesFromInspect 从 inspect 的 Mounts 抽出挂载点。
func volumeEntriesFromInspect(insp map[string]any) []VolumeEntry {
	raw, ok := insp["Mounts"].([]any)
	if !ok {
		return []VolumeEntry{}
	}
	out := []VolumeEntry{}
	for _, it := range raw {
		m, ok := it.(map[string]any)
		if !ok {
			continue
		}
		e := VolumeEntry{
			Name:        asString(m["Name"]),
			Type:        asString(m["Type"]),
			Source:      asString(m["Source"]),
			Destination: asString(m["Destination"]),
		}
		if e.Name == "" {
			e.Name = e.Destination
		}
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Destination < out[j].Destination })
	return out
}

// volumeHostDir 返回某个 named volume 在**容器内看得见**的目录。
// 第二个返回值为 false 表示「宿主上也许有，但我们看不见」。
func (s *Service) volumeHostDir(name string) (string, bool) {
	if name == "" {
		return "", false
	}
	local, ok := s.cfg.MapHostPath(filepath.Join("/var/lib/docker/volumes", name, "_data"))
	if !ok || !dirExists(local) {
		return "", false
	}
	return local, true
}

// PackVolumes 把 named volume 的数据打包到 dstDir 下（每个卷一个 <名字>.tar）。
//
// 返回的 VolumesMeta 会原样写进快照，UI 靠它告诉用户「这份快照到底含不含数据」。
func (s *Service) PackVolumes(ctx context.Context, insp map[string]any, dstDir string) VolumesMeta {
	meta := VolumesMeta{Included: true, Items: volumeEntriesFromInspect(insp)}
	if len(meta.Items) == 0 {
		meta.Warning = "该容器没有任何挂载点"
		return meta
	}
	if err := os.MkdirAll(dstDir, 0o700); err != nil {
		meta.Warning = "无法创建卷目录：" + err.Error()
		return meta
	}

	invisible := 0
	for i := range meta.Items {
		e := &meta.Items[i]
		switch e.Type {
		case "volume":
			src, visible := s.volumeHostDir(e.Name)
			if !visible {
				invisible++
				e.Note = "卷数据在宿主机的 /var/lib/docker/volumes 下，当前看不见 —— 快照只记录了它，没有打包"
				continue
			}
			size, err := tarDir(ctx, src, filepath.Join(dstDir, e.Name+".tar"))
			if err != nil {
				e.Note = "打包失败：" + err.Error()
				continue
			}
			e.Packed = true
			e.Bytes = size
			e.Note = "已打包卷数据"
		case "bind":
			e.Note = "绑定挂载的数据在宿主机目录 " + e.Source + "，Dockhelm 看不见 —— 只记录路径，不打包"
		case "tmpfs":
			e.Note = "tmpfs 是内存文件系统，重启即消失，无需备份"
		default:
			e.Note = "不支持的挂载类型，未打包"
		}
	}
	if invisible > 0 {
		meta.Warning = fmt.Sprintf("有 %d 个卷的数据在宿主机上看不见，这部分没有被备份", invisible)
	}
	return meta
}

// UnpackVolumes 把快照目录里的卷 tar 解回对应卷。
//
// 只在 opt.WithVolumes 打开时才调用 —— 覆盖卷数据是不可逆的，
// 所以每一步都写进 steps 里，让用户事后能逐条核对到底动了什么。
func (s *Service) UnpackVolumes(ctx context.Context, container, ts string, step func(string, ...any)) {
	dir := filepath.Join(s.cfg.ContainerBackupDir(), sanitize(container), sanitize(ts)+".volumes")
	entries, err := os.ReadDir(dir)
	if err != nil {
		step("跳过卷数据：这份快照没有打包任何卷（或目录不存在）")
		return
	}

	// 快照里的卷名 → 宿主机目录
	meta := s.loadVolumesMeta(container, ts)
	byName := map[string]VolumeEntry{}
	for _, e := range meta.Items {
		byName[e.Name] = e
	}

	for _, ent := range entries {
		if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".tar") {
			continue
		}
		volName := strings.TrimSuffix(ent.Name(), ".tar")
		// 路径穿越防护：卷名必须是快照里登记过的，且不含分隔符
		if strings.ContainsAny(volName, `/\`) || volName == "." || volName == ".." {
			step("✗ 跳过可疑的卷名 %q", volName)
			continue
		}
		if _, ok := byName[volName]; !ok {
			step("✗ 快照里没有登记 %s，跳过（防止把来路不明的 tar 解进宿主目录）", volName)
			continue
		}
		dst, visible := s.volumeHostDir(volName)
		if !visible {
			step("✗ 卷 %s 当前看不见，无法还原它的数据（需要在宿主机上手工处理）", volName)
			continue
		}
		n, err := untarDir(ctx, filepath.Join(dir, ent.Name()), dst)
		if err != nil {
			step("✗ 卷 %s 还原失败：%v", volName, err)
			continue
		}
		step("✓ 卷 %s 已还原（%d 个条目 → %s）", volName, n, dst)
	}
}

// loadVolumesMeta 读取快照里的 volumes 块。
func (s *Service) loadVolumesMeta(container, ts string) VolumesMeta {
	var out VolumesMeta
	doc, err := readSnapshot(filepath.Join(s.cfg.ContainerBackupDir(), sanitize(container), sanitize(ts)+".json"))
	if err != nil {
		return out
	}
	meta, ok := doc["_dockhelm"].(map[string]any)
	if !ok {
		return out
	}
	b, err := json.Marshal(meta["volumes"])
	if err != nil {
		return out
	}
	_ = json.Unmarshal(b, &out)
	return out
}

// LoadVolumes 暴露给 API：让前端知道某份快照里有没有卷数据。
func (s *Service) LoadVolumes(container, ts string) VolumesMeta {
	return s.loadVolumesMeta(container, ts)
}

// ---------- tar 读写 ----------

// tarDir 把 src 目录打成 dst（tar）。返回 tar 文件大小。
func tarDir(ctx context.Context, src, dst string) (int64, error) {
	f, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	tw := tar.NewWriter(f)
	defer tw.Close()

	root := filepath.Clean(src)
	err = filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		// 用正斜杠写进 tar，跨平台解包才不会错
		name := filepath.ToSlash(rel)

		// 符号链接：只记链接本身，不跟随（跟随会打出一份无限展开的树）
		if info.Mode()&os.ModeSymlink != 0 {
			link, err := os.Readlink(p)
			if err != nil {
				return nil
			}
			return tw.WriteHeader(&tar.Header{
				Typeflag: tar.TypeSymlink, Name: name, Linkname: link, Mode: 0o777,
			})
		}
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		hdr.Name = name
		if info.IsDir() {
			hdr.Name = name + "/"
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		if info.IsDir() || !info.Mode().IsRegular() {
			return nil
		}
		src, err := os.Open(p)
		if err != nil {
			return err
		}
		defer src.Close()
		_, err = io.Copy(tw, src)
		return err
	})
	if err != nil {
		return 0, err
	}
	if err := tw.Close(); err != nil {
		return 0, err
	}
	fi, err := os.Stat(dst)
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

// untarDir 把 tar 解到 dst 目录。返回解出的条目数。
func untarDir(ctx context.Context, tarPath, dst string) (int, error) {
	f, err := os.Open(tarPath)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	tr := tar.NewReader(f)
	root := filepath.Clean(dst)
	n := 0

	for {
		if ctx.Err() != nil {
			return n, ctx.Err()
		}
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return n, err
		}
		target := filepath.Join(root, filepath.FromSlash(hdr.Name))
		// 路径穿越防护：解出来的路径必须还在 root 里
		if !strings.HasPrefix(filepath.Clean(target), root+string(os.PathSeparator)) && filepath.Clean(target) != root {
			return n, fmt.Errorf("压缩包里有越界路径 %q", hdr.Name)
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(hdr.Mode)&0o777); err != nil {
				return n, err
			}
		case tar.TypeSymlink:
			_ = os.Remove(target)
			if err := os.Symlink(hdr.Linkname, target); err != nil {
				return n, err
			}
		default:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return n, err
			}
			out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode)&0o777)
			if err != nil {
				return n, err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return n, err
			}
			out.Close()
		}
		n++
	}
	return n, nil
}
