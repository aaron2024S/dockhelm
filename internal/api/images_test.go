package api

import "testing"

// inspectImageName 决定「已删除镜像 xxx」里的 xxx：
// 有 tag 用 tag，没 tag 用仓库名+短摘要，都没有退回短 ID。
func TestInspectImageName(t *testing.T) {
	cases := []struct {
		name string
		insp map[string]any
		want string
	}{
		{
			name: "有 tag 用第一个有效 tag",
			insp: map[string]any{
				"RepoTags":    []any{"nginx:1.27", "nginx:latest"},
				"RepoDigests": []any{},
			},
			want: "nginx:1.27",
		},
		{
			name: "无 tag 镜像用仓库名+短摘要",
			insp: map[string]any{
				"RepoTags":    []any{"<none>:<none>"},
				"RepoDigests": []any{"aaron2024s/dockhelm@sha256:2c75e866d5a5aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
			},
			want: "aaron2024s/dockhelm@2c75e866d5a5",
		},
		{
			name: "docker.io 官方镜像去掉前缀",
			insp: map[string]any{
				"RepoTags":    []any{},
				"RepoDigests": []any{"docker.io/library/redis@sha256:abcdef1234567890ffffffffffffffffffffffffffffffffffffffffffffff"},
			},
			want: "redis@abcdef123456",
		},
		{
			name: "全空退回空串由调用方兜底",
			insp: map[string]any{
				"RepoTags":    []any{},
				"RepoDigests": []any{},
			},
			want: "",
		},
		{
			name: "字段缺失也返回空串",
			insp: map[string]any{},
			want: "",
		},
	}
	for _, c := range cases {
		if got := inspectImageName(c.insp); got != c.want {
			t.Errorf("%s: got %q, want %q", c.name, got, c.want)
		}
	}
}
