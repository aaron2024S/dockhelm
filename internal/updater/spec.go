package updater

import (
	"maps"
	"strings"

	"github.com/aaron2024s/dockhelm/internal/dockerx"
)

// BuildCreateSpec 把一次 docker inspect 的结果转换成 POST /containers/create 的请求体。
//
// 这是「忠实克隆一个容器」的全部难点所在，逐条说明为什么这么做：
//
//  1. Config 与 HostConfig 整体拷贝，只删掉少数几个「读取时才存在、写回会出错」的字段
//     （Binds / Mounts / ContainerIDFile / MacAddress / Shell），其余原样保留 ——
//     环境变量、端口、重启策略、资源限制、capabilities、日志驱动、健康检查全都不需要枚举。
//  2. 挂载**必须从顶层 Mounts 重建**。inspect 会同时返回 HostConfig.Binds 与
//     HostConfig.Mounts，两者指向同一批挂载，一起回传会被守护进程判为「重复挂载点」；
//     而顶层 Mounts 是**解析后的真实结果**，命名卷、匿名卷、只读、传播模式都在里面。
//  3. 匿名卷要按 **Name** 复用（而不是让它重新生成）。这顺手修掉了经典的
//     「更新容器后匿名卷数据变孤儿」问题 —— 新容器接回同一个卷。
//  4. Hostname 只有在它等于容器短 ID（也就是从没被手工设置过）时才删掉，
//     这样既不会继承一个陈旧的 ID 当主机名，又不会丢掉用户显式设置的 hostname。
//  5. 网络别名与静态 IP 从 NetworkSettings.Networks 还原；但如果该网络已经被删掉，
//     就不要再传，否则创建会 404。host / none / container:* 网络模式不传 EndpointsConfig。
//
// existingNetworks 为 nil 表示「不校验网络是否存在」。
func BuildCreateSpec(insp map[string]any, existingNetworks map[string]bool) (config, hostConfig, networking map[string]any) {
	config = mapCopy(insp["Config"])
	hostConfig = mapCopy(insp["HostConfig"])

	// ---- Config 清理 ----
	if config != nil {
		if h, _ := config["Hostname"].(string); h != "" && strings.HasPrefix(inspectID(insp), h) {
			delete(config, "Hostname")
		}
		delete(config, "MacAddress") // 已废弃，回传会被忽略甚至报错
		delete(config, "Shell")      // Windows 专用
	}

	// ---- HostConfig 清理 ----
	if hostConfig != nil {
		delete(hostConfig, "Binds")
		delete(hostConfig, "Mounts")
		delete(hostConfig, "ContainerIDFile")
		// NetworkMode 保留：它决定了容器用哪个网络（也可能是 container:<id> / host / none）
	}
	mounts := buildMounts(insp)
	if hostConfig == nil {
		hostConfig = map[string]any{}
	}
	if len(mounts) > 0 {
		hostConfig["Mounts"] = mounts
	}

	// ---- 网络 ----
	networking = buildNetworking(insp, hostConfig, existingNetworks)
	return config, hostConfig, networking
}

// buildMounts 从顶层 Mounts 重建挂载列表。
func buildMounts(insp map[string]any) []map[string]any {
	raw, _ := insp["Mounts"].([]any)
	out := []map[string]any{}
	for _, item := range raw {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		typ, _ := m["Type"].(string)
		target, _ := m["Destination"].(string)
		if target == "" {
			continue
		}
		// tmpfs 不持久化，交给 HostConfig.Tmpfs（原样保留）处理
		if typ == "tmpfs" {
			continue
		}
		entry := map[string]any{"Target": target}
		switch typ {
		case "volume":
			name, _ := m["Name"].(string)
			if name == "" {
				// 没有 Name 说明是匿名卷且拿不到卷名（极端情况）——
				// 这种情况无法安全复用，跳过比悄悄换一个卷更安全
				continue
			}
			entry["Type"] = "volume"
			entry["Source"] = name
			entry["VolumeOptions"] = map[string]any{}
		case "bind":
			src, _ := m["Source"].(string)
			if src == "" {
				continue
			}
			entry["Type"] = "bind"
			entry["Source"] = src
			opts := map[string]any{}
			if p, _ := m["Propagation"].(string); p != "" && p != "rprivate" {
				opts["Propagation"] = p
			}
			if len(opts) > 0 {
				entry["BindOptions"] = opts
			}
		default:
			continue
		}
		if rw, ok := m["RW"].(bool); ok && !rw {
			entry["ReadOnly"] = true
		}
		out = append(out, entry)
	}
	return out
}

// buildNetworking 还原网络别名与静态 IP。
func buildNetworking(insp, hostConfig map[string]any, existing map[string]bool) map[string]any {
	mode := ""
	if hostConfig != nil {
		mode, _ = hostConfig["NetworkMode"].(string)
	}
	if mode == "host" || mode == "none" || strings.HasPrefix(mode, "container:") {
		return nil
	}
	ns, _ := insp["NetworkSettings"].(map[string]any)
	if ns == nil {
		return nil
	}
	nets, _ := ns["Networks"].(map[string]any)
	if len(nets) == 0 {
		return nil
	}
	endpoints := map[string]any{}
	for name, v := range nets {
		if existing != nil && !existing[name] {
			continue // 网络已被删除，传了会 404
		}
		nm, ok := v.(map[string]any)
		if !ok {
			continue
		}
		ep := map[string]any{}
		if aliases, ok := nm["Aliases"].([]any); ok && len(aliases) > 0 {
			ep["Aliases"] = aliases
		}
		if ipam, ok := nm["IPAMConfig"].(map[string]any); ok && len(ipam) > 0 {
			// 只有用户显式指定过静态 IP 时这里才有内容
			ep["IPAMConfig"] = ipam
		}
		if mac, ok := nm["MacAddress"].(string); ok && mac != "" {
			ep["MacAddress"] = mac
		}
		if len(ep) > 0 {
			endpoints[name] = ep
		}
	}
	if len(endpoints) == 0 {
		return nil
	}
	return map[string]any{"EndpointsConfig": endpoints}
}

// mapCopy 浅拷贝一个 map[string]any（nil 安全）。
func mapCopy(v any) map[string]any {
	src, ok := v.(map[string]any)
	if !ok || src == nil {
		return nil
	}
	dst := make(map[string]any, len(src))
	maps.Copy(dst, src)
	return dst
}

// NetworkNameSet 把网络列表转成集合（供 BuildCreateSpec 过滤已删除的网络）。
func NetworkNameSet(nets []map[string]any) map[string]bool {
	out := map[string]bool{}
	for _, n := range nets {
		if s, ok := n["Name"].(string); ok && s != "" {
			out[s] = true
		}
	}
	return out
}

// ImageRefOf 从容器 inspect 里取镜像引用（导出给其他包用）。
func ImageRefOf(insp map[string]any) string { return dockerx.ParseRef(inspectImageRef(insp)).String() }
