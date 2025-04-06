package dockerProv

import (
	"encoding/json"
	"fmt"

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
				Driver: "cid",
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

	bs, _ := json.Marshal(ops.ToHostConfig())
	fmt.Println(string(bs))
	return ops
}

func (ops StartContainerOptions) WithShmSize(MB int64) StartContainerOptions {
	ops.ShmSize = MB * 1024 * 1024
	return ops
}
