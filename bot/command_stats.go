package bot

import "gopkg.in/telebot.v3"

func commandStats(c telebot.Context) error {
	_, err := b.Send(c.Sender(), "stats", telebot.NoPreview)
	return err
}
