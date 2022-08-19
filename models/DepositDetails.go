package models

type DepositDetails struct {
	Amount      float32 `form:"amount"`
	PhoneNumber string  `form:"phonenumber"`
	UserName    string  `form:"username"`
}
