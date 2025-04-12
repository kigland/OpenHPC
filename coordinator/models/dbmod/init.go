package dbmod

import (
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	db.AutoMigrate(&User{}, &GC{}, &Token{})
}
