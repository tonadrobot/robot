package bot

import (
	"log"

	"gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

var conf *Config

var b *telebot.Bot

var db *gorm.DB

// Package init function
func init() {
	conf = initConfig()

	db = initDb()

	b = initTelegram(conf.TelegramKey)

	initMacaron()
}

func Start() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logs("Bot started successfully. 🚀")

	b.Start()
}
