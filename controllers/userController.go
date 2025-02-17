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

		var payload struct {
			ReceiverID uint   `json:"receiver_id" binding:"required"`
			Message    string `json:"message" binding:"required"`
		}

		if err := c.ShouldBind(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing receiver_id or message in request"})
			return
		}

		if strings.TrimSpace(payload.Message) == constants.Empty {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Message content cannot be empty"})
			return
		}

		if senderID == payload.ReceiverID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You cannot send a message to yourself"})
			return
		}

		if err := user.SendMessage(db, senderID, payload.ReceiverID, payload.Message); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not send message"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Message sent successfully"})
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
		err := db.Model(&models.Conversation{}).
			Joins("LEFT JOIN messages ON messages.conversation_id = conversations.id").
			Where("conversations.sender_id = ? OR conversations.receiver_id = ?", currentUserID, currentUserID).
			Group("conversations.id").
			Order("COALESCE(MAX(messages.created_at), conversations.created_at) desc").
			Preload("Messages", func(db *gorm.DB) *gorm.DB {
				return db.Order("created_at desc").Limit(1)
			}).
			Find(&conversations).Error

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load conversations"})
			return
		}

		c.JSON(http.StatusOK, conversations)
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

		receiverIDStr := c.Param("receiverID")
		senderIDStr := c.Param("senderID")

		receiverIDInt, err := strconv.Atoi(receiverIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receiver ID"})
			return
		}
		senderIDInt, err := strconv.Atoi(senderIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sender ID"})
			return
		}
		receiverID := uint(receiverIDInt)
		senderID := uint(senderIDInt)

		if currentUserID != receiverID && currentUserID != senderID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not a participant in this conversation"})
			return
		}

		var conversation models.Conversation
		if errDB := db.Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at desc")
		}).Where(
			"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
			senderID, receiverID, receiverID, senderID,
		).First(&conversation).Error; errDB != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
			return
		}

		c.JSON(http.StatusOK, conversation)
	}
}
