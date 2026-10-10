package dockerx

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
)

// newDistributionDaemon 起一个「按 tag 给摘要」的假守护进程。
//
// 真实守护进程的 /distribution/{name}/json 里，name 是一条**完整引用**：
// 不带 tag 就等于 :latest。所以这个假实现的摘要由「仓库:tag」决定 ——
// 一旦调用方把 tag 吞掉，拿回来的就是 :latest 的摘要，测试立刻失败。
//
// 2026-10-10 真机事故：DistributionInspect 只传了仓库名，redis:alpine 的本地摘要
// 被拿去跟 redis:latest 的远端摘要比，于是永久误报「有更新」；点更新又永远
// 「镜像未变化」（更新的拉取走的是正确引用）。当时假环境里的假守护进程
// 完全按仓库名匹配、根本不看 tag，一路绿灯 —— 所以这条假实现必须按 tag 区分。
func newDistributionDaemon(t *testing.T, seen *[]string) *httptest.Server {
	t.Helper()
	var mu sync.Mutex
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const prefix, suffix = "/distribution/", "/json"
		if !strings.HasPrefix(r.URL.Path, prefix) || !strings.HasSuffix(r.URL.Path, suffix) {
			http.NotFound(w, r)
			return
		}
		raw := r.URL.Path[len(prefix) : len(r.URL.Path)-len(suffix)]
		if dec, err := url.PathUnescape(raw); err == nil {
			raw = dec
		}
		mu.Lock()
		if seen != nil {
			*seen = append(*seen, raw)
		}
		mu.Unlock()

		// 模拟注册表：name 不带 tag 就按 :latest 解析
		name := raw
		if !strings.Contains(name, "@") && !strings.HasSuffix(name, ":latest") {
			if i := strings.LastIndex(name, ":"); i < 0 || strings.Contains(name[i+1:], "/") {
				name += ":latest"
			}
		}
		sum := sha256.Sum256([]byte(name))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"Descriptor": map[string]any{
				"mediaType": "application/vnd.oci.image.index.v1+json",
				"digest":    "sha256:" + hex.EncodeToString(sum[:]),
			},
			"Platforms": []map[string]any{{"architecture": "amd64", "os": "linux"}},
		})
	}))
}

func digestOf(s string) string {
	sum := sha256.Sum256([]byte(s))
	return "sha256:" + hex.EncodeToString(sum[:])
}

// TestDistributionInspectKeepsTag 是本仓库最该有、却一直缺的一条断言：
// 查远端摘要时必须把 tag（或 digest）带上，不能退回 :latest。
func TestDistributionInspectKeepsTag(t *testing.T) {
	var seen []string
	srv := newDistributionDaemon(t, &seen)
	defer srv.Close()
	c := clientFor(t, srv)

	cases := []struct {
		ref      string
		wantName string // 假守护进程应当收到的完整引用
	}{
		// 真机事故里的两个：非 latest 标签
		{"redis:alpine", "redis:alpine"},
		{"redis:7.4-alpine", "redis:7.4-alpine"},
		{"gdy666/lucky:v2", "gdy666/lucky:v2"},
		{"diygod/rsshub:chromium-bundled", "diygod/rsshub:chromium-bundled"},
		// 非 docker.io 的仓库
		{"lscr.io/linuxserver/msedge:latest", "lscr.io/linuxserver/msedge:latest"},
		{"ghcr.io/music-assistant/server:v1.2.3", "ghcr.io/music-assistant/server:v1.2.3"},
		// 不带 tag ⇒ 就是 latest（与 docker pull 的语义一致）
		{"redis", "redis:latest"},
		// 按 digest 钉住的引用，tag 位置必须是 digest
		{"redis@sha256:" + strings.Repeat("ab", 32), "redis@sha256:" + strings.Repeat("ab", 32)},
	}

	for _, tc := range cases {
		seen = seen[:0]
		got, err := c.DistributionInspect(context.Background(), tc.ref)
		if err != nil {
			t.Fatalf("%s：查远端摘要失败：%v", tc.ref, err)
		}
		if len(seen) != 1 {
			t.Fatalf("%s：期望恰好 1 次 /distribution 请求，实际 %d 次（%v）", tc.ref, len(seen), seen)
		}
		if seen[0] != tc.wantName {
			t.Errorf("%s：守护进程收到的引用是 %q，期望 %q（tag 被丢了？）", tc.ref, seen[0], tc.wantName)
		}
		if want := digestOf(tc.wantName); got != want {
			t.Errorf("%s：摘要 %s，期望 %s", tc.ref, got, want)
		}
	}
}

// TestDistributionInspectDoesNotCompareAgainstLatest 把事故场景钉死：
// 本地是 redis:alpine，远端给的必须是 alpine 的摘要，绝不能是 latest 的。
func TestDistributionInspectDoesNotCompareAgainstLatest(t *testing.T) {
	srv := newDistributionDaemon(t, nil)
	defer srv.Close()
	c := clientFor(t, srv)

	got, err := c.DistributionInspect(context.Background(), "redis:alpine")
	if err != nil {
		t.Fatalf("查远端摘要失败：%v", err)
	}
	if got == digestOf("redis:latest") {
		t.Fatalf("拿回了 :latest 的摘要 %s —— 非 latest 标签又被当成 latest 了", got)
	}
	if got != digestOf("redis:alpine") {
		t.Fatalf("摘要 %s，期望 %s", got, digestOf("redis:alpine"))
	}
}

// TestDistributionInspectSurfacesDaemonError 守护进程报错时必须返回错误，
// 调用方才能把容器标成「未知」而不是「有新版本」。
func TestDistributionInspectSurfacesDaemonError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"message":"manifest unknown"}`)
	}))
	defer srv.Close()
	c := clientFor(t, srv)

	if _, err := c.DistributionInspect(context.Background(), "redis:alpine"); err == nil {
		t.Fatal("守护进程 404 时必须返回错误，不能返回空摘要当成功")
	}
}
