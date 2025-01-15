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
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userIDStr := r.FormValue("userid")
	parentIDStr := r.FormValue("parent")
	quoteIDStr := r.FormValue("quote")
	body := r.FormValue("body")
	fmt.Println(body)

	if body == constants.EMPTY {
		http.Error(w, "Body cannot be empty", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid userid", http.StatusBadRequest)
		return
	}

	var parentID, quoteID *uint
	if parentIDStr != constants.EMPTY {
		parsedParentID, parentErr := strconv.ParseUint(parentIDStr, 10, 32)
		if parentErr != nil {
			http.Error(w, "Invalid parent ID", http.StatusBadRequest)
			return
		}
		tempParentID := uint(parsedParentID)
		parentID = &tempParentID
	}

	if quoteIDStr != constants.EMPTY {
		parsedQuoteID, parsedErr := strconv.ParseUint(quoteIDStr, 10, 32)
		if parsedErr != nil {
			http.Error(w, "Invalid quote ID", http.StatusBadRequest)
			return
		}
		tempQuoteID := uint(parsedQuoteID)
		quoteID = &tempQuoteID
	}

	err = user.CreatePost(db, uint(userID), parentID, quoteID, body)
	if err != nil {
		if err.Error() == constants.ERRNOUSER {
			http.Error(w, constants.ERRNOUSER, http.StatusBadRequest)
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

func GetAllPostsHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
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

	response, errMarshal := json.Marshal(posts)
	if errMarshal != nil {
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

func GetPostsByUserIDHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
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

	posts, errDB := user.GetAllPostsByUserID(db, finalID)
	if errDB != nil {
		if errors.Is(errDB, errors.New(constants.ERRNOPOST)) { // Directly compare with the constant
			http.Error(w, "No posts found with the given userID.", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Internal server error: %v", errDB), http.StatusInternalServerError)
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

func GetSpecificPostHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	postID, errPostID := extractPostID(r)
	if errPostID != nil {
		http.Error(w, errPostID.Error(), http.StatusBadRequest) // Handle invalid postID error
		return
	}

	// Fetch the post from the database using the postID
	var post models.Post
	if err := db.First(&post, postID).Error; err != nil {
		// Check for specific GORM error for "record not found"
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, constants.ERRNOPOST, http.StatusNotFound) // Post not found error
			return
		}

		// For other errors, log and return an internal server error
		http.Error(w, "An error occurred while fetching the post", http.StatusInternalServerError)
		return
	}

	// Marshal the post into JSON and send as the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if errEnc := json.NewEncoder(w).Encode(post); errEnc != nil {
		http.Error(w, "Failed to serialize response", http.StatusInternalServerError)
		return
	}
}

func EditPostHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	postID, err := extractPostID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newBody := r.FormValue("body")
	if newBody == constants.EMPTY {
		http.Error(w, "Body cannot be empty", http.StatusBadRequest)
		return
	}
	post, getPostErr := user.GetPostByID(db, postID)
	if getPostErr != nil {
		if errors.Is(getPostErr, gorm.ErrRecordNotFound) {
			http.Error(w, constants.ERRNOPOST, http.StatusNotFound) // Post not found error
			return
		}
		http.Error(w, "An error occurred while fetching the post", http.StatusInternalServerError)
		return
	}
	post.Body = newBody
	if db.Save(&post).Error != nil {
		http.Error(w, "Failed to update post", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Post updated successfully"))
}

func DeletePostHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	postID, err := extractPostID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	post, getPostErr := user.GetPostByID(db, postID)
	if getPostErr != nil {
		if errors.Is(getPostErr, gorm.ErrRecordNotFound) {
			http.Error(w, constants.ERRNOPOST, http.StatusNotFound) // Post not found error
			return
		}
		http.Error(w, "An error occurred while fetching the post", http.StatusInternalServerError)
		return
	}
	if deleteErr := db.Delete(&post).Error; deleteErr != nil {
		http.Error(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Post deleted successfully"))
}

// AUX.

// Extract postID from the URL.
const pathSize = 3 // It will always be /posts/{postid}. a way to overcome this is using MUX

func extractPostID(r *http.Request) (uint, error) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < pathSize {
		return 0, errors.New("invalid URL format")
	}
	postID, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, errors.New("invalid postid")
	}
	return uint(postID), nil
}

var GetAllPostsByUserIDEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "/profile/{userid}/posts",
	HandlerFunction: GetPostsByUserIDHandler,
}

var GetAllPostsEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "posts/all",
	HandlerFunction: GetAllPostsHandler,
}

var DeletePostEndpoint = models.Endpoint{
	Method:          models.DELETE,
	Path:            constants.BASEURL + "posts/{postid}/delete",
	HandlerFunction: DeletePostHandler,
}

var EditPostEndpoint = models.Endpoint{
	Method:          models.PUT,
	Path:            constants.BASEURL + "posts/{postid}/edit",
	HandlerFunction: EditPostHandler,
}

var GetSpecificPostEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "posts/{postid}",
	HandlerFunction: GetSpecificPostHandler,
}

var CreatePostEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "posts/create",
	HandlerFunction: CreatePostHandler,
}
