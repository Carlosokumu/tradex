package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/carlosokumu/dubbedapi/database"
	"github.com/carlosokumu/dubbedapi/models"
	"github.com/gin-gonic/gin"
)

func InsertPositionData(context *gin.Context) {
	bodyBytes, err := ioutil.ReadAll(context.Request.Body)

	if err != nil {
		fmt.Print(err.Error())
	}
	var responseObject models.OpenPosition

	parseError := json.Unmarshal(bodyBytes, &responseObject)

	if parseError != nil {
		fmt.Println(parseError)
	}

	record := database.Instance.Create(&responseObject)

	if record.Error != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"Database Error": record.Error.Error()})
		context.Abort()
		log.Fatal(record.Error)
		return
	}

}

func GetOpenPositions(context *gin.Context) {

	var openposition []models.OpenPosition
	var user models.User
	result := database.Instance.Table("open_positions").Find(&openposition)
	user.GetMtAccountBalance()

	if result.Error != nil {
		log.Fatal(result)
	}
	context.JSON(http.StatusOK, gin.H{"openpositions": openposition})

}
