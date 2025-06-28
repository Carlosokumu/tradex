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
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	UserName string `form:"username"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

type UserModel struct {
	gorm.Model
	UserName        string
	Email           string
	Password        string
	Avatar          string
	TradingAccounts []TradingAccount `gorm:"foreignKey:UserId"`
	RoleID          uint
	Role            Role
	Communities     []Community `gorm:"many2many:user_communities;"`
}

type TradingAccount struct {
	gorm.Model
	Platform  string
	AccountId string
	UserId    uint
}

type Role struct {
	gorm.Model
	Name string
}

type RunningPosition struct {
	gorm.Model
	UserID          uint
	Volume          *int64
	Price           *float64
	TradeSide       *int32
	SymbolId        *int64
	OpenTime        *int64
	Commission      *int64
	Swap            *int64
	MoneyDigits     *uint32
	PositionRisk    *float32
	PositionsReward *float32
}

func (user *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

func (userModel *UserModel) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(userModel.Password), 14)
	if err != nil {
		return err
	}
	userModel.Password = string(bytes)
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
	smarttader := "ubwmvktbfovpzonk"

	host := "smtp.gmail.com"

	gmailAuth := smtp.PlainAuth("", "swingwizardsinfo@gmail.com", smarttader, host)

	t, err := template.ParseFiles("html/registration.html")
	//address := host + ":" + os.Getenv("MAILPORT")
	address := host + ":" + "587"
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

	senderr := smtp.SendMail(address, gmailAuth, "swingwizardsinfo@gmail.com", []string{confirmationdata.Email}, body.Bytes())

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
	smarttader := "ubwmvktbfovpzonk"

	host := "smtp.gmail.com"

	// Configure hermes by setting a theme and your product info

	gmailAuth := smtp.PlainAuth("", "swingwizardsinfo@gmail.com", smarttader, host)

	t, err := template.ParseFiles(filename)
	//address := host + ":" + os.Getenv("MAILPORT")
	address := host + ":" + "587"

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

	senderr := smtp.SendMail(address, gmailAuth, "swingwizardsinfo@gmail.com", []string{email}, body.Bytes())

	if senderr != nil {
		log.Fatal(senderr)
	}
}

func GenerateCode() string {
	return fmt.Sprint(time.Now().Nanosecond())
}

func (user *User) GetMtAccountBalance() (*Mt4Account, error) {
	client := &http.Client{}

	/**
		    Fetch data from Mt4 api through nodejs sdk  provided.
	        Will switch to a message broker to make responses  faster.
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
	return &mt4account, nil
}
