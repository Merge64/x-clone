package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"main/constants"
	"main/models"
	"main/services/user"
	"net/http"
	"os"
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

		if err := user.CreateAccount(db, username, string(hashedPassword), mail, locationAux); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid parameters to create an account"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Account created successfully"})
	}
}

func LoginHandlerGin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		usernameOrEmail := c.PostForm("username-or-email")
		password := c.PostForm("password")

		// Validate input (either username or email must be provided)
		if usernameOrEmail == "" || password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}

		// Check if the user exists
		var u models.User
		if err := db.Where("username = ? OR mail = ?", usernameOrEmail, usernameOrEmail).First(&u).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Check if the password is correct
		if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}

		// Load secret key for JWT
		secretKey := os.Getenv("SECRET")
		if secretKey == "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Server configuration error"})
			return
		}

		// Generate JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": u.ID,
			"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		})

		tokenString, err := token.SignedString([]byte(secretKey))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		// Send response
		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}

// // TODO: In the future implement JWT.
//
//	func UserLoginHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
//		if r.Method != http.MethodPost {
//			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//			return
//		}
//
//		if err := r.ParseForm(); err != nil {
//			http.Error(w, "Invalid request body", http.StatusBadRequest)
//			return
//		}
//
//		inputUser := r.FormValue("username-or-email")
//		password := r.FormValue("password")
//
//		if inputUser == constants.EMPTY || password == constants.EMPTY {
//			http.Error(w, "Missing required fields", http.StatusBadRequest)
//			return
//		}
//
//		if !user.ValidateCredentials(db, inputUser, password) {
//			w.WriteHeader(http.StatusOK)
//			_, err := w.Write([]byte("username or password is incorrect"))
//			if err != nil {
//				log.Printf("Failed to write response: %v", err)
//			}
//			return
//		}
//
//		w.WriteHeader(http.StatusAccepted)
//		_, err := w.Write([]byte("Login successful"))
//		if err != nil {
//			log.Printf("Failed to write response: %v", err)
//		}
//	}
//
//	func FollowUserHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
//		if r.Method != http.MethodPost {
//			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//			return
//		}
//
//		followingID, getIDErr := getUserID(r)
//		if getIDErr != nil {
//			http.Error(w, "Invalid user", http.StatusBadRequest)
//			return
//		}
//
//		followedUserID, atoiErr := strconv.Atoi(r.PathValue("userid"))
//		if atoiErr != nil {
//			http.Error(w, "Invalid user ID", http.StatusBadRequest)
//			return
//		}
//
//		if followErr := user.FollowAccount(db, followingID, uint(followedUserID)); followErr != nil {
//			fmt.Println(followErr)
//			http.Error(w, "Failed to follow user", http.StatusInternalServerError)
//			return
//		}
//
//		w.WriteHeader(http.StatusOK)
//		_, err := w.Write([]byte("Follows user successfully"))
//		if err != nil {
//			log.Printf("Failed to write response: %v", err)
//		}
//	}
//
//	func UnfollowUserHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
//		if r.Method != http.MethodDelete {
//			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
//			return
//		}
//
//		followingID, getIDErr := getUserID(r)
//		if getIDErr != nil {
//			http.Error(w, "Invalid user", http.StatusBadRequest)
//			return
//		}
//
//		followedUserID, atoiErr := strconv.Atoi(r.PathValue("userid"))
//		if atoiErr != nil {
//			http.Error(w, "Invalid user ID", http.StatusBadRequest)
//			return
//		}
//
//		if unfollowErr := user.UnfollowAccount(db, followingID, uint(followedUserID)); unfollowErr != nil {
//			http.Error(w, "Failed to follow user", http.StatusInternalServerError)
//			return
//		}
//
//		w.WriteHeader(http.StatusOK)
//		_, err := w.Write([]byte("Unfollows user successfully"))
//		if err != nil {
//			log.Printf("Failed to write response: %v", err)
//		}
//	}
//
//	func getUserID(r *http.Request) (uint, error) {
//		var currentUser models.User
//		if decodeErr := json.NewDecoder(r.Body).Decode(&currentUser); decodeErr != nil {
//			return 0, decodeErr
//		}
//		return currentUser.ID, nil
//	}
//
//	var FollowUserEndpoint = models.Endpoint{
//		Method:          models.POST,
//		Path:            constants.BASEURL + "follow/{userid}",
//		HandlerFunction: FollowUserHandler,
//	}
//
//	var UnfollowUserEndpoint = models.Endpoint{
//		Method:          models.DELETE,
//		Path:            constants.BASEURL + "unfollow/{userid}",
//		HandlerFunction: UnfollowUserHandler,
//	}
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
