package bot

import (
	"github.com/go-macaron/cache"
	macaron "gopkg.in/macaron.v1"
)

func initMacaron() *macaron.Macaron {
	mac := macaron.Classic()

	mac.Use(macaron.Renderer())
	mac.Use(cache.Cacher())

	mac.Get("/data/:telegramid/:referral/:code", viewData)

	go mac.Run("0.0.0.0", 4040)

	return mac
}
