package main

import (
	"fmt"

	"net/http"

	"github.com/carlosokumu/dubbedapi/chat"
	"github.com/carlosokumu/dubbedapi/controllers"
	"github.com/carlosokumu/dubbedapi/database"
	"github.com/carlosokumu/dubbedapi/verification"
	"github.com/gin-gonic/gin"
)

func main() {

	fmt.Println("Fetching data  from api")

	//Connect to Postgres and migrate for the schemas
	//databaseUrl := os.Getenv("DATABASE_URL")
	databaseUrl := "postgres://swinngdata_user:4nZcOypBKc8E6RU96BftsnBMgClMGxqn@dpg-ci805l98g3n3vm2k9ifg-a.oregon-postgres.render.com/swinngdata"
	database.Connect(databaseUrl)
	database.Migrate()

	router := initRouter()
	//port := os.Getenv("PORT")
	port := "8080"
	router.Run(":" + port)

}

func initRouter() *gin.Engine {

	whitelist := make(map[string]bool)
	whitelist["0.0.0.0"] = true
	whitelist["https://swingwizards.vercel.app"] = true

	hub := chat.NewHub()
	go hub.Run()
	router := gin.Default()

	router.Use(allowAllDomains())
	router.LoadHTMLGlob("html/*")

	router.GET("/rascamps", func(c *gin.Context) {
		c.HTML(http.StatusOK, "rascampsprivacy.html", nil)
	})

	//[Websocket] Endpoindts ------
	router.GET("/ws", func(c *gin.Context) {
		chat.ServeWs(hub, c.Writer, c.Request)
	})

	router.GET("/ws/bot", func(c *gin.Context) {
		controllers.ReadBotEndpoint(c.Writer, c.Request)
	})

	//----------

	api := router.Group("/tradex")
	{

		api.POST("/user/register", controllers.RegisterUser)
		api.POST("/positiondata/add", controllers.InsertPositionData)
		api.PATCH("/user/phonenumber", controllers.UpdatePhoneNumber)
		api.GET("/positions/all", controllers.GetOpenPositions)
		api.POST("/user/login", controllers.LoginUser)
		api.POST("/user/email", controllers.SendOtp)
		api.POST("/user/confirmation", controllers.SendConfirmEmail)
		api.POST("/user/deposit", controllers.HandleDeposit)
		api.GET("/user/userinfo", controllers.GetUserInfo)
		api.GET("/user/verifytoken", verification.IsAuthorized(verification.UserIndex))
	}
	return router
}

func IPWhiteList(whitelist map[string]bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !whitelist[c.ClientIP()] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  http.StatusForbidden,
				"message": "Permission denied",
			})
			return
		}
	}
}

func allowAllDomains() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set the appropriate CORS headers to allow access from any domain
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type")

		// Handle the request
		c.Next()
	}
}
