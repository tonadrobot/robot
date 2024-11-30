package bot

import (
	"log"
	"time"

	"gopkg.in/telebot.v3"
)

var b2 *telebot.Bot

func initTelegram(key string) *telebot.Bot {
	b, err := telebot.NewBot(telebot.Settings{
		Token:     key,
		Poller:    &telebot.LongPoller{Timeout: 30 * time.Second},
		Verbose:   false,
		ParseMode: "html",
	})

	if err != nil {
		log.Fatal(err)
	}

	b.Handle("/start", commandStart)
	b.Handle("/stats", commandStats)
	b.Handle("/ranks", commandRanks)

	b2 = initTelegramOld("")

	return b
}

func initTelegramOld(key string) *telebot.Bot {
	b, err := telebot.NewBot(telebot.Settings{
		Token:     key,
		Poller:    &telebot.LongPoller{Timeout: 30 * time.Second},
		Verbose:   false,
		ParseMode: "html",
	})

	if err != nil {
		log.Fatal(err)
	}

	return b
}

func notify(msg string, tgid int64) {
	rec := &telebot.Chat{
		ID: tgid,
	}
	b.Send(rec, msg, telebot.NoPreview)
}

func notifyold(msg string, tgid int64) {
	ab := getAppButton()
	rec := &telebot.Chat{
		ID: tgid,
	}
	b2.Send(rec, msg, ab, telebot.NoPreview)
}
