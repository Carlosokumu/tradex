package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName              string   `gorm:"size:255;not null" form:"firstname"`
	LastName               string   `gorm:"size:255;not null" form:"lastname"`
	Username               string   `gorm:"size:150;not null;unique" form:"username"`
	Email                  string   `gorm:"size:100;not null;unique" form:"email"`
	Password               string   `gorm:"size:100;not null;unique" form:"password"`
	PhoneNumber            string   `gorm:"size:50;not null;unique" form:"phonenumber,omitempty"`
	Balance                *float32 `gorm:"default:0" form:"balance"`
	PercentageContribution *float32 `gorm:"default:0" form:"contribution,omitempty"`
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

func (user *User) GetMtAccountBalance() (*float32, error) {
	client := &http.Client{}

	/**
		    Fetch data from Mt4 api through nodejs sdk  provided.
	        Will switch to RabbitMq to make responses  fast.
	*/
	req, err := http.NewRequest("GET", "https://mt4functions.herokuapp.com/account", nil)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	//Set  headers to the requests
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	//Use the client to make the requests with the given [configurations]
	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}
	var mt4account Mt4Account

	err = json.Unmarshal(bodyBytes, &mt4account)

	if err != nil {
		return nil, err
	}

	fmt.Println(mt4account.Balance)
	return &mt4account.Balance, nil
}
