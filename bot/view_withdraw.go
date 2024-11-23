package bot

import (
	"time"

	"gopkg.in/macaron.v1"
)

func viewWithdraw(ctx *macaron.Context) {
	wr := &GeneralResponse{Success: true}
	tgid := getTgId(ctx)

	if tgid != 0 {
		u := getUserOrCreate2(tgid, "", "")
		send(int64(u.rewards()/10-5000000), u.AddressWithdraw, conf.Seed)
		u.LastUpdated = time.Now()
		u.delayedUpdateBalance()
		if err := db.Save(u).Error; err != nil {
			loge(err)
		}
	}

	ctx.Header().Add("Access-Control-Allow-Origin", "*")

	ctx.JSON(200, wr)
}
