package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/constants"
	"main/models"
	"main/services/user"
	"net/http"
	"strconv"
)

func ViewUserProfileHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, atoiErr := strconv.Atoi(c.Param("userid"))
		if atoiErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		currentUser, getUserErr := user.GetUserByID(db, uint(userID))
		if getUserErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": getUserErr.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "View Account successfully", "user": currentUser})
	}
}

// TODO: We need to finalize which method we are going to use to obtain this data.
func EditUserProfileHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		profileID, atoiErr := strconv.Atoi(c.Param("userid"))
		userID, _ := c.Get("userID")
		currentUserID, _ := userID.(uint)

		if atoiErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		if currentUserID != uint(profileID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not allowed to edit this profile"})
			return
		}

		var currentUser models.User
		if decodeErr := c.ShouldBindJSON(&currentUser); decodeErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
			return
		}

		if currentUser.ID != uint(profileID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User ID in URL does not match user ID in body"})
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

//// ----------------------------- AUX ----------------------------- //

var GetFollowersProfileEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "profile/followers/user/:userid",
	HandlerFunction: GetFollowersProfileHandler,
}

var GetFollowingProfileEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "profile/following/user/:userid",
	HandlerFunction: GetFollowingProfileHandler,
}

var ViewUserProfileEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "profile/:userid",
	HandlerFunction: ViewUserProfileHandler,
}

var EditUserProfileEndpoint = models.Endpoint{
	Method:          models.PUT,
	Path:            constants.BASEURL + "profile/:userid/edit",
	HandlerFunction: EditUserProfileHandler,
}
