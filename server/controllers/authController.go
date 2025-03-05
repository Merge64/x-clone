package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"main/constants"
	"main/models"
	"main/services/user"
	"net/http"
	"os"
	"time"
)

func SignUpHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Nickname string `json:"nickname" binding:"required"`
			Username string `json:"username" binding:"required"`
			Password string `json:"password" binding:"required"`
			Mail     string `json:"mail" binding:"required"`
			Location string `json:"location"` // optional field
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		var locationAux *string
		if req.Location != constants.Empty {
			locationAux = &req.Location
		}

		if req.Username == constants.Empty || req.Password == constants.Empty || req.Mail == constants.Empty {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		if !user.IsEmail(req.Mail) {
			c.JSON(http.StatusOK, gin.H{"error": "Invalid email"})
			return
		}

		if user.MailAlreadyUsed(db, req.Mail) {
			c.JSON(http.StatusOK, gin.H{"error": "Email already in use"})
			return
		}

		if user.UsernameAlreadyUsed(db, req.Username) {
			c.JSON(http.StatusOK, gin.H{"error": "Username already in use"})
			return
		}

		// Create the user account
		newUser := models.User{
			Nickname: req.Nickname,
			Username: req.Username,
			Password: string(hashedPassword),
			Mail:     req.Mail,
			Location: locationAux,
		}

		if err := db.Create(&newUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
			return
		}

		// Generate JWT token
		secretKey := os.Getenv("SECRET")
		if secretKey == constants.Empty {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": newUser.ID,
			"exp": time.Now().Add(time.Hour * constants.ExpDate).Unix(),
		})

		tokenString, err := token.SignedString([]byte(secretKey))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, constants.MaxCookieAge, "/", constants.Empty, false, true)

		// Return the token in the JSON response
		c.JSON(http.StatusCreated, gin.H{
			"message": "Account created successfully",
			"token":   tokenString,
		})
	}
}

func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			UsernameOrEmail string `json:"username_or_email" binding:"required"`
			Password        string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		var u models.User
		if err := db.Where("username = LOWER(?) OR mail = LOWER(?)",
			req.UsernameOrEmail, req.UsernameOrEmail).First(&u).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
			return
		}

		secretKey := os.Getenv("SECRET")
		if secretKey == constants.Empty {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": u.ID,
			"exp": time.Now().Add(time.Hour * constants.ExpDate).Unix(),
		})

		tokenString, err := token.SignedString([]byte(secretKey))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, constants.MaxCookieAge, "/", constants.Empty, false, true)
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}

func LogoutHandler(_ *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("Authorization", constants.Empty, -1, "/", constants.Empty, false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}
