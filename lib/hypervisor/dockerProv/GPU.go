package dockerProv

import (
	"github.com/docker/docker/api/types/container"
)

func GetGPUDeviceRequests(gpuCount int) []container.DeviceRequest {
	if gpuCount <= 0 {
		return nil
	}

	return []container.DeviceRequest{
		{
			Count:        gpuCount,
			Capabilities: [][]string{{"gpu"}},
			Options:      map[string]string{"gpu": "true"},
		},
	}
}

func (ops StartContainerOptions) WithGPU(gpuCount int) StartContainerOptions {
	switch ops.Provider {
	case ProviderPodman:
		ops.Resources.DeviceRequests = []container.DeviceRequest{
			{
				Driver: "cdi",
				DeviceIDs: []string{
					"nvidia.com/gpu=all", // TODO: Make this dynamic,
					// Related issue: https://github.com/containers/podman/pull/25171
					// Should works in podman v5.4.2+
				},
			},
		}
	default: // Docker
		ops.Resources.DeviceRequests = GetGPUDeviceRequests(gpuCount)
	}

	return ops
}

func (ops StartContainerOptions) WithGPUIds(gpuIds []string) StartContainerOptions {
	switch ops.Provider {
	case ProviderPodman:
		ids := []string{}
		for _, id := range gpuIds {
			ids = append(ids, "nvidia.com/gpu="+id)
		}
		ops.Resources.DeviceRequests = []container.DeviceRequest{
			{
				Driver:    "cdi",
				DeviceIDs: ids,
			},
		}
	default: // Docker
		ops.Resources.DeviceRequests = []container.DeviceRequest{
			{
				Driver:       "nvidia",
				Capabilities: [][]string{{"gpu"}},
				Options:      map[string]string{"gpu": "true"},
				DeviceIDs:    gpuIds,
			},
		}
	}

	return ops
}

func (ops StartContainerOptions) WithShmSize(MB int64) StartContainerOptions {
	ops.ShmSize = MB * 1024 * 1024
	return ops
}
