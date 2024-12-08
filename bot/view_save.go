package bot

import (
	"github.com/xssnick/tonutils-go/address"
	"gopkg.in/macaron.v1"
)

func viewSave(user UserForm, ctx *macaron.Context) {
	sr := &GeneralResponse{Success: true}
	tgid := getTgId(ctx)

	_, err := address.ParseAddr(user.AddressWithdraw)
	if err != nil {
		loge(err)
		sr.Success = false
	} else {
		if tgid != 0 {
			u := getUser(tgid)
			u.AddressWithdraw = user.AddressWithdraw
			if err := db.Save(u).Error; err != nil {
				loge(err)
			}
		}
	}

	ctx.Header().Add("Access-Control-Allow-Origin", "*")

	ctx.JSON(200, sr)
}

type UserForm struct {
	AddressWithdraw string `binding:"Required"`
}
