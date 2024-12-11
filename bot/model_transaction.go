package bot

import (
	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	TxHash    string `gorm:"size:255;uniqueIndex"`
	TxLT      uint64
	Processed bool `gorm:"default:1"`
}

func processTx(hash string, lt uint64) {
	t := &Transaction{}
	if res := db.Where(&Transaction{TxHash: hash, TxLT: lt}).FirstOrCreate(t); res.Error != nil {
		loge(res.Error)
	}
}

func isTxProcessed(hash string, lt uint64) bool {
	t := &Transaction{}
	if res := db.First(t, &Transaction{TxHash: hash, TxLT: lt}); res.Error != nil {
		return false
	}

	return t.ID != 0
}
