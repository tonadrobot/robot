package bot

import (
	"log"

	"gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

var conf *Config

var b *telebot.Bot

var db *gorm.DB

var cch *Cache

// Package init function
func init() {
	conf = initConfig()

	db = initDb()

	b = initTelegram(conf.TelegramKey)

	initMonitor()

	cch = initCache()

	initMacaron()
}

func Start() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logs("Bot started successfully. 🚀")

	// var users []*User
	// db.Find(&users)
	// counter := 0

	// for _, u := range users {
	// 	if !u.hasMigrated() {
	// 		counter++
	// 		log.Printf("%d Not: %s", counter, u.Name)
	// 	}
	// }

	// for _, u := range users {
	// 	u.processTmuPayments()
	// 	counter++
	// 	log.Printf("%d Not: %s", counter, u.Name)
	// }

	// notifytest(lRestartMining, BoardDev)

	// u := getUser(7967928871)
	// log.Println(u.processTmuPayments())

	b.Start()
}
