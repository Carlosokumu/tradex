package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/carlosokumu/dubbedapi/models"
	"github.com/gin-gonic/gin"
)

func PositionData(context *gin.Context) {

	bodyBytes, err := ioutil.ReadAll(context.Request.Body)

	if err != nil {
		fmt.Print(err.Error())
	}
	var responseObject models.PositionData

	parseError := json.Unmarshal(bodyBytes, &responseObject)

	if parseError != nil {
		fmt.Println(parseError)
	}

	fmt.Println("Ctrader Sendind data")
	context.String(http.StatusOK, "Hellow")
}
