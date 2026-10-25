package boot

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aaron2024s/dockhelm/internal/bus"
	"github.com/aaron2024s/dockhelm/internal/dockerx"
	"github.com/aaron2024s/dockhelm/internal/store"
)

// fakeDocker 一个够用的假守护进程：列表 / inspect / start。
// 只用标准库，不依赖真 Docker —— 本机没有引擎，这条链路只能这样验。
type fakeDocker struct {
	mu       sync.Mutex
	order    []string          // start 调用顺序
	state    map[string]string // name → running/exited
	policies map[string]string // name → restartPolicy
	failOn   map[string]string // name → 启动时报的错
	srv      *httptest.Server
}

func newFakeDocker(t *testing.T) *fakeDocker {
	f := &fakeDocker{
		state:    map[string]string{},
		policies: map[string]string{},
		failOn:   map[string]string{},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/_ping", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("OK")) })
	mux.HandleFunc("/containers/json", func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		out := make([]map[string]any, 0, len(f.state))
		for name, st := range f.state {
			out = append(out, map[string]any{
				"Id": "id-" + name, "Names": []string{"/" + name},
				"Image": "img:latest", "State": st, "Labels": map[string]string{},
			})
		}
		writeJSONTest(w, out)
	})
	mux.HandleFunc("/containers/", func(w http.ResponseWriter, r *http.Request) {
		// 路径形如 /containers/{name}/json 或 /containers/{name}/start
		rest := strings.TrimPrefix(r.URL.Path, "/containers/")
		parts := strings.Split(rest, "/")
		if len(parts) != 2 {
			http.NotFound(w, r)
			return
		}
		name, action := parts[0], parts[1]
		// 假守护进程里 ID 用 "id-" + 名字，收到 ID 时还原成名字再查表
		name = strings.TrimPrefix(name, "id-")
		f.mu.Lock()
		defer f.mu.Unlock()
		if action == "start" {
			if msg, bad := f.failOn[name]; bad {
				w.WriteHeader(http.StatusInternalServerError)
				writeJSONTest(w, map[string]string{"message": msg})
				return
			}
			f.order = append(f.order, name)
			f.state[name] = "running"
			w.WriteHeader(http.StatusNoContent)
			return
		}
		if action != "json" {
			http.NotFound(w, r)
			return
		}
		writeJSONTest(w, map[string]any{
			"Id": "id-" + name, "Name": "/" + name,
			"Config": map[string]any{"Image": "img:latest", "Labels": map[string]string{}},
			"State":  map[string]any{"Status": f.state[name], "Running": f.state[name] == "running"},
			"HostConfig": map[string]any{
				"RestartPolicy": map[string]any{"Name": f.policies[name], "MaximumRetryCount": 0},
			},
		})
	})
	f.srv = httptest.NewServer(mux)
	t.Cleanup(f.srv.Close)
	return f
}

func writeJSONTest(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

func newTestService(t *testing.T) (*Service, *store.Store, *fakeDocker) {
	t.Helper()
	fd := newFakeDocker(t)
	dc, err := dockerx.New("tcp://" + strings.TrimPrefix(fd.srv.URL, "http://"))
	if err != nil {
		t.Fatalf("创建客户端失败：%v", err)
	}
	st, err := store.Open(filepath.Join(t.TempDir(), "dockhelm.json"))
	if err != nil {
		t.Fatalf("打开 store 失败：%v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return New(dc, st, nil, func() string { return "dockhelm" }, nil), st, fd
}

// 一轮完整的执行：跳过的跳过、失败的继续、最后写一条自洽的记录。
func TestRunStartsSeriallyAndKeepsGoing(t *testing.T) {
	svc, st, fd := newTestService(t)
	fd.state = map[string]string{"alpha": "exited", "beta": "running", "boom": "exited", "zulu": "running"}
	fd.policies = map[string]string{
		"alpha": "on-failure", "beta": "no", "boom": "on-failure", "zulu": "unless-stopped",
	}
	fd.failOn = map[string]string{"boom": "port is already allocated"}

	plan := Plan{
		Enabled: true, InitialDelaySec: 0, WaitTimeoutSec: 5, DefaultDelaySec: 0,
		Items: []Item{
			{Name: "alpha", Auto: true, DelaySec: 0},
			{Name: "beta", Auto: true, DelaySec: 0},
			{Name: "boom", Auto: true, DelaySec: 0},
			{Name: "zulu", Auto: true, DelaySec: 0}, // unless-stopped：不归本功能管
		},
	}
	if err := svc.SavePlan(plan); err != nil {
		t.Fatalf("保存失败：%v", err)
	}

	if err := svc.Run(context.Background(), TriggerManual, true); err != nil {
		t.Fatalf("执行返回错误：%v", err)
	}

	if got := strings.Join(fd.order, ","); got != "alpha" {
		t.Errorf("只应启动 alpha（beta 在跑、boom 失败、zulu 不归本功能管），实际：%s", got)
	}

	stt := svc.Status()
	if stt.Running {
		t.Error("跑完之后 running 应为 false")
	}
	if stt.Total != 3 || stt.Success != 1 || stt.Skipped != 1 || stt.Failed != 1 {
		t.Errorf("统计不对：total=%d success=%d skipped=%d failed=%d", stt.Total, stt.Success, stt.Skipped, stt.Failed)
	}

	logs, _ := st.ListRunLogs(10)
	if len(logs) != 1 {
		t.Fatalf("应写 1 条运行记录，实际 %d 条", len(logs))
	}
	want := "顺序启动 3 个容器，成功 1、跳过 1、失败 1：boom"
	if logs[0].Message != want {
		t.Errorf("记录文案不对：\n  want %s\n  got  %s", want, logs[0].Message)
	}
	if logs[0].Kind != "boot" || logs[0].Status != "failed" {
		t.Errorf("记录类型/状态不对：%+v", logs[0])
	}
	if !strings.Contains(logs[0].Detail, "port is already allocated") {
		t.Errorf("详情里应带守护进程的原因：%q", logs[0].Detail)
	}

	lr, ok := svc.LastResult()
	if !ok || lr.Success != 1 || lr.Failed != 1 {
		t.Errorf("结构化回执不对：%+v（存在=%v）", lr, ok)
	}
}

// 全都在跑：不启动任何容器，记录换成「都已在运行」那句。
func TestRunAllSkipped(t *testing.T) {
	svc, st, fd := newTestService(t)
	fd.state = map[string]string{"a": "running", "b": "running"}
	fd.policies = map[string]string{"a": "on-failure", "b": "on-failure"}
	_ = svc.SavePlan(Plan{
		Enabled: true, WaitTimeoutSec: 5, DefaultDelaySec: 0,
		Items: []Item{{Name: "a", Auto: true}, {Name: "b", Auto: true}},
	})
	if err := svc.Run(context.Background(), TriggerManual, true); err != nil {
		t.Fatalf("执行返回错误：%v", err)
	}
	if len(fd.order) != 0 {
		t.Errorf("不该启动任何容器，实际：%v", fd.order)
	}
	logs, _ := st.ListRunLogs(5)
	if len(logs) != 1 || logs[0].Message != "顺序启动：2 个容器都已在运行，本次未启动任何容器" {
		t.Errorf("记录文案不对：%+v", logs)
	}
	if logs[0].Status != "success" {
		t.Errorf("全跳过应记成功，实际 %s", logs[0].Status)
	}
}

// 互斥：执行中再来一轮直接拒绝。
func TestRunRejectsConcurrent(t *testing.T) {
	svc, _, _ := newTestService(t)
	if !svc.begin(TriggerManual) {
		t.Fatal("第一次 begin 应成功")
	}
	defer svc.end()
	if err := svc.Run(context.Background(), TriggerManual, true); err != ErrBusy {
		t.Errorf("应返回 ErrBusy，实际 %v", err)
	}
}

// 总开关关掉时，开机触发不执行；手动触发照跑（用户是按了按钮才来的）。
func TestRunDisabledSkipsBootOnly(t *testing.T) {
	svc, _, fd := newTestService(t)
	fd.state = map[string]string{"a": "exited"}
	fd.policies = map[string]string{"a": "on-failure"}
	_ = svc.SavePlan(Plan{
		Enabled: false, WaitTimeoutSec: 5, DefaultDelaySec: 0,
		Items: []Item{{Name: "a", Auto: true}},
	})
	if err := svc.Run(context.Background(), TriggerBoot, true); err != nil {
		t.Fatalf("开关关掉时不应报错：%v", err)
	}
	if len(fd.order) != 0 {
		t.Errorf("开关关掉时不该启动容器：%v", fd.order)
	}
}

// 延时语义：这一行启动成功后等 N 秒才起下一个（不是启动前等）。
func TestRunDelayAfterEachStart(t *testing.T) {
	svc, _, fd := newTestService(t)
	fd.state = map[string]string{"a": "exited", "b": "exited"}
	fd.policies = map[string]string{"a": "on-failure", "b": "on-failure"}
	_ = svc.SavePlan(Plan{
		Enabled: true, WaitTimeoutSec: 5, DefaultDelaySec: 0,
		Items: []Item{{Name: "a", Auto: true, DelaySec: 1}, {Name: "b", Auto: true, DelaySec: 0}},
	})
	start := time.Now()
	if err := svc.Run(context.Background(), TriggerManual, true); err != nil {
		t.Fatalf("执行返回错误：%v", err)
	}
	if elapsed := time.Since(start); elapsed < 900*time.Millisecond {
		t.Errorf("a 的 1 秒延时应等满，实际只用了 %v", elapsed)
	}
	if got := strings.Join(fd.order, ","); got != "a,b" {
		t.Errorf("启动顺序应为 a,b，实际 %s", got)
	}
}

// 每启动一个容器都要推一条 container/changed —— 总览页靠它刷新「运行 / 停止」计数
// 和容器状态条；顺序启动不经过 API 的动作接口，没有这条事件页面就一直停在旧数字上。
func TestRunPublishesContainerChanged(t *testing.T) {
	svc, _, fd := newTestService(t)
	b := bus.New()
	svc.bus = b
	ch, cancel := b.Subscribe()
	defer cancel()

	fd.state = map[string]string{"alpha": "exited", "beta": "exited"}
	fd.policies = map[string]string{"alpha": "on-failure", "beta": "on-failure"}
	fd.failOn = map[string]string{"beta": "port is already allocated"}
	_ = svc.SavePlan(Plan{
		Enabled: true, WaitTimeoutSec: 5, DefaultDelaySec: 0,
		Items: []Item{{Name: "alpha", Auto: true}, {Name: "beta", Auto: true}},
	})
	if err := svc.Run(context.Background(), TriggerManual, true); err != nil {
		t.Fatalf("执行返回错误：%v", err)
	}

	got := map[string]string{}
	for {
		select {
		case ev := <-ch:
			if ev.Topic == "container" && ev.Kind == "changed" {
				got[ev.Data["name"].(string)] = ev.Status
			}
			continue
		default:
		}
		break
	}
	if got["alpha"] != "success" {
		t.Errorf("alpha 应推 success，实际 %q（全部：%v）", got["alpha"], got)
	}
	if got["beta"] != "failed" {
		t.Errorf("beta 应推 failed，实际 %q（全部：%v）", got["beta"], got)
	}
}
