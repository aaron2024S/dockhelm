package api

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/aaron2024s/dockhelm/internal/updater"
	"github.com/aaron2024s/dockhelm/internal/version"
)

// hostUsage 是「所有运行中容器」的 CPU / 内存合计。
type hostUsage struct {
	cpuPercent float64
	memUsed    int64
	memTotal   int64
	at         time.Time
	ok         bool
}

var (
	usageMu    sync.Mutex
	usageCache = hostUsage{}
)

// usageTTL 合计值缓存时长。容器统计要逐个向守护进程取，24 个容器约几百毫秒；
// 总览页与顶栏都会定期请求这个接口，缓存一下免得反复打扰守护进程。
const usageTTL = 10 * time.Second

// aggregateUsage 汇总所有运行中容器的 CPU 与内存占用。
//
// 单个容器的 cpuPercent 是「占满一个核 = 100%」，所以先把它们加起来、再除以核数，
// 才得到宿主机视角的总占用率。任何一个容器取不到统计都只是少算它一个，不影响整体。
func (s *Server) aggregateUsage(ctx context.Context, ids []string, cpus int) hostUsage {
	now := time.Now()
	usageMu.Lock()
	cached := usageCache
	usageMu.Unlock()
	if cached.ok && now.Sub(cached.at) < usageTTL {
		return cached
	}

	out := hostUsage{at: now, ok: true}
	if len(ids) == 0 {
		usageMu.Lock()
		usageCache = out
		usageMu.Unlock()
		return out
	}

	// 整批最多给 3 秒，慢了就让界面先出别的数字
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var (
		mu       sync.Mutex
		sumCPU   float64
		memUsed  int64
		memTotal int64
	)
	sem := make(chan struct{}, 8)
	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			raw, err := s.dc.StatsOneShot(ctx, id)
			if err != nil {
				return
			}
			sum := summarizeStats(raw)

			mu.Lock()
			sumCPU += num(sum["cpuPercent"])
			memUsed += int64(num(sum["memUsage"]))
			if l := int64(num(sum["memLimit"])); l > memTotal {
				memTotal = l
			}
			mu.Unlock()
		}(id)
	}
	wg.Wait()

	if cpus < 1 {
		cpus = 1
	}
	out.cpuPercent = sumCPU / float64(cpus)
	out.memUsed = memUsed
	out.memTotal = memTotal

	usageMu.Lock()
	usageCache = out
	usageMu.Unlock()
	return out
}

// hOverview 总览页所需的全部数据（一次请求拿齐，避免首屏打一堆请求）。
func (s *Server) hOverview(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := s.ctx(r)
	defer cancel()

	info, infoErr := s.dc.Info(ctx)
	containers, _ := s.dc.ListContainers(ctx)
	images, _ := s.dc.ListImages(ctx, false)

	running, stopped, paused, unhealthy := 0, 0, 0, 0
	runningIDs := make([]string, 0, len(containers))
	for _, c := range containers {
		switch c.State {
		case "running":
			running++
			runningIDs = append(runningIDs, c.ID)
		case "paused":
			paused++
		default:
			stopped++
		}
		if strings.Contains(c.Status, "unhealthy") {
			unhealthy++
		}
	}

	s.checkMu.RLock()
	cache := s.checkCache
	checkedAt := s.checkAt
	s.checkMu.RUnlock()

	withUpdates := []map[string]any{}
	unknown := 0
	for _, c := range cache {
		switch c.Status {
		case updater.StatusUpdateAvailable:
			withUpdates = append(withUpdates, map[string]any{
				"container":    c.Container,
				"image":        c.Image,
				"reason":       c.Reason,
				"localDigest":  c.LocalDigest,
				"remoteDigest": c.RemoteDigest,
			})
		case updater.StatusUnknown:
			unknown++
		}
	}

	var imageBytes, reclaimable int64
	for _, im := range images {
		imageBytes += im.Size
		// 没有 tag 的镜像就是悬空的，随时可以回收
		if len(im.RepoTags) == 0 {
			reclaimable += im.Size
		}
	}

	logs, _ := s.st.ListRunLogs(15)
	diskFree, diskTotal, _ := diskUsage(s.cfg.DataDir)

	out := map[string]any{
		"version": map[string]any{
			"version": version.Version,
			"name":    version.AppName,
			"commit":  version.Get().CommitShort,
		},
		"containers": map[string]any{
			"total":     len(containers),
			"running":   running,
			"stopped":   stopped,
			"paused":    paused,
			"unhealthy": unhealthy,
		},
		"images": map[string]any{
			"total":       len(images),
			"sizeBytes":   imageBytes,
			"reclaimable": reclaimable,
		},
		"updates": map[string]any{
			"available": len(withUpdates),
			"unknown":   unknown,
			"checkedAt": timeOrEmpty(checkedAt),
			"items":     withUpdates,
		},
		"disk": map[string]any{
			"free":  diskFree,
			"total": diskTotal,
		},
		"notify":   map[string]any{"sentToday": s.nt.GetSettings().SentToday},
		"recent":   logs,
		"excluded": keys(s.excluded()),
		"self":     s.up.SelfName(),
	}
	if len(runningIDs) > 0 {
		cpus := 0
		if infoErr == nil {
			cpus = info.CPUs
		}
		usage := s.aggregateUsage(ctx, runningIDs, cpus)
		out["usage"] = map[string]any{
			"cpuPercent": usage.cpuPercent,
			"memUsed":    usage.memUsed,
			"memTotal":   usage.memTotal,
		}
	}
	if infoErr == nil {
		out["docker"] = map[string]any{
			"name":     info.Name,
			"version":  info.ServerVersion,
			"os":       info.OperatingSystem,
			"arch":     info.Architecture,
			"kernel":   info.KernelVersion,
			"cpus":     info.CPUs,
			"memTotal": info.MemTotal,
			"rootDir":  info.DockerRootDir,
			"mirrors":  info.RegistryConfig.Mirrors,
		}
	} else {
		out["dockerError"] = infoErr.Error()
	}
	writeOK(w, out)
}
