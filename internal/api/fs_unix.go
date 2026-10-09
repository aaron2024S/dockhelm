//go:build linux || darwin || freebsd || openbsd || netbsd

package api

import "syscall"

// statfs 返回给定路径所在文件系统的 (可用字节, 总字节)。
// 只在 UI 上展示「数据盘还剩多少」，失败时返回 0 即可。
func statfs(path string) (free uint64, total uint64) {
	var st syscall.Statfs_t
	if err := syscall.Statfs(path, &st); err != nil {
		return 0, 0
	}
	bs := uint64(st.Bsize)
	return uint64(st.Bavail) * bs, uint64(st.Blocks) * bs
}
