package token

import (
	"fmt"
	"time"

	"github.com/carlosokumu/dubbedapi/models"
	"github.com/golang-jwt/jwt"
)

// global variable
var Secretkey = []byte("W7-67E0Zoi5V5RY2DYP6v1DM8Lp9JshVnw_B20EM7YI=")

const TOKEN_TTL = 1800

// generate JWT token
func GenerateJWT(email, username string) (string, error) {
	var mySigningKey = []byte(Secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = email
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println("something went wrong", err)
		return "", err
	}

	return tokenString, nil
}

func GenerateJWTWithUserModel(user models.UserModel) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID,
		"role": user.RoleID,
		"iat":  time.Now().Unix(),
		"eat":  time.Now().Add(time.Second * time.Duration(TOKEN_TTL)).Unix(),
	})
	return token.SignedString(Secretkey)
}
