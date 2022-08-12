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

func (user *User) SendMailConfirmation(confirmationdata *ConfirmationData) {

	//password := "hulisbfeulyecjpc"
	smarttader := "gntgkspsfqmkwech"

	host := "smtp.gmail.com"

	gmailAuth := smtp.PlainAuth("", "smarttraderkenya", smarttader, host)

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
		Name: confirmationdata.UserName,
	})

	senderr := smtp.SendMail(address, gmailAuth, "smarttraderkenya@gmail.com", []string{confirmationdata.Email}, body.Bytes())

	if senderr != nil {
		log.Fatal(senderr)
	}

}

func (user *User) SendOtpCode(email string) string {
	code := GenerateCode()
	getGmailAuth(email, "html/otp.html", struct {
		Code string
	}{
		Code: code[:6],
	})
	return code
}

func getGmailAuth(email, filename string, emailBody interface{}) {
	//password := "hulisbfeulyecjpc"
	smarttader := "gntgkspsfqmkwech"

	host := "smtp.gmail.com"

	// Configure hermes by setting a theme and your product info

	gmailAuth := smtp.PlainAuth("", "smarttraderkenya@gmail.com", smarttader, host)

	t, err := template.ParseFiles(filename)
	address := host + ":" + os.Getenv("MAILPORT")

	if err != nil {
		panic(err)
	}
	var body bytes.Buffer

	headers := "MIME-version : 1.0;\nContent-Type: text/html;"

	body.Write([]byte(fmt.Sprintf("Subject:Account Registration\n%s\n\n", headers)))

	Terr := t.Execute(&body, emailBody)

	if Terr != nil {
		log.Fatal(Terr)
	}

	senderr := smtp.SendMail(address, gmailAuth, "smarttraderkenya@gmail.com", []string{email}, body.Bytes())

	if senderr != nil {
		log.Fatal(senderr)
	}
}

func GenerateCode() string {
	return fmt.Sprint(time.Now().Nanosecond())
}
