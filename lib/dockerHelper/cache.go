package dockerHelper

import "path/filepath"

func (ops StartContainerOptions) WithPipCache(homeDir string) StartContainerOptions {
	if homeDir == "" {
		return ops
	}
	pipCacheDir := filepath.Join(homeDir, ".cache", "pip")
	if ops.Volumes == nil {
		ops.Volumes = map[string]struct{}{}
	}
	ops.Volumes[pipCacheDir] = struct{}{}
	ops.Binds = append(ops.Binds, pipCacheDir+":"+pipCacheDir)
	return ops
}
