package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DbUrl string

	SuperUsername     string
	SuperUserEmail    string
	SuperUserPassword string

	CloudinaryCloudName string
	CloudinaryAPIKey    string
	CloudinaryAPISecret string
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file", err)
	}

	DbUrl = os.Getenv("DATABASE_URL")

	//Super user
	SuperUsername = os.Getenv("SUPER_USER_USERNAME")
	SuperUserPassword = os.Getenv("SUPER_USER_PASSWORD")
	SuperUserEmail = os.Getenv("SUPER_USER_EMAIL")

	//Cloudinary
	CloudinaryCloudName = os.Getenv("CLOUDINARY_CLOUD_NAME")
	CloudinaryAPIKey = os.Getenv("CLOUDINARY_API_KEY")
	CloudinaryAPISecret = os.Getenv("CLOUDINARY_API_SECRET")

	if CloudinaryCloudName == "" || CloudinaryAPIKey == "" || CloudinaryAPISecret == "" || SuperUserEmail == "" || SuperUserPassword == "" || SuperUsername == "" || DbUrl == "" {
		log.Fatal("Missing required environment variables")
	}
}
