package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PositionData(context *gin.Context) {
	fmt.Println("Hitting by ctrader")
	context.String(http.StatusOK, "Hellow")
}
