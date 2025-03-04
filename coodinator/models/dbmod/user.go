package dbmod

type User struct {
	ID       string `json:"id" gorm:"primaryKey"`
	Nickname string `json:"nickname" gorm:"type:varchar(255);not null"`
	Password string `json:"password" gorm:"type:varchar(255);not null"`

	Admin bool `json:"admin" gorm:"default:false"`

	MaxVCPU   int `json:"max_vcpu" gorm:"default:-1"`   // -1 means no limit
	MaxVGPU   int `json:"max_vgpu" gorm:"default:-1"`   // -1 means no limit
	MaxMemory int `json:"max_memory" gorm:"default:-1"` // -1 means no limit
}
