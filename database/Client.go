package database

import (
	"fmt"
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
	}
	log.Println("Connected to Database!")
}

func Migrate() {

	err := Instance.AutoMigrate(
		&models.UserModel{},
		&models.TradingAccount{},
		&models.Role{},
		&models.Community{},
		&models.Post{},
		&models.Comment{},
		&models.Like{},
		&models.CommunityMember{},
	)

	if err != nil {
		log.Fatal(err)
	}
	SeedRoles()
	log.Println("Database Migration Completed!")
}

func SeedRoles() {

	//User Role
	userRole := models.Role{
		Name: "TradeShareUser",
	}
	// Trader Role
	traderRole := models.Role{
		Name: "TradeShareTrader",
	}
	//Admin Role
	adminRole := models.Role{
		Name: "TraderShareAdmin",
	}

	roles := []models.Role{userRole, traderRole, adminRole}

	for _, role := range roles {
		record := Instance.Create(&role)
		if record.Error != nil {
			log.Fatal("Failed to seed roles:", record.Error)
			return
		}
	}

	fmt.Println("Successfully initialized roles")

}
