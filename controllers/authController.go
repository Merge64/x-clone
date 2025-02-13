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
		if req.Location != "" {
			locationAux = &req.Location
		}

		if req.Username == "" || req.Password == "" || req.Mail == "" {
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

		if createErr := user.CreateAccount(db, req.Username, string(hashedPassword), req.Mail, locationAux); createErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid parameters to create an account"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Account created successfully"})
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
		if err := db.Where("username = ? OR mail = ?", req.UsernameOrEmail, req.UsernameOrEmail).First(&u).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		secretKey := os.Getenv("SECRET")
		if secretKey == "" {
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
		c.SetCookie("Authorization", tokenString, constants.MaxCookieAge, "/", "", false, true)

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}

func LogoutHandler(_ *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("Authorization", "", -1, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}
