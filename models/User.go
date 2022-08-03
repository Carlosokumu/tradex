package models

import (
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
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

	m := gomail.NewMessage()
	m.SetHeader("From", "carlosokumu254@gmail.com")
	m.SetHeader("To", "coderokush@gmail.com")
	m.SetHeader("RE", "Account Registration  Successful")
	m.SetHeader("Subject", "Account Registration")
	m.SetBody("text/html", `<center><u>Registration confirmation</u></center> <br> <p style = "margin-top: 1cm"> RE: Account Registration sucessful</p> <br> <img style = "margin-top: 1cm" src="https://i.pinimg.com/originals/aa/19/47/aa1947e08757e6a7d17724677ac850e6.jpg" alt="My image" /> <br> <p> Hi carlos</p><br> <p> Thank you for joining smartrader Community</p><br><p>It is recommended you take enough time to read through the terms and conditions before making any deposits</p><br> <p style = "margin-top: 1cm"> Download our app here</p><br><a href="https://www.qries.com/"> <img alt="Qries" src="https://www.freepnglogos.com/uploads/play-store-logo-png/play-store-logo-nisi-filters-australia-11.png"
	width=150" height="70"></a> `)

	d := gomail.NewPlainDialer(host, 587, "carlosokumu254@gmail.com", password)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
