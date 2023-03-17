package token

import (
	"fmt"
	"time"
	"github.com/golang-jwt/jwt"
)

//global variable
var Secretkey="Leonmjaluo"

//generate JWT token
func GenerateJWT(email,username string) (string, error) {
	var mySigningKey = []byte(Secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["email"] = email
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Println("something went wrong",err)
		return "", err
	}

	return tokenString, nil
}