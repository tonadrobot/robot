package bot

import (
	"fmt"

	"gopkg.in/telebot.v3"
)

func commandStats(c telebot.Context) error {
	msg := fmt.Sprintf(lStats, cch.StatsCache.Miners, cch.StatsCache.ActiveMiners, cch.StatsCache.TMU, cch.StatsCache.RewardTMU)

	_, err := b.Send(c.Chat(), msg, telebot.NoPreview)
	return err
}
