package backup

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestProjectInfoJSONShape 钉住 compose 项目接口的 JSON 契约。
//
// 背景（0.3.0 真实事故）：有 compose 项目、但一个文件都读不到时，
// 后端返回的 `readable` / `unreadable` 是 nil 切片 → JSON `null`，
// 前端 `p.readable.length` 抛 TypeError，整页渲染白屏（用户看到「点进去一片空白」）。
//
// 这里从**生产用的构造函数**出发断言：任何字段都不得序列化成 null。
func TestProjectInfoJSONShape(t *testing.T) {
	b, err := json.Marshal(newProjectInfo("demo"))
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	if strings.Contains(got, "null") {
		t.Fatalf("项目字段不得为 null（前端会取 .length）：%s", got)
	}
	for _, want := range []string{`"containers":[]`, `"configFiles":[]`, `"readable":[]`, `"unreadable":[]`} {
		if !strings.Contains(got, want) {
			t.Errorf("缺少 %s，实际：%s", want, got)
		}
	}
}

// TestListProjectsEmptyResultShape 列表本身也不能是 null。
func TestListProjectsEmptyResultShape(t *testing.T) {
	out := make([]ProjectInfo, 0, 0)
	b, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "[]" {
		t.Fatalf("空项目列表应序列化成 []，实际 %s", string(b))
	}
}
