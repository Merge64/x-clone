package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"main/constants"
	"main/mappers"
	"main/models"
	"main/services/user"
	"net/http"
	"strconv"
)

func GetAllPostsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawPosts, err := user.GetAllPosts(db)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "No posts found."})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		listPosts := user.ProcessPosts(rawPosts)
		c.JSON(http.StatusOK, listPosts) // <- Directly return the array
	}
}

func CreatePostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user details from context
		userID, _ := getUserIDFromContext(c)
		username, _ := getUsernameIDFromContext(c)
		nickname, _ := getNicknameFromContext(c)

		// Parse the request body
		var req struct {
			Body     string  `json:"body"`
			Quote    *string `json:"quote"`
			ParentID *uint   `json:"parent_id"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		// Validate the post body
		if err := validatePostBody(req.Body); err != nil {
			sendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		// Check if this is a repost (has a ParentID)
		var parentPost *models.Post
		if req.ParentID != nil {
			// Fetch the parent post
			if err := db.First(&parentPost, *req.ParentID).Error; err != nil {
				sendErrorResponse(c, http.StatusNotFound, "Parent post not found")
				return
			}

			// If the parent post itself is a repost, fetch the original post
			if parentPost.ParentID != nil {
				var originalPost models.Post
				if err := db.First(&originalPost, *parentPost.ParentID).Error; err != nil {
					sendErrorResponse(c, http.StatusNotFound, "Original post not found")
					return
				}
				parentPost = &originalPost
			}
		}

		// Create the post
		createdPost, err := user.CreatePost(db, userID, nickname, req.ParentID, username, req.Quote, req.Body, req.ParentID != nil)
		if err != nil {
			if err.Error() == constants.ErrNoUser {
				sendErrorResponse(c, http.StatusBadRequest, constants.ErrNoUser)
				return
			}
			sendErrorResponse(c, http.StatusInternalServerError, "Failed to create post")
			return
		}

		// Fetch the post with the parent post preloaded
		var postWithParent models.Post
		if errDB := db.Preload("ParentPost").First(&postWithParent, createdPost.ID).Error; errDB != nil {
			sendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch post details")
			return
		}

		// Process the post to include the nested parent_post in the response
		processedPost := mappers.ProcessPost(postWithParent)

		// Return the response
		c.JSON(http.StatusCreated, gin.H{
			"message": "Post created successfully",
			"post":    processedPost,
		})
	}
}

func GetPostsByUsernameHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		rawPosts, errDB := user.GetAllPostsByUsername(db, username)

		if errDB != nil {
			if errors.Is(errDB, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "No posts found with the given username."})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		listPosts := user.ProcessPosts(rawPosts)

		c.JSON(http.StatusOK, gin.H{"posts": listPosts})
	}
}

func GetSpecificPostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID, errPostID := strconv.Atoi(c.Param("postid"))
		if errPostID != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		var post models.Post
		if err := db.First(&post, postID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrNoPost})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while fetching the post"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"post": post})
	}
}

func EditPostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID, atoiErr := strconv.Atoi(c.Param("postid"))
		if atoiErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		postError := editPost(c, db, postID)

		c.JSON(postError.Status, postError.Message)
	}
}

func DeletePostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID, atoiErr := strconv.Atoi(c.Param("postid"))
		userID, _ := c.Get("userID")
		currentUserID, _ := userID.(uint)

		if atoiErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		postError := deletePost(db, currentUserID, postID)
		c.JSON(postError.Status, postError.Message)
	}
}

func CreateRepostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		parentIDStr := c.Param("postid")
		parentID, err := strconv.Atoi(parentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parent post id"})
			return
		}

		repostError := createRepost(c, db, parentID)
		c.JSON(repostError.Status, repostError.Message)
	}
}

func ToggleLikeHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")
		likerID, _ := userID.(uint)

		postID, atoiErr := strconv.Atoi(c.Param("postid"))
		if atoiErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		toggleResult, toggleErr := user.ToggleLike(db, likerID, uint(postID))
		if toggleErr != nil {
			log.Println("Toggle Like error:", toggleErr)
			if toggleResult.IsLiked {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like post"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike post"})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": toggleResult.MessageStatus})
	}
}

// AUX.
func getUserIDFromContext(c *gin.Context) (uint, error) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("unauthorized")
	}

	userID, ok := userIDStr.(uint)
	if !ok {
		return 0, errors.New("invalid user ID")
	}

	return userID, nil
}

func getUsernameIDFromContext(c *gin.Context) (string, error) {
	username, exists := c.Get("username")
	if !exists {
		return "", errors.New("unauthorized")
	}
	usernameStr, ok := username.(string)
	if !ok {
		return "", errors.New("invalid username type")
	}

	return usernameStr, nil
}

func getNicknameFromContext(c *gin.Context) (string, error) {
	nickname, exists := c.Get("nickname")
	if !exists {
		return "", errors.New("unauthorized")
	}
	nicknameStr, ok := nickname.(string)
	if !ok {
		return "", errors.New("invalid nickname type")
	}

	return nicknameStr, nil
}

func validatePostBody(body string) error {
	if body == constants.Empty {
		return errors.New("body cannot be empty")
	}
	return nil
}

// TODO: Implement Share post

func sendErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}

func editPost(c *gin.Context, db *gorm.DB, postID int) PostError {
	var req models.Post

	// Parse the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		return PostError{Message: gin.H{"error": "Invalid JSON: " + err.Error()}, Status: http.StatusBadRequest}
	}

	// Fetch the post to edit
	post, getPostErr := user.GetPostByID(db, uint(postID))
	if getPostErr != nil {
		if errors.Is(getPostErr, gorm.ErrRecordNotFound) {
			return PostError{Message: gin.H{"error": constants.ErrNoPost}, Status: http.StatusNotFound}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while fetching the post"})
		return PostError{
			Message: gin.H{"error": "An error occurred while fetching the post"},
			Status:  http.StatusInternalServerError,
		}
	}

	// Ensure the current user is the owner of the post
	userID, _ := c.Get("userID")
	currentUserID, _ := userID.(uint)
	if post.UserID != currentUserID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not the owner of this post"})
		return PostError{Message: gin.H{"error": "You are not the owner of this post"}, Status: http.StatusUnauthorized}
	}

	// Allow editing of either the body, the quote, or both
	if req.Body != "" {
		post.Body = req.Body
	}

	if req.Quote != nil {
		post.Quote = req.Quote
	}

	// Save the updated post to the database
	if db.Save(&post).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return PostError{Message: gin.H{"error": "Failed to update post"}, Status: http.StatusInternalServerError}
	}

	return PostError{Message: gin.H{"message": "Post updated successfully"}, Status: http.StatusOK}
}

func createRepost(c *gin.Context, db *gorm.DB, parentID int) PostError {
	userIDVal, _ := c.Get("userID")
	username, _ := getUsernameIDFromContext(c)
	nickname, _ := getNicknameFromContext(c)

	currentUserID, ok := userIDVal.(uint)
	if !ok {
		return PostError{Message: gin.H{"error": "Invalid user ID type"}, Status: http.StatusInternalServerError}
	}

	var payload struct {
		Quote string `json:"quote"`
	}

	if errJSON := c.ShouldBindJSON(&payload); errJSON != nil {
		return PostError{Message: gin.H{"error": "Invalid JSON payload"}, Status: http.StatusBadRequest}
	}

	var parentPost models.Post
	if errDB := db.First(&parentPost, parentID).Error; errDB != nil {
		return PostError{Message: gin.H{"error": "Original post not found"}, Status: http.StatusNotFound}
	}

	if parentPost.ParentID != nil && *parentPost.ParentID != 0 {
		var originalPost models.Post
		if errDB := db.First(&originalPost, *parentPost.ParentID).Error; errDB != nil {
			return PostError{Message: gin.H{"error": "Original post not found"}, Status: http.StatusNotFound}
		}
		parentPost = originalPost
	}

	r := models.Post{
		UserID:   currentUserID,
		ParentID: &parentPost.ID,
		Username: username,
		Body:     "",
		Quote:    &payload.Quote,
		Nickname: nickname,
		IsRepost: true,
	}

	// Create the post
	createdPost, errDB := user.CreatePost(db, r.UserID, r.Nickname, r.ParentID, r.Username, r.Quote, r.Body, r.IsRepost)
	if errDB != nil {
		return PostError{Message: gin.H{"error": "Failed to create repost"}, Status: http.StatusInternalServerError}
	}

	// Fetch the post with ParentPost preloaded
	var postWithParent models.Post
	if err := db.Preload("ParentPost").First(&postWithParent, createdPost.ID).Error; err != nil {
		return PostError{Message: gin.H{"error": "Failed to fetch repost"}, Status: http.StatusInternalServerError}
	}

	// Process the post to include parent_post in the response
	processedPost := mappers.ProcessPost(postWithParent)

	return PostError{Message: gin.H{"message": "Repost created successfully", "post": processedPost}, Status: http.StatusCreated}
}

func deletePost(db *gorm.DB, currentUserID uint, postID int) PostError {
	if !user.IsPostOwner(db, currentUserID, uint(postID)) {
		return PostError{Message: gin.H{"error": "You are not the owner of this post"}, Status: http.StatusUnauthorized}
	}

	post, getPostErr := user.GetPostByID(db, uint(postID))
	if getPostErr != nil {
		if errors.Is(getPostErr, gorm.ErrRecordNotFound) {
			return PostError{Message: gin.H{"error": constants.ErrNoPost}, Status: http.StatusNotFound}
		}
		return PostError{
			Message: gin.H{"error": "An error occurred while fetching the post"},
			Status:  http.StatusInternalServerError,
		}
	}

	if deleteErr := db.Delete(&post).Error; deleteErr != nil {
		return PostError{Message: gin.H{"error": "Failed to delete post"}, Status: http.StatusInternalServerError}
	}
	return PostError{Message: gin.H{"message": "Post deleted successfully"}, Status: http.StatusOK}
}

type PostError struct {
	Message gin.H
	Status  int
}
