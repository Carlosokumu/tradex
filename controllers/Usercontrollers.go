package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ajg/form"
	"github.com/carlosokumu/dubbedapi/database"
	"github.com/carlosokumu/dubbedapi/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(context *gin.Context) {
	var user models.User
	var password string

	d := form.NewDecoder(context.Request.Body)
	if err := d.Decode(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Parse Error": err.Error()})
		log.Fatal(err)
		context.Abort()
		return
	}
	password = user.Password

	if err := user.HashPassword(user.Password); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Internal server error": err.Error()})
		context.Abort()
		log.Fatal(err)
		return
	}

	//Create a new user record into the database

	record := database.Instance.Create(&user)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Database Error": record.Error.Error()})
		context.Abort()
		log.Fatal(record.Error)
		return
	}

	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "firstname": user.FirstName, "lastname": user.LastName, "email": user.Email, "username": user.Username, "password": password})
}

func LoginUser(context *gin.Context) {
	var credentials models.Credentials
	var user models.User

	d := form.NewDecoder(context.Request.Body)
	if err := d.Decode(&credentials); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Parse Error": err.Error()})
		log.Fatal(err)
		context.Abort()
		return
	}
	if result := database.Instance.Table("users").Where("username = ?", credentials.UserName).First(&user).Error; result != nil {
		context.JSON(http.StatusNotFound, gin.H{"response": result.Error()})
		fmt.Println(result)
		context.Abort()
		return
	}
	bytes, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), 14)

	fmt.Println("rawpassword", credentials.Password)

	result := CheckPasswordHash(credentials.Password, user.Password)
	fmt.Println("match", result)

	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"response": "something went wrong"})
		fmt.Println(err)
		context.Abort()
		return
	}

	password := string(bytes)

	fmt.Println(password)

	if result {
		context.JSON(http.StatusOK, gin.H{"response": "success"})
	} else {
		context.JSON(http.StatusForbidden, gin.H{"response": "unmatch"})
	}

}

func UpdateUser(context *gin.Context) {

	if result := database.Instance.Table("users").Model(&models.User{}).Where("username = ?", "webman").Update("username", "kalonje"); result.Error != nil {
		log.Fatal(result.Error)
		fmt.Println("Cannot find User")
	}
	context.JSON(http.StatusOK, "Done")
}

func UpdatePhoneNumber(context *gin.Context) {
	var phoneinfo models.PhoneInfo

	d := form.NewDecoder(context.Request.Body)
	if err := d.Decode(&phoneinfo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Parse Error": err.Error()})
		log.Fatal(err)
		context.Abort()
		return
	}

	if result := database.Instance.Table("users").Model(&models.User{}).Where("username = ?", phoneinfo.UserName).Update("phone_number", phoneinfo.PhoneNumber); result.Error != nil {
		log.Fatal(result.Error)
		fmt.Println("Cannot find User")
	}
	context.JSON(http.StatusOK, gin.H{"response": "Phone Number updated Sucessfully"})
}

func SendConfirmEmail(context *gin.Context) {
	var confirmData models.ConfirmationData
	var user models.User

	d := form.NewDecoder(context.Request.Body)
	if err := d.Decode(&confirmData); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Parse Error": err.Error()})
		log.Fatal(err)
		context.Abort()
		return
	}
	user.SendMailConfirmation(&confirmData)
	context.JSON(http.StatusOK, gin.H{"response": "Sucessfully sent mail"})
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SendOtp(context *gin.Context) {
	var user models.User
	var adressInfo models.AdressInfo

	d := form.NewDecoder(context.Request.Body)
	if err := d.Decode(&adressInfo); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Parse Error": err.Error()})
		log.Fatal(err)
		context.Abort()
		return
	}
	code := user.SendOtpCode(adressInfo.EmailAdress)
	context.JSON(http.StatusOK, gin.H{"code": code[:6]})

}

func HandleDeposit(context *gin.Context) {

	var depositDetails models.DepositDetails
	var user models.User
	var transactions models.Transactions

	//Handle decode for the user trying to deposit
	d := form.NewDecoder(context.Request.Body)

	if err := d.Decode(&depositDetails); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Parse Error": err.Error()})
		log.Fatal(err)
		context.Abort()
		return
	}
	transactions = models.Transactions{
		DepositedBy: depositDetails.UserName,
		PhoneNumber: depositDetails.PhoneNumber,
		Amount:      depositDetails.Amount,
	}

	if result := database.Instance.Table("users").Where("username = ?", depositDetails.UserName).First(&user).Error; result != nil {
		context.JSON(http.StatusNotFound, gin.H{"response": result.Error()})
		fmt.Println(result)
		context.Abort()
		return
	}
	//Create a new record for each deposit done by the user to the transactions table
	record := database.Instance.Create(&transactions)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Database Error": record.Error.Error()})
		context.Abort()
		log.Fatal(record.Error)
		return
	}

	if result := database.Instance.Table("users").Model(&models.User{}).Where("username = ?", depositDetails.UserName).Update("balance", *user.Balance+depositDetails.Amount); result.Error != nil {
		log.Fatal(result.Error)
		context.JSON(http.StatusNotAcceptable, gin.H{"Error": result.Error})
		context.Abort()
		fmt.Println("Cannot find User")
	}

	context.JSON(http.StatusCreated, gin.H{"response": *user.Balance})
}

func GetMtAccountBalance() {
	client := &http.Client{}

	/**
		    Fetch data from Mt4 api through nodejs sdk  provided.
	        Will switch to RabbitMq to make responses  fast.
	*/
	req, err := http.NewRequest("GET", "https://mt4functions.herokuapp.com/account", nil)

	if err != nil {
		fmt.Println(err)
	}

	//Set  headers to the requests
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	//Use the client to make the requests with the given [configurations]
	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
	}

	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Print(err.Error())
	}
	var mt4account models.Mt4Account

	json.Unmarshal(bodyBytes, &mt4account)

	fmt.Println(mt4account.Balance)

}
