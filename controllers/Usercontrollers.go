package controllers

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/ajg/form"
	"github.com/carlosokumu/dubbedapi/database"
	"github.com/carlosokumu/dubbedapi/dtos"
	"github.com/carlosokumu/dubbedapi/emailmethods"
	"github.com/carlosokumu/dubbedapi/models"
	"github.com/carlosokumu/dubbedapi/stringmethods"
	"github.com/carlosokumu/dubbedapi/token"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// global variable
var (
	Trader models.User
)

func RegisterUser(context *gin.Context) {
	var user models.User
	var userModel models.UserModel
	d := form.NewDecoder(context.Request.Body)
	if err := d.Decode(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Parse Error": err.Error()})
		log.Fatal(err)
		context.Abort()
		return
	}
	username := user.UserName
	email := user.Email
	password := user.Password
	if err := database.Instance.Table("user_models").Where("user_name = ?", username).Or("email = ? ", email).First(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if usernamelength := stringmethods.Charactercount(username); usernamelength < 4 {
				context.JSON(http.StatusBadRequest, gin.H{"username error": "username should be more of 4 or more characters"})
			} else if emailformatcredibility, _ := emailmethods.Emailformatverifier(email); !emailformatcredibility {
				context.JSON(http.StatusBadRequest, gin.H{"email error": "email address is invalid"})
			} else if passwordlength := stringmethods.Charactercount(password); passwordlength < 8 {
				context.JSON(http.StatusBadRequest, gin.H{"password error": "password is weak or invalid"})
			} else {

				if err := user.HashPassword(); err != nil {
					context.JSON(http.StatusInternalServerError, gin.H{"Internal server error": err.Error()})
					context.Abort()
					log.Fatal(err)
					return
				}

				//Create a new userModel entity
				userModel = models.UserModel{
					UserName: user.UserName,
					Email:    user.Email,
					Password: user.Password,
				}

				record := database.Instance.Create(&userModel)
				if record.Error != nil {
					context.JSON(http.StatusInternalServerError, gin.H{"Database Error": record.Error.Error()})
					context.Abort()
					return
				}
				_, tokenError := token.GenerateJWT(user.Email, user.UserName)
				if tokenError != nil {
					fmt.Println("failed to generate token:", tokenError)
					context.JSON(http.StatusInternalServerError, gin.H{"Token generation Error": tokenError})
					return
				}

				context.JSON(http.StatusCreated, gin.H{"user": models.User{
					UserName: user.UserName,
					Email:    user.Email,
					Password: user.Password,
				}},
				)
			}
		} else {
			context.JSON(http.StatusInternalServerError, gin.H{"error": err})
			context.Abort()
			return
		}
	} else {
		context.JSON(http.StatusConflict, gin.H{"error": "Provided username or email already exists"})
		context.Abort()
		return
	}
}

func LoginUser(context *gin.Context) {
	var credentials models.Credentials
	var userModel models.UserModel

	d := form.NewDecoder(context.Request.Body)
	if err := d.Decode(&credentials); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Parse Error": err.Error()})
		log.Fatal(err)
		context.Abort()
		return
	}
	if result := database.Instance.Table("user_models").Where("user_name = ?", credentials.UserName).First(&userModel).Error; result != nil {
		context.JSON(http.StatusNotFound, gin.H{"response": result.Error()})
		fmt.Println(result)
		context.Abort()
		return
	}
	_, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), 14)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"response": err})
		context.Abort()
		return
	}
	result := CheckPasswordHash(credentials.Password, userModel.Password)

	if result {
		context.JSON(http.StatusOK, gin.H{"user": userModel})
	} else {
		context.JSON(http.StatusUnauthorized, gin.H{"response": "password does not match username"})
	}
}

type PaginationData struct {
	NextPage     int
	PreviousPage *int
	CurrentPage  int
	HasMore      bool
}

func GetAllUsers(context *gin.Context) {
	var users []models.UserModel
	var totalRows int64
	pageStr := context.Query("page")
	page, _ := strconv.Atoi(pageStr)
	PageSize := 5
	fmt.Println("Page:", page)
	offset := (page - 1) * PageSize

	//Calculate total pages
	database.Instance.Table("user_models").Model(&models.UserModel{}).Count(&totalRows)
	totalPages := float64(totalRows / int64(PageSize))

	result := database.Instance.Table("user_models").Preload("TradingAccounts").Limit(PageSize).Offset(offset).Find(&users)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"server error": result.Error})
	}

	context.JSON(http.StatusOK, gin.H{"users": users, "pagination": PaginationData{
		CurrentPage: page,
		NextPage:    page + 1,
		PreviousPage: func() *int {
			if page <= 1 {
				return nil
			}
			previouspage := page - 1
			return &previouspage
		}(),
		HasMore: IsLastPage(page, int(totalPages), PageSize),
	},
	})
}

func IsLastPage(currentPage, totalRecords, pageSize int) bool {
	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))
	return currentPage == totalPages
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

func UpdateTradingAccount(context *gin.Context) {

	var tradingAccountDTO dtos.TradingAccountDTO
	var userModel models.UserModel

	// Bind JSON data from the request body to DTO
	if err := context.ShouldBindJSON(&tradingAccountDTO); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Please check your fields"})
		return
	}

	if result := database.Instance.Table("user_models").Where("user_name = ?", tradingAccountDTO.UserName).First(&userModel).Error; result != nil {
		context.JSON(http.StatusNotFound, gin.H{"response": result.Error()})
		fmt.Println(result)
		context.Abort()
		return
	} else {
		context.JSON(http.StatusOK, gin.H{"user": userModel})
	}

	// Convert DTO to GORM entity
	tradingAccount := models.TradingAccount{
		Platform:  tradingAccountDTO.Platform,
		AccountId: tradingAccountDTO.AccountId,
		UserId:    userModel.ID,
	}

	record := database.Instance.Create(&tradingAccount)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Database Error": record.Error.Error()})
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, gin.H{"trading_account": tradingAccount})
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

func EmailPassword(context *gin.Context) {
	var user models.EmailPassword
	fmt.Println("REQUESTURL:", context.Request.URL)
	d := form.NewDecoder(context.Request.Body)

	if err := d.Decode(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Parse Error": err.Error()})
		log.Fatal(err)
		context.Abort()
		return
	}
	email := user.Email
	password := user.Password

	if emailformatcredibility, _ := emailmethods.Emailformatverifier(email); !emailformatcredibility {
		context.JSON(http.StatusBadRequest, gin.H{"email error": "email address is invalid"})
	} else if passwordlength := stringmethods.Charactercount(password); passwordlength < 8 {
		context.JSON(http.StatusBadRequest, gin.H{"password error": "password is weak or invalid"})
	} else {
		context.JSON(http.StatusCreated, gin.H{"user": models.EmailPassword{
			Email:    user.Email,
			Password: user.Password,
		}},
		)
	}
}

// [Update]- refactor name  - store the refresh token,access token,client_id and secret
func Access_refresh_token_accout_id_secret(context *gin.Context) {
	var user models.AccessRefreshaccountsecret
	var userModel models.UserModel
	fmt.Println("REQUESTURL:", context.Request.URL)

	queryParams := context.Request.URL.Query()
	username := queryParams.Get("username")

	d := form.NewDecoder(context.Request.Body)

	if err := d.Decode(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Parse Error": err.Error()})
		log.Fatal(err)
		context.Abort()
		return
	}
	access := user.AccessToken
	refresh := user.RefreshToken
	client_id := user.Client_id
	secret := user.Secret

	if length := stringmethods.Charactercount(access); length < 10 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "empty access token"})
	} else if length := stringmethods.Charactercount(refresh); length < 10 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "empty refresh token"})
	} else if length := stringmethods.Charactercount(client_id); length < 1 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "empty account_id"})
	} else if length := stringmethods.Charactercount(secret); length < 10 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "empty secret"})
	} else {
		if result := database.Instance.Table("user_models").Model(&userModel).Select("AccessToken", "Secret", "RefreshToken").Where("user_name = ?", username).Updates(map[string]interface{}{"AccessToken": access, "Secret": secret, "RefreshToken": refresh}); result.Error != nil {
			log.Fatal(result.Error)
			fmt.Println("Cannot find User")
		}

		context.JSON(http.StatusCreated, gin.H{"user": userModel})
	}
}
func GetSpecificUser(context *gin.Context) {

	var userModel models.UserModel
	queryParams := context.Request.URL.Query()
	username := queryParams.Get("username")

	fmt.Println(username)
	if result := database.Instance.Table("user_models").Where("user_name = ?", username).First(&userModel).Error; result != nil {
		context.JSON(http.StatusNotFound, gin.H{"response": result.Error()})
		fmt.Println(result)
		context.Abort()
		return
	} else {
		context.JSON(http.StatusOK, gin.H{"user": userModel})
	}

}

func DeleteSpecificUser(context *gin.Context) {

	var userModel models.UserModel
	queryParams := context.Request.URL.Query()
	username := queryParams.Get("username")

	if result := database.Instance.Table("user_models").Where("user_name = ?", username).Delete(&userModel).Error; result != nil {
		context.JSON(http.StatusNotFound, gin.H{"response": result.Error()})
		fmt.Println(result)
		context.Abort()
		return
	} else {
		context.JSON(http.StatusOK, gin.H{"user": userModel})
	}

}
