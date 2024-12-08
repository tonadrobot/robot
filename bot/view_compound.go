package bot

import (
	"fmt"

	"gopkg.in/macaron.v1"
)

func viewCompound(ctx *macaron.Context) {
	cr := &GeneralResponse{Success: true}
	tgid := getTgId(ctx)

	if tgid != 0 {
		u := getUser(tgid)
		u.compound()

		r := u.Referrer

		if r != nil && r.ID != 0 && u.TMU >= 10100000 && !u.ReferralActive {
			r.TMU += 2500000
			if err := db.Save(r).Error; err != nil {
				loge(err)
			}
			msg := fmt.Sprintf(lNewRefTmu, float64(2500000)/float64(Mul9))
			notify(msg, r.TelegramId)

			u.ReferralActive = true
			if err := db.Save(u).Error; err != nil {
				loge(err)
			}
		}
	}

	ctx.Header().Add("Access-Control-Allow-Origin", "*")

	ctx.JSON(200, cr)
}

type GeneralResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
