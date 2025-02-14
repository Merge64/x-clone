package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/models"
	"main/services/user"
	"net/http"
	"strconv"
)

func ViewUserProfileHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")

		currentUser, getUserErr := user.GetUserByUsername(db, username)
		if getUserErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": getUserErr.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "View Account successfully", "user": currentUser})
	}
}

func EditUserProfileHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		profileID, _ := userID.(uint)

		var currentUser models.User
		currentUser.ID = profileID
		if decodeErr := c.ShouldBindJSON(&currentUser); decodeErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
			return
		}

		if updateProfileErr := user.UpdateProfile(db, &currentUser); updateProfileErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": updateProfileErr.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Edit Profile successfully", "user": currentUser})
	}
}

func GetFollowersProfileHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, atoiErr := strconv.Atoi(c.Param("userid"))
		if atoiErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		followers, getFollowersErr := user.GetFollowers(db, uint(userID))
		if getFollowersErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": getFollowersErr.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "View Followers Profile successfully", "followers": followers})
	}
}

func GetFollowingProfileHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, atoiErr := strconv.Atoi(c.Param("userid"))
		if atoiErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		following, getFollowingErr := user.GetFollowing(db, uint(userID))
		if getFollowingErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": getFollowingErr.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "View Profile Following successfully", "following": following})
	}
}
