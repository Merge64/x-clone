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

type ErrorMessage struct {
	Message gin.H
	Status  int
}
