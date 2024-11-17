package bot

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initDb() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("robot.db"), &gorm.Config{})

	if err != nil {
		log.Println(err)
	}

	sqlDB, err := db.DB()
	sqlDB.SetMaxOpenConns(1)

	if err != nil {
		log.Println(err)
	}

	// if err := db.AutoMigrate(&User{}); err != nil {
	// 	panic(err.Error())
	// }

	return db
}
