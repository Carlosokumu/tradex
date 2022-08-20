package models

type DepositDetails struct {
	Amount      float64 `form:"amount"`
	PhoneNumber string  `form:"phonenumber"`
	UserName    string  `form:"username"`
}
