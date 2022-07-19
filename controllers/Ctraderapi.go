package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PositionData(context *gin.Context) {
	fmt.Println("Ctrader Sendind data")
	context.String(http.StatusOK, "Hellow")
}
