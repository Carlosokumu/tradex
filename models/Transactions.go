package models

import "gorm.io/gorm"

type Transactions struct {
	gorm.Model
	DepositedBy string  `gorm:"size:255;not null" form:"username"`
	Amount      float32 `gorm:"size:255;not null" form:"amount"`
	PhoneNumber string  `gorm:"size:255;not null" form:"phonenumber"`
}
