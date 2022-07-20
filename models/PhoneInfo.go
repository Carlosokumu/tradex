package models

type PhoneInfo struct {
	UserName    string `form:"username"`
	PhoneNumber string `form:"phonenumber"`
}
