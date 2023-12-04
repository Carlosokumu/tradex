package models

type EmailPassword struct {
	Email string    `form:"email"`
	Password string `form:"password"`
}