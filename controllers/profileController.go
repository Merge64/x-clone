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
			Username string
			Mail     string
			Location string
		}

		currentUser, getUserErr := user.GetUserByUsername(db, username)
		if getUserErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": getUserErr.Error()})
			return
		}

		profile.Username = currentUser.Username
		profile.Mail = currentUser.Mail
		if currentUser.Location != nil {
			profile.Location = *currentUser.Location
		}
		c.JSON(http.StatusOK, gin.H{"message": "View Account successfully", "user": profile})
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
		username := c.Param("username")
		var listFollowers []struct {
			Username string  `json:"username"`
			Mail     string  `json:"mail"`
			Location *string `json:"location"`
		}

		followers, getFollowersErr := user.GetFollowers(db, username)
		if getFollowersErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": getFollowersErr.Error()})
			return
		}

		for _, currentUser := range followers {
			var profile struct {
				Username string  `json:"username"`
				Mail     string  `json:"mail"`
				Location *string `json:"location"`
			}

			profile.Username = currentUser.Username
			profile.Mail = currentUser.Mail
			profile.Location = currentUser.Location

			listFollowers = append(listFollowers, profile)
		}

		c.JSON(http.StatusOK, gin.H{"message": "View Followers Profile successfully", "followers": listFollowers})
	}
}

func GetFollowingProfileHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		var listFollowing []struct {
			Username string  `json:"username"`
			Mail     string  `json:"mail"`
			Location *string `json:"location"`
		}

		following, getFollowingErr := user.GetFollowing(db, username)
		if getFollowingErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": getFollowingErr.Error()})
			return
		}

		for _, currentUser := range following {
			var profile struct {
				Username string  `json:"username"`
				Mail     string  `json:"mail"`
				Location *string `json:"location"`
			}

			profile.Username = currentUser.Username
			profile.Mail = currentUser.Mail
			profile.Location = currentUser.Location

			listFollowing = append(listFollowing, profile)
		}

		c.JSON(http.StatusOK, gin.H{"message": "View Profile Following successfully", "following": listFollowing})
	}
}
