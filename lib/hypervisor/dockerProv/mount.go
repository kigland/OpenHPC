package dockerProv

import (
	"path/filepath"
)

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

func (ops StartContainerOptions) WithMountRW(from, to string) StartContainerOptions {
	if from == "" || to == "" {
		return ops
	}
	ops.Binds = append(ops.Binds, from+":"+to+":rw")
	return ops
}
