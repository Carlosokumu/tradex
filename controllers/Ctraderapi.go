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

func PositionData(context *gin.Context) {

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

	fmt.Println("Ctrader Sendind data:Entry Price ", responseObject.EntryPrice)
	fmt.Println("Ctrader Sendind data:Position Id", responseObject.PositionId)
	fmt.Println("Ctrader Sendind data:TradeType", responseObject.TradeType)
	fmt.Println("Ctrader Sendind data:EntryTime", responseObject.EntryTime)
	fmt.Println("Ctrader Sendind data:Quantity", responseObject.Quantity)

	context.String(http.StatusOK, "Hellow")
}
