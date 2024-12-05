package bot

import (
	"time"
)

type Cache struct {
	StatsCache *StatsCache
}

func (c *Cache) loadStatsCache() {
	tmu := float64(0)
	reward := float64(0)
	var users []*User
	db.Find(&users)
	count := len(users)
	countActive := 0

	for _, u := range users {
		tmu += (float64(u.TMU) / float64(Mul9))
		reward += (float64(u.rewards()) / float64(Mul9))
		if u.isActive() {
			countActive++
		}
	}

	c.StatsCache.Miners = count
	c.StatsCache.ActiveMiners = countActive
	c.StatsCache.TMU = tmu
	c.StatsCache.RewardTMU = reward
}

func (c *Cache) start() {
	for {
		c.loadStatsCache()

		time.Sleep(time.Second * 10)
	}
}

func initCache() *Cache {
	c := &Cache{}
	c.StatsCache = &StatsCache{}
	go c.start()

	return c
}

type StatsCache struct {
	Miners       int
	ActiveMiners int
	TMU          float64
	RewardTMU    float64
}
