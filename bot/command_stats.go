package bot

import (
	"fmt"

	"gopkg.in/telebot.v3"
)

func commandStats(c telebot.Context) error {
	count := int64(0)
	db.Find(&User{}).Count(&count)

	msg := fmt.Sprintf(lStats, count)

	_, err := b.Send(c.Chat(), msg, telebot.NoPreview)
	return err
}
