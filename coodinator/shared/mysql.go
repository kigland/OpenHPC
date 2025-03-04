package shared

import (
	"log"

	"github.com/kigland/HPC-Scheduler/coodinator/models/dbmod"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func initMySQL() {
	var err error
	DB, err = gorm.Open(mysql.Open(GetConfig().MySQL), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established successfully")

	dbmod.AutoMigrate(DB)
}
