package controllers

import (
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

func SignUpHandler(db *gorm.DB) gin.HandlerFunc {
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

func LoginHandler(db *gorm.DB) gin.HandlerFunc {
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
		c.SetCookie("Authorization", tokenString, constants.MAXCOOKIEAGE, "/", constants.EMPTY, false, true)

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}

func LogoutHandler(_ *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("Authorization", constants.EMPTY, -1, "/", constants.EMPTY, false, true)
		c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
	}
}

func FollowUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		followingID, _ := userID.(uint)

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
		userID, _ := c.Get("userID")
		followingID, _ := userID.(uint)

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

func ToggleLikeHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		likerID, _ := userID.(uint)

		postID, atoiErr := strconv.Atoi(c.Param("postid"))
		if atoiErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		toggleResult, toggleErr := user.ToggleLike(db, likerID, uint(postID))
		if toggleErr != nil {
			log.Println("Toggle Like error:", toggleErr)
			if toggleResult.IsLiked {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like post"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike post"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": toggleResult.MessageStatus})
	}
}

func SendMessageHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		senderVal, _ := c.Get("userID")
		senderID, ok := senderVal.(uint)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid userID in context"})
			return
		}

		receiverStr := c.Param("userid")
		receiverInt, err := strconv.Atoi(receiverStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}
		receiverID := uint(receiverInt)

		message := c.PostForm("message")
		if message == constants.EMPTY {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Message content cannot be empty"})
			return
		}

		if errMsg := user.SendMessage(db, senderID, receiverID, message); errMsg != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not send message"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
	}
}

func ListConversationsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserVal, _ := c.Get("userID")
		currentUserID, _ := currentUserVal.(uint)

		type ConvID struct {
			ID uint `json:"id"`
		}
		var conversationIDs []ConvID

		if err := db.Model(&models.Conversation{}).
			Select("ID").
			Where("sender_id = ? OR receiver_id = ?", currentUserID, currentUserID).
			Scan(&conversationIDs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load conversation IDs"})
			return
		}

		c.JSON(http.StatusOK, conversationIDs)
	}
}

func GetMessagesForConversationHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserVal, _ := c.Get("userID")
		currentUserID, _ := currentUserVal.(uint)

		convoIDStr := c.Param("conversationID")
		convoIDInt, err := strconv.Atoi(convoIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid conversation ID"})
			return
		}
		convoID := uint(convoIDInt)

		var conversation models.Conversation
		if errDB := db.First(&conversation, convoID).Error; errDB != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
			return
		}
		if conversation.SenderID != currentUserID && conversation.ReceiverID != currentUserID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not a participant in this conversation"})
			return
		}

		var messages []models.Message
		if errDB := db.Where("conversation_id = ?", conversation.ID).
			Order("created_at asc").
			Find(&messages).Error; errDB != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve messages"})
			return
		}

		c.JSON(http.StatusOK, messages)
	}
}
