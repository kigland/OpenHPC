package dbmod

type User struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`

	MaxVCPU   int `json:"max_vcpu"`   // -1 means no limit
	MaxVGPU   int `json:"max_vgpu"`   // -1 means no limit
	MaxMemory int `json:"max_memory"` // -1 means no limit
}
