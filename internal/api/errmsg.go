package api

import (
	"log"
	"strings"
)

// errText 把底层错误翻成一句人话再回给界面，原文只写进容器日志。
//
// 为什么不让 err.Error() 直出：面板上的报错会被用户截图、发到群里。
// `Get "http://127.0.0.1:1/containers/json": dial tcp 127.0.0.1:1:
// connectex: No connection could be made` 这种文本对用户毫无用处，还顺带
// 把内网地址、数据目录、Docker 套接字路径带了出去（界面文案铁律：说明、
// 结果列、toast、确认弹窗都不许出现内部参数）。
//
// 判定刻意保守：**只拦已知的技术错误特征**，其余原样透传 ——
// 业务校验文案（「两次输入的密码不一致」「cron 表达式不合法」「备份不存在」）
// 必须原样保留，那是用户能据此改动作的信息。
func errText(err error) string {
	if err == nil {
		return ""
	}
	raw := err.Error()
	switch {
	case isDockerUnreachable(raw):
		log.Printf("底层错误（Docker 不可达）：%v", err)
		return "连不上 Docker 守护进程。请确认容器挂载了 /var/run/docker.sock，且 Docker 服务正在运行"
	case isTechnical(raw):
		log.Printf("底层错误（已对界面隐藏细节）：%v", err)
		return "操作没有完成，详细信息已记录在容器日志里"
	}
	return raw
}

// isDockerUnreachable 判断是不是「连不上守护进程」这一类。
func isDockerUnreachable(raw string) bool {
	for _, s := range []string{
		"docker.sock",
		"connectex",
		"Cannot connect to the Docker daemon",
		"Is the docker daemon running",
		"no such host",
		"dial tcp",
		"dial unix",
		"connect: no such file or directory",
		"connection refused",
	} {
		if strings.Contains(raw, s) {
			return true
		}
	}
	return false
}

// isTechnical 判断是不是「Go 原文、用户看不懂」的错误。
func isTechnical(raw string) bool {
	for _, s := range []string{
		"dial tcp", "dial unix", "write tcp", "read tcp",
		"i/o timeout", "context deadline exceeded", "context canceled",
		"unexpected EOF", "connection reset", "broken pipe",
		"json: cannot unmarshal", "json: unsupported",
		"invalid character", "unexpected end of JSON",
		"permission denied", "operation not permitted",
		"The system cannot find the file", "The process cannot access",
		"/var/run/docker.sock", "\\\\?\\",
		"http://", "https://",
		"status code", "unexpected status",
		"rename ", "mkdir ", "open ",
	} {
		if strings.Contains(raw, s) {
			return true
		}
	}
	return false
}
