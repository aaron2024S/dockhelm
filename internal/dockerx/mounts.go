package dockerx

// SelfMountMap 从容器 inspect 结果里提取「宿主机路径 → 容器内路径」映射。
//
// 只用顶层的 Mounts —— 它是守护进程解析之后的最终结果（Binds、VolumesFrom、匿名卷都已展开），
// 比 HostConfig.Binds 之类原始声明确凿。有了它，Dockhelm 才能知道自己容器里
// 「宿主机那个目录」叫什么名字，于是 compose 里的挂载不必被迫写成冒号两边一致。
func SelfMountMap(insp map[string]any) map[string]string {
	out := map[string]string{}
	raw, ok := insp["Mounts"].([]any)
	if !ok {
		return out
	}
	for _, it := range raw {
		m, ok := it.(map[string]any)
		if !ok {
			continue
		}
		src, _ := m["Source"].(string)
		dst, _ := m["Destination"].(string)
		if src == "" || dst == "" {
			continue
		}
		out[src] = dst
	}
	return out
}
