package autoupdate

import (
	"reflect"
	"testing"
)

func TestPlan(t *testing.T) {
	containers := []Container{
		{Name: "qbittorrent", Image: "lscr.io/linuxserver/qbittorrent:latest", HasUpdate: true, Running: true},
		{Name: "jellyfin", Image: "jellyfin/jellyfin:10.9", HasUpdate: true, Running: true},
		{Name: "redis", Image: "redis:alpine", HasUpdate: true, Running: true},
		{Name: "postgres", Image: "postgres:16", HasUpdate: true, Running: true},
		{Name: "nginx", Image: "nginx:alpine", HasUpdate: false, Running: true},
		{Name: "dockhelm", Image: "aaron2024s/dockhelm:0.1.0", HasUpdate: true, Running: true},
	}
	opt := Options{SelfName: "dockhelm", Exclude: []string{"redis*", "postgres"}}

	got := Plan(containers, opt)

	want := []Candidate{
		{Name: "qbittorrent", Image: "lscr.io/linuxserver/qbittorrent:latest", HasUpdate: true, Running: true, WillUpdate: true},
		{Name: "jellyfin", Image: "jellyfin/jellyfin:10.9", HasUpdate: true, Running: true, WillUpdate: true},
		{Name: "redis", Image: "redis:alpine", HasUpdate: true, Running: true, Excluded: true, Reason: "命中排除列表"},
		{Name: "postgres", Image: "postgres:16", HasUpdate: true, Running: true, Excluded: true, Reason: "命中排除列表"},
		{Name: "nginx", Image: "nginx:alpine", HasUpdate: false, Running: true, Reason: "未检测到更新"},
		{Name: "dockhelm", Image: "aaron2024s/dockhelm:0.1.0", HasUpdate: true, Running: true, Protected: true, Reason: "Dockhelm 自身，永不自动更新"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Plan 结果不符\n got=%+v\nwant=%+v", got, want)
	}

	if names := Names(got); !reflect.DeepEqual(names, []string{"qbittorrent", "jellyfin"}) {
		t.Fatalf("Names 应为 [qbittorrent jellyfin]，实际 %v", names)
	}
	if n := Count(got); n != 2 {
		t.Fatalf("Count 应为 2，实际 %d", n)
	}
}

// 自身保护必须优先于排除列表：用户即使手动把自身加进排除列表之外，
// 也不该让它被选中 —— 更新到一半把自己停了是最难排查的故障。
func TestPlanSelfProtectionBeatsExclude(t *testing.T) {
	got := Plan([]Container{{Name: "dockhelm", Image: "x/dockhelm:1", HasUpdate: true}}, Options{})
	if len(got) != 1 || !got[0].Protected || got[0].WillUpdate {
		t.Fatalf("自身容器必须被保护，实际 %+v", got)
	}
}

func TestPlanSelfMatchesByKeyword(t *testing.T) {
	// 改名之后仍然认得出来
	got := Plan([]Container{{Name: "my-dockhelm-panel", HasUpdate: true}}, Options{SelfName: "whatever"})
	if !got[0].Protected {
		t.Fatalf("容器名含 dockhelm 应被保护，实际 %+v", got[0])
	}
}

func TestPlanEmptyExcludeAllowsAll(t *testing.T) {
	cands := Plan([]Container{
		{Name: "a", Image: "a:1", HasUpdate: true},
		{Name: "b", Image: "b:1", HasUpdate: false},
	}, Options{})
	if Count(cands) != 1 || cands[0].Name != "a" {
		t.Fatalf("无排除列表时应只挑有更新的，实际 %+v", cands)
	}
}

func TestMatchesExclude(t *testing.T) {
	cases := []struct {
		name     string
		exclude  []string
		cont     string
		image    string
		expected bool
	}{
		{"空列表不排除", nil, "redis", "redis:alpine", false},
		{"容器名精确", []string{"redis"}, "redis", "redis:alpine", true},
		{"容器名通配", []string{"redis*"}, "redis-exporter", "x/y:1", true},
		{"大小写不敏感", []string{"REDIS"}, "redis", "redis:alpine", true},
		{"镜像全名", []string{"ghcr.io/foo/bar:latest"}, "bar", "ghcr.io/foo/bar:latest", true},
		{"镜像仓库名", []string{"ghcr.io/foo/bar"}, "bar", "ghcr.io/foo/bar:1.2", true},
		{"带端口号的 registry", []string{"registry.local:5000/nginx"}, "nginx", "registry.local:5000/nginx:1.27", true},
		{"不命中", []string{"redis"}, "nginx", "nginx:alpine", false},
		{"非法通配表达式退化成精确匹配", []string{"[bad"}, "redis", "redis:alpine", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := MatchesExclude(c.exclude, c.cont, c.image); got != c.expected {
				t.Fatalf("MatchesExclude(%v, %q, %q) = %v，期望 %v", c.exclude, c.cont, c.image, got, c.expected)
			}
		})
	}
}

func TestStripTag(t *testing.T) {
	cases := map[string]string{
		"redis:alpine":                  "redis",
		"nginx":                         "nginx",
		"lscr.io/linuxserver/plex:1.2":  "lscr.io/linuxserver/plex",
		"registry.local:5000/nginx:1.2": "registry.local:5000/nginx",
		"registry.local:5000/nginx":     "registry.local:5000/nginx",
	}
	for in, want := range cases {
		if got := StripTag(in); got != want {
			t.Fatalf("StripTag(%q) = %q，期望 %q", in, got, want)
		}
	}
}

func TestIsSelf(t *testing.T) {
	if !IsSelf("dockhelm", "dockhelm") {
		t.Fatal("精确同名应判为自身")
	}
	if !IsSelf("Dockhelm", "") {
		t.Fatal("关键字匹配不区分大小写")
	}
	if IsSelf("", "dockhelm") {
		t.Fatal("空名字不是自身")
	}
	if IsSelf("nginx", "dockhelm") {
		t.Fatal("nginx 不该被判为自身")
	}
}
