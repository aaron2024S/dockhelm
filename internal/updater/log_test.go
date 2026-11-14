package updater

import (
	"strings"
	"testing"
)

type logRecord struct{ kind, ref, status, message, detail string }

// collect 用零值 Updater + 日志钩子接住 logBatch 写出的记录。
func collect(t *testing.T, out []Result) []logRecord {
	t.Helper()
	var got []logRecord
	u := &Updater{}
	u.SetLogSink(func(kind, ref, status, message, detail string) {
		got = append(got, logRecord{kind, ref, status, message, detail})
	})
	u.logBatch(out)
	return got
}

// 一轮更新 2 个容器：只许出现 1 条记录，且说明里点名两个容器。
//
// 2026-11-14 用户反馈：同一分钟里出现「ntfy 镜像已更新」「rsshub 镜像已更新」
// 「自动更新完成」三行，问「一个时间更新了 2 个镜像，为什么不写成一条记录」。
func TestLogBatchMergesIntoOneRecord(t *testing.T) {
	got := collect(t, []Result{
		{Container: "ntfy", Status: ResultUpdated, Steps: []string{"开始处理 ntfy"}},
		{Container: "rsshub", Status: ResultUpdated, Steps: []string{"开始处理 rsshub"}},
	})
	if len(got) != 1 {
		t.Fatalf("一轮更新应只写 1 条记录，实际 %d 条：%+v", len(got), got)
	}
	r := got[0]
	if r.kind != "update" || r.status != "success" {
		t.Errorf("kind/status = %q/%q", r.kind, r.status)
	}
	if r.ref != "batch" {
		t.Errorf("多台记录的对象列应为 batch 占位，实际 %q", r.ref)
	}
	if r.message != "共 2 个容器：更新了 ntfy、rsshub" {
		t.Errorf("message = %q", r.message)
	}
	// 逐台过程不能丢，排查时还要能看
	if !strings.Contains(r.detail, "【ntfy】") || !strings.Contains(r.detail, "【rsshub】") {
		t.Errorf("detail 应含每台的段落：%q", r.detail)
	}
}

// 单台成功保持历史文案（对象列是容器名，说明是「镜像已更新」）。
func TestLogBatchSingleKeepsOldWording(t *testing.T) {
	got := collect(t, []Result{{Container: "ntfy", Status: ResultUpdated}})
	if len(got) != 1 {
		t.Fatalf("want 1 record, got %d", len(got))
	}
	if got[0].ref != "ntfy" || got[0].message != "镜像已更新" || got[0].status != "success" {
		t.Errorf("单台文案变了：%+v", got[0])
	}
}

// 部分失败：一条记录说清总数、成功名单、失败名单（原因压成短语）。
func TestLogBatchFailureNamesTheFailed(t *testing.T) {
	got := collect(t, []Result{
		{Container: "ntfy", Status: ResultUpdated},
		{Container: "rsshub", Status: ResultFailed, Message: "拉取失败：pull failed: read tcp [240e::1]:443: connection reset"},
	})
	if len(got) != 1 {
		t.Fatalf("want 1 record, got %d", len(got))
	}
	r := got[0]
	if r.status != "failed" {
		t.Errorf("status = %q, want failed", r.status)
	}
	if r.message != "共 2 个容器：更新了 ntfy，失败 rsshub（拉取失败）" {
		t.Errorf("message = %q", r.message)
	}
}

// 全部已是最新：一条记录，不逐台刷「镜像已是最新」。
func TestLogBatchAllUpToDate(t *testing.T) {
	got := collect(t, []Result{
		{Container: "a", Status: ResultUpToDate},
		{Container: "b", Status: ResultUpToDate},
		{Container: "c", Status: ResultUpToDate},
	})
	if len(got) != 1 {
		t.Fatalf("want 1 record, got %d", len(got))
	}
	if got[0].status != "up_to_date" || got[0].message != "共 3 个容器：镜像均是最新" {
		t.Errorf("记录 = %+v", got[0])
	}
}

// 更新 + 跳过混合：点名更新的，另外折算个数。
func TestLogBatchMixedSkipped(t *testing.T) {
	got := collect(t, []Result{
		{Container: "ntfy", Status: ResultUpdated},
		{Container: "x", Status: ResultUpToDate},
		{Container: "y", Status: ResultSkipped},
	})
	if len(got) != 1 {
		t.Fatalf("want 1 record, got %d", len(got))
	}
	if got[0].message != "共 3 个容器：更新了 ntfy，另有 2 个已是最新" {
		t.Errorf("message = %q", got[0].message)
	}
}

// 空结果集不写记录（比如全部被排除，调用方提前返回）。
func TestLogBatchEmpty(t *testing.T) {
	if got := collect(t, nil); len(got) != 0 {
		t.Errorf("空结果集不该写记录，实际 %d 条", len(got))
	}
}
