package notify

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
)

// ChannelField 渠道的一个配置字段。
type ChannelField struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Type        string `json:"type"` // text | password | textarea
	Required    bool   `json:"required"`
	Placeholder string `json:"placeholder"`
	Help        string `json:"help"`
}

// Preset 一个渠道预设 = 一组字段 + 一个「URL / 方法 / 请求头 / 请求体」模板。
//
// 加一个新渠道只需要往 Presets 里加一条数据，不需要动任何代码 ——
// 这是为了避开「每接一个通知服务就要写一套适配器」的泥潭。
type Preset struct {
	Type        string            `json:"type"`
	Label       string            `json:"label"`
	Description string            `json:"description"`
	Fields      []ChannelField    `json:"fields"`
	Method      string            `json:"method"`
	URL         string            `json:"url"`
	ContentType string            `json:"contentType"` // application/json | text/plain | application/x-www-form-urlencoded
	Headers     map[string]string `json:"headers"`
	Body        string            `json:"body"`
	// Config 里是否还能自定义请求头（仅自定义 webhook 为 true）
	AllowCustomHeaders bool `json:"allowCustomHeaders"`
}

// Presets 全部渠道预设。
var Presets = []Preset{
	{
		Type: "telegram", Label: "Telegram", Description: "Telegram 机器人",
		Fields: []ChannelField{
			{Key: "token", Label: "Bot Token", Type: "password", Required: true, Placeholder: "123456:ABC-DEF…", Help: "找 @BotFather 创建机器人后获得"},
			{Key: "chatId", Label: "Chat ID", Type: "text", Required: true, Placeholder: "-1001234567890", Help: "给机器人发条消息后访问 /getUpdates 可看到"},
		},
		Method:      "POST",
		URL:         "https://api.telegram.org/bot{{token}}/sendMessage",
		ContentType: "application/json",
		Body:        `{"chat_id":"{{chatId}}","text":"{{text}}","parse_mode":"HTML","disable_web_page_preview":true}`,
	},
	{
		Type: "bark", Label: "Bark", Description: "iOS 推送（Bark）",
		Fields: []ChannelField{
			{Key: "server", Label: "服务器", Type: "text", Placeholder: "https://api.day.app", Help: "自建 Bark 就填自己的域名；官方留空用默认"},
			{Key: "key", Label: "Key", Type: "password", Required: true, Placeholder: "你的 Bark Key"},
			{Key: "group", Label: "分组", Type: "text", Placeholder: "dockhelm"},
		},
		Method:      "POST",
		URL:         "{{server}}/{{key}}",
		ContentType: "application/json",
		Body:        `{"title":"{{title}}","body":"{{text}}","group":"{{group}}","level":"{{sound}}"}`,
	},
	{
		Type: "ntfy", Label: "ntfy", Description: "ntfy 推送（可自建）",
		Fields: []ChannelField{
			{Key: "server", Label: "服务器", Type: "text", Placeholder: "https://ntfy.sh"},
			{Key: "topic", Label: "主题 Topic", Type: "text", Required: true, Placeholder: "dockhelm-alerts"},
			{Key: "token", Label: "访问令牌", Type: "password", Help: "自建且开了鉴权时填写"},
		},
		Method:      "POST",
		URL:         "{{server}}/{{topic}}",
		ContentType: "text/plain",
		Headers:     map[string]string{"Title": "{{title}}", "Priority": "{{priority}}", "Tags": "{{tags}}"},
		Body:        "{{text}}",
	},
	{
		Type: "wecom", Label: "企业微信机器人", Description: "企业微信群机器人 Webhook",
		Fields: []ChannelField{
			{Key: "key", Label: "Webhook Key", Type: "password", Required: true, Placeholder: "群机器人 URL 里 key= 后面的部分"},
			{Key: "mentioned", Label: "@成员手机号", Type: "text", Placeholder: "13800138000，多人用逗号分隔"},
		},
		Method:      "POST",
		URL:         "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key={{key}}",
		ContentType: "application/json",
		Body:        `{"msgtype":"text","text":{"content":"{{text}}","mentioned_mobile_list":{{mentionedJson}}}}`,
	},
	{
		Type: "dingtalk", Label: "钉钉机器人", Description: "钉钉群自定义机器人",
		Fields: []ChannelField{
			{Key: "token", Label: "access_token", Type: "password", Required: true, Placeholder: "机器人 Webhook 里的 access_token"},
			{Key: "secret", Label: "加签密钥", Type: "password", Help: "安全设置选了「加签」才需要填"},
		},
		Method:      "POST",
		URL:         "https://oapi.dingtalk.com/robot/send?access_token={{token}}&timestamp={{timestamp}}&sign={{sign}}",
		ContentType: "application/json",
		Body:        `{"msgtype":"text","text":{"content":"{{text}}"}}`,
	},
	{
		Type: "feishu", Label: "飞书机器人", Description: "飞书群自定义机器人",
		Fields: []ChannelField{
			{Key: "token", Label: "Webhook Token", Type: "password", Required: true, Placeholder: "…/hook/ 后面的部分"},
			{Key: "secret", Label: "签名校验密钥", Help: "开启了签名校验才需要填"},
		},
		Method:      "POST",
		URL:         "https://open.feishu.cn/open-apis/bot/v2/hook/{{token}}",
		ContentType: "application/json",
		Body:        `{"msg_type":"text","content":{"text":"{{text}}"}}`,
	},
	{
		Type: "serverchan", Label: "Server 酱", Description: "方糖 Server 酱",
		Fields: []ChannelField{
			{Key: "sendkey", Label: "SendKey", Type: "password", Required: true, Placeholder: "SCT…"},
		},
		Method:      "POST",
		URL:         "https://sctapi.ftqq.com/{{sendkey}}.send",
		ContentType: "application/x-www-form-urlencoded",
		Body:        "title={{title}}&desp={{text}}",
	},
	{
		Type: "pushplus", Label: "PushPlus", Description: "PushPlus 微信推送",
		Fields: []ChannelField{
			{Key: "token", Label: "Token", Type: "password", Required: true},
			{Key: "topic", Label: "群组编码", Help: "一对多推送时填写"},
		},
		Method:      "POST",
		URL:         "https://www.pushplus.plus/send",
		ContentType: "application/json",
		Body:        `{"token":"{{token}}","title":"{{title}}","content":"{{text}}","template":"markdown","topic":"{{topic}}"}`,
	},
	{
		Type: "smtp", Label: "邮件 SMTP", Description: "发送到邮箱",
		Fields: []ChannelField{
			{Key: "host", Label: "SMTP 服务器", Type: "text", Required: true, Placeholder: "smtp.qq.com"},
			{Key: "port", Label: "端口", Type: "text", Required: true, Placeholder: "465（SSL）或 587（STARTTLS）"},
			{Key: "username", Label: "用户名", Type: "text", Required: true, Placeholder: "you@qq.com"},
			{Key: "password", Label: "密码 / 授权码", Type: "password", Required: true},
			{Key: "from", Label: "发件人", Type: "text", Placeholder: "留空则用用户名"},
			{Key: "to", Label: "收件人", Type: "text", Required: true, Placeholder: "多人用逗号分隔"},
			{Key: "tls", Label: "加密方式", Type: "text", Placeholder: "ssl / starttls / none"},
		},
		Method: "SMTP",
	},
	{
		Type: "webhook", Label: "自定义 Webhook", Description: "任意 URL + 自定义请求头与请求体",
		Fields: []ChannelField{
			{Key: "url", Label: "请求 URL", Type: "text", Required: true, Placeholder: "https://example.com/hook"},
			{Key: "method", Label: "方法", Type: "text", Placeholder: "POST / PUT"},
			{Key: "contentType", Label: "Content-Type", Type: "text", Placeholder: "application/json"},
			{Key: "headers", Label: "自定义请求头", Type: "textarea", Placeholder: `{"Authorization":"Bearer xxx"}`, Help: "JSON 对象，值里同样可以用模板变量"},
			{Key: "body", Label: "请求体模板", Type: "textarea", Placeholder: `{"text":"{{text}}"}`, Help: "可用变量见下方说明"},
		},
		Method:             "POST",
		AllowCustomHeaders: true,
	},
}

// PresetMap 按 type 索引。
var PresetMap = func() map[string]Preset {
	m := make(map[string]Preset, len(Presets))
	for _, p := range Presets {
		m[p.Type] = p
	}
	return m
}()

// DefaultConfig 返回某预设的默认配置（带合理默认值）。
func DefaultConfig(typ string) map[string]any {
	switch typ {
	case "bark":
		return map[string]any{"server": "https://api.day.app", "group": "dockhelm"}
	case "ntfy":
		return map[string]any{"server": "https://ntfy.sh"}
	case "webhook":
		return map[string]any{"method": "POST", "contentType": "application/json"}
	case "smtp":
		return map[string]any{"port": "465", "tls": "ssl"}
	}
	return map[string]any{}
}

// TemplateVar 一个模板变量。
type TemplateVar struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

// TemplateVars 渲染模板时可用的变量（README 与通知页都会展示这份清单）。
var TemplateVars = []TemplateVar{
	{"event", "事件代码，如 update_failed"},
	{"eventLabel", "事件中文名，如「容器更新失败」"},
	{"level", "级别：urgent / normal"},
	{"title", "标题（已含事件名与容器名）"},
	{"text", "完整正文（标题 + 详情 + 时间 + 主机）"},
	{"message", "详情文字"},
	{"container", "容器名"},
	{"image", "镜像引用"},
	{"result", "结果短语，如「更新成功」「已回滚」"},
	{"time", "事件时间（本地时区）"},
	{"host", "Dockhelm 主机名"},
	{"url", "面板地址（若已在设置里填写）"},
	{"sound", "Bark 用：urgent 事件为 «» ，否则空"},
	{"priority", "ntfy 用：urgent → 4，normal → 3"},
	{"tags", "ntfy 用：urgent → rotating_light，normal → whale"},
	{"titleEncoded", "URL 编码后的标题"},
	{"textEncoded", "URL 编码后的正文"},
}

// render 把模板里的 {{key}} 替换成 values 里对应的值。
// inURL=true 时对每个值做 URL 转义（用于路径/查询参数中的变量）。
func render(tpl string, values map[string]string, inURL bool) string {
	if tpl == "" {
		return ""
	}
	var b strings.Builder
	i := 0
	for i < len(tpl) {
		open := strings.Index(tpl[i:], "{{")
		if open < 0 {
			b.WriteString(tpl[i:])
			break
		}
		b.WriteString(tpl[i : i+open])
		rest := tpl[i+open+2:]
		closeIdx := strings.Index(rest, "}}")
		if closeIdx < 0 {
			b.WriteString(tpl[i+open:])
			break
		}
		key := strings.TrimSpace(rest[:closeIdx])
		val, ok := values[key]
		if !ok {
			val = ""
		}
		if inURL {
			b.WriteString(url.QueryEscape(val))
		} else {
			b.WriteString(val)
		}
		i = i + open + 2 + closeIdx + 2
	}
	return b.String()
}

// buildValues 汇总一次通知的全部可用变量。
func buildValues(event string, vars map[string]string, host, panelURL, title, text string) map[string]string {
	level := LevelOf(event)
	if v, ok := vars["level"]; ok && v != "" {
		level = v
	}
	out := map[string]string{}
	for k, v := range vars {
		out[k] = v
	}
	out["event"] = event
	out["eventLabel"] = Label(event)
	out["level"] = level
	out["title"] = title
	out["text"] = text
	out["host"] = host
	if panelURL != "" {
		out["url"] = panelURL
	}
	if _, ok := out["container"]; !ok {
		out["container"] = ""
	}
	if _, ok := out["image"]; !ok {
		out["image"] = ""
	}
	if _, ok := out["message"]; !ok {
		out["message"] = ""
	}
	if _, ok := out["result"]; !ok {
		out["result"] = ""
	}
	sound := ""
	if level == LevelUrgent {
		sound = "«"
	}
	out["sound"] = sound
	if level == LevelUrgent {
		out["priority"] = "4"
		out["tags"] = "rotating_light"
	} else {
		out["priority"] = "3"
		out["tags"] = "whale"
	}
	out["titleEncoded"] = url.QueryEscape(title)
	out["textEncoded"] = url.QueryEscape(text)
	return out
}

// RenderRequest 把渠道配置 + 一次事件渲染成一次 HTTP 请求。
// 返回 (method, url, headers, body, contentType, error)。
func RenderRequest(chType string, cfg map[string]any, values map[string]string) (string, string, map[string]string, string, string, error) {
	p, ok := PresetMap[chType]
	if !ok {
		return "", "", nil, "", "", fmt.Errorf("未知的渠道类型 %q", chType)
	}
	if chType == "smtp" {
		return "SMTP", "", nil, "", "", nil
	}

	// 把渠道配置本身也作为变量参与渲染（{{token}} / {{chatId}} …）
	merged := map[string]string{}
	for k, v := range values {
		merged[k] = v
	}
	for k, v := range cfg {
		if k == "headers" {
			continue
		}
		merged[k] = anyToString(v)
	}
	// 常见缺省
	if merged["server"] == "" {
		switch chType {
		case "bark":
			merged["server"] = "https://api.day.app"
		case "ntfy":
			merged["server"] = "https://ntfy.sh"
		}
	}
	merged["timestamp"] = values["timestamp"]
	merged["sign"] = values["sign"]
	if merged["group"] == "" {
		merged["group"] = "dockhelm"
	}
	if merged["mentionedJson"] == "" {
		merged["mentionedJson"] = buildMentionedJSON(merged["mentioned"])
	}
	if merged["topic"] == "" && chType == "pushplus" {
		merged["topic"] = ""
	}

	method := p.Method
	if m, ok := cfg["method"].(string); ok && strings.TrimSpace(m) != "" {
		method = strings.ToUpper(strings.TrimSpace(m))
	}
	rawURL := p.URL
	if u, ok := cfg["url"].(string); ok && chType == "webhook" && strings.TrimSpace(u) != "" {
		rawURL = strings.TrimSpace(u)
	}
	ct := p.ContentType
	if c, ok := cfg["contentType"].(string); ok && chType == "webhook" && strings.TrimSpace(c) != "" {
		ct = strings.TrimSpace(c)
	}
	bodyTpl := p.Body
	if b, ok := cfg["body"].(string); ok && chType == "webhook" && b != "" {
		bodyTpl = b
	}

	headers := map[string]string{}
	for k, v := range p.Headers {
		headers[k] = render(v, merged, false)
	}
	if p.AllowCustomHeaders {
		if raw, ok := cfg["headers"].(string); ok && strings.TrimSpace(raw) != "" {
			var custom map[string]string
			if err := json.Unmarshal([]byte(raw), &custom); err == nil {
				for k, v := range custom {
					headers[k] = render(v, merged, false)
				}
			}
		}
	}

	finalURL := render(rawURL, merged, true)
	finalBody := render(bodyTpl, merged, false)
	if ct == "" {
		ct = "application/json"
	}
	return method, finalURL, headers, finalBody, ct, nil
}

// buildMentionedJSON 把逗号分隔的手机号数组渲染成 JSON 数组字符串（企业微信用）。
func buildMentionedJSON(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "[]"
	}
	parts := []string{}
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			parts = append(parts, p)
		}
	}
	b, _ := json.Marshal(parts)
	return string(b)
}

func anyToString(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case bool:
		if t {
			return "true"
		}
		return "false"
	case float64:
		if t == float64(int64(t)) {
			return fmt.Sprintf("%d", int64(t))
		}
		return fmt.Sprintf("%v", t)
	default:
		b, err := json.Marshal(t)
		if err != nil {
			return fmt.Sprintf("%v", t)
		}
		return string(b)
	}
}

// PresetsSorted 预设按显示名排序（给前端下拉用）。
func PresetsSorted() []Preset {
	out := make([]Preset, len(Presets))
	copy(out, Presets)
	sort.Slice(out, func(i, j int) bool { return out[i].Label < out[j].Label })
	return out
}
