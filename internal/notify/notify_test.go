package notify

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/aaron2024s/dockhelm/internal/store"
)

// 2026-10-11 用户反馈沉淀：
//  1. 主机显示的是容器随机 ID（eba791308b83），要换成机器型号（dx4600）；
//  2. 主机要放头部，不能挂在最后一行 —— 消息被截断时最后一行最先消失；
//  3. 批量失败推送被腰斩过一次，正文要有源头截断兜底。

func TestGenericBoardName(t *testing.T) {
	for _, s := range []string{"System Product Name", "To be filled by O.E.M.", "Default string", "None", "TYPE1"} {
		if !genericBoardName(s) {
			t.Errorf("%q 应识别为占位值", s)
		}
	}
	for _, s := range []string{"DX4600", "DS923+", "ROCK 5B"} {
		if genericBoardName(s) {
			t.Errorf("%q 不是占位值，不该被过滤", s)
		}
	}
}

func TestHostLabelEnvOverride(t *testing.T) {
	t.Setenv("DOCKHELM_HOST_LABEL", "dx4600")
	if got := HostLabel(); got != "dx4600" {
		t.Fatalf("环境变量应最优先，实际 %q", got)
	}
}

func TestHostLabelFallbackNonEmpty(t *testing.T) {
	os.Unsetenv("DOCKHELM_HOST_LABEL")
	if got := HostLabel(); got == "" {
		t.Fatal("任何环境下都不该返回空串")
	}
}

// compose 的行序：主机/时间必须钉在头部。
func TestComposePutsHostFirst(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir + "/db.json")
	if err != nil {
		t.Fatal(err)
	}
	m := New(st, nil)
	m.host = "dx4600"
	_, text := m.compose("batch_update_failed", "", "", map[string]string{
		"container": "3 个容器",
		"result":    "自动更新有失败",
		"message":   "失败明细：\n· web：拉取失败",
	})
	lines := strings.Split(text, "\n")
	if len(lines) < 4 {
		t.Fatalf("正文行数不对：%q", text)
	}
	if !strings.HasPrefix(lines[1], "主机：dx4600") {
		t.Fatalf("第 2 行应是主机，实际 %q", lines[1])
	}
	if !strings.HasPrefix(lines[2], "时间：") {
		t.Fatalf("第 3 行应是时间，实际 %q", lines[2])
	}
	if strings.Contains(text, "主机：eba791308b83") {
		t.Fatal("不该再出现容器随机 ID 式主机名")
	}
}

func TestTruncateNotifyText(t *testing.T) {
	// 不超长 ⇒ 原样
	short := "一行\n两行"
	if got := truncateNotifyText(short); got != short {
		t.Fatalf("短文本应原样返回：%q", got)
	}
	// 超长 ⇒ 按行截断，头部信息保留
	head := strings.Repeat("头", 100)
	var b strings.Builder
	b.WriteString(head + "\n")
	for i := 0; i < 400; i++ {
		b.WriteString(strings.Repeat("x", 60) + "\n")
	}
	got := truncateNotifyText(b.String())
	if len(got) > maxNotifyTextBytes+100 {
		t.Fatalf("截断后仍超长：%d 字节", len(got))
	}
	if !strings.HasPrefix(got, head) {
		t.Fatal("头部信息必须保留")
	}
	if !strings.HasSuffix(got, "……（内容过长已截断）") {
		t.Fatalf("缺截断提示：%q", got[len(got)-60:])
	}
}

// 2026-10-11 用户反馈：ntfy 配置失败、没有 Gotify 渠道。参照 course-learn
// 的成熟做法把 ntfy 改成 JSON body，并修掉三个渲染层缺陷。
func TestRenderNtfy(t *testing.T) {
	vals := map[string]string{
		"title":    "🚨 Dockhelm 容器更新失败 · app",
		"text":     "详情：拉取失败\n时间：2026-10-11 09:00:00",
		"priority": "4",
		"tags":     "rotating_light",
	}
	cfg := map[string]any{"server": "https://ntfy.sh", "topic": "dockhelm-alerts", "token": "tk_xxx"}
	method, u, headers, body, ct, err := RenderRequest("ntfy", cfg, vals)
	if err != nil {
		t.Fatalf("RenderRequest: %v", err)
	}
	if method != "POST" || ct != "application/json" {
		t.Fatalf("method=%s ct=%s", method, ct)
	}
	// 修复①：server 变量不能被 QueryEscape，否则 https:// 变 %3A%2F%2F
	if !strings.HasPrefix(u, "https://ntfy.sh") {
		t.Fatalf("URL 被转义坏了：%q", u)
	}
	// 修复②：标题里的中文+emoji 走 JSON body，绝不能出现在请求头里
	// （Go 的 http 客户端会拒绝非 ASCII 头值）。
	for k, v := range headers {
		if !isASCII(v) {
			t.Fatalf("请求头 %s 含非 ASCII 值：%q", k, v)
		}
	}
	if headers["Authorization"] != "Bearer tk_xxx" {
		t.Fatalf("修复③：token 没进 Authorization 头：%v", headers)
	}
	for _, want := range []string{`"topic":"dockhelm-alerts"`, `"priority":4`, `"tags":["rotating_light"]`} {
		if !strings.Contains(body, want) {
			t.Fatalf("body 缺 %s：%q", want, body)
		}
	}

	// token 留空 ⇒ 不带 Authorization 头（别发空 "Bearer " 出去）
	_, _, headers2, _, _, _ := RenderRequest("ntfy", map[string]any{"server": "https://ntfy.sh", "topic": "t"}, vals)
	if _, ok := headers2["Authorization"]; ok {
		t.Fatalf("token 为空不该有 Authorization 头：%v", headers2)
	}
}

func TestRenderGotify(t *testing.T) {
	vals := map[string]string{
		"title": "✅ Dockhelm 批量更新完成", "text": "详情：全部成功",
		"gotifyPriority": "5",
	}
	// token 里的特殊字符必须转义（它在 query string 里）
	cfg := map[string]any{"server": "https://push.example.com", "token": "Aabc 123/xyz"}
	_, u, _, body, ct, err := RenderRequest("gotify", cfg, vals)
	if err != nil {
		t.Fatalf("RenderRequest: %v", err)
	}
	if ct != "application/json" {
		t.Fatalf("ct=%s", ct)
	}
	if !strings.HasPrefix(u, "https://push.example.com/message?token=") || !strings.Contains(u, "Aabc+123%2Fxyz") {
		t.Fatalf("URL 不对：%q", u)
	}
	for _, want := range []string{`"title":"✅ Dockhelm 批量更新完成"`, `"priority":5`} {
		if !strings.Contains(body, want) {
			t.Fatalf("body 缺 %s：%q", want, body)
		}
	}
}

func TestPresetHasGotify(t *testing.T) {
	if _, ok := PresetMap["gotify"]; !ok {
		t.Fatal("gotify 预设缺失")
	}
}

func isASCII(s string) bool {
	for _, r := range s {
		if r > 127 {
			return false
		}
	}
	return true
}

// 2026-10-12 用户反馈：推送标题与「结果：」那一行完全重复（Docker 连接断开 …）。
// 规则：结果短语与事件名相同时不另起一行；只有补充了新信息（「全量备份有失败」）
// 才输出「结果：」。
func TestComposeDropsDuplicateResult(t *testing.T) {
	dir := t.TempDir()
	st, err := store.Open(dir + "/db.json")
	if err != nil {
		t.Fatal(err)
	}
	m := New(st, nil)
	m.host = "dx4600"

	// 结果 == 事件名 ⇒ 不出现「结果：」行，详情照常保留
	_, text := m.compose("docker_disconnected", "", "", map[string]string{
		"result": "Docker 连接断开", "message": "Docker 守护进程当前不可达",
	})
	if strings.Contains(text, "结果：") {
		t.Fatalf("结果与标题重复时不该出现「结果：」行：\n%s", text)
	}
	if !strings.Contains(text, "详情：Docker 守护进程当前不可达") {
		t.Fatalf("详情行应保留：\n%s", text)
	}

	// 结果补充了新信息 ⇒ 保留
	_, text = m.compose("backup_failed", "", "", map[string]string{
		"container": "24 个容器", "result": "全量备份有失败", "message": "成功 1、失败 23",
	})
	if !strings.Contains(text, "结果：全量备份有失败") {
		t.Fatalf("结果与事件名不同时必须保留：\n%s", text)
	}
}

// JSON 类预设：正文里的真实换行、引号、反斜杠必须被转义成合法 JSON 字符串内容。
// 曾经的缺陷：原样嵌入 ⇒ 发出去的是非法 JSON，对方 400。
func TestJSONBodyEscapesControlChars(t *testing.T) {
	vals := map[string]string{
		"title":    `带"引号"的标题`,
		"text":     "第一行\n第二行 带\"引号\" 和 \\ 反斜杠\t制表",
		"priority": "4",
		"tags":     "rotating_light",
	}
	cfg := map[string]any{"server": "https://ntfy.sh", "topic": "t", "token": ""}
	_, _, _, body, _, err := RenderRequest("ntfy", cfg, vals)
	if err != nil {
		t.Fatalf("RenderRequest: %v", err)
	}
	if !json.Valid([]byte(body)) {
		t.Fatalf("body 不是合法 JSON：%q", body)
	}
	var parsed struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal([]byte(body), &parsed); err != nil {
		t.Fatalf("解析失败: %v", err)
	}
	want := "第一行\n第二行 带\"引号\" 和 \\ 反斜杠\t制表"
	if parsed.Message != want {
		t.Fatalf("转义后还原的正文不一致：\n got %q\nwant %q", parsed.Message, want)
	}
}
