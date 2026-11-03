package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestSecurityHeaders 全站安全头。这几行是「面板能被别的站点 iframe 嵌进去」
// 与「浏览器嗅探 Content-Type」这两类问题的唯一防线。
func TestSecurityHeaders(t *testing.T) {
	s := &Server{}
	h := s.securityMW(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))

	want := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"Referrer-Policy":        "no-referrer",
	}
	for k, v := range want {
		if got := rec.Header().Get(k); got != v {
			t.Errorf("%s = %q，期望 %q", k, got, v)
		}
	}
	csp := rec.Header().Get("Content-Security-Policy")
	for _, need := range []string{"default-src 'self'", "frame-ancestors 'none'", "object-src 'none'"} {
		if !strings.Contains(csp, need) {
			t.Errorf("CSP 缺少 %q：%s", need, csp)
		}
	}
	// 脚本必须严格要求 'self'：一旦放行 'unsafe-inline'，这层就等于没有
	if strings.Contains(csp, "script-src 'self' 'unsafe-inline'") || strings.Contains(csp, "script-src 'unsafe-inline'") {
		t.Errorf("CSP 不应允许内联脚本：%s", csp)
	}
}

// TestAboutRequiresAuth 部署细节（数据目录 / Docker 地址 / 监听端口）只给登录用户看。
func TestAboutRequiresAuth(t *testing.T) {
	s := &Server{mux: http.NewServeMux()}
	s.routes()

	rec := httptest.NewRecorder()
	s.mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/about", nil))
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("未登录访问 /api/about 应 401，实际 %d：%s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "dataDir") {
		t.Fatal("未登录响应里出现了部署细节")
	}

	// 必须确认它确实挂在了「需要登录」的分支上，而不是压根没注册（404 也算错）
	if rec.Code == http.StatusNotFound {
		t.Fatal("/api/about 路由没有注册")
	}
}

// TestErrTextHidesInternals 界面文案铁律：不出现接口路径、字段名、连接串、异常原文。
func TestErrTextHidesInternals(t *testing.T) {
	dockerErr := errors.New(`Get "http://127.0.0.1:1/containers/json": dial tcp 127.0.0.1:1: connectex: No connection could be made`)
	got := errText(dockerErr)
	if strings.Contains(got, "127.0.0.1") || strings.Contains(got, "containers/json") {
		t.Fatalf("Docker 不可达时泄露了内部参数：%q", got)
	}
	if !strings.Contains(got, "Docker") {
		t.Fatalf("应当告诉用户是 Docker 连不上，实际：%q", got)
	}

	ioErr := errors.New("open D:\\data\\dockhelm.db: The process cannot access the file because it is being used by another process")
	got = errText(ioErr)
	if strings.Contains(got, "dockhelm.db") || strings.Contains(got, "D:\\") {
		t.Fatalf("文件类错误泄露了路径：%q", got)
	}

	// 业务文案必须原样透传 —— 用户要靠它决定下一步怎么做
	for _, msg := range []string{"两次输入的密码不一致", "cron 表达式不合法：预期 5 段", "备份不存在"} {
		if got := errText(errors.New(msg)); got != msg {
			t.Errorf("业务文案被改写：%q → %q", msg, got)
		}
	}
	if errText(nil) != "" {
		t.Error("nil error 应返回空串")
	}
}
