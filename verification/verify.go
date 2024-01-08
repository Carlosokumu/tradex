package verification

import (
	"net/http"

	"github.com/carlosokumu/dubbedapi/controllers"
	"github.com/carlosokumu/dubbedapi/token"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// verify token
func IsAuthorized(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(context *gin.Context) {

		if context.Request.Header["Token"] == nil {
			context.JSON(http.StatusInternalServerError, gin.H{"Error": "token not found"})
			return
		}

		var mySigningKey = []byte(token.Secretkey)

		token, err := jwt.Parse(context.Request.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, nil
			}
			return mySigningKey, nil
		})

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"Error": "token has been expired"})
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["username"] == controllers.Trader.UserName {
				context.Request.Header.Set("Username", controllers.Trader.UserName)
				handler(context)
				return

			}
		}
		context.JSON(http.StatusForbidden, gin.H{"response": "Not Authorized"})
	}
}

// verify user
func UserIndex(context *gin.Context) {
	if context.Request.Header.Get("Username") != controllers.Trader.UserName {
		context.JSON(http.StatusForbidden, gin.H{"response": "Not Authorized"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"response welcome": controllers.Trader.UserName})
}
