package models

type Credentials struct {
	UserName string `form:"username"`
	Password string `form:"password"`
}
