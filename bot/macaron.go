package bot

import (
	"github.com/go-macaron/binding"
	"github.com/go-macaron/cache"
	macaron "gopkg.in/macaron.v1"
)

func initMacaron() *macaron.Macaron {
	mac := macaron.Classic()

	mac.Use(macaron.Renderer())
	mac.Use(cache.Cacher())

	mac.Get("/data/:telegramid/:referral/:code/:name", viewData)
	mac.Get("/paid/:telegramid", viewPayment)

	mac.Post("/save/:telegramid", binding.Bind(UserForm{}), viewSave)
	mac.Post("/compound/:telegramid", viewCompound)
	mac.Post("/withdraw/:telegramid", viewWithdraw)
	mac.Post("/restart/:telegramid", viewRestart)

	go mac.Run("0.0.0.0", 4040)

	return mac
}
