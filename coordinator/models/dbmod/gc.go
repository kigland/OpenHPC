package dbmod

type GC struct {
	ID    string `json:"id" gorm:"primaryKey"`
	Value string `json:"value" gorm:"type:text"`
}
