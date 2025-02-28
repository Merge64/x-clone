package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"main/constants"
	"main/models"
	"main/services/user"
	"net/http"
	"strconv"
	"strings"
)

func FollowUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		toggleFollow(db, c, true)
	}
}

func UnfollowUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		toggleFollow(db, c, false)
	}
}

func GetMessagesForConversationHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserVal, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		currentUserID, ok := currentUserVal.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid userID in context"})
			return
		}

		receiverIDStr, senderIDStr := c.Param("receiverID"), c.Param("senderID")
		receiverIDInt, errReciver := strconv.Atoi(receiverIDStr)
		senderIDInt, errSender := strconv.Atoi(senderIDStr)

		if errReciver != nil || errSender != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receiver or sender ID"})
			return
		}
		receiverID := uint(receiverIDInt)
		senderID := uint(senderIDInt)

		if currentUserID != receiverID && currentUserID != senderID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not a participant in this conversation"})
			return
		}

		var conversation models.Conversation
		if errDB := preloadMessages(db, senderID, receiverID, conversation); errDB != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
			return
		}

		c.JSON(http.StatusOK, conversation)
	}
}

func SendMessageHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		senderVal, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		senderID, ok := senderVal.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid userID in context"})
			return
		}

		errorMessage := sendMessage(c, senderID, db)
		c.JSON(errorMessage.Status, errorMessage.Message)
	}
}

func ListConversationsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		currentUserVal, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		currentUserID, ok := currentUserVal.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid userID in context"})
			return
		}

		var conversations []models.Conversation
		err := getConversation(db, currentUserID, conversations)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load conversations"})
			return
		}

		c.JSON(http.StatusOK, conversations)
	}
}

// Aux

func toggleFollow(db *gorm.DB, c *gin.Context, isFollowing bool) {
	userID, _ := c.Get("userID")
	followingID, _ := userID.(uint)
	var (
		successMessage  string
		logErrorMessage string
		expr            string
	)

	followedUserID, atoiErr := strconv.Atoi(c.Param("userid"))
	if atoiErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	if isFollowing {
		if followErr := user.FollowAccount(db, followingID, uint(followedUserID)); followErr != nil {
			log.Println("Follow error:", followErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
			return
		}
	} else {
		if unfollowErr := user.UnfollowAccount(db, followingID, uint(followedUserID)); unfollowErr != nil {
			log.Println("Unfollow error:", unfollowErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unfollow user"})
			return
		}
	}

	if isFollowing {
		successMessage = "Followed user successfully"
		logErrorMessage = "increment"
		expr = "+"
	} else {
		successMessage = "Unfollowed user successfully"
		logErrorMessage = "decrement"
		expr = "-"
	}

	if err := db.Model(&models.User{}).
		Where("id = ?", followedUserID).
		UpdateColumn("follower_count", gorm.Expr("follower_count "+expr+" 1")).Error; err != nil {
		log.Println("Failed to "+logErrorMessage+" follower_count:", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": successMessage})
}

func sendMessage(c *gin.Context, senderID uint, db *gorm.DB) ErrorMessage {
	var payload struct {
		ReceiverID uint   `json:"receiver_id" binding:"required"`
		Message    string `json:"message" binding:"required"`
	}

	if err := c.ShouldBind(&payload); err != nil {
		return ErrorMessage{
			Message: gin.H{"error": "Missing receiver_id or message in request"},
			Status:  http.StatusBadRequest,
		}
	}

	if strings.TrimSpace(payload.Message) == constants.Empty {
		return ErrorMessage{Message: gin.H{"error": "Message content cannot be empty"}, Status: http.StatusBadRequest}
	}

	if senderID == payload.ReceiverID {
		return ErrorMessage{Message: gin.H{"error": "You cannot send a message to yourself"}, Status: http.StatusBadRequest}
	}

	if err := user.SendMessage(db, senderID, payload.ReceiverID, payload.Message); err != nil {
		return ErrorMessage{Message: gin.H{"error": "Could not send message"}, Status: http.StatusInternalServerError}
	}
	return ErrorMessage{Message: gin.H{"message": "Message sent successfully"}, Status: http.StatusOK}
}

func preloadMessages(db *gorm.DB, senderID uint, receiverID uint, conversation models.Conversation) error {
	return db.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc")
	}).Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		senderID, receiverID, receiverID, senderID,
	).First(&conversation).Error
}

func getConversation(db *gorm.DB, currentUserID uint, conversations []models.Conversation) error {
	return db.Model(&models.Conversation{}).
		Joins("LEFT JOIN messages ON messages.conversation_id = conversations.id").
		Where("conversations.sender_id = ? OR conversations.receiver_id = ?", currentUserID, currentUserID).
		Group("conversations.id").
		Order("COALESCE(MAX(messages.created_at), conversations.created_at) desc").
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc").Limit(1)
		}).
		Find(&conversations).Error
}

func GetUserInfoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

		var u models.User
		if err := db.First(&u, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"username":  u.Username,
			"email":     u.Mail,
			"createdAt": u.CreatedAt,
		})
	}
}

func UpdateUsernameHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

		var request struct {
			Username string `json:"username" binding:"required"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		var u models.User
		if err := db.First(&u, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		if u.Username == request.Username {
			c.JSON(http.StatusOK, gin.H{"message": "Username unchanged"})
			return
		}

		var existingUser models.User
		if err := db.Where("username = ?", request.Username).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already taken"})
			return
		}

		u.Username = request.Username
		if err := db.Save(&u).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update username"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":  "Username updated successfully",
			"username": u.Username,
		})
	}
}

type ErrorMessage struct {
	Message gin.H
	Status  int
}
