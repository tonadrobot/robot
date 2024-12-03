package bot

import (
	"log"
	"time"

	"gopkg.in/telebot.v3"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	TelegramId       int64  `gorm:"size:255;uniqueIndex"`
	Code             string `gorm:"size:255;uniqueIndex"`
	AddressWithdraw  string `gorm:"size:255;uniqueIndex"`
	AddressDeposit   string `gorm:"size:255;uniqueIndex"`
	Seed             string `gorm:"size:255"`
	TMU              uint64
	Balance          uint64
	LastUpdated      time.Time
	TimeLock         *time.Time
	ReferrerID       *uint
	Referrer         *User
	Name             string `gorm:"size:255"`
	ReferralActive   bool   `gorm:"default:false"`
	CompoundCount    uint64
	MiningTime       time.Time `gorm:"default:'2024-12-03 23:00:00.390330053+01:00'"`
	LastNotification time.Time `gorm:"default:'2024-12-03 23:00:00.390330053+01:00'"`
}

func (u *User) rewards() uint64 {
	r := uint64(0)

	if !u.isFollower() {
		return r
	}

	r = uint64(time.Since(u.LastUpdated).Seconds() * float64(u.TMU) / (2400 * 3600))

	return r
}

func (u *User) compound() {
	u.TMU += u.rewards()
	u.CompoundCount++
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

func (u *User) isFollower() bool {
	ut, err := b.ChatByID(u.TelegramId)
	if err != nil {
		// loge(err)
		log.Println(err)
		return false
	}

	cb, err := b.ChatByID(Board)
	if err != nil {
		loge(err)
		return false
	}

	cm, err := b.ChatMemberOf(cb, ut)
	if err != nil {
		loge(err)
		return false
	}

	if cm.Role == "member" ||
		cm.Role == "administrator" ||
		cm.Role == "creator" {
		return true
	}

	return false
}

func (u *User) isMember() bool {
	ut, err := b.ChatByID(u.TelegramId)
	if err != nil {
		loge(err)
		return false
	}

	cb, err := b.ChatByID(Group)
	if err != nil {
		loge(err)
		return false
	}

	cm, err := b.ChatMemberOf(cb, ut)
	if err != nil {
		loge(err)
		return false
	}

	if cm.Role == "member" ||
		cm.Role == "administrator" ||
		cm.Role == "creator" {
		return true
	}

	return false
}

func (u *User) hasMigrated() bool {
	_, err := b.ChatByID(u.TelegramId)
	if err != nil {
		// loge(err)
		log.Println(err)
		return false
	}

	return true
}

func getUserOrCreate(c telebot.Context) *User {
	u := &User{}

	code := c.Sender().Username
	if len(code) == 0 {
		code = generateCode()
	}

	if res := db.Where(&User{TelegramId: c.Sender().ID}).Attrs(
		&User{
			TMU:             10000000,
			Code:            code,
			AddressWithdraw: code,
			AddressDeposit:  code,
			LastUpdated:     time.Now(),
			Name:            c.Sender().FirstName,
		}).FirstOrCreate(u); res.Error != nil {

		loge(res.Error)
	}

	if u.AddressDeposit == u.Code {
		notify(lNewUser, Group)
		s, a := generateSeedAddress()
		u.AddressDeposit = a
		u.Seed = s
		if err := db.Save(u).Error; err != nil {
			loge(err)
		}
	}

	p := c.Message().Payload

	if u.ReferrerID == nil && len(p) > 0 {
		r := getUserByCode(p)
		if r.ID != 0 && r.ID != u.ID {
			u.ReferrerID = &r.ID
			if err := db.Save(u).Error; err != nil {
				loge(err)
			}
		}
	}

	return u
}

func getUserOrCreate2(tgid int64, code string, name string) *User {
	u := &User{}

	if code == "undefined" {
		code = generateCode()
	}

	if res := db.Preload("Referrer").Where(&User{TelegramId: tgid}).Attrs(
		&User{
			TMU:             10000000,
			Code:            code,
			AddressWithdraw: code,
			AddressDeposit:  code,
			LastUpdated:     time.Now(),
			Name:            name,
		}).FirstOrCreate(u); res.Error != nil {

		loge(res.Error)
	}

	if u.AddressDeposit == u.Code {
		notify(lNewUser, Group)
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
