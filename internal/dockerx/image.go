package dockerx

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Ref 是解析后的镜像引用。
type Ref struct {
	Original   string `json:"original"`
	Registry   string `json:"registry"`   // docker.io / ghcr.io / ...
	Repository string `json:"repository"` // library/redis
	Tag        string `json:"tag"`        // alpine
	Digest     string `json:"digest"`     // 形如 sha256:...，可空
	// LocalName 是与 RepoDigests 比较时用的仓库名（docker.io 会去掉 library/ 前缀）。
	LocalName string `json:"localName"`
}

// String 返回可用于 pull 的引用。
func (r Ref) String() string {
	s := r.Repository
	if r.Registry != "" && r.Registry != "docker.io" {
		s = r.Registry + "/" + r.Repository
	}
	if r.Tag != "" {
		s += ":" + r.Tag
	}
	if r.Digest != "" {
		s += "@" + r.Digest
	}
	return s
}

// IsOfficial 是否是 docker.io 官方镜像。
func (r Ref) IsOfficial() bool {
	return r.Registry == "docker.io" && strings.HasPrefix(r.Repository, "library/")
}

// defaultTag 补全默认标签。
func defaultTag(t string) string {
	if t == "" {
		return "latest"
	}
	return t
}

// ParseRef 解析镜像引用。
//
// 规则与 Docker CLI 一致：第一段包含 '.' 或 ':' 或是 "localhost" ⇒ 它是 registry 域名，
// 否则视为 Docker Hub 镜像（官方镜像补 library/ 前缀）。
func ParseRef(s string) Ref {
	orig := strings.TrimSpace(s)
	r := Ref{Original: orig}
	rest := orig

	// 分离 digest
	if i := strings.Index(rest, "@"); i >= 0 {
		r.Digest = rest[i+1:]
		rest = rest[:i]
	}

	// 分离 registry
	first := rest
	remainder := ""
	if i := strings.Index(rest, "/"); i >= 0 {
		first = rest[:i]
		remainder = rest[i+1:]
	}
	if remainder != "" && (strings.ContainsAny(first, ".:") || first == "localhost") {
		r.Registry = first
		rest = remainder
	} else {
		r.Registry = "docker.io"
	}

	// 分离 tag（注意不能把 registry 端口当成 tag，此处 rest 已不含 registry）
	if i := strings.LastIndex(rest, ":"); i >= 0 && !strings.Contains(rest[i+1:], "/") {
		r.Tag = rest[i+1:]
		rest = rest[:i]
	}
	r.Tag = defaultTag(r.Tag)
	r.Repository = rest

	// 官方镜像补 library/；同时计算本地比对名
	if r.Registry == "docker.io" {
		if !strings.Contains(r.Repository, "/") {
			r.Repository = "library/" + r.Repository
		}
		r.LocalName = strings.TrimPrefix(r.Repository, "library/")
	} else {
		r.LocalName = r.Registry + "/" + r.Repository
	}
	if r.Tag == "" && r.Digest == "" {
		r.Tag = "latest"
	}
	return r
}

// DigestOnly 返回 sha256:...（去掉算法前缀则返回空）。
func DigestOnly(d string) string {
	d = strings.TrimSpace(d)
	if strings.HasPrefix(d, "sha256:") {
		return d
	}
	return d
}

// ---------- 本地镜像信息 ----------

// LocalDigests 从镜像 inspect 结果里取出 RepoDigests 对应的 sha256 列表。
func LocalDigests(inspect map[string]any) []string {
	out := []string{}
	raw, ok := inspect["RepoDigests"]
	if !ok {
		return out
	}
	list, ok := raw.([]any)
	if !ok {
		return out
	}
	for _, item := range list {
		s, _ := item.(string)
		if i := strings.LastIndex(s, "@"); i >= 0 {
			s = s[i+1:]
		}
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}

// ImageID 从镜像 inspect 结果里取 Id。
func ImageID(inspect map[string]any) string {
	s, _ := inspect["Id"].(string)
	return s
}

// ShortID 把 sha256:abcd... 截成 abcd 12 位。
func ShortID(id string) string {
	id = strings.TrimPrefix(id, "sha256:")
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

// ---------- 拉取 ----------

// PullEvent 拉取过程中的一条进度。
type PullEvent struct {
	Status string `json:"status"`
	ID     string `json:"id"`
	Detail string `json:"detail"`
	// Progress 进度字符串（如 "1.2MB/3.4MB"）
	Progress string `json:"progress"`
}

// PullResult 拉取结果。
type PullResult struct {
	Ref        string   `json:"ref"`
	ImageID    string   `json:"imageId"`
	Digest     string   `json:"digest"`
	UpToDate   bool     `json:"upToDate"` // 守护进程明确回了 "Image is up to date"
	Downloaded bool     `json:"downloaded"` // 守护进程明确回了 "Downloaded newer image"
	Messages   []string `json:"messages"`
}

// Pull 拉取镜像。
//
// 关键点：这是**唯一**会让镜像真正变化的地方 —— 走的是 Docker 守护进程自己的
// ImagePull，因此用户的 registry-mirrors / 私有仓库凭据全部自动生效。
// 检测侧与拉取侧因此永远同源，不会出现 dockerCopilot 那种「检测说有新版本、
// 拉取说已是最新」的永久误报。
//
// onEvent 可为 nil。
func (c *Client) Pull(ctx context.Context, ref string, onEvent func(PullEvent)) (*PullResult, error) {
	r := ParseRef(ref)
	q := url.Values{}
	q.Set("fromImage", r.Repository)
	if r.Registry != "" && r.Registry != "docker.io" {
		// 非 Docker Hub 必须带上完整域名，守护进程才知道去哪个仓库拉
		q.Set("fromImage", r.Registry+"/"+r.Repository)
	}
	if r.Tag != "" {
		q.Set("tag", r.Tag)
	}
	if r.Digest != "" {
		q.Del("tag")
		q.Set("fromImage", q.Get("fromImage")+"@"+r.Digest)
	}

	// 私有仓库凭据：交给守护进程用 ~/.docker/config.json；若用户在设置里配了
	// 显式凭据则通过 X-Registry-Auth 传递（见 authRegistryHeader）。
	rc, err := c.stream(ctx, http.MethodPost, "/images/create", q, nil, nil)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	res := &PullResult{Ref: r.String()}
	sc := bufio.NewScanner(rc)
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		var m struct {
			Status         string `json:"status"`
			ID             string `json:"id"`
			Progress       string `json:"progress"`
			ProgressDetail struct {
				Current int64 `json:"current"`
				Total   int64 `json:"total"`
			} `json:"progressDetail"`
			Error       string `json:"error"`
			ErrorDetail struct {
				Message string `json:"message"`
			} `json:"errorDetail"`
		}
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			continue
		}
		if m.Error != "" || m.ErrorDetail.Message != "" {
			msg := m.ErrorDetail.Message
			if msg == "" {
				msg = m.Error
			}
			return res, fmt.Errorf("pull failed: %s", msg)
		}
		if m.Status != "" {
			res.Messages = append(res.Messages, m.Status)
			low := strings.ToLower(m.Status)
			switch {
			case strings.Contains(low, "image is up to date"):
				res.UpToDate = true
			case strings.Contains(low, "downloaded newer image"):
				res.Downloaded = true
			}
			if strings.HasPrefix(m.Status, "Digest:") {
				res.Digest = strings.TrimSpace(strings.TrimPrefix(m.Status, "Digest:"))
			}
			if onEvent != nil {
				onEvent(PullEvent{
					Status:   m.Status,
					ID:       m.ID,
					Progress: m.Progress,
					Detail: fmt.Sprintf("%d/%d", m.ProgressDetail.Current, m.ProgressDetail.Total),
				})
			}
		}
	}
	if err := sc.Err(); err != nil {
		return res, err
	}

	// 拉完再 inspect 一次，拿到权威的本地镜像 ID 与 digest
	if insp, err := c.ImageInspect(ctx, r.String()); err == nil {
		res.ImageID = ImageID(insp)
		if res.Digest == "" {
			if ds := LocalDigests(insp); len(ds) > 0 {
				res.Digest = ds[0]
			}
		}
	}
	return res, nil
}

// ---------- 远端检测（同源） ----------

// DistributionInspect 询问守护进程「这个镜像在它对应的仓库里是什么」。
//
// 与 dockerCopilot 的「自拼 registry HTTP 请求 + 硬编码加速站表」不同：
// 这个接口由守护进程自己解析仓库端点，因此与 docker pull 走的是同一套配置。
// 如果注册表需要认证（401）或接口不可用，返回错误 —— 调用方必须把该容器标成
// 「未知」而不是「有新版本」，这条规则是从 dockerCopilot 的误报里学到的。
func (c *Client) DistributionInspect(ctx context.Context, ref string) (digest string, err error) {
	r := ParseRef(ref)
	name := r.LocalName
	if r.Registry != "docker.io" {
		name = r.Registry + "/" + r.Repository
	}
	var out struct {
		Descriptor struct {
			Digest string `json:"digest"`
			MediaType string `json:"mediaType"`
		} `json:"Descriptor"`
		Platforms []map[string]any `json:"Platforms"`
	}
	if err := c.do(ctx, http.MethodGet, "/distribution/"+url.PathEscape(name)+"/json", nil, nil, nil, &out); err != nil {
		return "", err
	}
	return out.Descriptor.Digest, nil
}

// PullStreamReader 把容器日志的 multiplexed stream 去帧。
// Docker 的日志流是 8 字节头（stream type + 3 字节 padding + 4 字节大端长度）后跟负载。
type stdWriter struct {
	src io.Reader
	buf []byte
}

// DemuxLogs 把 docker logs 的多路复用流还原成纯文本（follow 模式也可用）。
func DemuxLogs(r io.Reader) io.Reader {
	pr, pw := io.Pipe()
	go func() {
		defer pw.Close()
		br := bufio.NewReaderSize(r, 32*1024)
		hdr := make([]byte, 8)
		for {
			if _, err := io.ReadFull(br, hdr); err != nil {
				return
			}
			// 若非标准帧（某些驱动返回裸文本），把已读的当文本输出
			if hdr[0] > 2 {
				_, _ = pw.Write(hdr)
				_, _ = io.Copy(pw, br)
				return
			}
			n := int(uint32(hdr[4])<<24 | uint32(hdr[5])<<16 | uint32(hdr[6])<<8 | uint32(hdr[7]))
			if n <= 0 {
				continue
			}
			if _, err := io.CopyN(pw, br, int64(n)); err != nil {
				return
			}
		}
	}()
	return pr
}
