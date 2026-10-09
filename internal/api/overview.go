package api

import (
	"net/http"
	"strings"

	"github.com/aaron2024s/dockhelm/internal/updater"
	"github.com/aaron2024s/dockhelm/internal/version"
)

// hOverview 总览页所需的全部数据（一次请求拿齐，避免首屏打一堆请求）。
func (s *Server) hOverview(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := s.ctx(r)
	defer cancel()

	info, infoErr := s.dc.Info(ctx)
	containers, _ := s.dc.ListContainers(ctx)
	images, _ := s.dc.ListImages(ctx, false)

	running, stopped, paused, unhealthy := 0, 0, 0, 0
	for _, c := range containers {
		switch c.State {
		case "running":
			running++
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
				"container": c.Container, "image": c.Image, "reason": c.Reason,
			})
		case updater.StatusUnknown:
			unknown++
		}
	}

	var imageBytes int64
	for _, im := range images {
		imageBytes += im.Size
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
			"total":     len(images),
			"sizeBytes": imageBytes,
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
	if infoErr == nil {
		out["docker"] = map[string]any{
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
