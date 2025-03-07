package controllers

import (
	"errors"
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
		req, err := parseSignUpRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if errSignup := validateSignUpRequest(db, req); errSignup != nil {
			c.JSON(http.StatusOK, gin.H{"error": errSignup.Error()})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		newUser := models.User{
			Nickname: req.Nickname,
			Username: req.Username,
			Password: string(hashedPassword),
			Mail:     req.Mail,
			Location: getOptionalLocation(req.Location),
		}

		if errDB := db.Create(&newUser).Error; errDB != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
			return
		}

		tokenString, err := generateJWT(newUser.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		setAuthCookie(c, tokenString)

		c.JSON(http.StatusCreated, gin.H{
			"message": "Account created successfully",
			"token":   tokenString,
		})
	}
}

func LoginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		req, err := parseLoginRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		u, err := findUserByUsernameOrEmail(db, req.UsernameOrEmail)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
			return
		}

		if errHashPassword := bcrypt.CompareHashAndPassword([]byte(u.Password),
			[]byte(req.Password)); errHashPassword != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid credentials"})
			return
		}

		tokenString, err := generateJWT(u.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		setAuthCookie(c, tokenString)
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}

func LogoutHandler(_ *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		setAuthCookie(c, constants.Empty) // Expire the cookie
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}

// AUX.

func findUserByUsernameOrEmail(db *gorm.DB, usernameOrEmail string) (*models.User, error) {
	var u models.User
	err := db.Where("username = LOWER(?) OR mail = LOWER(?)", usernameOrEmail, usernameOrEmail).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func parseLoginRequest(c *gin.Context) (*LoginRequest, error) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("invalid JSON: %s")
	}
	return &req, nil
}

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

func parseSignUpRequest(c *gin.Context) (*SignUpRequest, error) {
	var req SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, errors.New("invalid JSON: %s")
	}
	if req.Username == constants.Empty || req.Password == constants.Empty || req.Mail == constants.Empty {
		return nil, errors.New("missing required fields")
	}
	return &req, nil
}

type SignUpRequest struct {
	Nickname string `json:"nickname" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Mail     string `json:"mail" binding:"required"`
	Location string `json:"location"` // optional
}

func validateSignUpRequest(db *gorm.DB, req *SignUpRequest) error {
	if !user.IsEmail(req.Mail) {
		return errors.New("invalid email")
	}
	if user.MailAlreadyUsed(db, req.Mail) {
		return errors.New("email already in use")
	}
	if user.UsernameAlreadyUsed(db, req.Username) {
		return errors.New("username already in use")
	}
	return nil
}

func getOptionalLocation(location string) *string {
	if location == constants.Empty {
		return nil
	}
	return &location
}

func generateJWT(userID uint) (string, error) {
	secretKey := os.Getenv("SECRET")
	if secretKey == constants.Empty {
		return constants.Empty, errors.New("server configuration error")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * constants.ExpDate).Unix(),
	})

	return token.SignedString([]byte(secretKey))
}

func setAuthCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"Authorization",
		token, constants.MaxCookieAge,
		"/", constants.Empty,
		false,
		true)
}
