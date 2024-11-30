package bot

import (
	"fmt"

	"gopkg.in/telebot.v3"
)

func commandRanks(c telebot.Context) error {
	var users []*User
	msg := lRanks

	db.Order("tmu desc").Limit(10).Find(&users)

	for _, u := range users {
		tmu := float64(u.TMU) / float64(Mul9)
		msg += fmt.Sprintf("\n<b>%s</b> - <code>%.9f TMU (%d)</code>", u.Name, tmu, u.CompoundCount)
	}

	_, err := b.Send(c.Chat(), msg, telebot.NoPreview)
	return err
}
