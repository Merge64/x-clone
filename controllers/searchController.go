package controllers

import (
	"encoding/json"
	"gorm.io/gorm"
	"main/constants"
	"main/models"
	"main/services/user"
	"net/http"
	"strings"
)

func SearchUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Path
	prefix := "/search/"
	if !strings.HasPrefix(path, prefix) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	username := strings.TrimPrefix(path, prefix)
	if username == "" {
		http.Error(w, "Missing 'username' parameter", http.StatusBadRequest)
		return
	}

	users, err := user.SearchUsersByUsername(db, username)
	if err != nil {
		http.Error(w, "No users found", http.StatusNotFound)
		return
	}

	// Convert the result to JSON
	response, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	// Set the response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

var SearchEndPoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "search/{username}",
	HandlerFunction: SearchUser,
}
