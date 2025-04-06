package dockerProv

import "github.com/docker/docker/api/types/container"

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
	ops.Resources.DeviceRequests = GetGPUDeviceRequests(gpuCount)
	return ops
}

func (ops StartContainerOptions) WithShmSize(MB int64) StartContainerOptions {
	ops.ShmSize = MB * 1024 * 1024
	return ops
}
