package bot

import (
	"gopkg.in/telebot.v3"
)

func commandStart(c telebot.Context) error {
	getUserOrCreate(c)
	ab := getAppButton()

	b.Send(c.Sender(), lStart, ab)

	// notify(lNewUser, Admin)

	return nil
}

func getAppButton() *telebot.ReplyMarkup {
	rm := &telebot.ReplyMarkup{}
	btn := rm.URL("⚪️ Launch Ton City Miner", "https://t.me/TonCityRobot/miner")

	rm.Inline(
		rm.Row(btn),
	)

	return rm
}
