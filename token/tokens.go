package token

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/carlosokumu/dubbedapi/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type Token struct {
	User        models.User `json:"user"`
	TokenString string      `json:"token"`
}

// extract token from request Authorization header
func getTokenFromRequest(context *gin.Context) string {
	bearerToken := context.Request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) == 2 {
		return splitToken[1]
	}
	return ""
}

// check token validity
func getToken(context *gin.Context) (*jwt.Token, error) {
	tokenString := getTokenFromRequest(context)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return Secretkey, nil
	})
	return token, err
}

// validate JWT token
func ValidateJWT(context *gin.Context) error {
	token, err := getToken(context)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return nil
	}
	return errors.New("invalid token provided")
}

// validate Trader JWT token
func ValidateTraderRoleJWT(context *gin.Context) error {
	token, err := getToken(context)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	userRole := uint(claims["role"].(float64))
	if ok && token.Valid && userRole == 2 || userRole == 1 {
		return nil
	}
	return errors.New("invalid author token provided")
}

// check for a valid trader token
func JWTAuthTrader() gin.HandlerFunc {
	return func(context *gin.Context) {
		err := ValidateJWT(context)
		if err != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			context.Abort()
			return
		}
		error := ValidateTraderRoleJWT(context)
		if error != nil {
			context.JSON(http.StatusUnauthorized, gin.H{"error": "Only verified  traders are allowed to perform this action"})
			context.Abort()
			return
		}
		context.Next()
	}
}
