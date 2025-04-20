package stat

import (
	"github.com/gin-gonic/gin"
	"github.com/kigland/OpenHPC/coordinator/models/apimod"
	"github.com/kigland/OpenHPC/lib/nv"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

func StateHandler(c *gin.Context) {
	c.JSON(200, apimod.StatInfo{
		Gpu: getGPUState(),
		Mem: getMemState(),
		Cpu: getCpuState(),
	})
}

func getGPUState() []apimod.GpuStatInfo {
	info, err := nv.GetNvidiaSmiLog()
	if err != nil {
		return nil
	}
	infoParsed, err := info.Parse()
	if err != nil {
		return nil
	}
	var res []apimod.GpuStatInfo
	for _, gpu := range infoParsed.GPUs {
		res = append(res, apimod.GpuStatInfo{
			Name:    gpu.Name,
			MemUtil: int32(gpu.Utilization.Memory),
			Util:    int32(gpu.Utilization.GPU),
			Mem: apimod.MemStatInfo{
				Total: int32(gpu.Mem.TotalMB),
				Used:  int32(gpu.Mem.UsedMB),
				Unit:  "MiB",
			},
		})
	}

	return res
}

func getMemState() apimod.MemStatInfo {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return apimod.MemStatInfo{}
	}
	return apimod.MemStatInfo{
		Total: int32(memInfo.Total / (1024 * 1024)),
		Used:  int32(memInfo.Used / (1024 * 1024)),
		Unit:  "MiB",
	}
}

func getCpuState() apimod.CpuStatInfo {
	percentages, err := cpu.Percent(0, true)
	if err != nil {
		return apimod.CpuStatInfo{}
	}
	sumPerc := 0.0
	count := 0
	for _, percentage := range percentages {
		sumPerc += percentage
		count++
	}
	avgPerc := sumPerc / float64(count)
	return apimod.CpuStatInfo{
		AvgLoad: float32(avgPerc),
	}
}
