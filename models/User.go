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
	//from := "carlosokumu254@gmail.com"
	password := "hulisbfeulyecjpc"

	//toEmailAddress := "coderokush@gmail.com"
	//to := []string{toEmailAddress}

	host := "smtp.gmail.com"
	//port := "587"
	//port := os.Getenv("MAILPORT")
	//address := host + ":" + port

	// subject := "Subject: This is the subject of the mail\n"
	// //body := "This is the body of the mail"
	// mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	// body := "<html><body><img src= "httpa" + "://i.pinimg.com/originals/aa/19/47/aa1947e08757e6a7d17724677ac850e6.jpg" alt="My image" /></body></html>"
	// message := []byte(subject + mime + body)

	// auth := smtp.PlainAuth("", from, password, host)

	// err := smtp.SendMail(address, auth, from, to, message)
	// if err != nil {
	// 	panic(err)
	// }

	m := gomail.NewMessage()
	m.SetHeader("From", "carlosokumu254@gmail.com")
	m.SetHeader("To", "coderokush@gmail.com")
	m.SetHeader("Subject", "Hello!")
	m.Embed("images/android.png")
	m.SetBody("text/html", `<img src="cid:images/android.png" alt="My image" />`)

	d := gomail.NewPlainDialer(host, 587, "carlosokumu254@gmail.com", password)

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
