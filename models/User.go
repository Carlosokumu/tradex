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
	//FirstName   string `form:"firstname"`
	FirstName   string   `gorm:"size:255;not null" form:"firstname"`
	LastName    string   `gorm:"size:255;not null;unique" form:"lastname"`
	Username    string   `gorm:"size:150;not null;unique" form:"username"`
	Email       string   `gorm:"size:100;not null;unique" form:"email"`
	Password    string   `gorm:"size:100;not null;unique" form:"password"`
	PhoneNumber string   `gorm:"size:50;not null;unique" form:"phonenumber,omitempty"`
	Balance     *float32 `gorm:"default:0" form:"balance"`

	//LastName    string `form:"lastname"`
	//Username    string `form:"username"`
	//Email       string `form:"email"`
	//Password    string `form:"password"`
	//PhoneNumber string `form:"phonenumber,omitempty"`
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
	smarttader := "jxkpndrkvjbceokd"

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

	terr := t.Execute(&body, emailBody)

	if terr != nil {
		log.Fatal(terr)
	}

	senderr := smtp.SendMail(address, gmailAuth, "smarttraderkenya@gmail.com", []string{email}, body.Bytes())

	if senderr != nil {
		log.Fatal(senderr)
	}
}

func GenerateCode() string {
	return fmt.Sprint(time.Now().Nanosecond())
}
