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

	// Configure hermes by setting a theme and your product info

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

	// m := gomail.NewMessage()
	// m.SetHeader("From", "carlosokumu254@gmail.com")
	// m.SetHeader("To", "coderokush@gmail.com")
	// m.SetHeader("RE", "Account Registration  Successful")
	// m.SetHeader("Subject", "Account Registration")
	// m.SetBody("text/html",emailBody)
	// m.SetBody("text/html",
	// 	`
	// <center>
	// <u>Registration confirmation</u>
	// </center>
	// <br>
	// <h2 style = "margin-top: 1cm"> RE: Account Registration sucessful</h2>
	// <br>
	// <img style = "margin-top: 1cm" src="https://www.freepnglogos.com/uploads/play-store-logo-png/play-store-logo-nisi-filters-australia-11.png" alt="My image" width="150" height="70" />
	// <br>
	// <p> Hi <var>username</var></p>
	// <br>
	// <p> Thank you for joining smart trader Community</p>
	// <br>
	// <p>It is recommended  you take enough time to read through the terms and conditions before making any deposits</p>
	// <br>
	// <p style = "margin-top: 1cm"> Download our app here</p>
	// <br>
	// <a href="https://www.qries.com/"> <img alt="Qries" src="https://www.freepnglogos.com/uploads/play-store-logo-png/play-store-logo-nisi-filters-australia-11.png" width=150" height="70"></a>
	// <br>
	// <p>You can view your account here</p>
	// <br>
	// <a href = "url">https:linktoaccount.com</a>
	// `)

	// d := gomail.NewPlainDialer(host, 587, "carlosokumu254@gmail.com", password)

	// if err := d.DialAndSend(m); err != nil {
	// 	panic(err)
	// }
}

func SendOtpCode() {
	code := GenerateCode()
	getGmailAuth("html/otp.html", struct {
		Code string
	}{
		Code: code[:6],
	})
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
	return time.Nanosecond.String()
}
