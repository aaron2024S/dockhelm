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

// HumanBytes 是运行记录 / 通知 / 计划任务共用的体积说法，四舍五入口径要一致。
func TestHumanBytes(t *testing.T) {
	cases := []struct {
		n    int64
		want string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1024 * 1024, "1.0 MB"},
		{3 * 1024 * 1024 * 1024, "3.0 GB"},
	}
	for _, c := range cases {
		if got := HumanBytes(c.n); got != c.want {
			t.Errorf("HumanBytes(%d) = %q, want %q", c.n, got, c.want)
		}
	}
}

// PruneImageResult.Summary 是手动清理与计划任务清理共用的那一句。
// 关键：清理对象是悬空镜像（按定义没有名字），所以形状是「几个 + 多大」，
// 不点名；只有 Docker 回了 Untagged 条目（带名字）才点。
func TestPruneImageSummary(t *testing.T) {
	if got := (PruneImageResult{}).Summary(); got != "没有可清理的未使用镜像" {
		t.Errorf("空结果应给一句人话，got %q", got)
	}

	// 典型的悬空镜像清理：只有 Deleted（裸 ID），不该出现在文案里。
	res := PruneImageResult{
		Deleted: []string{"sha256:2c75e866d5a5", "sha256:aaaa11112222"},
		Freed:   1536 * 1024 * 1024,
	}
	got := res.Summary()
	want := "清理未使用镜像 2 个，释放 1.5 GB"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if strings.Contains(got, "2c75e866") {
		t.Errorf("裸 ID 不该出现在文案里：%q", got)
	}
	if res.Count() != 2 {
		t.Errorf("Count = %d, want 2", res.Count())
	}

	// 带 Untagged（顺手摘掉的标签）时，它也计入数量 —— 但名字不出现在摘要里，
	// 名字走 Names（运行记录的对象行），这是 2026-10-13 定下的分工。
	res2 := PruneImageResult{
		Deleted:  []string{"sha256:aaa"},
		Untagged: []string{"nginx:1.26", "redis:7.2"},
		Freed:    1024,
	}
	if got2 := res2.Summary(); strings.Contains(got2, "nginx:1.26") || !strings.Contains(got2, "3 个") {
		t.Errorf("摘要只报数量与体积，got %q", got2)
	}
}

// PruneImageResult.Names 输出「这次删了哪些镜像」—— 运行记录的对象行。
//
// 用户 2026-10-13：「原 ID 行就要显示出删除了哪些镜像」。所以这里必须给名字，
// 且绝不出现裸 ID：查不到名字的也要说清它是什么 —— 无标签的顶层镜像说成
// 「无标签镜像 <短ID>」，prune 混进 Deleted 的中间层 ID 说成「中间层 <短ID>」
// （用户 2026-10-13 追问「镜像页能读到名字，怎么记录里不行」——因为那是层，
// 镜像页根本不显示，Docker 也没给过它任何名字）。
func TestPruneImageNames(t *testing.T) {
	res := PruneImageResult{
		// 三个删掉的对象：一个查得到名字，一个是无标签的顶层镜像，
		// 一个不在清理前的镜像列表里（随镜像一起删掉的中间层）。
		Deleted:  []string{"sha256:2c75e866d5a5aaaaaaaa", "sha256:beef12345678bbbb", "sha256:c0ffee0000001111"},
		Untagged: []string{"nginx:1.26"},
	}
	idx := PruneNameIndex{
		Names:    map[string]string{"2c75e866d5a5": "aaron2024s/dockhelm@2c75e866d5a5"},
		ImageIDs: map[string]bool{"2c75e866d5a5": true, "beef12345678": true},
		Digests:  map[string]string{},
	}

	got := res.Names(idx)
	want := []string{
		"aaron2024s/dockhelm@2c75e866d5a5",
		"nginx:1.26",
		"中间层 c0ffee000000",
		"无标签镜像 beef12345678",
	}
	if len(got) != len(want) {
		t.Fatalf("名单条数 = %d（%v），want %d", len(got), got, len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("名单[%d] = %q, want %q（已排序）", i, got[i], want[i])
		}
	}
	for _, n := range got {
		if n == "2c75e866d5a5" || n == "beef12345678" || n == "c0ffee000000" {
			t.Errorf("名单里不该出现裸 ID：%q", n)
		}
	}

	// Untagged 里的摘要引用：归一成 仓库@短摘要；指向的镜像已在 Deleted
	// 里点过名时去重，不让同一个镜像占两行。
	res2 := PruneImageResult{
		Deleted:  []string{"sha256:2c75e866d5a5aaaaaaaa"},
		Untagged: []string{"aaron2024s/dockhelm@sha256:2c75e866d5a5ffffffffffffffffffffffffffffffff", "redis@sha256:fedcba987654ffffffffffffffffffffffffffffffff"},
	}
	idx2 := PruneNameIndex{
		Names:    map[string]string{"2c75e866d5a5": "aaron2024s/dockhelm@2c75e866d5a5"},
		ImageIDs: map[string]bool{"2c75e866d5a5": true},
		Digests:  map[string]string{"2c75e866d5a5": "2c75e866d5a5", "fedcba987654": "fedcba987654"},
	}
	got2 := res2.Names(idx2)
	want2 := []string{
		"aaron2024s/dockhelm@2c75e866d5a5",
		"redis@fedcba987654",
	}
	if len(got2) != len(want2) {
		t.Fatalf("去重后名单条数 = %d（%v），want %d", len(got2), got2, len(want2))
	}
	for i := range want2 {
		if got2[i] != want2[i] {
			t.Errorf("名单[%d] = %q, want %q（已排序）", i, got2[i], want2[i])
		}
	}

	// 对照表整体缺失（清理前那次列举失败）也不能崩，且仍不出现裸 ID。
	for _, n := range (PruneImageResult{Deleted: []string{"sha256:feedface0000"}}).Names(PruneNameIndex{}) {
		if !strings.HasPrefix(n, "无标签镜像 ") {
			t.Errorf("取不到对照时应兜底成「无标签镜像 <短ID>」，got %q", n)
		}
	}
}

// FoldNames 是一行版的名单：超出上限必须折叠并写明总数，否则会把记录撑爆。
func TestFoldNames(t *testing.T) {
	if got := FoldNames(nil, 3); got != "" {
		t.Errorf("空名单应为空串，got %q", got)
	}
	if got := FoldNames([]string{"a", "b"}, 3); got != "a、b" {
		t.Errorf("未超上限应全列，got %q", got)
	}
	got := FoldNames([]string{"a", "b", "c", "d", "e"}, 3)
	if got != "a、b、c…等 5 个" {
		t.Errorf("超上限应折叠并写明总数，got %q", got)
	}
}

// ImageDisplayName 是「清理前抓名字」的取名规则：有标签用标签，
// 没标签用 仓库@短摘要，匿名层返回空串（由调用方兜底）。
func TestImageDisplayName(t *testing.T) {
	cases := []struct {
		name string
		im   ImageSummary
		want string
	}{
		{"有标签", ImageSummary{RepoTags: []string{"nginx:1.27"}, RepoDigests: []string{"x@sha256:aa"}}, "nginx:1.27"},
		{"跳过 none 标签", ImageSummary{RepoTags: []string{"<none>:<none>", "redis:7.2"}}, "redis:7.2"},
		{"无标签但有仓库摘要", ImageSummary{RepoDigests: []string{"aaron2024s/dockhelm@sha256:2c75e866d5a5ffffffffffffffffffffffffffffffffffffffffffffffff"}},
			"aaron2024s/dockhelm@2c75e866d5a5"},
		{"匿名层", ImageSummary{}, ""},
	}
	for _, c := range cases {
		if got := ImageDisplayName(c.im); got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}
