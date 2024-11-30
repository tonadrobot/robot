package bot

import (
	"fmt"

	"gopkg.in/telebot.v3"
)

func commandStats(c telebot.Context) error {
	tmu := float64(0)
	reward := float64(0)
	var users []*User
	db.Find(&users)
	count := len(users)

	for _, u := range users {
		tmu += (float64(u.TMU) / float64(Mul9))
		reward += (float64(u.rewards()) / float64(Mul9))
	}

	msg := fmt.Sprintf(lStats, count, tmu, reward)

	_, err := b.Send(c.Chat(), msg, telebot.NoPreview)
	return err
}
