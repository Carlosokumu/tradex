package models

import "gorm.io/gorm"

type Transactions struct {
	gorm.Model
	DepositedBy string  `gorm:"size:255;not null;unique" form:"username"`
	Amount      float32 `gorm:"size:255;not null" form:"amount"`
}
