package models

import "gorm.io/gorm"

type Credentials struct {
	gorm.Model
	UserName string `form:"username"`
	Password string `form:"password"`
}
