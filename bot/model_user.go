package bot

import (
	"fmt"
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
	CycleCount       uint64
	CycleCountTotal  uint64
	MiningTime       time.Time `gorm:"default:'2024-12-03 16:00:00.390330053+01:00'"`
	LastNotification time.Time `gorm:"default:'2024-12-03 16:00:00.390330053+01:00'"`
	LastTxLT         uint64    `gorm:"default:0"`
	LastTxHash       string    `gorm:"default:''"`
}

func (u *User) rewards() uint64 {
	r := uint64(0)

	if !u.isFollower() {
		return r
	}

	r = uint64(time.Since(u.LastUpdated).Seconds() * float64(u.TMU) / (2400 * 3600))

	cycleIndex := float64(u.CycleCount) / float64(time.Since(u.MiningTime).Hours()/24)
	if cycleIndex > 1 {
		cycleIndex = 1
	}

	// log.Printf("cycle index: %s %.9f", u.Name, cycleIndex)

	r = uint64(float64(r) * cycleIndex)

	return r
}

func (u *User) compound() {
	u.TMU += u.rewards()
	u.CompoundCount++
	if u.CycleCount > 0 {
		u.CycleCountTotal += (u.CycleCount - 1)
	}
	u.CycleCount = 1
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
		// log.Println(err)
		return false
	}

	cb, err := b.ChatByID(Board)
	if err != nil {
		// loge(err)
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

func (u *User) isActive() bool {
	return time.Since(u.MiningTime).Minutes() <= 2280
}

func (u *User) processTmuPayments() bool {
	new := checkNewTmu(u)
	// checkNewTmu(u)

	if new >= 50000000 {
		new *= 10
		u.TMU += new
		now := time.Now()
		u.TimeLock = &now

		if time.Since(u.LastUpdated).Hours() > (60 * 24) {
			u.LastUpdated = now
		}

		if err := db.Save(u).Error; err != nil {
			loge(err)
		}

		if u.ReferrerID != nil {
			r := u.Referrer
			r.TMU += (uint64(new) * 25 / 100)
			if err := db.Save(r).Error; err != nil {
				loge(err)
			}
			notify(fmt.Sprintf(lNewRefTmu, float64((new*25/100))/float64(Mul9)), r.TelegramId)
		}

		notify(fmt.Sprintf(lNewMint, float64(new)/float64(Mul9)), GroupHall)

		go splitPayment(u)

		return true
	}

	return false
}

func getUserOrCreate(c telebot.Context) (*User, error) {
	u := &User{}

	code := c.Sender().Username
	if len(code) == 0 {
		code = generateCode()
	}

	if res := db.Where(&User{TelegramId: c.Sender().ID}).Attrs(
		&User{
			TMU:              10000000,
			Code:             code,
			AddressWithdraw:  code,
			AddressDeposit:   code,
			LastUpdated:      time.Now(),
			LastNotification: time.Now(),
			MiningTime:       time.Now(),
			Name:             c.Sender().FirstName,
		}).FirstOrCreate(u); res.Error != nil {

		loge(res.Error)
		return u, res.Error
	}

	if u.AddressDeposit == u.Code {
		notify(lNewUser, Group)
		s, a, err := generateSeedAddress()
		if err != nil {
			return u, err
		}
		u.AddressDeposit = a
		u.Seed = s
		if err := db.Save(u).Error; err != nil {
			loge(err)
			return u, err
		}
	}

	p := c.Message().Payload

	if u.ReferrerID == nil && len(p) > 0 {
		r := getUserByCode(p)
		if r.ID != 0 && r.ID != u.ID {
			u.ReferrerID = &r.ID
			if err := db.Save(u).Error; err != nil {
				loge(err)
				return u, err
			}
		}
	}

	return u, nil
}

func getUserOrCreate2(tgid int64, code string, name string) (*User, error) {
	u := &User{}

	if code == "undefined" {
		code = generateCode()
	}

	if res := db.Preload("Referrer").Where(&User{TelegramId: tgid}).Attrs(
		&User{
			TMU:              10000000,
			Code:             code,
			AddressWithdraw:  code,
			AddressDeposit:   code,
			LastUpdated:      time.Now(),
			LastNotification: time.Now(),
			MiningTime:       time.Now(),
			Name:             name,
		}).FirstOrCreate(u); res.Error != nil {

		loge(res.Error)
		return u, res.Error
	}

	if u.AddressDeposit == u.Code {
		s, a, err := generateSeedAddress()
		if err != nil {
			return u, err
		}
		u.AddressDeposit = a
		u.Seed = s
		if err := db.Save(u).Error; err != nil {
			loge(err)
			if err != nil {
				return u, err
			}
		}
		notify(lNewUser, Group)
	}

	return u, nil
}

func getUserByCode(code string) *User {
	u := &User{}

	db.First(u, &User{Code: code})

	return u
}

func getUser(tgid int64) *User {
	u := &User{}

	db.Preload("Referrer").First(u, &User{TelegramId: tgid})

	return u
}
