package dbmod

type Token struct {
	Token  string `json:"token" gorm:"primaryKey"`
	UserId string `json:"user_id" gorm:"type:text;not null;index"`
}
