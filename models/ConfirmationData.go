package models

type ConfirmationData struct {
	UserName string `form:"username"`
	Email    string `form:"email"`
}
