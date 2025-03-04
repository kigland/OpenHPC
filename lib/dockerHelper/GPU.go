package dockerHelper

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
