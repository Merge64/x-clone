package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/models"
	"main/services/user"
	"net/http"
)

func ViewUserProfileHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		var profile struct {
			// NameTag   string
			Username  string
			Location  string
			CreatedAt string
		}

		currentUser, getUserErr := user.GetUserByUsername(db, username)
		if getUserErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": getUserErr.Error()})
			return
		}

		profile.Username = currentUser.Username
		profile.CreatedAt = currentUser.CreatedAt.String()
		if currentUser.Location != nil {
			profile.Location = *currentUser.Location
		}
		c.JSON(http.StatusOK, gin.H{"message": "View Account successfully", "profile": profile})
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

		c.JSON(http.StatusOK, gin.H{"message": "Edit Profile successfully"})
	}
}

func GetFollowersProfileHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")

		followers, getFollowersErr := user.GetFollowers(db, username)
		if getFollowersErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": getFollowersErr.Error()})
			return
		}

		listFollowers := user.EnlistUsers(followers)

		c.JSON(http.StatusOK, gin.H{"message": "View Followers Profile successfully", "followers": listFollowers})
	}
}

func GetFollowingProfileHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")

		following, getFollowingErr := user.GetFollowing(db, username)
		if getFollowingErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": getFollowingErr.Error()})
			return
		}

		listFollowing := user.EnlistUsers(following)

		c.JSON(http.StatusOK, gin.H{"message": "View Profile Following successfully", "following": listFollowing})
	}
}

func IsAlreadyFollowingHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		followedUsername := c.Param("username")
		usernameAux, _ := c.Get("username")
		username, _ := usernameAux.(string)

		isFollowing, isFollowingErr := user.IsFollowing(db, followedUsername, username)
		if isFollowingErr != nil {
			c.JSON(http.StatusOK, gin.H{"message": "Check Following successfully", "isFollowing": isFollowing})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Check Following successfully", "isFollowing": isFollowing})
	}
}
