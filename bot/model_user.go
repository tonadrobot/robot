package bot

import (
	"time"

	"gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	TelegramId      int64  `gorm:"size:255;uniqueIndex"`
	Code            string `gorm:"size:255;uniqueIndex"`
	AddressWithdraw string `gorm:"size:255;uniqueIndex"`
	AddressDeposit  string `gorm:"size:255;uniqueIndex"`
	Seed            string `gorm:"size:255"`
	TMU             uint64
	Balance         uint64
	LastUpdated     time.Time
	TimeLock        *time.Time
	ReferrerID      uint
	Referrer        *User
}

func (u *User) rewards() uint64 {
	r := uint64(0)

	r = uint64(time.Since(u.LastUpdated).Seconds() * float64(u.TMU) / (3000 * 3600))

	return r
}

func (u *User) compound() {
	u.TMU += u.rewards()
	u.LastUpdated = time.Now()
	if err := db.Save(u).Error; err != nil {
		loge(err)
	}
}

func (u *User) delayedUpdateBalance() {
	go func(u *User) {
		time.Sleep(time.Minute * 1)
		b := getBalance(u.AddressDeposit)
		if b != u.Balance {
			u.Balance = b
			if err := db.Save(u).Error; err != nil {
				loge(err)
			}
		}
	}(u)
}

func getUserOrCreate(c telebot.Context) *User {
	u := &User{}

	code := c.Sender().Username
	if len(code) == 0 {
		code = generateCode()
	}

	if res := db.Where(&User{TelegramId: c.Sender().ID}).Attrs(
		&User{
			Code:            code,
			AddressWithdraw: code,
			AddressDeposit:  code,
		}).FirstOrCreate(u); res.Error != nil {

		loge(res.Error)
	}

	if u.AddressDeposit == u.Code {
		s, a := generateSeedAddress()
		u.AddressDeposit = a
		u.Seed = s
		if err := db.Save(u).Error; err != nil {
			loge(err)
		}
	}

	p := c.Message().Payload

	if u.ReferrerID == 0 && len(p) > 0 {
		r := getUserByCode(p)
		if r.ID != 0 && r.ID != u.ID {
			u.ReferrerID = r.ID
			if err := db.Save(u).Error; err != nil {
				loge(err)
			}
		}
	}

	return u
}

func getUserOrCreate2(tgid int64) *User {
	u := &User{}

	code := generateCode()

	if res := db.Preload("Referrer").Where(&User{TelegramId: tgid}).Attrs(
		&User{
			Code:            code,
			AddressWithdraw: code,
			AddressDeposit:  code,
		}).FirstOrCreate(u); res.Error != nil {

		loge(res.Error)
	} else if res.RowsAffected > 0 {
		notify(lNewUser, Admin)
	}

	if u.AddressDeposit == u.Code {
		s, a := generateSeedAddress()
		u.AddressDeposit = a
		u.Seed = s
		if err := db.Save(u).Error; err != nil {
			loge(err)
		}
	}

	return u
}

func getUserByCode(code string) *User {
	u := &User{}

	db.First(u, &User{Code: code})

	return u
}
