package controllers

import (
	"errors"
	"fmt"
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

func GetSpecificPostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postID, errPostID := strconv.Atoi(c.Param("postid"))
		if errPostID != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		var post models.Post
		// Use Preload to eagerly load the ParentPost relationship
		if err := db.Preload("ParentPost").First(&post, postID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrNoPost})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while fetching the post"})
			return
		}

		// Process the single post using the ProcessPost function
		processedPost := user.ProcessPost(post)

		c.JSON(http.StatusOK, gin.H{"post": processedPost})
	}
}

func CreatePostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract user details
		userID, _ := user.GetUserIDFromContext(c)
		username, _ := user.GetUsernameIDFromContext(c)
		nickname, _ := user.GetNicknameFromContext(c)

		// Parse and validate request
		req, err := parsePostRequest(c)
		if err != nil {
			sendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		// Validate post body
		if errValidate := validatePostBody(req.Body); errValidate != nil {
			sendErrorResponse(c, http.StatusBadRequest, errValidate.Error())
			return
		}

		// Ensure the parent post exists (if provided)
		if req.ParentID != nil {
			if _, errFetch := fetchParentPost(db, req.ParentID); errFetch != nil {
				sendErrorResponse(c, http.StatusNotFound, errFetch.Error())
				return
			}
		}

		// Create post
		createdPost, err := user.CreatePost(db,
			userID,
			nickname,
			req.ParentID,
			username,
			req.Quote,
			req.Body,
			req.ParentID != nil)
		if err != nil {
			handlePostCreationError(c, err)
			return
		}

		// Fetch and process the post
		processedPost, err := fetchAndProcessPost(db, createdPost.ID)
		if err != nil {
			sendErrorResponse(c, http.StatusInternalServerError, "Failed to fetch post details")
			return
		}

		// Return response
		c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully", "post": processedPost})
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
		userID, _ := user.GetUserIDFromContext(c)

		// Validate userID (check if empty/zero)
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized: User ID is missing or invalid"})
			return
		}

		parentIDStr := c.Param("postid")
		parentID, err := strconv.Atoi(parentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parent post ID"})
			return
		}

		// Convert parentID (int) to uint
		parentIDUint := uint(parentID)

		// Check if the user already has a repost of this post
		var existingRepost models.Post
		errCheck := db.Where("user_id = ? AND parent_id = ? AND is_repost = ?",
			userID,
			parentIDUint,
			true).First(&existingRepost).Error

		if errCheck == nil {
			// Existing repost found → Delete it
			db.Delete(&existingRepost)
			db.Model(&models.Post{}).Where("id = ?", parentIDUint).Update("reposts_count", gorm.Expr("reposts_count - 1"))
			c.JSON(http.StatusOK, gin.H{"message": "Repost deleted successfully"})
			return
		}

		postErr := createRepost(c, db, parentID)
		if postErr.Status != http.StatusCreated {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Repost could not be created"})
			return
		}

		db.Model(&models.Post{}).Where("id = ?", parentIDUint).Update("reposts_count", gorm.Expr("reposts_count + 1"))
		c.JSON(http.StatusCreated, gin.H{"message": "Repost created successfully"})
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

func CheckRepostedHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		userID, err := user.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

		// Get post ID from URL parameter
		postIDStr := c.Param("postid")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		// Check for existing repost
		var existingRepost models.Post
		err = db.Where("user_id = ? AND parent_id = ? AND is_repost = ?",
			userID,
			postID,
			true).
			First(&existingRepost).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"reposted": false})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// If no error, repost exists
		c.JSON(http.StatusOK, gin.H{"reposted": true})
	}
}

func CheckIfLikedHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		userID, err := user.GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			return
		}

		// Get post ID from URL parameter
		postIDStr := c.Param("postid")
		postID, err := strconv.Atoi(postIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		// Check if like exists in database
		var like models.Like
		err = db.Where("user_id = ? AND post_id = ?", userID, postID).First(&like).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"liked": false})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// If no error, like exists
		c.JSON(http.StatusOK, gin.H{"liked": true})
	}
}

func GetCommentsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		postIDStr := c.Param("postid")
		postID, err := strconv.ParseUint(postIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		var comments []models.Post
		result := db.Where("parent_id = ? AND is_repost = ?", uint(postID), false).Find(&comments)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"comments": comments})
	}
}

func CreateCommentHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract parent post ID from URL
		postIDStr := c.Param("postid")
		parentPostID, err := strconv.ParseUint(postIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		// Check if parent post exists
		var parentPost models.Post
		if err := db.First(&parentPost, parentPostID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Parent post not found"})
			return
		}

		// Extract user details
		userID, _ := user.GetUserIDFromContext(c)
		username, _ := user.GetUsernameIDFromContext(c)
		nickname, _ := user.GetNicknameFromContext(c)

		// Parse request body
		var req struct {
			Body  string  `json:"body"`
			Quote *string `json:"quote"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate body
		if err := validatePostBody(req.Body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Create the comment (set ParentID and IsRepost=false)
		parentIDUint := uint(parentPostID)
		createdPost, err := user.CreatePost(
			db,
			userID,
			nickname,
			&parentIDUint, // ParentID from URL
			username,
			req.Quote,
			req.Body,
			false, // isRepost explicitly set to false
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
			return
		}

		// Return the created comment
		c.JSON(http.StatusCreated, gin.H{
			"message": "Comment created successfully",
			"comment": createdPost,
		})
	}
}

// AUX.

// Parses the request body and returns the structured request.
func parsePostRequest(c *gin.Context) (*struct {
	Body     string  `json:"body"`
	Quote    *string `json:"quote"`
	ParentID *uint   `json:"parent_id"`
}, error) {
	var req struct {
		Body     string  `json:"body"`
		Quote    *string `json:"quote"`
		ParentID *uint   `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		return nil, fmt.Errorf("invalid JSON: %s", err.Error())
	}
	return &req, nil
}

var ErrParentPostNotFound = errors.New("parent post ID not provided")

// Fetches the parent post if it exists and resolves original posts for reposts.
func fetchParentPost(db *gorm.DB, parentID *uint) (*models.Post, error) {
	if parentID == nil {
		return nil, ErrParentPostNotFound
	}

	var parentPost models.Post
	if err := db.First(&parentPost, *parentID).Error; err != nil {
		return nil, errors.New("parent post not found")
	}

	if parentPost.ParentID != nil {
		var originalPost models.Post
		if err := db.First(&originalPost, *parentPost.ParentID).Error; err != nil {
			return nil, errors.New("original post not found")
		}
		return &originalPost, nil
	}
	return &parentPost, nil
}

// Handles post creation errors in a centralized function.
func handlePostCreationError(c *gin.Context, err error) {
	if err.Error() == constants.ErrNoUser {
		sendErrorResponse(c, http.StatusBadRequest, constants.ErrNoUser)
		return
	}
	sendErrorResponse(c, http.StatusInternalServerError, "failed to create post")
}

// Fetches the created post and processes it into the API response format.
func fetchAndProcessPost(db *gorm.DB, postID uint) (mappers.PostResponse, error) {
	var postWithParent models.Post
	if err := db.Preload("ParentPost").First(&postWithParent, postID).Error; err != nil {
		return mappers.PostResponse{}, err
	}
	return mappers.ProcessPost(postWithParent), nil
}

func validatePostBody(body string) error {
	if body == constants.Empty {
		return errors.New("body cannot be empty")
	}
	return nil
}

func sendErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}

func editPost(c *gin.Context, db *gorm.DB, postID int) PostError {
	var req models.Post

	// Parse the request body
	if err := c.ShouldBindJSON(&req); err != nil {
		return PostError{Message: gin.H{"error": "invalid JSON: " + err.Error()},
			Status: http.StatusBadRequest}
	}

	// Fetch the post to edit
	post, getPostErr := user.GetPostByID(db, uint(postID))
	if getPostErr != nil {
		if errors.Is(getPostErr, gorm.ErrRecordNotFound) {
			return PostError{Message: gin.H{"error": constants.ErrNoPost},
				Status: http.StatusNotFound}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "an error occurred while fetching the post"})
		return PostError{
			Message: gin.H{"error": "an error occurred while fetching the post"},
			Status:  http.StatusInternalServerError,
		}
	}

	// Ensure the current user is the owner of the post
	userID, _ := c.Get("userID")
	currentUserID, _ := userID.(uint)
	if post.UserID != currentUserID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "you are not the owner of this post"})
		return PostError{Message: gin.H{"error": "you are not the owner of this post"},
			Status: http.StatusUnauthorized}
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update post"})
		return PostError{Message: gin.H{"error": "failed to update post"},
			Status: http.StatusInternalServerError}
	}

	return PostError{Message: gin.H{"message": "post updated successfully"},
		Status: http.StatusOK}
}

func createRepost(c *gin.Context, db *gorm.DB, parentID int) PostError {
	userIDVal, _ := c.Get("userID")
	username, _ := user.GetUsernameIDFromContext(c)
	nickname, _ := user.GetNicknameFromContext(c)

	currentUserID, ok := userIDVal.(uint)
	if !ok {
		return PostError{Message: gin.H{"error": "invalid user ID type"},
			Status: http.StatusInternalServerError}
	}

	var payload struct {
		Quote string `json:"quote"`
	}

	if errJSON := c.ShouldBindJSON(&payload); errJSON != nil {
		return PostError{Message: gin.H{"error": "invalid JSON payload"},
			Status: http.StatusBadRequest}
	}

	var parentPost models.Post
	if errDB := db.First(&parentPost, parentID).Error; errDB != nil {
		return PostError{Message: gin.H{"error": "original post not found"},
			Status: http.StatusNotFound}
	}

	if parentPost.ParentID != nil && *parentPost.ParentID != 0 {
		var originalPost models.Post
		if errDB := db.First(&originalPost, *parentPost.ParentID).Error; errDB != nil {
			return PostError{Message: gin.H{"error": "original post not found"},
				Status: http.StatusNotFound}
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
		return PostError{Message: gin.H{"error": "failed to create repost"},
			Status: http.StatusInternalServerError}
	}

	// Fetch the post with ParentPost preloaded
	var postWithParent models.Post
	if err := db.Preload("ParentPost").First(&postWithParent, createdPost.ID).Error; err != nil {
		return PostError{Message: gin.H{"error": "failed to fetch repost"},
			Status: http.StatusInternalServerError}
	}

	// Process the post to include parent_post in the response
	processedPost := mappers.ProcessPost(postWithParent)

	return PostError{Message: gin.H{"message": "repost created successfully", "post": processedPost},
		Status: http.StatusCreated}
}

func deletePost(db *gorm.DB, currentUserID uint, postID int) PostError {
	if !user.IsPostOwner(db, currentUserID, uint(postID)) {
		return PostError{Message: gin.H{"error": "you are not the owner of this post"},
			Status: http.StatusUnauthorized}
	}

	post, getPostErr := user.GetPostByID(db, uint(postID))
	if getPostErr != nil {
		if errors.Is(getPostErr, gorm.ErrRecordNotFound) {
			return PostError{Message: gin.H{"error": constants.ErrNoPost}, Status: http.StatusNotFound}
		}
		return PostError{
			Message: gin.H{"error": "an error occurred while fetching the post"},
			Status:  http.StatusInternalServerError,
		}
	}

	if deleteErr := db.Delete(&post).Error; deleteErr != nil {
		return PostError{Message: gin.H{"error": "failed to delete post"}, Status: http.StatusInternalServerError}
	}
	return PostError{Message: gin.H{"message": "post deleted successfully"}, Status: http.StatusOK}
}

type PostError struct {
	Message gin.H
	Status  int
}
