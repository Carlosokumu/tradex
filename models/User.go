package models

import (
	"net/smtp"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName   string `form:"firstname"`
	LastName    string `form:"lastname"`
	Username    string `form:"username"`
	Email       string `form:"email"`
	Password    string `form:"password"`
	PhoneNumber string `form:"phonenumber,omitempty"`
}

func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}

func (user *User) SendMail() {
	from := "carlosokumu254@gmail.com"
	password := "hulisbfeulyecjpc"

	toEmailAddress := "coderokush@gmail.com"
	to := []string{toEmailAddress}

	host := "smtp.gmail.com"
	//port := "587"
	port := os.Getenv("MAILPORT")
	address := host + ":" + port

	subject := "Subject: This is the subject of the mail\n"
	body := "This is the body of the mail"
	message := []byte(subject + body)

	auth := smtp.PlainAuth("", from, password, host)

	err := smtp.SendMail(address, auth, from, to, message)
	if err != nil {
		panic(err)
	}
}
