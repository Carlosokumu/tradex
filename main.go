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

	
	path := "postgres://lmbdoeimaxunbj:5aaf1d7c9fec6e330ec64a15059c8b7f66db09c2711e87cc0d56a4cffec03b0d@ec2-44-206-11-200.compute-1.amazonaws.com:5432/d99joi1q9b5r3t"
	database.Connect(path)
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
