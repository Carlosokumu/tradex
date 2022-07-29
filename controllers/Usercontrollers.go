package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ajg/form"
	"github.com/carlosokumu/dubbedapi/database"
	"github.com/carlosokumu/dubbedapi/models"
	"github.com/gin-gonic/gin"
)

func RegisterUser(context *gin.Context) {
	var user models.User

	d := form.NewDecoder(context.Request.Body)
	if err := d.Decode(&user); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Parse Error": err.Error()})
		log.Fatal(err)
		context.Abort()
		return
	}

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
	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "firstname": user.FirstName, "lastname": user.LastName, "email": user.Email, "username": user.Username})
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
		context.JSON(http.StatusNotFound, gin.H{"Http404": result.Error()})
		fmt.Println(result)
		context.Abort()
		return
	}
	context.JSON(http.StatusOK, gin.H{"response": user.FirstName})
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

func TestRouter(context *gin.Context) {
	context.String(http.StatusOK, "Hellow")
}
