package controllers

import (
	"log"
	"net/http"

	"github.com/carlosokumu/dubbedapi/database"
	"github.com/carlosokumu/dubbedapi/models"
	"github.com/carlosokumu/dubbedapi/utils"
	"github.com/gin-gonic/gin"
)

func VerifyTrader(context *gin.Context) {
	username := context.Query("username")
	//TODO : Add check for username validity
	if result := database.Instance.Table("user_models").Model(&models.UserModel{}).Where("user_name = ?", username).Update("role_id", utils.TRADER); result.Error != nil {
		log.Fatal(result.Error)
	}

	context.JSON(http.StatusOK, gin.H{"response": "Congratulations.You have been verified as a trader"})
}
