package middleware

import (
	"ApiGin/initializers"
	"ApiGin/models"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func Auth(cgin *gin.Context) {
	//Extract the "Bearer Token"
	authorizationHeader := cgin.GetHeader("Authorization")
	bearerToken := strings.Split(authorizationHeader, " ")

	if len(bearerToken) < 2 {
		cgin.AbortWithStatus(http.StatusUnauthorized)

		return
	}
	tokenString := bearerToken[1]
	if tokenString == "" {
		cgin.AbortWithStatus(http.StatusUnauthorized)
		cgin.Abort()
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		//Check the expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			//Token expired
			cgin.JSON(http.StatusUnauthorized, gin.H{"unauthorized": "Token Expired"})
			cgin.Abort()
		}

		//Check the user in the database.
		var user models.User
		initializers.DB.First(&user, claims["sub"])
		if user.Id == 0 {
			cgin.AbortWithStatus(http.StatusUnauthorized)
		}

	} else {

		cgin.JSON(http.StatusUnauthorized, gin.H{"unauthorized": err.Error()})
		cgin.Abort()
	}
	//Continue with the execution
	cgin.Next()
}
