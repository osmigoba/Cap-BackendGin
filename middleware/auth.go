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

func Auth(admin bool) gin.HandlerFunc {
	//(bool admin)
	//Extract the "Bearer Token"
	return func(cgin *gin.Context) {
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
		fmt.Printf("\nLos valores del claim %v %T", claims, claims)
		fmt.Println("\nclaim de admin:", claims["admin"])
		var claimAdmin bool
		claimAdmin = claims["admin"].(bool)
		fmt.Printf("\nTipo del claim Admin: %T", claimAdmin)

		if ok && token.Valid {
			fmt.Println("\n entre en el IF")

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
			if admin {
				//ADMIN USER
				if user.IsAdmin && claimAdmin {
					fmt.Println("Authorized  protected Route")
					cgin.Next()
				} else {
					fmt.Println(" Not Authorized, Protected Route")
					cgin.AbortWithStatus(http.StatusUnauthorized)
				}
			} else {
				// No Admin User
				fmt.Println("Authorized Not protected Route")
				cgin.Next()
			}

		} else {

			cgin.JSON(http.StatusUnauthorized, gin.H{"unauthorized": err.Error()})
			cgin.Abort()
		}
		//Continue with the execution

	}

}
