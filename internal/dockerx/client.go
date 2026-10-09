// Package dockerx 是对 Docker Engine API 的轻量封装。
//
// 为什么不用 github.com/docker/docker SDK：
//  1. 依赖树极浅（本包只用标准库），配合纯 Go 的 SQLite 可以轻松 CGO_ENABLED=0 交叉编译；
//  2. 「克隆已有容器」这个核心动作，直接拿 inspect 返回的 Config/HostConfig 原样回传
//     比在 SDK 类型之间来回映射更忠实、更不容易漏字段；
//  3. 拉取进度、事件流这些本来就要自己解析流式 JSON。
package dockerx

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Client Docker Engine API 客户端。
type Client struct {
	hc      *http.Client
	scheme  string // http / https
	addr    string // 用于拼 URL 的 host 部分
	rawHost string

	mu         sync.RWMutex
	apiVersion string
}

// ErrNotFound 表示 Docker 守护进程返回 404。
var ErrNotFound = errors.New("docker: not found")

// ErrConflict 表示 409（比如容器名已存在、容器正在运行）。
var ErrConflict = errors.New("docker: conflict")

// StatusError 携带 HTTP 状态码与守护进程返回的 message。
type StatusError struct {
	Code    int
	Message string
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("docker api %d: %s", e.Code, e.Message)
}

// IsNotFound 判定 404。
func IsNotFound(err error) bool {
	var se *StatusError
	return errors.As(err, &se) && se.Code == http.StatusNotFound
}

// IsConflict 判定 409。
func IsConflict(err error) bool {
	var se *StatusError
	return errors.As(err, &se) && se.Code == http.StatusConflict
}

// New 根据 DOCKER_HOST 风格地址创建客户端（unix:///path 或 tcp://host:port）。
func New(host string) (*Client, error) {
	host = strings.TrimSpace(host)
	if host == "" {
		host = "unix:///var/run/docker.sock"
	}
	c := &Client{rawHost: host}

	switch {
	case strings.HasPrefix(host, "unix://"):
		sock := strings.TrimPrefix(host, "unix://")
		tr := &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				var d net.Dialer
				return d.DialContext(ctx, "unix", sock)
			},
			MaxIdleConns:          10,
			IdleConnTimeout:       30 * time.Second,
			ResponseHeaderTimeout: 0, // pull/events 可能很久
			DisableCompression:    true,
		}
		c.hc = &http.Client{Transport: tr}
		c.scheme = "http"
		c.addr = "docker"
	case strings.HasPrefix(host, "tcp://"), strings.HasPrefix(host, "http://"), strings.HasPrefix(host, "https://"):
		u, err := url.Parse(host)
		if err != nil {
			return nil, fmt.Errorf("invalid DOCKER_HOST %q: %w", host, err)
		}
		c.scheme = u.Scheme
		if c.scheme == "tcp" {
			c.scheme = "http"
		}
		c.addr = u.Host
		c.hc = &http.Client{Transport: &http.Transport{DisableCompression: true}}
	default:
		return nil, fmt.Errorf("unsupported DOCKER_HOST %q (expected unix:// or tcp://)", host)
	}
	c.hc.Timeout = 0 // 由每个请求的 context 控制超时
	return c, nil
}

// Host 返回原始地址（用于展示）。
func (c *Client) Host() string { return c.rawHost }

// APIVersion 返回协商到的 API 版本（未知时为空，此时不带版本前缀）。
func (c *Client) APIVersion() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.apiVersion
}

// Negotiate 读取守护进程版本并决定 API 版本前缀。
func (c *Client) Negotiate(ctx context.Context) error {
	var v struct {
		APIVersion string `json:"ApiVersion"`
		Version    string `json:"Version"`
	}
	// 先不带版本前缀探测
	if err := c.do(ctx, http.MethodGet, "/version", nil, nil, nil, &v); err != nil {
		return err
	}
	ver := strings.TrimSpace(v.APIVersion)
	if ver == "" {
		ver = "1.41"
	}
	c.mu.Lock()
	c.apiVersion = ver
	c.mu.Unlock()
	return nil
}

func (c *Client) urlFor(path string) string {
	ver := c.APIVersion()
	p := path
	if ver != "" {
		p = "/v" + ver + path
	}
	return c.scheme + "://" + c.addr + p
}

// do 执行一次请求。out 非 nil 时把响应体解为 JSON（204/空体忽略）。
func (c *Client) do(ctx context.Context, method, path string, q url.Values, body any, hdr map[string]string, out any) error {
	rc, err := c.stream(ctx, method, path, q, body, hdr)
	if err != nil {
		return err
	}
	defer rc.Close()
	if out == nil {
		_, _ = io.Copy(io.Discard, rc)
		return nil
	}
	dec := json.NewDecoder(rc)
	if err := dec.Decode(out); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		return fmt.Errorf("decode %s %s: %w", method, path, err)
	}
	return nil
}

// stream 发起请求并返回响应体流（调用方负责 Close）。非 2xx 会被读成 StatusError。
func (c *Client) stream(ctx context.Context, method, path string, q url.Values, body any, hdr map[string]string) (io.ReadCloser, error) {
	full := c.urlFor(path)
	if len(q) > 0 {
		full += "?" + q.Encode()
	}
	var rdr io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, full, rdr)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	for k, v := range hdr {
		req.Header.Set(k, v)
	}
	resp, err := c.hc.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		defer resp.Body.Close()
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 8192))
		msg := strings.TrimSpace(string(raw))
		var em struct {
			Message string `json:"message"`
		}
		if json.Unmarshal(raw, &em) == nil && em.Message != "" {
			msg = em.Message
		}
		if msg == "" {
			msg = http.StatusText(resp.StatusCode)
		}
		return nil, &StatusError{Code: resp.StatusCode, Message: msg}
	}
	return resp.Body, nil
}

// Ping 探活。
func (c *Client) Ping(ctx context.Context) error {
	rc, err := c.stream(ctx, http.MethodGet, "/_ping", nil, nil, nil)
	if err != nil {
		return err
	}
	defer rc.Close()
	_, _ = io.Copy(io.Discard, rc)
	return nil
}

// ---------- 基础信息 ----------

// Info 是 GET /info 里我们关心的字段子集。
type Info struct {
	ID              string `json:"ID"`
	Name            string `json:"Name"`
	ServerVersion   string `json:"ServerVersion"`
	OperatingSystem string `json:"OperatingSystem"`
	OSType          string `json:"OSType"`
	Architecture    string `json:"Architecture"`
	KernelVersion   string `json:"KernelVersion"`
	CPUs            int    `json:"NCPU"`
	MemTotal        int64  `json:"MemTotal"`
	Containers      int    `json:"Containers"`
	ContainersRunning int  `json:"ContainersRunning"`
	ContainersPaused  int  `json:"ContainersPaused"`
	ContainersStopped int  `json:"ContainersStopped"`
	Images          int    `json:"Images"`
	DockerRootDir   string `json:"DockerRootDir"`
	Driver          string `json:"Driver"`
	RegistryConfig  struct {
		Mirrors []string `json:"Mirrors"`
	} `json:"RegistryConfig"`
}

// Info 拉取守护进程信息（包含它实际生效的 registry-mirrors，非常关键）。
func (c *Client) Info(ctx context.Context) (*Info, error) {
	var out Info
	if err := c.do(ctx, http.MethodGet, "/info", nil, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ---------- 事件 ----------

// Event 是 GET /events 的一条事件。
type Event struct {
	Type   string            `json:"Type"`
	Action string            `json:"Action"`
	Actor  struct {
		ID         string            `json:"ID"`
		Attributes map[string]string `json:"Attributes"`
	} `json:"Actor"`
	Time int64 `json:"time"`
}

// Events 订阅事件流，直到 ctx 取消。handler 在每个事件上被调用。
func (c *Client) Events(ctx context.Context, handler func(Event)) error {
	q := url.Values{}
	q.Set("filters", `{"type":["container"]}`)
	rc, err := c.stream(ctx, http.MethodGet, "/events", q, nil, nil)
	if err != nil {
		return err
	}
	defer rc.Close()
	sc := bufio.NewScanner(rc)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := bytes.TrimSpace(sc.Bytes())
		if len(line) == 0 {
			continue
		}
		var ev Event
		if err := json.Unmarshal(line, &ev); err != nil {
			continue
		}
		handler(ev)
	}
	if err := sc.Err(); err != nil && ctx.Err() == nil {
		return err
	}
	return ctx.Err()
}

// ---------- 容器 ----------

// ContainerSummary 对应 GET /containers/json 的条目。
type ContainerSummary struct {
	ID      string   `json:"Id"`
	Names   []string `json:"Names"`
	Image   string   `json:"Image"`
	ImageID string   `json:"ImageID"`
	Command string   `json:"Command"`
	Created int64    `json:"Created"`
	State   string   `json:"State"`
	Status  string   `json:"Status"`
	Labels  map[string]string `json:"Labels"`
	Ports   []struct {
		IP          string `json:"IP"`
		PrivatePort uint16 `json:"PrivatePort"`
		PublicPort  uint16 `json:"PublicPort"`
		Type        string `json:"Type"`
	} `json:"Ports"`
}

// Name 返回去掉前导斜线的容器名。
func (s ContainerSummary) Name() string {
	if len(s.Names) == 0 {
		return ""
	}
	return strings.TrimPrefix(s.Names[0], "/")
}

// ComposeProject 返回 compose 项目名（若不是 compose 管理的容器则为空）。
func (s ContainerSummary) ComposeProject() string {
	if s.Labels == nil {
		return ""
	}
	return s.Labels["com.docker.compose.project"]
}

// ListContainers 列出所有容器（含已停止）。
func (c *Client) ListContainers(ctx context.Context) ([]ContainerSummary, error) {
	q := url.Values{}
	q.Set("all", "1")
	var out []ContainerSummary
	if err := c.do(ctx, http.MethodGet, "/containers/json", q, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// Inspect 返回容器的原始 inspect 结果（保留全部字段以便原样回传创建新容器）。
func (c *Client) Inspect(ctx context.Context, id string) (map[string]any, error) {
	var out map[string]any
	if err := c.do(ctx, http.MethodGet, "/containers/"+url.PathEscape(id)+"/json", nil, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ContainerAction 执行 start/stop/restart/kill/pause/unpause。
func (c *Client) ContainerAction(ctx context.Context, id, action string, timeout *int) error {
	q := url.Values{}
	if timeout != nil {
		q.Set("t", fmt.Sprint(*timeout))
	}
	return c.do(ctx, http.MethodPost, "/containers/"+url.PathEscape(id)+"/"+action, q, nil, nil, nil)
}

// RenameContainer 改名（更新流程里把旧容器改成 __bak_ 前缀）。
func (c *Client) RenameContainer(ctx context.Context, id, newName string) error {
	q := url.Values{}
	q.Set("name", newName)
	return c.do(ctx, http.MethodPost, "/containers/"+url.PathEscape(id)+"/rename", q, nil, nil, nil)
}

// RemoveContainer 删除容器。v=true 同时删除匿名卷（危险，需显式确认）。
func (c *Client) RemoveContainer(ctx context.Context, id string, force, removeVolumes bool) error {
	q := url.Values{}
	if force {
		q.Set("force", "1")
	}
	if removeVolumes {
		q.Set("v", "1")
	}
	return c.do(ctx, http.MethodDelete, "/containers/"+url.PathEscape(id), q, nil, nil, nil)
}

// CreateContainer 创建容器。config / hostConfig 直接来自 Inspect 的对应字段，
// 这样挂载、端口、环境、网络、重启策略等全部原样保留 —— 这是「忠实克隆」的关键。
func (c *Client) CreateContainer(ctx context.Context, name string, config, hostConfig map[string]any, networking map[string]any) (string, error) {
	q := url.Values{}
	if name != "" {
		q.Set("name", name)
	}
	body := map[string]any{
		"HostConfig": hostConfig,
	}
	if config != nil {
		body["Config"] = config
	}
	if networking != nil {
		body["NetworkingConfig"] = networking
	}
	var out struct {
		ID       string   `json:"Id"`
		Warnings []string `json:"Warnings"`
	}
	if err := c.do(ctx, http.MethodPost, "/containers/create", q, body, nil, &out); err != nil {
		return "", err
	}
	return out.ID, nil
}

// UpdateContainer 更新容器的资源限制（还原/调整时用）。
func (c *Client) UpdateContainer(ctx context.Context, id string, update map[string]any) error {
	return c.do(ctx, http.MethodPost, "/containers/"+url.PathEscape(id)+"/update", nil, update, nil, nil)
}

// Logs 返回容器日志流。follow=true 时保持连接。
func (c *Client) Logs(ctx context.Context, id string, tail int, follow bool, timestamps bool) (io.ReadCloser, error) {
	q := url.Values{}
	q.Set("stdout", "1")
	q.Set("stderr", "1")
	if tail > 0 {
		q.Set("tail", fmt.Sprint(tail))
	} else {
		q.Set("tail", "all")
	}
	if follow {
		q.Set("follow", "1")
	}
	if timestamps {
		q.Set("timestamps", "1")
	}
	return c.stream(ctx, http.MethodGet, "/containers/"+url.PathEscape(id)+"/logs", q, nil, nil)
}

// StatsOneShot 取一次资源占用快照。
func (c *Client) StatsOneShot(ctx context.Context, id string) (map[string]any, error) {
	q := url.Values{}
	q.Set("stream", "0")
	q.Set("one-shot", "1")
	var out map[string]any
	if err := c.do(ctx, http.MethodGet, "/containers/"+url.PathEscape(id)+"/stats", q, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ---------- 镜像 ----------

// ImageSummary 对应 GET /images/json 的条目。
type ImageSummary struct {
	ID          string            `json:"Id"`
	ParentID    string            `json:"ParentId"`
	RepoTags    []string          `json:"RepoTags"`
	RepoDigests []string          `json:"RepoDigests"`
	Created     int64             `json:"Created"`
	Size        int64             `json:"Size"`
	Labels      map[string]string `json:"Labels"`
	Containers  int64             `json:"Containers"`
}

// ListImages 列出镜像。all=false 时隐藏中间层镜像。
func (c *Client) ListImages(ctx context.Context, all bool) ([]ImageSummary, error) {
	q := url.Values{}
	if all {
		q.Set("all", "1")
	}
	var out []ImageSummary
	if err := c.do(ctx, http.MethodGet, "/images/json", q, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ImageInspect 返回镜像 inspect 原始结果。
func (c *Client) ImageInspect(ctx context.Context, ref string) (map[string]any, error) {
	var out map[string]any
	if err := c.do(ctx, http.MethodGet, "/images/"+url.PathEscape(ref)+"/json", nil, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// RemoveImage 删除镜像。
func (c *Client) RemoveImage(ctx context.Context, ref string, force, noprune bool) error {
	q := url.Values{}
	if force {
		q.Set("force", "1")
	}
	if noprune {
		q.Set("noprune", "1")
	}
	return c.do(ctx, http.MethodDelete, "/images/"+url.PathEscape(ref), q, nil, nil, nil)
}

// PruneImages 清理悬空镜像，返回释放的字节数。
func (c *Client) PruneImages(ctx context.Context) (int64, error) {
	var out struct {
		ImagesDeleted []map[string]string `json:"ImagesDeleted"`
		SpaceReclaimed int64              `json:"SpaceReclaimed"`
	}
	if err := c.do(ctx, http.MethodPost, "/images/prune", nil, map[string]any{}, nil, &out); err != nil {
		return 0, err
	}
	return out.SpaceReclaimed, nil
}

// TagImage 给本地镜像再打一个标签。
// 用途：通过加速站拉到的镜像，需要重新打回原始引用，这样 compose 与其它工具
// 仍然按原来的名字找得到它。
func (c *Client) TagImage(ctx context.Context, source, repo, tag string) error {
	q := url.Values{}
	q.Set("repo", repo)
	if tag != "" {
		q.Set("tag", tag)
	}
	return c.do(ctx, http.MethodPost, "/images/"+url.PathEscape(source)+"/tag", q, nil, nil, nil)
}

// RemoveContainerByName 便捷删除。
func (c *Client) RemoveContainerByName(ctx context.Context, name string, force, vols bool) error {
	return c.RemoveContainer(ctx, name, force, vols)
}

// ---------- 网络 / 卷 ----------

// ListNetworks 列出网络。
func (c *Client) ListNetworks(ctx context.Context) ([]map[string]any, error) {
	var out []map[string]any
	if err := c.do(ctx, http.MethodGet, "/networks", nil, nil, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// ListVolumes 列出卷。
func (c *Client) ListVolumes(ctx context.Context) ([]map[string]any, error) {
	var out struct {
		Volumes []map[string]any `json:"Volumes"`
	}
	if err := c.do(ctx, http.MethodGet, "/volumes", nil, nil, nil, &out); err != nil {
		return nil, err
	}
	return out.Volumes, nil
}

// RemoveVolume 删除卷。
func (c *Client) RemoveVolume(ctx context.Context, name string, force bool) error {
	q := url.Values{}
	if force {
		q.Set("force", "1")
	}
	return c.do(ctx, http.MethodDelete, "/volumes/"+url.PathEscape(name), q, nil, nil, nil)
}
