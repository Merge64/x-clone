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

func SearchUserHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
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
	if username == constants.EMPTY {
		http.Error(w, "Missing 'username' parameter", http.StatusBadRequest)
		return
	}

	users, err := user.SearchUserByUsername(db, username)
	if err != nil {
		http.Error(w, constants.ERRNOUSER, http.StatusNotFound)
		return
	}

	response, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		return
	}
}

func SearchPostHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	keyword := r.URL.Query().Get("search")
	if keyword == constants.EMPTY {
		http.Error(w, "Missing 'search' query parameter", http.StatusBadRequest)
		return
	}

	posts, err := user.SearchPostsByKeywords(db, keyword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		return
	}
}

var SearchUserEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "search/{username}",
	HandlerFunction: SearchUserHandler,
}

var SearchPostEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "posts",
	HandlerFunction: SearchPostHandler,
}
