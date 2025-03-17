package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"x-clone/server/constants"
	"x-clone/server/mappers"
	"x-clone/server/models"
	"x-clone/server/services/user"
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

func GetAllRepliesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")

		rawPosts, err := user.GetAllRepliesByUsername(db, username)

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

func PostsWLikesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")

		rawPosts, err := user.PostsWLikesByUsername(db, username)

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
		c.JSON(http.StatusOK, listPosts) //
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
		parentIDUint := uint(parentID)

		parentPost, err := resolveParentPost(db, parentIDUint)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Original post or comment not found"})
			return
		}

		handled, err := handleExistingRepost(db, userID, parentPost.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to handle repost"})
			return
		}
		if handled {
			c.JSON(http.StatusOK, gin.H{"message": "Repost deleted successfully"})
			return
		}

		postErr := createRepost(c, db, int(parentPost.ID))
		if postErr.Status != http.StatusCreated {
			c.JSON(postErr.Status, postErr.Message)
			return
		}

		if errIncrement := incrementRepostCount(db, parentPost.ID); errIncrement != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update repost count"})
			return
		}

		c.JSON(http.StatusCreated, postErr.Message)
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
		// Preload ParentPost to include it in the processing
		result := db.Preload("ParentPost").Where("parent_id = ? AND is_repost = ?", uint(postID), false).Find(&comments)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
			return
		}

		// Process each comment using ProcessPost
		processedComments := make([]mappers.PostResponse, len(comments))
		for i, comment := range comments {
			processedComments[i] = user.ProcessPost(comment) // Use user.ProcessPost if part of a package
		}

		c.JSON(http.StatusOK, gin.H{"comments": processedComments})
	}
}

func CountRepostsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract parent post ID from URL
		postIDStr := c.Param("postid")
		postID, err := strconv.ParseUint(postIDStr, 10, 32) // uint64
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		// Convert uint64 to uint
		postIDUint := uint(postID)

		// Fetch the repost count
		count, err := GetRepostCount(db, postIDUint)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve repost count"})
			return
		}

		// Return the count in JSON format
		c.JSON(http.StatusOK, gin.H{"reposts_count": count})
	}
}

func CountLikesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract parent post ID from URL
		postIDStr := c.Param("postid")
		postID, err := strconv.ParseUint(postIDStr, 10, 32) // uint64
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		// Convert uint64 to uint
		postIDUint := uint(postID)

		// Fetch the repost count
		count, err := GetLikesCount(db, postIDUint)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve repost count"})
			return
		}

		// Return the count in JSON format
		c.JSON(http.StatusOK, gin.H{"likes_count": count})
	}
}

func CountCommentsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract parent post ID from URL
		postIDStr := c.Param("postid")
		postID, err := strconv.ParseUint(postIDStr, 10, 32) // uint64
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		// Convert uint64 to uint
		postIDUint := uint(postID)

		// Fetch the repost count
		count, err := GetCommentsCount(db, postIDUint)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve repost count"})
			return
		}

		// Return the count in JSON format
		c.JSON(http.StatusOK, gin.H{"comments_count": count})
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
		if errDB := db.First(&parentPost, parentPostID).Error; errDB != nil {
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
		if errJSON := c.ShouldBindJSON(&req); errJSON != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errJSON.Error()})
			return
		}

		// Validate body
		if errValidate := validatePostBody(req.Body); errValidate != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": errValidate.Error()})
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

func resolveParentPost(db *gorm.DB, parentID uint) (models.Post, error) {
	var parentPost models.Post
	if err := db.First(&parentPost, parentID).Error; err != nil {
		return models.Post{}, err
	}

	if parentPost.IsRepost && parentPost.ParentID != nil {
		var originalPost models.Post
		if err := db.First(&originalPost, *parentPost.ParentID).Error; err != nil {
			return models.Post{}, err
		}
		parentPost = originalPost
	}

	return parentPost, nil
}

func handleExistingRepost(db *gorm.DB, userID, parentPostID uint) (bool, error) {
	var existingRepost models.Post
	err := db.Where("user_id = ? AND parent_id = ? AND is_repost = ?", userID, parentPostID, true).
		First(&existingRepost).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if errDB := db.Delete(&existingRepost).Error; errDB != nil {
		return false, errDB
	}

	if errIncrement := decrementRepostCount(db, parentPostID); errIncrement != nil {
		return false, errIncrement
	}

	return true, nil
}

func incrementRepostCount(db *gorm.DB, parentPostID uint) error {
	return db.Model(&models.Post{ID: parentPostID}).Update("reposts_count", gorm.Expr("reposts_count + 1")).Error
}

func decrementRepostCount(db *gorm.DB, parentPostID uint) error {
	return db.Model(&models.Post{ID: parentPostID}).Update("reposts_count", gorm.Expr("reposts_count - 1")).Error
}

func createRepost(c *gin.Context, db *gorm.DB, parentID int) PostError {
	userIDVal, _ := c.Get("userID")
	username, _ := user.GetUsernameIDFromContext(c)
	nickname, _ := user.GetNicknameFromContext(c)

	currentUserID, ok := userIDVal.(uint)
	if !ok {
		return PostError{
			Message: gin.H{"error": "invalid user ID type"},
			Status:  http.StatusInternalServerError,
		}
	}

	var payload struct {
		Quote string `json:"quote"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		return PostError{
			Message: gin.H{"error": "invalid JSON payload"},
			Status:  http.StatusBadRequest,
		}
	}

	// Create the repost with resolved parent ID
	parentIDUint := uint(parentID)
	createdPost, err := user.CreatePost(
		db,
		currentUserID,
		nickname,
		&parentIDUint,
		username,
		&payload.Quote,
		constants.Empty,
		true, // isRepost is always true
	)
	if err != nil {
		return PostError{
			Message: gin.H{"error": "failed to create repost"},
			Status:  http.StatusInternalServerError,
		}
	}

	// Fetch the created repost with ParentPost preloaded
	var postWithParent models.Post
	if errPreloadDB := db.Preload("ParentPost").First(&postWithParent, createdPost.ID).Error; errPreloadDB != nil {
		return PostError{
			Message: gin.H{"error": "failed to fetch repost"},
			Status:  http.StatusInternalServerError,
		}
	}

	return PostError{
		Message: gin.H{"message": "repost created successfully"},
		Status:  http.StatusCreated,
	}
}

func GetRepostCount(db *gorm.DB, postID uint) (int64, error) {
	var repostCount int64
	err := db.Model(&models.Post{}).
		Where("id = ?", postID).
		Pluck("reposts_count", &repostCount).Error // Use Pluck instead of Select + Count
	return repostCount, err
}

func GetLikesCount(db *gorm.DB, postID uint) (int64, error) {
	var likesCount int64
	err := db.Model(&models.Post{}).
		Where("id = ?", postID).
		Pluck("likes_count", &likesCount).Error // Use Pluck instead of Select + Count
	return likesCount, err
}

func GetCommentsCount(db *gorm.DB, postID uint) (int64, error) {
	var commentsCount int64

	// Query for counting comments with the specific parent_id, which is postID
	err := db.Model(&models.Post{}).
		Where("parent_id = ? AND is_repost = ?", postID, false). // Fixing the WHERE clause
		Count(&commentsCount).Error                              // Use Count directly instead of Pluck for counting rows

	return commentsCount, err
}

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
	post, getPostErr := user.GetSimplePostByID(db, uint(postID))
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

func deletePost(db *gorm.DB, currentUserID uint, postID int) PostError {
	if !user.IsPostOwner(db, currentUserID, uint(postID)) {
		return PostError{Message: gin.H{"error": "you are not the owner of this post"},
			Status: http.StatusUnauthorized}
	}

	post, getPostErr := user.GetSimplePostByID(db, uint(postID))
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
