package controllers

import (
	"gorm.io/gorm"
	"main/constants"
	"main/models"
	"main/services/user"
	"net/http"
	"strconv"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
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
		parsedParentID, err := strconv.ParseUint(parentIDStr, 10, 32)
		if err != nil {
			http.Error(w, "Invalid parent ID", http.StatusBadRequest)
			return
		}
		tempParentID := uint(parsedParentID)
		parentID = &tempParentID
	}

	if quoteIDStr != "" {
		parsedQuoteID, err := strconv.ParseUint(quoteIDStr, 10, 32)
		if err != nil {
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
	w.Write([]byte("Post created successfully"))
}

var CreatePostEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "posts/create",
	HandlerFunction: CreatePostHandler,
}
