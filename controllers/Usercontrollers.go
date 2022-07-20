package controllers

import (
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
	context.JSON(http.StatusCreated, gin.H{"userId": user.ID, "firstname": user.FirstName, "lastName": user.LastName, "email": user.Email, "username": user.Username})
}

func UpdateUser(context *gin.Context) {
	//var user models.User
	// if result := database.Instance.Table("users").Where("username = ?", "carlos").First(&user); result.Error != nil {
	// 	log.Fatal(result)
	// }
	// user.Email = "coderokush@gmail.com"
	// database.Instance.Table("users").Save(&user)
	if result := database.Instance.Table("users").Model(&models.User{}).Where("username = ?", "carlo").Update("username", "webman"); result.Error != nil {
		log.Fatal(result.Error)
		// fmt.Println(user.LastName)
	}
	context.JSON(http.StatusOK, "Done")
}

func TestRouter(context *gin.Context) {
	context.String(http.StatusOK, "Hellow")
}
