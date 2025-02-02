package authentication

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"main/models"
	"net/http"
	"os"
	"time"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("Authorization")

		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("SECRET")), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			expDate, _ := claims["exp"].(float64)
			if float64(time.Now().Unix()) > expDate {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			var currentUser models.User
			db.First(&currentUser, claims["sub"])

			if currentUser.ID == 0 {
				c.AbortWithStatus(http.StatusUnauthorized)
			}
			userID, _ := claims["sub"].(float64)
			c.Set("userID", uint(userID))
			c.Next()
		} else {
			fmt.Println(err)
		}
	}
}
