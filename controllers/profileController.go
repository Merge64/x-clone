package controllers

import (
	"encoding/json"
	"gorm.io/gorm"
	"main/constants"
	"main/models"
	"main/services/user"
	"net/http"
	"strconv"
)

func ViewUserProfileHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	userID, err := strconv.Atoi(r.PathValue("userid"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, getUserErr := user.GetUserByID(db, uint(userID))
	if getUserErr != nil {
		http.Error(w, getUserErr.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(user)
}

// ----------------------------- AUX ----------------------------- //

var ViewUserProfileEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "/profile/{userid}",
	HandlerFunction: ViewUserProfileHandler,
}
