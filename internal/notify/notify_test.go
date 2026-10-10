package notify

import (
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
