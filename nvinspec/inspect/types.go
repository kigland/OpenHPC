package inspect

import "time"

type Options struct {
	NvidiaSMI string
	Podman    string
	ProcRoot  string
	UseSudo   bool
	Timeout   time.Duration
}

type GPUProcess struct {
	GPUUUID       string `json:"gpu_uuid"`
	PID           int    `json:"pid"`
	ProcessName   string `json:"process_name"`
	UsedMemoryMiB int    `json:"used_memory_mib"`
}

type ContainerInfo struct {
	ID    string `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Image string `json:"image,omitempty"`
}

type Finding struct {
	GPUUUID        string `json:"gpu_uuid"`
	PID            int    `json:"pid"`
	ProcessName    string `json:"process_name"`
	UsedMemoryMiB  int    `json:"used_memory_mib"`
	ContainerID    string `json:"container_id,omitempty"`
	ContainerName  string `json:"container_name,omitempty"`
	ContainerImage string `json:"container_image,omitempty"`
	Source         string `json:"source"`
	Reason         string `json:"reason,omitempty"`
}

type Report struct {
	Findings []Finding `json:"findings"`
	Warnings []string  `json:"warnings,omitempty"`
}

func (o *Options) setDefaults() {
	if o.NvidiaSMI == "" {
		o.NvidiaSMI = "nvidia-smi"
	}
	if o.Podman == "" {
		o.Podman = "podman"
	}
	if o.ProcRoot == "" {
		o.ProcRoot = "/proc"
	}
	if o.Timeout <= 0 {
		o.Timeout = 8 * time.Second
	}
}
