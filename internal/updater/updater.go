// Package updater 是 Dockhelm 的核心：镜像更新。
//
// 它与 dockerCopilot 最本质的区别只有一条，但决定了「没更新的容器会不会被停掉」：
//
//	dockerCopilot：pull → stop → rename → create → start   （无条件停容器）
//	Dockhelm     ：pull → 比对镜像 ID → 没变就【直接结束，容器一个字节都不碰】
//	                                → 变了才 stop → rename → create → start
//
// 第二条同样重要：**检测（有没有新版本）与拉取永远同源**。检测侧走的是守护进程的
// /distribution 接口（与 docker pull 同一套仓库端点解析），异常时一律标记为「未知」，
// 绝不退化成「有新版本」—— 这正是 dockerCopilot 永久误报的成因。
package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aaron2024s/dockhelm/internal/bus"
	"github.com/aaron2024s/dockhelm/internal/dockerx"
)

// Notifier 通知钩子（用接口避免与 notify 包循环依赖）。
type Notifier interface {
	// Emit 投递一个事件；实现必须保证非阻塞、失败只记日志。
	Emit(event string, vars map[string]string)
}

// noopNotifier 空实现。
type noopNotifier struct{}

func (noopNotifier) Emit(string, map[string]string) {}

// Options 更新引擎参数。
type Options struct {
	// Concurrency 批量更新时的并发度（默认 2）。注册表限流时调小。
	Concurrency int
	// HealthTimeout 重建后等待健康的超时（默认 90s）。
	HealthTimeout time.Duration
	// StableSeconds 没有 healthcheck 时，启动后需保持运行多久才算成功（默认 6s）。
	StableSeconds int
	// KeepBackups 每个容器保留多少个 __bak_ 快照容器（默认 0 = 成功即删）。
	KeepBackups int
	// BackupDir 配置快照目录。
	BackupDir string

	// PullOnce 同一轮批量更新里，同一个镜像只真正拉取一次（默认开）。
	// 4 个容器共用 nginx:alpine 时下载 1 次、重建 4 个 —— 镜像 ID 是内容寻址的，
	// 复用第一次的结果完全安全。
	PullOnce bool
	// BackupBefore 重建前先写一份配置快照（默认开）。这是失败回滚与人肉排查的底牌。
	BackupBefore bool
	// CleanupAfter 更新成功后，确认没有任何容器再引用旧镜像时把它删掉（默认开）。
	CleanupAfter bool
	// DirectFirst 显式域名（ghcr.io、私有仓库等）不套用加速源（默认开）。
	// 关掉它意味着「所有镜像都先走加速站」，只有把加速站当成通用代理时才该这么做。
	DirectFirst bool
}

// DefaultOptions 返回默认参数。
func DefaultOptions(backupDir string) Options {
	return Options{
		Concurrency:   2,
		HealthTimeout: 90 * time.Second,
		StableSeconds: 6,
		KeepBackups:   0,
		BackupDir:     backupDir,
		PullOnce:      true,
		BackupBefore:  true,
		CleanupAfter:  true,
		DirectFirst:   true,
	}
}

// Policy 设置页里那一组更新开关。零值（false）是有意义的取值，
// 所以调用方要显式传全部字段。
type Policy struct {
	Concurrency  int
	PullOnce     bool
	BackupBefore bool
	CleanupAfter bool
	DirectFirst  bool
}

// Apply 把设置页的策略应用到引擎参数上。并发度越界时退回 2（与设置页校验一致）。
func (o *Options) Apply(p Policy) {
	if p.Concurrency >= 1 && p.Concurrency <= 8 {
		o.Concurrency = p.Concurrency
	} else {
		o.Concurrency = 2
	}
	o.PullOnce = p.PullOnce
	o.BackupBefore = p.BackupBefore
	o.CleanupAfter = p.CleanupAfter
	o.DirectFirst = p.DirectFirst
}

// Updater 更新引擎。
type Updater struct {
	dc       *dockerx.Client
	bus      *bus.Bus
	notify   Notifier
	selfName string
	selfID   string

	// opts 用原子指针存：设置页随时可以改（并发度、那几个策略开关），
	// 而更新流程有多条 goroutine 在读它 —— 普通字段会构成数据竞争。
	opts atomic.Pointer[Options]

	mu    sync.Mutex
	locks map[string]*sync.Mutex // 逐容器串行锁（键 = 容器名，重建前后都稳定）
	// resolveMu 只包住「把名字/ID 解析成容器名」这一次 Inspect，
	// 保证两个并发调用拿到的是同一把锁（见 lockForContainer）。
	resolveMu sync.Mutex
	onLog     func(kind, ref, status, message, detail string)
	nowStr    func() string
	// mirrorFn 返回用户配置的「拉取加速源」；空字符串表示完全交给守护进程。
	mirrorFn func() string

	// batchMu/batchPull 是「同一轮批量更新里同一镜像只拉一次」的缓存。
	// 只在 UpdateMany 的范围内有效，批量结束就清空 —— 不能让上一次的结论
	// 漏到下一次，那会让「刚推了新版本却检测不到」变成偶发问题。
	batchMu   sync.Mutex
	batchPull map[string]*dockerx.PullResult
}

// New 创建更新引擎。
func New(dc *dockerx.Client, b *bus.Bus, opts Options) *Updater {
	u := &Updater{
		dc:     dc,
		bus:    b,
		notify: noopNotifier{},
		locks:  map[string]*sync.Mutex{},
		nowStr: func() string { return time.Now().Format("2006-01-02 15:04:05") },
	}
	u.opts.Store(&opts)
	return u
}

// opt 取当前参数的一份快照。Options 是小结构体，复制比加锁便宜也更不容易出错。
func (u *Updater) opt() Options {
	if p := u.opts.Load(); p != nil {
		return *p
	}
	return Options{}
}

// SetNotifier 注入通知器。
func (u *Updater) SetNotifier(n Notifier) {
	if n == nil {
		n = noopNotifier{}
	}
	u.notify = n
}

// SetLogSink 注入运行记录写入器。
func (u *Updater) SetLogSink(fn func(kind, ref, status, message, detail string)) {
	u.onLog = fn
}

// SetSelf 告知引擎「哪个容器是 Dockhelm 自己」，永不更新。
func (u *Updater) SetSelf(name, id string) {
	u.selfName = name
	u.selfID = id
}

// Options 返回当前参数的一份快照。
func (u *Updater) Options() Options { return u.opt() }

// SetOptions 换一组新参数。设置页改完开关立刻生效，不必重启进程。
//
// Concurrency 会被正在跑的 UpdateMany 在启动时读走，所以改并发度只影响
// 「下一次批量更新」—— 这也是设置页里文案的原话。
func (u *Updater) SetOptions(opt Options) { u.opts.Store(&opt) }

// SetMirror 注入「拉取加速源」解析函数。
func (u *Updater) SetMirror(fn func() string) { u.mirrorFn = fn }

// pullImage 拉取镜像。
//
// 若用户在设置里指定了加速源、且镜像来自 Docker Hub、且引用里没写死 digest，
// 就用「<加速站>/<仓库>:<标签>」去拉，拉完再打回原始标签 —— 这样即使用户
// 没法改 daemon.json，也能让 Dockhelm 走自己选的加速站，而 compose 与其它工具
// 仍然按原来的名字找得到镜像。
//
// 无论走哪条路，**最终判定「有没有变化」用的都是内容寻址的镜像 ID**，
// 所以检测与拉取永远同源，不会出现 dockerCopilot 那种永久误报。
func (u *Updater) pullImage(ctx context.Context, imageRef string, onEvent func(dockerx.PullEvent)) (*dockerx.PullResult, error) {
	mirror := ""
	if u.mirrorFn != nil {
		mirror = strings.TrimSpace(u.mirrorFn())
	}
	ref := dockerx.ParseRef(imageRef)
	// DirectFirst：显式写了域名的镜像（ghcr.io、私有仓库…）直连，不套加速。
	// 关掉这个开关等于把加速站当成通用代理，对非 Docker Hub 的镜像也去试一把。
	if mirror == "" || (u.opt().DirectFirst && ref.Registry != "docker.io") || ref.Digest != "" {
		return u.dc.Pull(ctx, imageRef, onEvent)
	}
	host := mirrorHost(mirror)
	if host == "" {
		return u.dc.Pull(ctx, imageRef, onEvent)
	}
	target := fmt.Sprintf("%s/%s:%s", host, ref.LocalName, ref.Tag)
	res, err := u.dc.Pull(ctx, target, onEvent)
	if err != nil {
		return res, err
	}
	if res.ImageID != "" {
		if terr := u.dc.TagImage(ctx, target, ref.LocalName, ref.Tag); terr != nil {
			return res, fmt.Errorf("镜像已拉取，但回写原始标签 %s 失败：%w", imageRef, terr)
		}
	}
	res.Ref = imageRef
	return res, nil
}

// mirrorHost 去掉协议前缀与尾部斜杠，得到可以直接拼进镜像名的域名。
func mirrorHost(mirror string) string {
	m := strings.TrimSpace(strings.TrimSuffix(mirror, "/"))
	m = strings.TrimPrefix(m, "https://")
	m = strings.TrimPrefix(m, "http://")
	return m
}

// pullForBatch 在一轮批量更新内对同一个镜像只真正拉取一次。
//
// 返回的 fromCache 表示这次是复用了别的容器拉到的结果 —— 那种情况下
// **不能**采信 pr.UpToDate：那是「第一次拉取时镜像是否已经最新」的结论，
// 对另一个还跑着旧镜像的容器并不成立，采信它就会漏掉一次本该做的重建。
func (u *Updater) pullForBatch(
	ctx context.Context,
	image string,
	onEvent func(dockerx.PullEvent),
) (pr *dockerx.PullResult, fromCache bool, err error) {
	if !u.opt().PullOnce {
		pr, err = u.pullImage(ctx, image, onEvent)
		return pr, false, err
	}

	u.batchMu.Lock()
	cached, hit := u.batchPull[image]
	u.batchMu.Unlock()
	if hit {
		return cached, true, nil
	}

	pr, err = u.pullImage(ctx, image, onEvent)
	if err != nil {
		return nil, false, err
	}
	u.batchMu.Lock()
	if u.batchPull == nil {
		u.batchPull = map[string]*dockerx.PullResult{}
	}
	u.batchPull[image] = pr
	u.batchMu.Unlock()
	return pr, false, nil
}

// resetBatchPull 清空批量拉取缓存。批量开始与结束都要调。
func (u *Updater) resetBatchPull() {
	u.batchMu.Lock()
	u.batchPull = nil
	u.batchMu.Unlock()
}

// SelfName 返回自身容器名。
func (u *Updater) SelfName() string { return u.selfName }

func (u *Updater) log(kind, ref, status, message, detail string) {
	if u.onLog != nil {
		u.onLog(kind, ref, status, message, detail)
	}
}

// lockFor 取某容器的串行锁（同一容器不会被两个任务同时更新）。
//
// key 必须是**容器名**。调用方不要直接用它 —— 用 lockForContainer，
// 它负责把「名字或 ID」统一换成名字，原因见那里的注释。
func (u *Updater) lockFor(key string) *sync.Mutex {
	u.mu.Lock()
	defer u.mu.Unlock()
	m, ok := u.locks[key]
	if !ok {
		m = &sync.Mutex{}
		u.locks[key] = m
	}
	return m
}

// lockForContainer 解析出容器名并返回它对应的串行锁。
//
// **锁键必须是容器名，不能是 ID**，两个原因：
//
//  1. 同一个容器既可能被名字引用（界面上的批量更新按名字），也可能被 ID 引用
//     （计划任务用的是容器列表里的 ID）。以前直接把调用方传进来的字符串当锁键，
//     于是 "web" 与 "a1b2c3…" 会拿到两把互不相干的锁 —— 等于没有互斥。真撞上时
//     （计划任务到点 + 用户手点更新），同一个容器会被并发执行 stop→rename→create：
//     轻则 rename 撞名失败，重则留下 `xxx__bak_` 悬挂容器 + 没启起来的新容器。
//
//  2. 更新本身就是「重建」：更新完容器会拿到一个**新的 ID**。若用 ID 当锁键，
//     第二个调用只要解析得晚一点就会拿到新 ID 的锁，与仍在收尾的第一个调用
//     并发操作同一个容器 —— 名字是重建前后都不变的那一个标识。
//
// 解析只用一把很短的全局锁包住一次 Inspect（本地 socket，几毫秒）。
//
// 注意：拿到锁之后调用方**必须重新 Inspect 一次**。这里读到的状态是取锁之前那一刻的，
// 前一个持锁者可能已经把容器改名/重建过了。
func (u *Updater) lockForContainer(ctx context.Context, nameOrID string) (*sync.Mutex, error) {
	u.resolveMu.Lock()
	defer u.resolveMu.Unlock()
	insp, err := u.dc.Inspect(ctx, nameOrID)
	if err != nil {
		return nil, err
	}
	key := inspectName(insp)
	if key == "" {
		key = nameOrID // 理论上不会发生；真发生了至少保证不与别人共享锁
	}
	return u.lockFor(key), nil
}

// IsSelf 判断是否是 Dockhelm 自身。
func (u *Updater) IsSelf(nameOrID string) bool {
	if nameOrID == "" {
		return false
	}
	if u.selfName != "" && (nameOrID == u.selfName) {
		return true
	}
	if u.selfID != "" && strings.HasPrefix(u.selfID, nameOrID) {
		return true
	}
	if u.selfName != "" && strings.HasPrefix(u.selfName, strings.TrimPrefix(nameOrID, "/")) {
		return true
	}
	return false
}

// ---------- 检测 ----------

// CheckStatus 检测结论。
type CheckStatus string

const (
	StatusUpToDate        CheckStatus = "up_to_date"       // 已是最新
	StatusUpdateAvailable CheckStatus = "update_available" // 有新版本
	StatusUnknown         CheckStatus = "unknown"          // 无法判定（网络/认证/接口不可用）
	StatusNoUpstream      CheckStatus = "no_upstream"      // 本地构建的镜像，无远端可比
)

// CheckResult 单容器检测结果。
type CheckResult struct {
	Container    string      `json:"container"`
	ID           string      `json:"id"`
	Image        string      `json:"image"`
	LocalDigest  string      `json:"localDigest"`
	RemoteDigest string      `json:"remoteDigest"`
	Status       CheckStatus `json:"status"`
	Reason       string      `json:"reason"`
	CheckedAt    string      `json:"checkedAt"`
}

// Check 检测单个容器是否有新版本（只读，不拉取、不动容器）。
func (u *Updater) Check(ctx context.Context, nameOrID string) (*CheckResult, error) {
	insp, err := u.dc.Inspect(ctx, nameOrID)
	if err != nil {
		return nil, err
	}
	name := inspectName(insp)
	res := &CheckResult{
		Container: name,
		ID:        inspectID(insp),
		Image:     inspectImageRef(insp),
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
	}

	if res.Image == "" {
		res.Status = StatusNoUpstream
		res.Reason = "容器没有记录镜像引用"
		return res, nil
	}
	ref := dockerx.ParseRef(res.Image)

	// 本地 digest 集合
	local := map[string]bool{}
	for _, d := range containerRepoDigests(insp) {
		if d != "" {
			local[d] = true
			if res.LocalDigest == "" {
				res.LocalDigest = d
			}
		}
	}
	if len(local) == 0 {
		// 镜像 inspect 一次拿到 RepoDigests（有的容器 inspect 里带了 Image 的 RepoDigests）
		if imgInsp, err := u.dc.ImageInspect(ctx, ref.String()); err == nil {
			for _, d := range dockerx.LocalDigests(imgInsp) {
				local[d] = true
				if res.LocalDigest == "" {
					res.LocalDigest = d
				}
			}
		}
	}
	if len(local) == 0 {
		res.Status = StatusNoUpstream
		res.Reason = "该镜像没有仓库摘要（本地构建或已丢失 tag），无法与远端比较"
		return res, nil
	}

	remote, err := u.dc.DistributionInspect(ctx, ref.String())
	if err != nil {
		// 🚨 关键：检测失败绝不能当成「有更新」。
		res.Status = StatusUnknown
		res.Reason = "无法向仓库查询最新摘要：" + err.Error()
		return res, nil
	}
	res.RemoteDigest = remote
	if remote == "" {
		res.Status = StatusUnknown
		res.Reason = "仓库未返回摘要"
		return res, nil
	}
	if local[remote] {
		res.Status = StatusUpToDate
		res.Reason = "本地摘要与仓库一致"
		return res, nil
	}
	res.Status = StatusUpdateAvailable
	res.Reason = "仓库摘要与本地不一致"
	return res, nil
}

// CheckAll 并发检测全部容器（跳过自身与被排除的）。
//
// 检测只走守护进程的 /distribution（只下 manifest，不拉层、不动容器）。
// 曾经还有一个「深度检测」——真的 pull 一次再比镜像 ID；它只是拿带宽与磁盘
// 换一个几乎总是相同的答案，且做成页面级按钮时会一次拉满所有容器，已删。
// 拉取后比对镜像 ID 的那套逻辑仍在更新流程里（UpdateMany），那才是它该待的地方。
func (u *Updater) CheckAll(ctx context.Context, exclude map[string]bool) []CheckResult {
	list, err := u.dc.ListContainers(ctx)
	if err != nil {
		return []CheckResult{}
	}
	type job struct{ id, name string }
	jobs := []job{}
	for _, c := range list {
		n := c.Name()
		if u.IsSelf(n) || exclude[n] {
			continue
		}
		jobs = append(jobs, job{c.ID, n})
	}

	sem := make(chan struct{}, max(1, u.opt().Concurrency*2))
	var wg sync.WaitGroup
	var mu sync.Mutex
	out := []CheckResult{}
	for _, j := range jobs {
		wg.Add(1)
		go func(j job) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			sub, cancel := context.WithTimeout(ctx, 60*time.Second)
			defer cancel()
			r, err := u.Check(sub, j.id)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				out = append(out, CheckResult{
					Container: j.name, ID: j.id, Status: StatusUnknown,
					Reason: err.Error(), CheckedAt: time.Now().UTC().Format(time.RFC3339),
				})
				return
			}
			out = append(out, *r)
		}(j)
	}
	wg.Wait()
	sort.Slice(out, func(i, j int) bool {
		rank := map[CheckStatus]int{StatusUpdateAvailable: 0, StatusUnknown: 1, StatusUpToDate: 2, StatusNoUpstream: 3}
		if rank[out[i].Status] != rank[out[j].Status] {
			return rank[out[i].Status] < rank[out[j].Status]
		}
		return out[i].Container < out[j].Container
	})
	return out
}

// ---------- 更新 ----------

// ResultStatus 更新动作的结果状态。
type ResultStatus string

const (
	ResultUpToDate ResultStatus = "up_to_date" // 拉取后镜像未变，已跳过（容器未被动过）
	ResultUpdated  ResultStatus = "updated"    // 真的更新了
	ResultFailed   ResultStatus = "failed"     // 失败且回滚成功
	ResultBroken   ResultStatus = "broken"     // 失败且回滚也失败（需人工介入）
	ResultSkipped  ResultStatus = "skipped"    // 被排除/自身
	ResultPulling  ResultStatus = "pulling"
	ResultStopped  ResultStatus = "stopped"
	ResultCreated  ResultStatus = "created"
)

// Result 单容器更新结果。
type Result struct {
	Container  string       `json:"container"`
	Image      string       `json:"image"`
	OldImageID string       `json:"oldImageId"`
	NewImageID string       `json:"newImageId"`
	Status     ResultStatus `json:"status"`
	Message    string       `json:"message"`
	RolledBack bool         `json:"rolledBack"`
	Duration   int64        `json:"durationMs"`
	Steps      []string     `json:"steps"`
	// PullReused 本次拉取是复用了同一轮里别的容器拉到的结果（跳过下载）。
	PullReused bool `json:"pullReused"`
	// ReclaimedBytes 更新成功后清理旧镜像回收到的字节数（0 表示没删）。
	ReclaimedBytes int64 `json:"reclaimedBytes"`
}

func (u *Updater) step(res *Result, format string, a ...any) {
	line := fmt.Sprintf(format, a...)
	res.Steps = append(res.Steps, line)
	u.bus.Publish("update", "step", "running", map[string]any{
		"container": res.Container,
		"message":   line,
	})
}

func (u *Updater) status(res *Result, st ResultStatus, msg string) {
	res.Status = st
	res.Message = msg
	u.bus.Publish("update", "container_status", string(st), map[string]any{
		"container": res.Container,
		"image":     res.Image,
		"status":    string(st),
		"message":   msg,
		"oldImage":  dockerx.ShortID(res.OldImageID),
		"newImage":  dockerx.ShortID(res.NewImageID),
	})
}

// Update 更新单个容器。
//
// 语义：**镜像没变就完全不动容器**。这是对 dockerCopilot 那个
// 「一键更新 10 个，没更新的也被停掉」缺陷的直接修复。
func (u *Updater) Update(ctx context.Context, nameOrID string, force bool) *Result {
	started := time.Now()
	res := &Result{Container: nameOrID}

	lk, err := u.lockForContainer(ctx, nameOrID)
	if err != nil {
		u.status(res, ResultFailed, "读取容器信息失败："+err.Error())
		u.log("update", nameOrID, "failed", res.Message, "")
		return res
	}
	lk.Lock()
	defer lk.Unlock()

	// 锁拿到手之后必须重新读一次：上面那次 Inspect 只是为了拿 ID 做锁键，
	// 读到的状态可能是前一个持锁者动手之前的样子。
	insp, err := u.dc.Inspect(ctx, nameOrID)
	if err != nil {
		u.status(res, ResultFailed, "读取容器信息失败："+err.Error())
		u.log("update", nameOrID, "failed", res.Message, "")
		return res
	}
	name := inspectName(insp)
	res.Container = name
	res.Image = inspectImageRef(insp)
	res.OldImageID, _ = insp["Image"].(string)

	if u.IsSelf(name) {
		u.status(res, ResultSkipped, "这是 Dockhelm 自身容器，永不自动更新")
		return res
	}
	if res.Image == "" {
		u.status(res, ResultSkipped, "容器没有镜像引用，跳过")
		return res
	}

	wasRunning := inspectRunning(insp)
	u.step(res, "开始处理 %s（镜像 %s，旧 ID %s）", name, res.Image, dockerx.ShortID(res.OldImageID))

	// ---- 1. 拉取（走守护进程，registry-mirrors 自动生效）----
	u.status(res, ResultPulling, "正在拉取镜像…")
	pr, fromCache, err := u.pullForBatch(ctx, res.Image, func(ev dockerx.PullEvent) {
		if ev.Status != "" {
			u.bus.Publish("update", "pull_progress", "running", map[string]any{
				"container": name,
				"status":    ev.Status,
				"progress":  ev.Progress,
			})
		}
	})
	if err != nil {
		u.status(res, ResultFailed, "拉取失败："+err.Error())
		u.log("update", name, "failed", res.Message, strings.Join(res.Steps, "\n"))
		u.notify.Emit("update_failed", map[string]string{
			"container": name, "image": res.Image, "result": "拉取失败", "message": err.Error(),
		})
		res.Duration = time.Since(started).Milliseconds()
		return res
	}
	res.NewImageID = pr.ImageID
	res.PullReused = fromCache
	if fromCache {
		u.step(res, "本轮已拉取过 %s，直接复用（新 ID %s）", res.Image, dockerx.ShortID(pr.ImageID))
	} else {
		u.step(res, "拉取完成，新 ID %s（digest %s）", dockerx.ShortID(pr.ImageID), shortDigest(pr.Digest))
	}

	// ---- 2. 核心判定：镜像没变 ⇒ 到此为止 ----
	if res.NewImageID != "" && res.NewImageID == res.OldImageID && !force {
		u.status(res, ResultUpToDate, "镜像已是最新，容器保持原样（未停止、未重建）")
		u.step(res, "镜像 ID 未变化（%s），跳过重建", dockerx.ShortID(res.OldImageID))
		u.log("update", name, "up_to_date", res.Message, strings.Join(res.Steps, "\n"))
		res.Duration = time.Since(started).Milliseconds()
		return res
	}
	if res.NewImageID == "" {
		u.status(res, ResultFailed, "拉取后无法读取镜像 ID，拒绝在不确定状态下重建容器")
		u.log("update", name, "failed", res.Message, strings.Join(res.Steps, "\n"))
		res.Duration = time.Since(started).Milliseconds()
		return res
	}
	if pr.UpToDate && !force && !fromCache {
		u.status(res, ResultUpToDate, "守护进程报告镜像已是最新，容器保持原样")
		u.log("update", name, "up_to_date", res.Message, strings.Join(res.Steps, "\n"))
		res.Duration = time.Since(started).Milliseconds()
		return res
	}
	u.step(res, "检测到镜像变化，开始重建容器")

	// ---- 3. 写配置快照（回滚与还原的依据，可用设置关掉）----
	if u.opt().BackupBefore {
		snapPath, snapErr := u.saveSnapshot(name, insp)
		if snapErr != nil {
			u.step(res, "⚠ 配置快照写入失败：%v（继续，但不写快照不影响回滚）", snapErr)
		} else {
			u.step(res, "已写入配置快照 %s", filepath.Base(snapPath))
		}
	} else {
		u.step(res, "已在设置里关闭「更新前自动备份配置」，跳过写快照")
	}

	// ---- 4~8. 停止 / 改名 / 重建 / 启动 / 自检 ----
	newID, failMsg, rb := u.recreate(ctx, insp, name, res, wasRunning)
	if failMsg != "" {
		res.RolledBack = rb
		if rb {
			u.status(res, ResultFailed, failMsg+"（已成功回滚到旧容器）")
		} else {
			u.status(res, ResultBroken, failMsg+"（回滚失败，请手工处理）")
		}
		u.log("update", name, "failed", res.Message, strings.Join(res.Steps, "\n"))
		u.notify.Emit("update_failed", map[string]string{
			"container": name, "image": res.Image,
			"result":  map[bool]string{true: "已回滚", false: "回滚失败"}[rb],
			"message": failMsg,
		})
		res.Duration = time.Since(started).Milliseconds()
		return res
	}
	res.NewImageID = newID
	u.status(res, ResultUpdated, "更新成功"+map[bool]string{true: "", false: "（原容器为停止状态，保持停止）"}[wasRunning])
	u.log("update", name, "success", "已更新到 "+dockerx.ShortID(newID), strings.Join(res.Steps, "\n"))

	// 9. 清理旧镜像：只在确认没有任何容器再引用它时才动手。
	res.ReclaimedBytes = u.cleanupOldImage(ctx, res.OldImageID, res.NewImageID, res)

	u.notify.Emit("update_success", map[string]string{
		"container": name, "image": res.Image,
		"result": "更新成功", "message": "新镜像 " + dockerx.ShortID(newID),
	})
	res.Duration = time.Since(started).Milliseconds()
	return res
}

// cleanupOldImage 更新成功后，把已经被取代的旧镜像删掉，回收磁盘。
//
// 三条硬约束：
//  1. 关闭开关时直接返回（返回 0 表示没有回收）。
//  2. 新旧 ID 相同（强制重建场景）绝对不删 —— 那正是当前正在跑的那一份。
//  3. 先自己数一遍还有没有容器引用它。只依赖 Docker 的 409 也能防住，
//     但那样「为什么没删」就只有一个冷冰冰的报错，不如自己判一次说明白。
func (u *Updater) cleanupOldImage(ctx context.Context, oldID, newID string, res *Result) int64 {
	if !u.opt().CleanupAfter {
		return 0
	}
	if oldID == "" || oldID == newID {
		return 0
	}
	list, err := u.dc.ListContainers(ctx)
	if err != nil {
		u.step(res, "跳过旧镜像清理：读不到容器列表（%v）", err)
		return 0
	}
	for _, c := range list {
		if c.ImageID != "" && c.ImageID == oldID {
			u.step(res, "旧镜像 %s 仍被容器 %s 引用，保留不删", dockerx.ShortID(oldID), c.Name())
			return 0
		}
	}

	// 先量一下体积，删完再问就来不及了
	var size int64
	if imgs, err := u.dc.ListImages(ctx, false); err == nil {
		for _, im := range imgs {
			if im.ID == oldID {
				size = im.Size
				break
			}
		}
	}
	if err := u.dc.RemoveImage(ctx, oldID, false, false); err != nil {
		u.step(res, "旧镜像 %s 未能删除（可能仍被引用）：%v", dockerx.ShortID(oldID), err)
		return 0
	}
	u.step(res, "已清理旧镜像 %s，回收 %s", dockerx.ShortID(oldID), humanBytes(size))
	return size
}

// humanBytes 把字节数说成人话（只用于步骤日志）。
func humanBytes(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for m := n / unit; m >= unit; m /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGT"[exp])
}

// recreate 执行「停旧 → 改名保留 → 建新 → 启动 → 自检」，失败自动回滚。
// 返回 (新容器ID, 失败信息, 是否已回滚)。失败信息为空表示成功。
func (u *Updater) recreate(ctx context.Context, insp map[string]any, name string, res *Result, wasRunning bool) (string, string, bool) {
	oldID := inspectID(insp)
	bakName := fmt.Sprintf("%s__bak_%s", truncName(name, 40), time.Now().Format("20060102-150405"))

	// 先把创建请求体准备好并校验 —— **必须在动这个容器之前**。
	// 否则会走成「停旧容器 → 改名 → 创建失败 → 回滚」：容器白停一次，
	// 日志里还会留下一次吓人的失败，而问题其实在第一步就看得见。
	// 先取一次现存网络，避免把已被删除的网络别名传回去导致 404。
	var existingNets map[string]bool
	if nets, err := u.dc.ListNetworks(ctx); err == nil {
		existingNets = NetworkNameSet(nets)
	}
	cfg, hostCfg, netCfg := BuildCreateSpec(insp, existingNets)
	if fail := u.validateCreateSpec(cfg, insp, res); fail != "" {
		return "", fail, false // false：一个容器都没动过，没什么可回滚的
	}

	// AutoRemove 的容器一停就自动删除，改名会失败 —— 先关掉它（新容器仍按原配置创建）。
	// 这个改动必须能被回滚还原（否则更新失败回滚之后，容器会永久丢掉 --rm 语义）。
	autoRemoveOn := false
	if hc, ok := insp["HostConfig"].(map[string]any); ok {
		if ar, _ := hc["AutoRemove"].(bool); ar {
			if err := u.dc.UpdateContainer(ctx, oldID, map[string]any{"AutoRemove": false}); err == nil {
				autoRemoveOn = true
				u.step(res, "原容器开启了自动删除，已临时关闭以保证可回滚")
			}
		}
	}

	// 5. 停止
	if wasRunning {
		u.status(res, ResultStopped, "正在停止旧容器…")
		timeout := 30
		if err := u.dc.ContainerAction(ctx, oldID, "stop", &timeout); err != nil {
			return "", "停止旧容器失败：" + err.Error(), true
		}
		u.step(res, "旧容器已停止")
	} else {
		u.step(res, "旧容器原本就是停止状态，无需停止")
	}

	// 6. 改名保留（绝不删除 —— 旧容器的可写层与运行期身份因此保留下来）
	if err := u.dc.RenameContainer(ctx, oldID, bakName); err != nil {
		if wasRunning {
			// 这一步的「把旧容器重新启起来」也是补救动作，用独立 ctx
			cctx, cancel := cleanupCtx(ctx)
			_ = u.dc.ContainerAction(cctx, oldID, "start", nil)
			cancel()
		}
		return "", "重命名旧容器失败：" + err.Error(), true
	}
	u.step(res, "旧容器已改名为 %s（保留以便回滚）", bakName)

	// 7. 用原配置创建同名新容器（请求体在函数开头就备好了）
	u.status(res, ResultCreated, "正在创建新容器…")
	newID, err := u.dc.CreateContainer(ctx, name, cfg, hostCfg, netCfg)
	if err != nil {
		u.step(res, "创建新容器失败：%v，开始回滚", err)
		return "", "创建新容器失败：" + err.Error(), u.rollback(ctx, oldID, bakName, name, wasRunning, autoRemoveOn, res)
	}
	u.step(res, "新容器已创建（%s）", dockerx.ShortID(newID))

	// 8. 启动（原本停止的容器保持停止）
	if wasRunning {
		if err := u.dc.ContainerAction(ctx, newID, "start", nil); err != nil {
			u.discardContainer(ctx, newID)
			return "", "启动新容器失败：" + err.Error(), u.rollback(ctx, oldID, bakName, name, wasRunning, autoRemoveOn, res)
		}
		u.step(res, "新容器已启动，开始健康检查")
		// 9. 健康判定
		if err := u.waitHealthy(ctx, newID, insp); err != nil {
			u.discardContainer(ctx, newID)
			return "", "新容器健康检查未通过：" + err.Error(), u.rollback(ctx, oldID, bakName, name, wasRunning, autoRemoveOn, res)
		}
		u.step(res, "健康检查通过")
	}

	// 10. 清理备份容器（保留策略由 KeepBackups 决定）
	if u.opt().KeepBackups <= 0 {
		if err := u.dc.RemoveContainer(ctx, oldID, false, false); err != nil {
			u.step(res, "⚠ 备份容器 %s 删除失败（可稍后手动清理）：%v", bakName, err)
		} else {
			u.step(res, "已清理备份容器 %s", bakName)
		}
	} else {
		u.step(res, "按保留策略保留备份容器 %s", bakName)
	}
	return newID, "", false
}

// validateCreateSpec 在动容器之前判断「这份配置能不能建出一个容器」。
// 返回空字符串表示可以继续；否则返回一句给人看的失败原因。
//
// 存在的理由：创建请求体是从 inspect 拼出来的，而 inspect 来自网络 ——
// 一旦它不完整（没有 Config），等真到了 create 那一步才失败，代价是
// 「旧容器已经停了、改了名，再回滚」。而这个判断在任何时候都能做，
// 且做完就知道没法继续，所以必须在动手之前做。
func (u *Updater) validateCreateSpec(cfg map[string]any, insp map[string]any, res *Result) string {
	patched, fail := EnsureCreateSpec(cfg, insp)
	if patched != "" {
		u.step(res, "容器没有记录镜像引用，改用镜像 ID %s 重建", dockerx.ShortID(patched))
	}
	if fail == "" {
		return ""
	}
	return fail + "，拒绝重建，容器保持原样"
}

// rollback 回滚：删新容器 → 把 __bak_ 改回原名 → 按原状态启动 →
// 把为改名而临时关掉的 AutoRemove 还原回去。
//
// autoRemoveOn 表示「本次确实动过 AutoRemove」。不还原的话，一次失败的更新
// 会永久改掉这个容器的语义：它本该在退出时被自动删除，回滚之后却不会了。
// cleanupCtx 返回一个**不受调用方取消影响**的短超时上下文，专供回滚/清理使用。
//
// 回滚存在的意义就是「前面已经失败了」，而最常见的失败原因恰恰是调用方的
// ctx 超时或取消（批量更新各有 20 / 60 分钟预算）。若这些补救动作继续用同一个
// ctx，它们必然跟着一起失败 —— 于是承诺过的「失败自动回滚」在最需要它的时候
// 恰好失效：留下一个已改名停着的旧容器 + 一个建了但没启起来的新容器，
// 而返回的 RolledBack=true 还让调用方以为已经退回去了。
func cleanupCtx(ctx context.Context) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.WithoutCancel(ctx), 90*time.Second)
}

// discardContainer 清掉一个不该留下的容器（用独立 ctx，别被调用方的取消带走）。
func (u *Updater) discardContainer(ctx context.Context, id string) {
	cctx, cancel := cleanupCtx(ctx)
	defer cancel()
	_ = u.dc.ContainerAction(cctx, id, "stop", nil)
	_ = u.dc.RemoveContainer(cctx, id, true, false)
}

// rollback 把旧容器改回原名并恢复运行。返回是否真的回滚成功。
func (u *Updater) rollback(ctx context.Context, oldID, bakName, name string, wasRunning, autoRemoveOn bool, res *Result) bool {
	// 回滚全程用独立 ctx，理由见 cleanupCtx
	ctx, cancel := cleanupCtx(ctx)
	defer cancel()

	u.step(res, "开始回滚：删除新容器并把旧容器改回原名")
	if err := u.dc.RenameContainer(ctx, oldID, name); err != nil {
		u.step(res, "✗ 回滚失败（改名）：%v", err)
		return false
	}
	if wasRunning {
		if err := u.dc.ContainerAction(ctx, oldID, "start", nil); err != nil {
			u.step(res, "✗ 回滚失败（启动）：%v", err)
			return false
		}
	}
	if autoRemoveOn {
		if err := u.dc.UpdateContainer(ctx, oldID, map[string]any{"AutoRemove": true}); err != nil {
			u.step(res, "⚠ 已回滚，但恢复自动删除（--rm）失败：%v", err)
		} else {
			u.step(res, "已还原原来的自动删除（--rm）设置")
		}
	}
	u.step(res, "✓ 已回滚，容器 %s 恢复为更新前的状态", name)
	return true
}

// waitHealthy 等待容器健康。
func (u *Updater) waitHealthy(ctx context.Context, id string, before map[string]any) error {
	hasHealth := false
	if c, ok := before["Config"].(map[string]any); ok {
		if _, ok := c["Healthcheck"]; ok {
			hasHealth = true
		}
	}
	deadline := time.Now().Add(u.opt().HealthTimeout)
	stableSince := time.Time{}
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("超时（%s）", u.opt().HealthTimeout)
		}
		insp, err := u.dc.Inspect(ctx, id)
		if err != nil {
			return err
		}
		st, _ := insp["State"].(map[string]any)
		if st == nil {
			return fmt.Errorf("无法读取容器状态")
		}
		running, _ := st["Running"].(bool)
		if !running {
			exit, _ := st["ExitCode"].(float64)
			return fmt.Errorf("容器已退出（退出码 %d）", int(exit))
		}
		if hasHealth {
			h, _ := st["Health"].(map[string]any)
			if h != nil {
				hs, _ := h["Status"].(string)
				switch hs {
				case "healthy":
					return nil
				case "unhealthy":
					return fmt.Errorf("健康检查状态为 unhealthy")
				}
			}
		} else {
			if stableSince.IsZero() {
				stableSince = time.Now()
			} else if time.Since(stableSince) >= time.Duration(u.opt().StableSeconds)*time.Second {
				return nil
			}
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(1500 * time.Millisecond):
		}
	}
}

// UpdateMany 批量更新（有界并发），并汇总成一条通知。
func (u *Updater) UpdateMany(ctx context.Context, names []string, force bool) []Result {
	if len(names) == 0 {
		return nil
	}
	conc := max(1, u.opt().Concurrency)
	sem := make(chan struct{}, conc)
	var wg sync.WaitGroup
	var mu sync.Mutex
	out := []Result{}

	// 同一轮里同一镜像只拉一次。缓存的生命周期就是这一批，结束必须清掉 ——
	// 留着会让「刚推的新版本」在下一轮里被上一次的旧结论吃掉。
	u.resetBatchPull()
	defer u.resetBatchPull()

	batch := u.bus.Publish("update", "batch_start", "running", map[string]any{
		"total": len(names),
	})
	_ = batch

	for _, n := range names {
		wg.Add(1)
		go func(n string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			r := u.Update(ctx, n, force)
			mu.Lock()
			out = append(out, *r)
			mu.Unlock()
		}(n)
	}
	wg.Wait()

	updated, skipped, failed := 0, 0, 0
	var reclaimed int64
	reused := 0
	for _, r := range out {
		reclaimed += r.ReclaimedBytes
		if r.PullReused {
			reused++
		}
		switch r.Status {
		case ResultUpdated:
			updated++
		case ResultUpToDate, ResultSkipped:
			skipped++
		default:
			failed++
		}
	}
	u.bus.Publish("update", "batch_done", "success", map[string]any{
		"total": len(names), "updated": updated, "skipped": skipped, "failed": failed,
		"reclaimedBytes": reclaimed, "pullReused": reused,
	})
	summary := fmt.Sprintf("共 %d 个容器：更新 %d、已是最新/跳过 %d、失败 %d", len(names), updated, skipped, failed)
	if reclaimed > 0 {
		summary += "，清理旧镜像回收 " + humanBytes(reclaimed)
	}
	if reused > 0 {
		summary += fmt.Sprintf("，%d 个容器复用了同轮已拉取的镜像", reused)
	}
	// 汇总记录/通知只在**多个容器**时写：单个容器的批次里，下面每条结果本来就各有一条
	// 运行记录（更新成功 / 镜像已是最新…），再来一条「共 1 个容器：更新 0…」纯属重复 ——
	// 同一次更新在总览页出现两行，用户会以为执行了两遍。
	if len(names) > 1 {
		u.log("update", "batch", "done", summary, "")
		u.notify.Emit("batch_update_done", map[string]string{
			"container": fmt.Sprintf("%d 个容器", len(names)),
			"result":    "批量更新完成",
			"message":   summary,
		})
	}
	return out
}

// saveSnapshot 把容器配置快照写到 <BackupDir>/<name>/<ts>.json。
//
// 文件格式必须与 backup.Service.Snapshot 写出来的完全一致（外层包一层 _dockhelm
// 元信息、inspect 放里面），否则这些「更新前快照」在备份页里既看不到镜像与状态，
// 也**还原不了** —— 而「更新失败还能退回去」正是写它的唯一理由。
func (u *Updater) saveSnapshot(name string, insp map[string]any) (string, error) {
	if u.opt().BackupDir == "" {
		return "", fmt.Errorf("未配置备份目录")
	}
	dir := filepath.Join(u.opt().BackupDir, name)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	ts := time.Now().Format("20060102-150405")
	doc := map[string]any{
		"_dockhelm": map[string]any{
			"snapshotAt": time.Now().UTC().Format(time.RFC3339),
			"reason":     "pre_update",
			"version":    1,
		},
		"inspect": insp,
	}
	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", err
	}
	p := filepath.Join(dir, ts+".json")
	// 同一秒内重复写时加后缀，别互相覆盖
	for i := 1; ; i++ {
		if _, statErr := os.Stat(p); os.IsNotExist(statErr) {
			break
		}
		p = filepath.Join(dir, fmt.Sprintf("%s-%d.json", ts, i))
	}
	// 0600：配置里可能含密码（Env）
	if err := os.WriteFile(p, b, 0o600); err != nil {
		return "", err
	}
	return p, nil
}

// ---------- inspect 取值助手 ----------

func inspectName(insp map[string]any) string {
	if s, ok := insp["Name"].(string); ok {
		return strings.TrimPrefix(s, "/")
	}
	return ""
}

func inspectID(insp map[string]any) string {
	s, _ := insp["Id"].(string)
	return s
}

func inspectImageRef(insp map[string]any) string {
	if c, ok := insp["Config"].(map[string]any); ok {
		if s, ok := c["Image"].(string); ok {
			return s
		}
	}
	return ""
}

func inspectRunning(insp map[string]any) bool {
	if st, ok := insp["State"].(map[string]any); ok {
		b, _ := st["Running"].(bool)
		return b
	}
	return false
}

// containerRepoDigests 从容器 inspect 里尽力取镜像摘要。
func containerRepoDigests(insp map[string]any) []string {
	out := []string{}
	// 有的版本把 Image 的 RepoDigests 放在 ImageManifestDescriptor 里
	if imd, ok := insp["ImageManifestDescriptor"].(map[string]any); ok {
		if d, ok := imd["digest"].(string); ok && d != "" {
			out = append(out, d)
		}
	}
	return out
}

func shortDigest(d string) string {
	d = strings.TrimPrefix(d, "sha256:")
	if len(d) > 12 {
		return d[:12]
	}
	if d == "" {
		return "-"
	}
	return d
}

func truncName(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
