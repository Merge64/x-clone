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
		userID, atoiErr := strconv.Atoi(c.Param("userid"))
		if atoiErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var currentUser models.User
		if decodeErr := c.ShouldBindJSON(&currentUser); decodeErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user data"})
			return
		}

		if currentUser.ID != uint(userID) {
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

//func GetFollowersProfileHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
//	userID, atoiErr := strconv.Atoi(r.PathValue("userid"))
//	if atoiErr != nil {
//		http.Error(w, "Invalid user ID", http.StatusBadRequest)
//		return
//	}
//
//	followers, getFollowersErr := user.GetFollowers(db, uint(userID))
//	if getFollowersErr != nil {
//		http.Error(w, getFollowersErr.Error(), http.StatusBadRequest)
//		return
//	}
//
//	if encodeError := json.NewEncoder(w).Encode(followers); encodeError != nil {
//		http.Error(w, encodeError.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	w.Header().Set("Content-Type", "application/json")
//	if _, err := w.Write([]byte("View Followers Profile successfully")); err != nil {
//		return
//	}
//}
//
//func GetFollowingProfileHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
//	userID, atoiErr := strconv.Atoi(r.PathValue("userid"))
//	if atoiErr != nil {
//		http.Error(w, "Invalid user ID", http.StatusBadRequest)
//		return
//	}
//
//	following, getFollowingErr := user.GetFollowing(db, uint(userID))
//	if getFollowingErr != nil {
//		http.Error(w, getFollowingErr.Error(), http.StatusBadRequest)
//		return
//	}
//
//	if encodeError := json.NewEncoder(w).Encode(following); encodeError != nil {
//		http.Error(w, encodeError.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.WriteHeader(http.StatusOK)
//	w.Header().Set("Content-Type", "application/json")
//	if _, err := w.Write([]byte("View Profile Following successfully")); err != nil {
//		return
//	}
//}
//
//// ----------------------------- AUX ----------------------------- //
//
//var GetFollowersProfileEndpoint = models.Endpoint{
//	Method:          models.GET,
//	Path:            constants.BASEURL + "profile/followers/user/{userid}",
//	HandlerFunction: GetFollowersProfileHandler,
//}
//
//var GetFollowingProfileEndpoint = models.Endpoint{
//	Method:          models.GET,
//	Path:            constants.BASEURL + "profile/following/user/{userid}",
//	HandlerFunction: GetFollowingProfileHandler,
//}

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
