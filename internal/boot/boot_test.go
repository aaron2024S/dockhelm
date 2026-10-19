package boot

import "testing"

func infos(names ...string) []ContainerInfo {
	out := make([]ContainerInfo, 0, len(names))
	for _, n := range names {
		out = append(out, ContainerInfo{Name: n, State: "exited", Policy: "on-failure"})
	}
	return out
}

func TestNormalizeClamps(t *testing.T) {
	p := Plan{
		InitialDelaySec: -5,
		WaitTimeoutSec:  0,
		DefaultDelaySec: 999,
		Items: []Item{
			{Name: "  mysql  ", Auto: true, DelaySec: -1},
			{Name: "mysql", Auto: true, DelaySec: 5}, // 重名 → 只留第一条
			{Name: "", Auto: true, DelaySec: 3},      // 空名 → 丢掉
			{Name: "redis", Auto: true, DelaySec: 999},
		},
	}
	got := p.Normalize()
	if got.InitialDelaySec != DefaultInitialDelaySec {
		t.Errorf("初始等待应回落到默认值，得到 %d", got.InitialDelaySec)
	}
	if got.WaitTimeoutSec != DefaultWaitTimeoutSec {
		t.Errorf("超时应回落到默认值，得到 %d", got.WaitTimeoutSec)
	}
	if got.DefaultDelaySec != DefaultDelaySec {
		t.Errorf("默认延时应回落到 %d，得到 %d", DefaultDelaySec, got.DefaultDelaySec)
	}
	if len(got.Items) != 2 {
		t.Fatalf("条目应去重去空后剩 2 条，得到 %d：%+v", len(got.Items), got.Items)
	}
	if got.Items[0].Name != "mysql" {
		t.Errorf("名字应去掉首尾空白，得到 %q", got.Items[0].Name)
	}
	if got.Items[0].DelaySec != got.DefaultDelaySec {
		t.Errorf("负延时应收敛到默认值，得到 %d", got.Items[0].DelaySec)
	}
	if got.Items[1].DelaySec != maxDelaySec {
		t.Errorf("超大延时应收敛到 %d，得到 %d", maxDelaySec, got.Items[1].DelaySec)
	}
}

// 首次打开页面：没有任何条目时，按「compose 项目分组 + 项目内按名称」生成初始顺序，
// 没有项目的散装容器排在最后（方案 §八 第 3 条）。
func TestMergeInitialOrder(t *testing.T) {
	p := DefaultPlan()
	list := []ContainerInfo{
		{Name: "emby", Policy: "on-failure", Project: "media"},
		{Name: "qbittorrent", Policy: "on-failure", Project: "media"},
		{Name: "nginx", Policy: "on-failure"}, // 散装
		{Name: "mysql", Policy: "on-failure", Project: "db"},
		{Name: "adguard", Policy: "on-failure", Project: "db"},
	}
	got := p.Merge(list)
	want := []string{"adguard", "mysql", "emby", "qbittorrent", "nginx"}
	if len(got.Items) != len(want) {
		t.Fatalf("条目数不对：%+v", got.Items)
	}
	for i, n := range want {
		if got.Items[i].Name != n {
			t.Errorf("第 %d 位应为 %s，得到 %s（整体：%+v）", i+1, n, got.Items[i].Name, got.Items)
		}
	}
	for _, it := range got.Items {
		if !it.Auto {
			t.Errorf("新容器默认应开机启动：%+v", it)
		}
		if it.DelaySec != DefaultDelaySec {
			t.Errorf("新容器默认延时应为 %d：%+v", DefaultDelaySec, it)
		}
	}
}

// 已存在的条目顺序、开关、延时都要原样保留；被删掉的剔除；新出现的追加末尾。
func TestMergeKeepsOrder(t *testing.T) {
	p := DefaultPlan()
	p.Items = []Item{
		{Name: "mysql", Auto: true, DelaySec: 20},
		{Name: "redis", Auto: false, DelaySec: 0},
		{Name: "gone", Auto: true, DelaySec: 3}, // 已被删除
	}
	got := p.Merge([]ContainerInfo{
		{Name: "redis", State: "exited", Policy: "on-failure"},
		{Name: "mysql", State: "exited", Policy: "on-failure"},
		{Name: "newbie", State: "exited", Policy: "on-failure"},
	})
	if len(got.Items) != 3 {
		t.Fatalf("应为 3 条，得到 %+v", got.Items)
	}
	if got.Items[0].Name != "mysql" || got.Items[0].DelaySec != 20 {
		t.Errorf("原有条目应保留顺序与延时：%+v", got.Items[0])
	}
	if got.Items[1].Name != "redis" || got.Items[1].Auto {
		t.Errorf("关闭的开关不能被重置：%+v", got.Items[1])
	}
	if got.Items[2].Name != "newbie" {
		t.Errorf("新容器应追加末尾：%+v", got.Items[2])
	}
}

func TestResolveKinds(t *testing.T) {
	p := DefaultPlan()
	p.Items = []Item{
		{Name: "mysql", Auto: true, DelaySec: 20}, // 已停、受管 → 启动
		{Name: "redis", Auto: true, DelaySec: 3},  // 已在跑 → 跳过（计入）
		{Name: "qb", Auto: false, DelaySec: 3},    // 开关关着 → 不计入
		{Name: "emby", Auto: true, DelaySec: 3},   // unless-stopped → 不碰、不计入
		{Name: "ghost", Auto: true, DelaySec: 3},  // 不存在 → 记失败
	}
	list := []ContainerInfo{
		{Name: "mysql", State: "exited", Policy: "on-failure"},
		{Name: "redis", State: "running", Policy: "no"},
		{Name: "qb", State: "exited", Policy: "on-failure"},
		{Name: "emby", State: "running", Policy: "unless-stopped"},
	}
	steps := Resolve(p, list)
	want := []StepKind{StepStart, StepRunning, StepDisabled, StepUnmanaged, StepMissing}
	if len(steps) != len(want) {
		t.Fatalf("步数不对：%+v", steps)
	}
	for i, k := range want {
		if steps[i].Kind != k {
			t.Errorf("%s 应为 %s，得到 %s", steps[i].Name, k, steps[i].Kind)
		}
	}

	// 统计口径：只有 start / running / missing 计入「N 个容器」
	counted := 0
	for _, st := range steps {
		if st.Counted() {
			counted++
		}
	}
	if counted != 3 {
		t.Errorf("应计入 3 个（mysql / redis / ghost），得到 %d", counted)
	}
}

func TestEstimateSeconds(t *testing.T) {
	p := DefaultPlan()
	p.InitialDelaySec = 30
	p.Items = []Item{
		{Name: "a", Auto: true, DelaySec: 10},
		{Name: "b", Auto: true, DelaySec: 3},
		{Name: "c", Auto: true, DelaySec: 0},
	}
	steps := Resolve(p, []ContainerInfo{
		{Name: "a", State: "exited", Policy: "on-failure"},
		{Name: "b", State: "exited", Policy: "on-failure"},
		{Name: "c", State: "running", Policy: "on-failure"}, // 已在跑，不计耗时
	})
	// 30 + (8+10) + (8+3)
	if got := EstimateSeconds(p, steps); got != 59 {
		t.Errorf("预计耗时应为 59 秒，得到 %d", got)
	}
}

// 运行记录的三句话（方案 §六）。本功能在 23:xx 从「编排」改名为「顺序启动」，
// 文案随之统一，不再出现「编排」二字。
func TestSummary(t *testing.T) {
	cases := []struct {
		name            string
		total, ok, skip int
		failed          []string
		want            string
	}{
		{"全部正常", 50, 47, 3, nil, "顺序启动 50 个容器，成功 47、跳过 3"},
		{"有失败", 50, 45, 3, []string{"nginx", "foo"}, "顺序启动 50 个容器，成功 45、跳过 3、失败 2：foo、nginx"},
		{"全跳过", 50, 0, 50, nil, "顺序启动：50 个容器都已在运行，本次未启动任何容器"},
		{"没有受管容器", 0, 0, 0, nil, "顺序启动：没有需要启动的容器"},
	}
	for _, c := range cases {
		if got := Summary(c.total, c.ok, c.skip, c.failed); got != c.want {
			t.Errorf("%s：\n  want %s\n  got  %s", c.name, c.want, got)
		}
	}
}

func TestSummaryFoldsLongList(t *testing.T) {
	failed := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	got := Summary(10, 0, 0, failed)
	want := "顺序启动 10 个容器，成功 0、跳过 0、失败 10：a、b、c、d、e、f、g、h…等 10 个"
	if got != want {
		t.Errorf("\n  want %s\n  got  %s", want, got)
	}
}
