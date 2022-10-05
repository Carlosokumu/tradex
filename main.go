package main

import (
	"fmt"
	"os"

	"github.com/carlosokumu/dubbedapi/controllers"
	"github.com/carlosokumu/dubbedapi/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	fmt.Println("Fetching data  from api")

	//Connect to Postgres and migrate for the schemas
	databaseUrl := os.Getenv("DATABASE_URL")
	database.Connect(databaseUrl)
	database.Migrate()

	
	router := initRouter()
	port := os.Getenv("PORT")
	router.Run(":" + port)

	

}

func initRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("html/*")
	router.GET("/rascamps", func(c *gin.Context) {
		c.HTML(http.StatusOK, "rascampsprivacy.html", nil)
	})

	api := router.Group("/tradex")
	{

		api.POST("/user/register", controllers.RegisterUser)
		api.POST("/positiondata/add", controllers.InsertPositionData)
		api.PATCH("/user", controllers.UpdateUser)
		api.PATCH("/user/phonenumber", controllers.UpdatePhoneNumber)
		api.GET("/positions/all", controllers.GetOpenPositions)
		api.POST("/user/login", controllers.LoginUser)
		api.POST("/user/email", controllers.SendOtp)
		api.POST("/user/confirmation", controllers.SendConfirmEmail)
		api.POST("/user/deposit", controllers.HandleDeposit)
		api.GET("/user/userinfo", controllers.GetUserInfo)
		
	}
	return router
}
