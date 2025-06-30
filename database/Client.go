package database

import (
	"errors"
	"fmt"
	"log"

	"github.com/carlosokumu/dubbedapi/config"
	"github.com/carlosokumu/dubbedapi/dtos"
	"github.com/carlosokumu/dubbedapi/models"
	"github.com/carlosokumu/dubbedapi/token"
	"github.com/carlosokumu/dubbedapi/utils"
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
		}
	}

	fmt.Println("Successfully initialized roles")

}

func CreateSuperUser() (string, error) {

	userdto := dtos.UserDto{
		UserName: config.SuperUsername,
		Email:    config.SuperUserEmail,
		Password: config.SuperUserPassword,
	}

	err := utils.ValidateUserInput(&userdto)
	if err != nil {
		log.Printf("Failed to validate super user details: %v", err)
		return "", err
	}

	tx := Instance.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var existingUser models.UserModel
	if err := tx.Where("user_name = ? OR role_id = ?", userdto.UserName, utils.ADMIN).First(&existingUser).Error; err == nil {
		tx.Rollback()
		return "", fmt.Errorf("superuser already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return "", fmt.Errorf("database error: %v", err)
	}

	hashedPassword, err := utils.HashPassword(userdto.Password)
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to hash password: %v", err)
	}

	// Generate JWT
	token, err := token.GenerateJWTWithUserModel(models.UserModel{
		UserName: userdto.UserName,
		Email:    userdto.Email,
		Password: string(hashedPassword),
		RoleID:   utils.ADMIN,
	})
	if err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	// Create super user
	superUser := models.UserModel{
		UserName: userdto.UserName,
		Email:    userdto.Email,
		Password: string(hashedPassword),
		RoleID:   utils.ADMIN,
		Token:    token,
	}

	if err := tx.Create(&superUser).Error; err != nil {
		tx.Rollback()
		return "", fmt.Errorf("failed to create user: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return "", fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Superuser account created for %s", userdto.UserName)

	return token, nil
}
