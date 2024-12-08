package bot

import (
	"time"

	"gopkg.in/macaron.v1"
)

func viewRestart(ctx *macaron.Context) {
	rr := &GeneralResponse{Success: true}
	tgid := getTgId(ctx)

	if tgid != 0 {
		u := getUser(tgid)

		if time.Since(u.MiningTime).Minutes() > 1410 {
			u.MiningTime = time.Now()
			u.CycleCount++
			if err := db.Save(u).Error; err != nil {
				loge(err)
			}
		} else {
			rr.Success = false
		}
	}

	ctx.Header().Add("Access-Control-Allow-Origin", "*")

	ctx.JSON(200, rr)
}
