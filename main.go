// Dockhelm —— 自托管的 Docker 容器管理面板。
//
// 单二进制 + //go:embed 内嵌前端，只依赖标准库（外加 bcrypt 与 cron 两个小库），
// 因此 CGO_ENABLED=0 就能交叉编译 amd64 / arm64。
package main

import (
	"context"
	"embed"
	"errors"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/aaron2024s/dockhelm/internal/api"
	"github.com/aaron2024s/dockhelm/internal/auth"
	"github.com/aaron2024s/dockhelm/internal/backup"
	"github.com/aaron2024s/dockhelm/internal/bus"
	"github.com/aaron2024s/dockhelm/internal/config"
	"github.com/aaron2024s/dockhelm/internal/dockerx"
	"github.com/aaron2024s/dockhelm/internal/notify"
	"github.com/aaron2024s/dockhelm/internal/scheduler"
	"github.com/aaron2024s/dockhelm/internal/store"
	"github.com/aaron2024s/dockhelm/internal/updater"
	"github.com/aaron2024s/dockhelm/internal/version"
	"github.com/aaron2024s/dockhelm/internal/watch"
)

//go:embed all:web/dist
var webFS embed.FS

func main() {
	log.SetFlags(log.LstdFlags)
	log.Printf("%s %s (%s) 启动中…", version.AppName, version.Version, version.Get().CommitShort)

	cfg := config.Load()
	if err := cfg.EnsureDirs(); err != nil {
		log.Fatalf("无法创建数据目录 %s：%v", cfg.DataDir, err)
	}

	// —— 持久化 ——
	st, err := store.Open(cfg.DBPath())
	if err != nil {
		log.Fatalf("无法打开数据文件 %s：%v", cfg.DBPath(), err)
	}
	defer st.Close()

	// —— 事件总线 ——
	b := bus.New()

	// —— Docker 客户端 ——
	dc, err := dockerx.New(cfg.DockerHost)
	if err != nil {
		log.Fatalf("无法解析 DOCKER_HOST %q：%v", cfg.DockerHost, err)
	}
	{
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		if err := dc.Negotiate(ctx); err != nil {
			// 启动时连不上不致命：面板仍然要能打开并告诉用户哪里不对
			log.Printf("⚠ 暂时无法连接 Docker 守护进程（%s）：%v", cfg.DockerHost, err)
		} else {
			log.Printf("已连接 Docker 守护进程 %s（API v%s）", dc.Host(), dc.APIVersion())
		}
		cancel()
	}

	// —— 认证 ——
	authMgr, err := auth.New(cfg.AuthPath(), st)
	if err != nil {
		log.Fatalf("无法读取认证数据：%v", err)
	}
	if w := authMgr.LoadWarning(); w != "" {
		log.Printf("⚠ %s", w)
	}
	if cfg.ForcePassword != "" {
		if err := authMgr.ForceSetPassword(cfg.ForcePassword); err != nil {
			log.Printf("⚠ DOCKHELM_PASSWORD 设置失败：%v", err)
		} else {
			log.Printf("已按 DOCKHELM_PASSWORD 重置登录密码（建议用完就删掉这个环境变量）")
		}
	}
	if !authMgr.Initialized() {
		log.Printf("首次启动：请打开面板设置登录密码")
	}

	// —— 通知 ——
	nt := notify.New(st, b)
	nt.Start()
	defer nt.Stop()

	// —— 更新引擎 ——
	opts := updater.DefaultOptions(cfg.ContainerBackupDir())
	if n := config.AtoiDefault(st.GetSetting("update.concurrency", ""), 0); n > 0 {
		opts.Concurrency = n
	}
	up := updater.New(dc, b, opts)
	up.SetNotifier(nt)
	up.SetLogSink(func(kind, ref, status, message, detail string) {
		st.AddRunLog(kind, ref, status, message, detail)
	})

	// 自身容器识别（永不更新自己）+ 读自己的挂载映射
	{
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		name, id := watch.DetectSelf(ctx, dc, cfg.SelfContainer)
		up.SetSelf(name, id)
		if name != "" {
			log.Printf("识别到 Dockhelm 自身容器：%s（%s），已加入永久排除", name, id[:min(12, len(id))])
			target := id
			if target == "" {
				target = name
			}
			if insp, err := dc.Inspect(ctx, target); err == nil {
				cfg.SetMounts(dockerx.SelfMountMap(insp))
			} else {
				log.Printf("读取自身挂载失败，路径映射退回 DOCKHELM_HOST_ROOTS：%v", err)
			}
		}
		cancel()
		log.Printf("宿主路径映射：%s", cfg.Summary())
	}

	// 加速源配置注入
	up.SetMirror(func() string {
		var c struct {
			PullMirror string `json:"pullMirror"`
		}
		st.GetJSON("registry.settings", &c)
		return c.PullMirror
	})

	// —— 备份 ——
	bk := backup.New(cfg, dc, st, nt)

	// —— 计划任务 ——
	sch := scheduler.New(st, up, dc, bk, nt, b, func() map[string]bool {
		out := map[string]bool{}
		if self := up.SelfName(); self != "" {
			out[self] = true
		}
		var list []string
		st.GetJSON("exclude.containers", &list)
		for _, n := range list {
			if n = strings.TrimSpace(n); n != "" {
				out[n] = true
			}
		}
		return out
	})
	if err := sch.Start(); err != nil {
		log.Printf("⚠ 计划任务调度器启动失败：%v", err)
	}
	defer sch.Stop()

	// —— 静态资源 ——
	static := buildStaticHandler()

	// —— HTTP 服务 ——
	srv := api.New(api.Deps{
		Cfg: cfg, Store: st, Docker: dc, Auth: authMgr,
		Update: up, Sched: sch, Backup: bk, Notify: nt, Bus: b, Static: static,
	})
	httpSrv := &http.Server{
		Addr:              cfg.Listen,
		Handler:           srv.Handler(),
		ReadHeaderTimeout: 20 * time.Second,
		// 不设 WriteTimeout：SSE 与日志流要长时间挂着
		IdleTimeout: 120 * time.Second,
	}

	// —— 后台任务 ——
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	// Docker 事件观察（容器意外退出 / 崩溃循环 / 健康失败 / 连接断开）
	watcher := watch.New(dc, nt, st, b, up.SelfName)
	go watcher.Run(appCtx)

	// 定期清理
	go housekeeping(appCtx, st, bk, nt)

	// 周期检测 + 自动更新。
	//
	// 启动后首次巡检由 StartAutoLoop 的首轮承担（checkOnStart 开着则 8 秒后跑，
	// 否则 2 分钟）—— 这里曾经还额外挂了一个 8 秒的 goroutine 直接调 CheckAll，
	// 与首轮重复：启动瞬间两套并行打同一个镜像仓库，白白触发限流。
	go srv.StartAutoLoop(appCtx)
	{
		interval := config.AtoiDefault(st.GetSetting("update.checkInterval", ""), 6)
		if interval > 0 {
			log.Printf("周期更新检测已启用：每 %d 小时巡检一次（自动更新 %s）",
				interval, onOff(st.GetSetting("update.autoApply", "0") == "1"))
		} else {
			log.Printf("周期更新检测已关闭（设置页可开启）")
		}
	}

	// —— 启动 HTTP ——
	errCh := make(chan error, 1)
	go func() {
		log.Printf("HTTP 服务监听 %s（来自 %s）", cfg.Listen, cfg.ListenSource)
		if err := httpSrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	// —— 优雅退出 ——
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-sigCh:
		log.Printf("收到信号 %v，正在退出…", sig)
	case err := <-errCh:
		log.Printf("HTTP 服务异常退出：%v", err)
	}
	appCancel()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = httpSrv.Shutdown(shutdownCtx)
	_ = st.Close()
	log.Printf("已退出")
}

// buildStaticHandler 构造前端静态资源处理器，未命中的路径回退到 index.html（SPA 路由）。
func buildStaticHandler() http.Handler {
	sub, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		log.Printf("⚠ 前端资源不可用：%v", err)
		return nil
	}
	fileServer := http.FileServer(http.FS(sub))
	index, indexErr := fs.ReadFile(sub, "index.html")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}
		if f, err := sub.Open(path); err == nil {
			_ = f.Close()
			// 带哈希的静态资源可以长缓存
			if strings.HasPrefix(path, "assets/") {
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}
			fileServer.ServeHTTP(w, r)
			return
		}
		// SPA 回退
		if indexErr != nil {
			http.Error(w, "前端资源未构建", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-cache, must-revalidate")
		_, _ = w.Write(index)
	})
}

// housekeeping 定期清理过期会话、登录记录、去重键与旧备份。
func housekeeping(ctx context.Context, st *store.Store, bk *backup.Service, nt *notify.Manager) {
	purge := func() {
		st.PurgeExpiredSessions()
		st.PurgeOldLoginAttempts()
		for _, k := range st.DedupeKeys() {
			if t := st.LastNotifyTime(strings.TrimPrefix(k, "dedupe:")); !t.IsZero() && time.Since(t) > 24*time.Hour {
				st.DeleteSetting(k)
			}
		}
	}
	purge()
	t := time.NewTicker(time.Hour)
	defer t.Stop()
	lastPrune := time.Now()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			purge()
			if time.Since(lastPrune) > 20*time.Hour {
				// 快照保留策略来自设置页；更新前快照默认受保护，
				// 免得「保留最近 N 份」把唯一能回滚的那一份挤掉。
				readInt := func(key string, def int) int {
					return config.AtoiDefault(st.GetSetting(key, ""), def)
				}
				res := bk.Prune(backup.PruneOptions{
					KeepPerContainer: readInt("backup.keepPerContainer", 10),
					MaxAgeDays:       readInt("backup.maxAgeDays", 30),
					MaxTotalMB:       readInt("backup.maxTotalMB", 2048),
					KeepPreUpdate:    st.GetSetting("backup.keepPreUpdate", "1") == "1",
				})
				if res.Removed > 0 {
					log.Printf("已清理 %d 份过期配置快照，释放 %.1f MB", res.Removed, float64(res.FreedBytes)/1024/1024)
				}
				lastPrune = time.Now()
			}
		}
	}
}

// onOff 把 "1"/"0" 渲染成中文开关，只为启动日志好读。
func onOff(enabled bool) string {
	if enabled {
		return "已开启"
	}
	return "已关闭"
}
