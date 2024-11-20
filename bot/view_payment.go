package bot

import (
	"fmt"
	"time"

	"gopkg.in/macaron.v1"
)

func viewPayment(ctx *macaron.Context) {
	pr := &GeneralResponse{Success: false}
	tgid := getTgId(ctx)

	if tgid != 0 {
		u := getUserOrCreate2(tgid, "", "")
		b := getBalance(u.AddressDeposit)
		if u.Balance != b {
			if b > u.Balance {
				new := (b - u.Balance)
				new *= 10
				u.TMU += new
				now := time.Now()
				u.TimeLock = &now
				pr.Success = true

				if time.Since(u.LastUpdated).Hours() > (60 * 24) {
					u.LastUpdated = now
				}

				if u.ReferrerID != nil {
					r := u.Referrer
					r.TMU += (new * 25 / 100)
					if err := db.Save(r).Error; err != nil {
						loge(err)
					}
					notify(fmt.Sprintf(lNewRefTmu, float64((new*25/100))/float64(Mul9)), r.TelegramId)
				}

				notify(fmt.Sprintf(lNewMint, float64(new)/float64(Mul9)), Group)
			}
			u.Balance = b
			if err := db.Save(u).Error; err != nil {
				loge(err)
			}
		}
	}

	ctx.Header().Add("Access-Control-Allow-Origin", "*")

	ctx.JSON(200, pr)
}
