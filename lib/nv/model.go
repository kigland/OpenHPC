package nv

import (
	"encoding/xml"
	"os/exec"
	"strconv"
	"strings"
)

type NvidiaSmiLog struct {
	// XMLName       xml.Name `xml:"nvidia_smi_log"`
	// Text          string   `xml:",chardata"`
	// Timestamp     string   `xml:"timestamp"`
	// DriverVersion string   `xml:"driver_version"`
	// CudaVersion   string   `xml:"cuda_version"`
	AttachedGpus string `xml:"attached_gpus"`
	Gpu          []struct {
		Text             string `xml:",chardata"`
		ID               string `xml:"id,attr"`
		ProductName      string `xml:"product_name"`
		PerformanceState string `xml:"performance_state"`
		FbMemoryUsage    struct {
			Text     string `xml:",chardata"`
			Total    string `xml:"total"`
			Reserved string `xml:"reserved"`
			Used     string `xml:"used"`
			Free     string `xml:"free"`
		} `xml:"fb_memory_usage"`
		Utilization struct {
			Text       string `xml:",chardata"`
			GpuUtil    string `xml:"gpu_util"`
			MemoryUtil string `xml:"memory_util"`
		} `xml:"utilization"`
		Temperature struct {
			Text                   string `xml:",chardata"`
			GpuTemp                string `xml:"gpu_temp"`
			GpuTempTlimit          string `xml:"gpu_temp_tlimit"`
			GpuTempMaxThreshold    string `xml:"gpu_temp_max_threshold"`
			GpuTempSlowThreshold   string `xml:"gpu_temp_slow_threshold"`
			GpuTempMaxGpuThreshold string `xml:"gpu_temp_max_gpu_threshold"`
			GpuTargetTemperature   string `xml:"gpu_target_temperature"`
			MemoryTemp             string `xml:"memory_temp"`
			GpuTempMaxMemThreshold string `xml:"gpu_temp_max_mem_threshold"`
		} `xml:"temperature"`
		GpuPowerReadings struct {
			Text       string `xml:",chardata"`
			PowerState string `xml:"power_state"`
			PowerDraw  string `xml:"power_draw"`
			// CurrentPowerLimit   string `xml:"current_power_limit"`
			// RequestedPowerLimit string `xml:"requested_power_limit"`
			// DefaultPowerLimit   string `xml:"default_power_limit"`
			// MinPowerLimit       string `xml:"min_power_limit"`
			// MaxPowerLimit       string `xml:"max_power_limit"`
		} `xml:"gpu_power_readings"`
	} `xml:"gpu"`
}

func GetNvidiaSmiLog() (*NvidiaSmiLog, error) {
	data, err := exec.Command("nvidia-smi", "-q", "-x").Output()
	if err != nil {
		return nil, err
	}
	return ParseNvidiaSmiLog(data)
}

func ParseNvidiaSmiLog(data []byte) (*NvidiaSmiLog, error) {
	var log NvidiaSmiLog
	if err := xml.Unmarshal(data, &log); err != nil {
		return nil, err
	}
	return &log, nil
}

func (log *NvidiaSmiLog) Parse() (*NVInfo, error) {
	info := &NVInfo{
		GPUCount: len(log.Gpu),
		GPUs:     make([]GPU, len(log.Gpu)),
	}
	for i, gpu := range log.Gpu {
		info.GPUs[i] = GPU{
			Name:  gpu.ProductName,
			State: gpu.PerformanceState,
			Mem: GPUMemory{
				TotalMB:    parseMiB(gpu.FbMemoryUsage.Total),
				ReservedMB: parseMiB(gpu.FbMemoryUsage.Reserved),
				UsedMB:     parseMiB(gpu.FbMemoryUsage.Used),
				FreeMB:     parseMiB(gpu.FbMemoryUsage.Free),
			},
			Utilization: GPUUtilization{
				GPU:    parsePercent(gpu.Utilization.GpuUtil),
				Memory: parsePercent(gpu.Utilization.MemoryUtil),
			},
		}
	}
	return info, nil
}

func parseMiB(str string) int {
	parts := strings.Split(str, "MiB")
	return parseInt(parts[0], -1)
}

func parsePercent(str string) int {
	parts := strings.Split(str, "%")
	return parseInt(parts[0], -1)
}

func parseInt(str string, def int) int {
	num, err := strconv.Atoi(strings.TrimSpace(str))
	if err != nil {
		return def
	}
	return num
}

type NVInfo struct {
	GPUCount int
	GPUs     []GPU
}

type GPU struct {
	Name        string
	State       string
	Mem         GPUMemory
	Utilization GPUUtilization
}

type GPUMemory struct {
	TotalMB    int
	ReservedMB int
	UsedMB     int
	FreeMB     int
}

type GPUUtilization struct {
	GPU    int
	Memory int
}
