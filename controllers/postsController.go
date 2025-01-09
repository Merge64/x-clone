package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"main/constants"
	"main/models"
	"main/services/user"
	"net/http"
	"strconv"
	"strings"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userIDStr := r.FormValue("userid")
	parentIDStr := r.FormValue("parent")
	quoteIDStr := r.FormValue("quote")
	body := r.FormValue("body")

	if body == "" {
		http.Error(w, "Body cannot be empty", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid userid", http.StatusBadRequest)
		return
	}

	var parentID, quoteID *uint
	if parentIDStr != "" {
		parsedParentID, err2 := strconv.ParseUint(parentIDStr, 10, 32)
		if err2 != nil {
			http.Error(w, "Invalid parent ID", http.StatusBadRequest)
			return
		}
		tempParentID := uint(parsedParentID)
		parentID = &tempParentID
	}

	if quoteIDStr != "" {
		parsedQuoteID, err3 := strconv.ParseUint(quoteIDStr, 10, 32)
		if err3 != nil {
			http.Error(w, "Invalid quote ID", http.StatusBadRequest)
			return
		}
		tempQuoteID := uint(parsedQuoteID)
		quoteID = &tempQuoteID
	}

	err = user.CreatePost(db, uint(userID), parentID, quoteID, body)
	if err != nil {
		if err.Error() == "user does not exist" {
			http.Error(w, "user does not exist", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte("Post created successfully"))
	if err != nil {
		return
	}
}

func GetEveryPostHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	posts, err := user.GetAllPosts(db)
	if err != nil {
		if errors.Is(err, errors.New("no posts found")) {
			http.Error(w, "No posts found.", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	response, err2 := json.Marshal(posts)
	if err2 != nil {
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

func GetPostsByUserID(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := r.URL.Path
	prefix := "/posts/"
	if !strings.HasPrefix(path, prefix) {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}

	urlID := strings.TrimPrefix(path, prefix)
	if urlID == constants.EMPTY {
		http.Error(w, "Missing 'ID' parameter", http.StatusBadRequest)
		return
	}
	parsedID, err := strconv.ParseUint(urlID, 10, 32)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	finalID := uint(parsedID)

	posts, err := user.GetAllPostsByUserID(db, finalID)
	if err != nil {
		if errors.Is(err, errors.New("no posts found")) {
			http.Error(w, "No posts found with the given userID.", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Internal server error: %v", err), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(posts)
	if err != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		return
	}
}

var GetPostsByUserIDEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "posts/{userid}",
	HandlerFunction: GetPostsByUserID,
}

var GetEveryPostEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "posts/all",
	HandlerFunction: GetEveryPostHandler,
}

var CreatePostEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "posts/create",
	HandlerFunction: CreatePostHandler,
}
