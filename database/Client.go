package database

import (
	"log"

	"github.com/carlosokumu/dubbedapi/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Instance *gorm.DB
var dbError error

func Connect(connectionString string) {

	Instance, dbError = gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if dbError != nil {
		log.Fatal(dbError)
		panic("Cannot connect to DB")
	}
	log.Println("Connected to Database!")
}

func Migrate() {

	err := Instance.AutoMigrate(
		&models.User{},
		&models.OpenPosition{},
	)

	if err != nil {
		log.Println(err)
	}
	log.Println("Database Migration Completed!")
}
