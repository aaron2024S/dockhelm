package dockerx

import (
	"encoding/json"
	"testing"
)

func TestSelfMountMap(t *testing.T) {
	raw := `{"Mounts":[
		{"Type":"bind","Source":"/volume1/docker","Destination":"/host/docker"},
		{"Type":"volume","Source":"/var/lib/docker/volumes/appdata/_data","Destination":"/data"},
		{"Source":"","Destination":"/nope"},
		{"Source":"/only-source"}
	]}`
	var insp map[string]any
	if err := json.Unmarshal([]byte(raw), &insp); err != nil {
		t.Fatal(err)
	}

	m := SelfMountMap(insp)
	if len(m) != 2 {
		t.Fatalf("应当只提取到 2 条有效映射，实际 %d 条：%+v", len(m), m)
	}
	if m["/volume1/docker"] != "/host/docker" {
		t.Fatalf("绑定挂载映射不对：%+v", m)
	}
	if m["/var/lib/docker/volumes/appdata/_data"] != "/data" {
		t.Fatalf("具名卷映射不对：%+v", m)
	}
}

func TestSelfMountMapWithoutMounts(t *testing.T) {
	if m := SelfMountMap(map[string]any{}); len(m) != 0 {
		t.Fatalf("没有 Mounts 时应返回空映射，实际 %+v", m)
	}
}
