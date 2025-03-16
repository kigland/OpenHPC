package dockerHelper

import "path/filepath"

const PIP_CACHE_VOLUME = "pip-cache"

func (ops StartContainerOptions) WithPipCache(homeDir string) StartContainerOptions {
	if homeDir == "" {
		return ops
	}
	pipCacheDir := filepath.Join(homeDir, ".cache", "pip")
	if ops.Volumes == nil {
		ops.Volumes = map[string]struct{}{}
	}
	ops.Volumes[pipCacheDir] = struct{}{}
	ops.Binds = append(ops.Binds, PIP_CACHE_VOLUME+":"+pipCacheDir)
	return ops
}

func (ops StartContainerOptions) WithRDS(rdsDir, mountAt string) StartContainerOptions {
	if rdsDir == "" || mountAt == "" {
		return ops
	}
	if ops.Volumes == nil {
		ops.Volumes = map[string]struct{}{}
	}
	ops.Volumes[rdsDir] = struct{}{}
	ops.Binds = append(ops.Binds, rdsDir+":"+mountAt+":rw")
	return ops
}
