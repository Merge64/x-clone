package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"main/constants"
	"main/models"
	"main/services/user"
	"net/http"
	"sort"
	"strings"
	"time"
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
		// Retrieve the current username
		currentUsernameAux, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Ensure the currentUsername is a string
		currentUsername, ok := currentUsernameAux.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid username in context"})
			return
		}

		// Get sender and receiver usernames from URL parameters
		senderUsername := c.Param("senderUsername")
		receiverUsername := c.Param("receiverUsername")

		// Check if the current user is part of the conversation
		if currentUsername != senderUsername && currentUsername != receiverUsername {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not a participant in this conversation"})
			return
		}

		// Preload messages for the conversation
		var conversation models.Conversation
		if errDB := preloadMessages(db, senderUsername, receiverUsername, &conversation); errDB != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Conversation not found"})
			return
		}

		c.JSON(http.StatusOK, conversation)
	}
}

func SendMessageHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		senderValAux, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		senderValStr, ok := senderValAux.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid senderValStr in context"})
			return
		}

		receiverValStr := c.Param("rUsername")
		if receiverValStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Receiver username is required"})
			return
		}

		if !checkUsernameExists(c, db, senderValStr, receiverValStr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sender or receiver"})
			return
		}

		// Proceed with sending the message
		errorMessage := sendMessage(c, senderValStr, receiverValStr, db)
		c.JSON(errorMessage.Status, errorMessage.Message)
	}
}

func ListConversationsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get current user from context
		currentUserVal, exists := c.Get("username")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		currentUsername, ok := currentUserVal.(string)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid username in context"})
			return
		}

		// Get conversations with latest messages
		var conversations []models.Conversation
		err := getConversation(db, currentUsername, &conversations)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load conversations"})
			return
		}

		// Sort conversations by latest message timestamp (descending)
		sort.Slice(conversations, func(i, j int) bool {
			var iTime, jTime time.Time

			// Get latest message time for conversation i
			if len(conversations[i].Messages) > 0 {
				iTime = conversations[i].Messages[0].CreatedAt
			} else {
				iTime = conversations[i].UpdatedAt // Fallback to conversation update time
			}

			// Get latest message time for conversation j
			if len(conversations[j].Messages) > 0 {
				jTime = conversations[j].Messages[0].CreatedAt
			} else {
				jTime = conversations[j].UpdatedAt
			}

			return iTime.After(jTime) // Descending order
		})

		// Format response
		formattedConversations := []gin.H{}
		for _, conv := range conversations {
			// Get message details
			lastMessage := ""
			timestamp := ""
			if len(conv.Messages) > 0 {
				lastMessage = conv.Messages[0].Content
				timestamp = conv.Messages[0].CreatedAt.Format("Jan 02 15:04")
			}

			// Determine other participant
			var partnerUsername, partnerNickname string
			if conv.SenderUsername == currentUsername {
				partnerUsername = conv.ReceiverUsername
				partnerNickname = conv.ReceiverNickname
			} else {
				partnerUsername = conv.SenderUsername
				partnerNickname = conv.SenderNickname
			}

			formattedConversations = append(formattedConversations, gin.H{
				"id":        conv.ID,
				"username":  partnerUsername,
				"nickname":  partnerNickname,
				"content":   lastMessage,
				"timestamp": timestamp,
			})
		}

		c.JSON(http.StatusOK, formattedConversations)
	}
}

// Aux

func checkUsernameExists(c *gin.Context, db *gorm.DB, sender string, receiver string) bool {
	var senderCount, receiverCount int64
	db.Model(&models.User{}).Where("username = ?", sender).Count(&senderCount)
	db.Model(&models.User{}).Where("username = ?", receiver).Count(&receiverCount)

	if senderCount == 0 || receiverCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sender or receiver does not exist"})
		return false
	} else {
		return true
	}
}

func toggleFollow(db *gorm.DB, c *gin.Context, isFollowing bool) {
	followingUsernameAux, _ := c.Get("username")
	followingUsername, _ := followingUsernameAux.(string)

	var (
		successMessage  string
		logErrorMessage string
		expr            string
	)

	followedUsername := c.Param("username")

	if isFollowing {
		if followErr := user.FollowAccount(db, followingUsername, followedUsername); followErr != nil {
			log.Println("Follow error:", followErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
			return
		}
	} else {
		if unfollowErr := user.UnfollowAccount(db, followingUsername, followedUsername); unfollowErr != nil {
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
		Where("username = ?", followedUsername).
		UpdateColumn("follower_count", gorm.Expr("follower_count "+expr+" 1")).Error; err != nil {
		log.Println("Failed to "+logErrorMessage+" follower_count:", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": successMessage})
}

func sendMessage(c *gin.Context, senderStr string, receiverStr string, db *gorm.DB) ErrorMessage {
	var payload struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBind(&payload); err != nil {
		return ErrorMessage{
			Message: gin.H{"error": "Missing message in request"},
			Status:  http.StatusBadRequest,
		}
	}

	if strings.TrimSpace(payload.Message) == constants.Empty {
		return ErrorMessage{Message: gin.H{"error": "Message content cannot be empty"}, Status: http.StatusBadRequest}
	}

	if receiverStr == senderStr {
		return ErrorMessage{Message: gin.H{"error": "You cannot send a message to yourself"}, Status: http.StatusBadRequest}
	}

	if err := user.SendMessage(db, senderStr, receiverStr, payload.Message); err != nil {
		return ErrorMessage{Message: gin.H{"error": "Could not send message"}, Status: http.StatusInternalServerError}
	}
	return ErrorMessage{Message: gin.H{"message": "Message sent successfully"}, Status: http.StatusOK}
}

func preloadMessages(db *gorm.DB, senderUsername string, receiverUsername string, conversation *models.Conversation) error {
	return db.Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at asc") // Changed to ascending order
	}).Where(
		"(sender_username = ? AND receiver_username = ?) OR (sender_username = ? AND receiver_username = ?)",
		senderUsername, receiverUsername, receiverUsername, senderUsername,
	).First(conversation).Error
}

func getConversation(db *gorm.DB, currentUsername string, conversations *[]models.Conversation) error {
	// First, load conversations **without messages**
	err := db.
		Where("sender_username = ? OR receiver_username = ?", currentUsername, currentUsername).
		Order("updated_at DESC").
		Find(conversations).Error

	if err != nil {
		return err
	}

	// Then, for each conversation, fetch the **latest** message separately
	for i := range *conversations {
		var latestMessage models.Message
		err := db.
			Where("conversation_id = ?", (*conversations)[i].ID).
			Order("created_at DESC").
			Limit(1).
			Find(&latestMessage).Error

		if err != nil {
			return err
		}

		// If a message was found, add it to the conversation
		if latestMessage.ID != 0 {
			(*conversations)[i].Messages = []models.Message{latestMessage}
		}
	}

	return nil
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

		// Start a transaction
		tx := db.Begin()
		if tx.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
			return
		}

		// Update the username in the user table
		u.Username = request.Username
		if err := tx.Save(&u).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update username"})
			return
		}

		// Update the username in all posts by this user
		if err := tx.
			Model(&models.Post{}).
			Where("user_id = ?", userID).
			Update("username", request.Username).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update posts with new username"})
			return
		}

		// Commit the transaction
		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
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
