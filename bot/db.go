package bot

import (
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Initializes DB object
func initDb() *gorm.DB {
	var db *gorm.DB
	var err error

	if conf.Dev {
		db, err = gorm.Open(sqlite.Open(conf.DbURI), &gorm.Config{})
	} else {
		db, err = gorm.Open(postgres.Open(conf.DbURI), &gorm.Config{})
	}

	if err != nil {
		loge(err)
	}

	if err := db.AutoMigrate(&User{}, &Transaction{}); err != nil {
		panic(err.Error())
	}

	return db
}
