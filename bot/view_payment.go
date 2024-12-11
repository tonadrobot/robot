package bot

import (
	"gopkg.in/macaron.v1"
)

func viewPayment(ctx *macaron.Context) {
	pr := &GeneralResponse{Success: false}
	tgid := getTgId(ctx)

	if tgid != 0 {
		u := getUser(tgid)
		if u.processTmuPayments() {
			pr.Success = true
		}
	}

	ctx.Header().Add("Access-Control-Allow-Origin", "*")

	ctx.JSON(200, pr)
}
