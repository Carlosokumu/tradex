package main

import (
	"net/http"
	"os"

	"github.com/carlosokumu/dubbedapi/chat"
	"github.com/carlosokumu/dubbedapi/controllers"
	"github.com/carlosokumu/dubbedapi/database"
	"github.com/carlosokumu/dubbedapi/token"
	"github.com/gin-gonic/gin"
)

func main() {
	databaseUrl := os.Getenv("DATABASE_URL")
	database.Connect(databaseUrl)
	database.Migrate()
	router := initRouter()
	port := os.Getenv("PORT")

	router.Run(port)

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

	authRoutes := router.Group("/auth/user")
	// registration route
	authRoutes.POST("/register", controllers.RegisterUser)
	authRoutes.POST("/login", controllers.LoginUser)

	//Admin route
	adminRotes := router.Group("/admin")
	adminRotes.Use(token.JWTAuthAdmin())
	adminRotes.POST("/verify/trader", controllers.VerifyTrader)

	//public user level access routes
	publicUserRoutes := router.Group("/api/v1/user")
	publicUserRoutes.GET("/traders", controllers.GetTraders)

	//protected trader access level  routes
	protectedTraderRoutes := router.Group("/api/v1/trader")
	protectedTraderRoutes.Use(token.JWTAuthTrader())
	protectedTraderRoutes.PATCH("/connect", controllers.ConnectTradingAccount)

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
