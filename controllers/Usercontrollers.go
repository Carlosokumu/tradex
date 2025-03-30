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
	"github.com/carlosokumu/dubbedapi/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// global variable
var (
	Trader models.User
)

func RegisterUser(context *gin.Context) {
	//var user models.User
	var userModel models.UserModel

	var userDto dtos.UserDto

	if err := context.ShouldBind(&userDto); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate user input
	if err := validateUserInput(&userDto); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username or email already exists
	if userExists := checkUserExists(userDto.UserName, userDto.Email); userExists {
		context.JSON(http.StatusConflict, gin.H{"error": "Provided username  already exists"})
		return
	}

	userModel = models.UserModel{
		UserName: userDto.UserName,
		Email:    userDto.Email,
		Password: userDto.Password,
		RoleID:   utils.USER,
	}
	if err := userModel.HashPassword(); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		context.Abort()
		log.Fatal(err)
		return
	}

	record := database.Instance.Create(&userModel)
	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": record.Error.Error()})
		context.Abort()
		return
	}

	token, tokenError := token.GenerateJWTWithUserModel(models.UserModel{
		UserName: userModel.UserName,
		Password: userModel.Password,
		RoleID:   utils.USER,
	})

	if tokenError != nil {
		fmt.Println("failed to generate token:", tokenError)
		context.JSON(http.StatusInternalServerError, gin.H{"error": tokenError})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"user": models.User{
		UserName: userModel.UserName,
		Email:    userModel.Email,
		Password: userModel.Password,
	}, "token": token},
	)
}

func LoginUser(context *gin.Context) {
	var userModel models.UserModel
	var UserLoginDto dtos.UserLoginDto

	if err := context.ShouldBind(&UserLoginDto); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if result := database.Instance.Table("user_models").Where("user_name = ?", UserLoginDto.UserName).Preload("TradingAccounts").Preload("Role").First(&userModel).Error; result != nil {
		context.JSON(http.StatusNotFound, gin.H{"response": result.Error()})
		return
	}
	_, err := bcrypt.GenerateFromPassword([]byte(UserLoginDto.Password), 14)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"response": err})
		return
	}
	result := CheckPasswordHash(UserLoginDto.Password, userModel.Password)

	token, tokenError := token.GenerateJWTWithUserModel(models.UserModel{
		UserName: userModel.UserName,
		Password: userModel.Password,
		RoleID:   userModel.RoleID,
	})

	if tokenError != nil {
		fmt.Println("failed to generate token:", tokenError)
		context.JSON(http.StatusInternalServerError, gin.H{"error": tokenError})
		return
	}

	if result {
		context.JSON(http.StatusOK, gin.H{"user": userModel, "token": token})
	} else {
		context.JSON(http.StatusUnauthorized, gin.H{"response": "password does not match username"})
	}
}

type PaginationData struct {
	NextPage     *int
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
		NextPage: func() *int {
			if true {
				nextPage := page + 1
				return &nextPage
			}
			return nil
		}(),
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

func GetTraders(context *gin.Context) {
	var users []models.UserModel
	var totalRows int64
	pageStr := context.Query("page")
	page, _ := strconv.Atoi(pageStr)
	PageSize := 5
	offset := (page - 1) * PageSize

	//Calculate total pages
	database.Instance.Table("user_models").Model(&models.UserModel{}).Count(&totalRows)
	totalPages := float64(totalRows / int64(PageSize))

	hasNextPage := (offset + PageSize) < int(totalPages)

	result := database.Instance.Table("user_models").Where("role_id = ?", utils.TRADER).Preload("TradingAccounts").Preload("Role").Limit(PageSize).Offset(offset).Find(&users)
	if result.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"server error": result.Error})
	}

	context.JSON(http.StatusOK, gin.H{"users": users, "pagination": PaginationData{
		CurrentPage: page,
		NextPage: func() *int {
			if hasNextPage {
				nextPage := page + 1
				return &nextPage
			}
			return nil
		}(),
		PreviousPage: func() *int {
			if page <= 1 {
				return nil
			}
			previouspage := page - 1
			return &previouspage
		}(),
		HasMore: hasNextPage,
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
	}
	context.JSON(http.StatusOK, gin.H{"response": "Phone Number updated Sucessfully"})
}

func ConnectTradingAccount(context *gin.Context) {

	var tradingAccountDTO dtos.TradingAccountDto
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

func checkUserExists(username, email string) bool {
	var userModel models.UserModel
	if err := database.Instance.Table("user_models").
		Where("user_name = ?", username).
		Or("email = ?", email).
		First(&userModel).
		Error; err != nil {
		return false
	} else {
		return true
	}
}

func validateUserInput(user *dtos.UserDto) error {
	if usernamelength := len(user.UserName); usernamelength < 4 {
		return errors.New("username should be more of 4 or more characters")
	}
	if emailformatcredibility, _ := emailmethods.Emailformatverifier(user.Email); !emailformatcredibility {
		return errors.New("please check that your email is correctly formatted")
	}
	if passwordlength := stringmethods.Charactercount(user.Password); passwordlength < 8 {
		return errors.New("a password should be of 8 or more characters")
	}
	return nil
}
