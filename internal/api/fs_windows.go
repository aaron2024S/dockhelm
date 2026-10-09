//go:build windows

package api

// statfs 在 Windows 上没有直接对应实现；宿主机磁盘信息对本面板没有意义
// （Dockhelm 实际跑在 Linux 容器里），这里返回 0 让 UI 隐藏该字段。
func statfs(string) (free uint64, total uint64) { return 0, 0 }
