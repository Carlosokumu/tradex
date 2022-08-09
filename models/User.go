package models

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"time"

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

	password := "hulisbfeulyecjpc"

	host := "smtp.gmail.com"

	gmailAuth := smtp.PlainAuth("", "carlosokumu254@gmail.com", password, host)

	t, err := template.ParseFiles("html/registration.html")
	address := host + ":" + os.Getenv("MAILPORT")

	if err != nil {
		panic(err)
	}
	var body bytes.Buffer

	headers := "MIME-version : 1.0;\nContent-Type: text/html;"

	body.Write([]byte(fmt.Sprintf("Subject:Account Registration\n%s\n\n", headers)))

	t.Execute(&body, struct {
		Name string
	}{
		Name: "Carlos",
	})

	senderr := smtp.SendMail(address, gmailAuth, "carlosokumu254@gmail.com", []string{"coderokush@gmail.com"}, body.Bytes())

	if senderr != nil {
		log.Fatal(senderr)
	}

}

func (user *User) SendOtpCode() string {
	code := GenerateCode()
	getGmailAuth("html/otp.html", struct {
		Code string
	}{
		Code: code[:6],
	})
	return code
}

func getGmailAuth(filename string, emailBody interface{}) {
	password := "hulisbfeulyecjpc"

	host := "smtp.gmail.com"

	// Configure hermes by setting a theme and your product info

	gmailAuth := smtp.PlainAuth("", "carlosokumu254@gmail.com", password, host)

	t, err := template.ParseFiles(filename)
	address := host + ":" + os.Getenv("MAILPORT")

	if err != nil {
		panic(err)
	}
	var body bytes.Buffer

	headers := "MIME-version : 1.0;\nContent-Type: text/html;"

	body.Write([]byte(fmt.Sprintf("Subject:Account Registration\n%s\n\n", headers)))

	t.Execute(&body, emailBody)

	senderr := smtp.SendMail(address, gmailAuth, "carlosokumu254@gmail.com", []string{"coderokush@gmail.com"}, body.Bytes())

	if senderr != nil {
		log.Fatal(senderr)
	}
}

func GenerateCode() string {
	return fmt.Sprint(time.Now().Nanosecond())
}
