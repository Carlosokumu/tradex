package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/carlosokumu/dubbedapi/database"
	"github.com/carlosokumu/dubbedapi/models"
	"github.com/gin-gonic/gin"
)

func VerifyTrader(context *gin.Context) {

	id, _ := strconv.Atoi(context.Param("id"))

	fmt.Println("Id:", id)

	if result := database.Instance.Table("user_models").Model(&models.UserModel{}).Where("user_name = ?", "carlos").Update("role_id", 2); result.Error != nil {
		log.Fatal(result.Error)
		fmt.Println("Cannot find User")
	}

	context.JSON(http.StatusOK, gin.H{"response": "Congratulations.You have been verified as a trader"})
}
