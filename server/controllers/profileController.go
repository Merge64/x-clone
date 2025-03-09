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

		var u models.User
		if err := db.Where("username = ?", username).First(&u).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Fetch follower and following counts
		var followerCount int64
		db.Model(&models.Follow{}).Where("followed_username = ?", username).Count(&followerCount)

		var followingCount int64
		db.Model(&models.Follow{}).Where("following_username = ?", username).Count(&followingCount)

		// Return profile data with counts
		c.JSON(http.StatusOK, gin.H{
			"profile": gin.H{
				"username":        u.Username,
				"nickname":        u.Nickname,
				"mail":            u.Mail,
				"location":        u.Location,
				"bio":             u.Bio,
				"follower_count":  followerCount,
				"following_count": followingCount,
				"created_at":      u.CreatedAt, // Include the CreatedAt field
			},
		})
	}
}

func EditUserProfileHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		usernameAux, _ := c.Get("username")
		nicknameAux, _ := c.Get("nickname")
		nickname, _ := nicknameAux.(string)
		username, _ := usernameAux.(string)

		var currentUser models.User
		if decodeErr := c.ShouldBindJSON(&currentUser); decodeErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
			return
		}

		if currentUser.Nickname != nickname {
			updateErr := user.UpdateNicknamePosts(db, username, currentUser.Nickname)
			if updateErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": updateErr.Error()})
				return
			}
		}

		if updateProfileErr := user.UpdateProfile(db, username, &currentUser); updateProfileErr != nil {
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

		followerCount := len(followers)

		c.JSON(http.StatusOK, gin.H{
			"users":           followers,
			"following_count": followerCount,
		})
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

		//listFollowing := user.EnlistUsers(following)
		followingCount := len(following)

		c.JSON(http.StatusOK, gin.H{
			"users":           following,
			"following_count": followingCount,
		})
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
