package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"main/constants"
	"main/models"
	"main/services/user"
	"net/http"
	"os"
	"strconv"
	"time"
)

func SignUpHandlerGin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var locationAux *string
		username := c.PostForm("username")
		password := c.PostForm("password")
		mail := c.PostForm("mail")
		location := c.PostForm("location")

		if location != constants.EMPTY {
			locationAux = &location
		}

		if username == constants.EMPTY || password == constants.EMPTY || mail == constants.EMPTY {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		if !user.IsEmail(mail) {
			c.JSON(http.StatusOK, gin.H{"error": "Invalid email"})
			return
		}

		if user.MailAlreadyUsed(db, mail) {
			c.JSON(http.StatusOK, gin.H{"error": "Email already in use"})
			return
		}

		if user.UsernameAlreadyUsed(db, username) {
			c.JSON(http.StatusOK, gin.H{"error": "Username already in use"})
			return
		}

		if createErr := user.CreateAccount(db, username, string(hashedPassword), mail, locationAux); createErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid parameters to create an account"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Account created successfully"})
	}
}

func LoginHandlerGin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		usernameOrEmail := c.PostForm("username-or-email")
		password := c.PostForm("password")

		if usernameOrEmail == constants.EMPTY || password == constants.EMPTY {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		var u models.User
		if err := db.Where("username = ? OR mail = ?", usernameOrEmail, usernameOrEmail).First(&u).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		secretKey := os.Getenv("SECRET")
		if secretKey == constants.EMPTY {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": u.ID,
			"exp": time.Now().Add(time.Hour * constants.EXPDATE).Unix(),
		})

		tokenString, err := token.SignedString([]byte(secretKey))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, constants.MAXCOOKIEAGE, "/", "", false, true)

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}

func FollowUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		followingID, getIDErr := getUserID(c)
		if getIDErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		followedUserID, atoiErr := strconv.Atoi(c.Param("userid"))
		if atoiErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		if followErr := user.FollowAccount(db, followingID, uint(followedUserID)); followErr != nil {
			log.Println("Follow error:", followErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Followed user successfully"})
	}
}

func UnfollowUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		followingID, err := getUserID(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		followedUserID, atoiErr := strconv.Atoi(c.Param("userid"))
		if atoiErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		if unfollowErr := user.UnfollowAccount(db, followingID, uint(followedUserID)); unfollowErr != nil {
			log.Println("Unfollow error:", unfollowErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Unfollowed user successfully"})
	}
}

func getUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("user ID not found")
	}

	if userIDUint, ok := userID.(uint); ok {
		return userIDUint, nil
	}

	return 0, errors.New("invalid user ID format")
}

var FollowUserEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "follow/:userid",
	HandlerFunction: FollowUserHandler,
}

var UnfollowUserEndpoint = models.Endpoint{
	Method:          models.DELETE,
	Path:            constants.BASEURL + "unfollow/:userid",
	HandlerFunction: UnfollowUserHandler,
}

var UserSignUpEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "signup",
	HandlerFunction: SignUpHandlerGin,
}

var UserLoginEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "login",
	HandlerFunction: LoginHandlerGin,
}
