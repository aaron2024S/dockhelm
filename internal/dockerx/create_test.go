package dockerx

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newCreateDaemon 起一个「照真实守护进程的规则解析请求体」的假守护进程。
//
// 规则来自 moby：请求体被解进 types.ContainerCreateConfig，其中
//
//	type ContainerCreateConfig struct {
//		Name string
//		*container.Config   // ← 内嵌指针，JSON 里被摊平
//		HostConfig       *HostConfig
//		NetworkingConfig *NetworkConfig
//	}
//
// 内嵌指针只有在顶层出现它的字段时才会被分配。所以：
//   - 顶层有 Image ⇒ 正常创建；
//   - 把配置塞进 "Config" 子对象 ⇒ 守护进程当作「没给配置」，
//     回 400 "config cannot be empty in order to create a container"。
//
// 2026-10-09 就是栽在这里：本仓库此前发的正是嵌套形状，而自带的假守护进程
// 压根不读请求体，于是一路绿灯、到真机上才炸。这个假守护进程的存在意义
// 就是让那种形状再也过不去。
func newCreateDaemon(t *testing.T, sink *map[string]any) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/containers/create") {
			http.NotFound(w, r)
			return
		}
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			t.Errorf("读请求体失败：%v", err)
		}
		var body map[string]any
		if err := json.Unmarshal(raw, &body); err != nil {
			t.Errorf("请求体不是合法 JSON：%v（原文 %s）", err, raw)
		}
		if sink != nil {
			*sink = body
		}
		w.Header().Set("Content-Type", "application/json")
		if img, _ := body["Image"].(string); img == "" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = io.WriteString(w, `{"message":"config cannot be empty in order to create a container"}`)
			return
		}
		_, _ = io.WriteString(w, `{"Id":"`+strings.Repeat("a", 64)+`","Warnings":[]}`)
	}))
}

func clientFor(t *testing.T, srv *httptest.Server) *Client {
	t.Helper()
	c, err := New("tcp://" + strings.TrimPrefix(srv.URL, "http://"))
	if err != nil {
		t.Fatalf("构造客户端失败：%v", err)
	}
	return c
}

// TestCreateContainerSendsFlatBody 是本仓库最该有、却一直缺的一条断言：
// 创建容器的请求体必须是**摊平**的（配置字段与 HostConfig 同级），
// 顶层不能有 "Config" 键。
func TestCreateContainerSendsFlatBody(t *testing.T) {
	var got map[string]any
	srv := newCreateDaemon(t, &got)
	defer srv.Close()
	c := clientFor(t, srv)

	id, err := c.CreateContainer(context.Background(), "nginx",
		map[string]any{
			"Image": "nginx:1.27-alpine",
			"Env":   []any{"TZ=Asia/Shanghai"},
			"Cmd":   []any{"nginx", "-g", "daemon off;"},
		},
		map[string]any{"NetworkMode": "bridge"},
		map[string]any{"EndpointsConfig": map[string]any{"bridge": map[string]any{}}},
	)
	if err != nil {
		t.Fatalf("创建容器失败（真实守护进程正是这样拒绝嵌套形状的）：%v", err)
	}
	if id != strings.Repeat("a", 64) {
		t.Errorf("没读到返回的容器 ID：%q", id)
	}

	if v, _ := got["Image"].(string); v != "nginx:1.27-alpine" {
		t.Errorf("顶层 Image 丢了：%v", got["Image"])
	}
	if _, ok := got["Env"]; !ok {
		t.Error("顶层 Env 丢了")
	}
	if _, ok := got["Cmd"]; !ok {
		t.Error("顶层 Cmd 丢了")
	}
	if _, ok := got["HostConfig"]; !ok {
		t.Error("HostConfig 丢了")
	}
	if _, ok := got["NetworkingConfig"]; !ok {
		t.Error("NetworkingConfig 丢了")
	}
	if _, ok := got["Config"]; ok {
		t.Error(`请求体里出现了 "Config" 键 —— 真实守护进程会把它当未知字段静默忽略，然后回 400`)
	}
}

// TestCreateContainerRefusesEmptyConfig 配置为空时必须在本地就拦下，
// 因为交给守护进程只会换来一句看不懂的 400，而且那时容器可能已经被停了。
func TestCreateContainerRefusesEmptyConfig(t *testing.T) {
	srv := newCreateDaemon(t, nil)
	defer srv.Close()
	c := clientFor(t, srv)

	if _, err := c.CreateContainer(context.Background(), "x", nil, map[string]any{}, nil); err == nil {
		t.Error("config 为 nil 时必须直接报错")
	}
	if _, err := c.CreateContainer(context.Background(), "x", map[string]any{}, nil, nil); err == nil {
		t.Error("config 为空 map 时也必须直接报错")
	}
}
