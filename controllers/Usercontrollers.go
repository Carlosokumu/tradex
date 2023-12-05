package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/ajg/form"
	"github.com/carlosokumu/dubbedapi/database"
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
	username := user.Username
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
					UserName: user.Username,
					Email:    user.Email,
					Password: user.Password,
				}

				record := database.Instance.Create(&userModel)
				if record.Error != nil {
					context.JSON(http.StatusInternalServerError, gin.H{"Database Error": record.Error.Error()})
					context.Abort()
					return
				}

				_, tokenError := token.GenerateJWT(user.Email, user.Username)
				if tokenError != nil {
					fmt.Println("failed to generate token:", tokenError)
					context.JSON(http.StatusInternalServerError, gin.H{"Token generation Error": tokenError})
					return
				}

				context.JSON(http.StatusCreated, gin.H{"user": models.User{
					Username: user.Username,
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
		context.JSON(http.StatusOK, gin.H{"error": "Provided username or email already exists"})
		context.Abort()
		return
	}
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
	_, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), 14)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"response": err})
		context.Abort()
		return
	}
	result := CheckPasswordHash(credentials.Password, user.Password)

	if result {
		context.JSON(http.StatusOK, gin.H{"user": user})
	} else {
		context.JSON(http.StatusUnauthorized, gin.H{"response": "password does not match username"})
	}

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

func GetUserInfo(context *gin.Context) {
	//var user models.User
	var stagedUser models.User
	username := context.Query("username")

	// mt4Account, err := user.GetMtAccountBalance()

	// if err != nil {
	// 	return
	// }
	//floatingprofit := mt4Account.Equity - mt4Account.Balance
	// if err := database.Instance.Where("username = ?", username).First(&user).Error; err != nil {
	// 	fmt.Println(err)
	// 	context.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
	// 	context.Abort()
	// 	return
	// }
	// //divide the main accounts data into  data for an individual user
	// individualprofit := (*user.PercentageContribution / 100) * floatingprofit
	// individualaccountBalance := (*user.PercentageContribution / 100) * mt4Account.Balance
	// individualEquity := (*user.PercentageContribution / 100) * mt4Account.Equity

	// //Update the data in an individual user before feeding it to the user
	// if result := database.Instance.Table("users").Model(&models.User{}).Where("username = ?", username).Updates(models.User{
	// 	FloatingProfit: &individualprofit,
	// 	Balance:        &individualaccountBalance,
	// 	Equity:         &individualEquity,
	// }); result.Error != nil {
	// 	context.JSON(http.StatusNotAcceptable, gin.H{"Error": result.Error})
	// 	context.Abort()
	// 	fmt.Println("Cannot find User")
	// 	return
	// }
	if err := database.Instance.Where("username = ?", username).Preload("Positions").First(&stagedUser).Error; err != nil {
		fmt.Println(err)
		context.JSON(http.StatusNotFound, gin.H{"Error": err.Error()})
		context.Abort()
		return
	}

	context.JSON(http.StatusOK, gin.H{"user": stagedUser})
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
func Access_refresh_token_accout_id_secret(context *gin.Context) {
	var user models.AccessRefreshaccountsecret
	fmt.Println("REQUESTURL:", context.Request.URL)

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
		context.JSON(http.StatusCreated, gin.H{"user": models.AccessRefreshaccountsecret{
			AccessToken:  user.AccessToken,
			RefreshToken: user.RefreshToken,
			Client_id:    user.Client_id,
			Secret:       user.Secret,
		}},
		)
	}
}
