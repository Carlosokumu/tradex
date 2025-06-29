package controllers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/carlosokumu/dubbedapi/database"
	"github.com/carlosokumu/dubbedapi/models"
	"github.com/carlosokumu/dubbedapi/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func VerifyTrader(ctx *gin.Context) {
	username := strings.TrimSpace(ctx.Query("username"))
	if username == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "username parameter is required",
		})
		return
	}

	tx := database.Instance.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Recovered from panic in VerifyTrader: %v", r)
		}
	}()

	var user models.UserModel
	if err := tx.Where("user_name = ?", username).First(&user).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error":   "user not found",
				"details": fmt.Sprintf("username: %s", username),
			})
			return
		}
		log.Printf("Database error fetching user: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to verify user",
		})
		return
	}

	if user.RoleID == utils.TRADER {
		tx.Rollback()
		ctx.JSON(http.StatusConflict, gin.H{
			"error": "user is already a trader",
		})
		return
	}

	if err := tx.Model(&user).Update("role_id", utils.TRADER).Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to update user role: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to verify user",
		})
		return
	}

	if err := tx.Commit().Error; err != nil {
		log.Printf("Transaction commit failed: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to complete verification",
		})
		return
	}

	log.Printf("User %s (ID: %d) verified as trader", user.UserName, user.ID)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User successfully verified as trader",
		"data": gin.H{
			"user_id":   user.ID,
			"user_name": user.UserName,
			"new_role":  "trader",
			"role_id":   utils.TRADER,
			"timestamp": time.Now().Format(time.RFC3339),
		},
	})
}
