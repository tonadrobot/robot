package bot

import (
	"log"
	"time"

	"gopkg.in/telebot.v3"
)

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

	return b
}

func notify(msg string, tgid int64) {
	rec := &telebot.Chat{
		ID: tgid,
	}
	b.Send(rec, msg, telebot.NoPreview)
}

func notifytest(msg string, tgid int64) {
	mb := getMiningButton()
	rec := &telebot.Chat{
		ID: tgid,
	}
	b.Send(rec, msg, mb, telebot.NoPreview)
}

func notifystart(msg string, tgid int64) {
	sb := getStartButton()
	rec := &telebot.Chat{
		ID: tgid,
	}
	_, err := b.Send(rec, msg, sb, telebot.NoPreview)
	if err != nil {
		loge(err)
	}
}

func getMiningButton() *telebot.ReplyMarkup {
	rm := &telebot.ReplyMarkup{}
	btn := rm.URL("⚪️ Restart Mining", "https://t.me/TonCityRobot/miner?startapp=restart")

	rm.Inline(
		rm.Row(btn),
	)

	return rm
}

func getRestartButton() *telebot.ReplyMarkup {
	rm := &telebot.ReplyMarkup{}
	btn := rm.URL("⚪️ Restart Mining", "https://t.me/DevTonCityBoard/6")

	rm.Inline(
		rm.Row(btn),
	)

	return rm
}

func getStartButton() *telebot.ReplyMarkup {
	rm := &telebot.ReplyMarkup{}
	btn := rm.URL("⚪️ Start Mining TON", "https://t.me/TonCityRobot/miner")

	rm.Inline(
		rm.Row(btn),
	)

	return rm
}

func notifyEnd(tgid int64) {
	rb := getRestartButton()

	rec := &telebot.Chat{
		ID: tgid,
	}

	_, err := b.Send(rec, lCycleFinished, rb, telebot.NoPreview)
	if err != nil {
		loge(err)
	}
}
