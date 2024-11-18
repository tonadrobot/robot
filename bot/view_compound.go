package bot

import "gopkg.in/macaron.v1"

func viewCompound(ctx *macaron.Context) {
	cr := &GeneralResponse{Success: true}
	tgid := getTgId(ctx)

	if tgid != 0 {
		u := getUserOrCreate2(tgid, "")
		u.compound()
	}

	ctx.Header().Add("Access-Control-Allow-Origin", "*")

	ctx.JSON(200, cr)
}

type GeneralResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
