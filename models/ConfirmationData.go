package models

type ConfirmationData struct {
	UserName string `form:"username"`
	Email    string `form:"email"`
}
type MasterAccount struct {
	AccountLogin uint
	Balance      *int64
}
