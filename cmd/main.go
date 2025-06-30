package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/carlosokumu/dubbedapi/chat"
	"github.com/carlosokumu/dubbedapi/config"
	"github.com/carlosokumu/dubbedapi/controllers"
	"github.com/carlosokumu/dubbedapi/database"
	"github.com/carlosokumu/dubbedapi/token"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()
	migrate := flag.Bool("migrate", false, "Run database migrations")
	createSuperUser := flag.Bool("create-superuser", false, "Create a superuser account")
	flag.Parse()

	database.Connect(config.DbUrl)

	if *migrate {
		log.Println("Running database migrations...")
		database.Migrate()
	}

	if *createSuperUser {
		log.Println("Creating superuser account...")
		_, err := database.CreateSuperUser()
		if err != nil {
			log.Fatal("Failed to create superuser:", err)
		}
		log.Println("Superuser created successfully")
	}
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

	//[Websocket] Endpoindts ------
	router.GET("/ws", func(c *gin.Context) {
		chat.ServeWs(hub, c.Writer, c.Request)
	})

	router.GET("/ws/bot", func(c *gin.Context) {
		controllers.ReadBotEndpoint(c.Writer, c.Request)
	})

	authRoutes := router.Group("api/v1/auth/user")
	// registration route
	authRoutes.POST("/register", controllers.RegisterUser)
	authRoutes.POST("/login", controllers.LoginUser)

	//Admin route
	adminRotes := router.Group("api/v1/admin")
	adminRotes.Use(token.JWTAuthAdmin())
	adminRotes.POST("/verify/trader", controllers.VerifyTrader)

	//public user level access routes
	publicUserRoutes := router.Group("/api/v1/user")
	publicUserRoutes.GET("/traders", controllers.GetTraders)
	publicUserRoutes.PATCH("/community/join", controllers.AddNewMemberToCommunity)
	publicUserRoutes.GET("/communities", controllers.GetAllCommunities)
	publicUserRoutes.GET("/community", controllers.GetCommunityByName)

	//protected trader access level  routes
	protectedTraderRoutes := router.Group("/api/v1/trader")
	protectedTraderRoutes.Use(token.JWTAuthTrader())
	protectedTraderRoutes.GET("/connect", controllers.ConnectTradingAccount)
	protectedTraderRoutes.POST("/community/create", controllers.CreateCommunity)
	protectedTraderRoutes.POST("/community/post", controllers.PostToCommunity)

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
