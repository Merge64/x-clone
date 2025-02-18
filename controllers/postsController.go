package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"main/constants"
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
		c.JSON(http.StatusOK, gin.H{"posts": listPosts})
	}
}

func CreatePostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := getUserIDFromContext(c)
		var req models.Post

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		if err := validatePostBody(req.Body); err != nil {
			sendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		var quote *string // its empty when created
		var parentID *uint
		if err := user.CreatePost(db, userID, parentID, quote, req.Body); err != nil {
			if err.Error() == constants.ErrNoUser {
				sendErrorResponse(c, http.StatusBadRequest, constants.ErrNoUser)
				return
			}
			sendErrorResponse(c, http.StatusInternalServerError, "Failed to create post")
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully"})
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
		userID, _ := c.Get("userID")
		currentUserID, _ := userID.(uint)

		if atoiErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		if !user.IsPostOwner(db, currentUserID, uint(postID)) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not the owner of this post"})
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

		if !user.IsPostOwner(db, currentUserID, uint(postID)) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "You are not the owner of this post"})
			return
		}

		post, getPostErr := user.GetPostByID(db, uint(postID))
		if getPostErr != nil {
			if errors.Is(getPostErr, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrNoPost})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while fetching the post"})
			return
		}

		if deleteErr := db.Delete(&post).Error; deleteErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
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

// TODO: Implement Share post

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

	if err := c.ShouldBindJSON(&req); err != nil {
		return PostError{Message: gin.H{"error": "Invalid JSON: " + err.Error()}, Status: http.StatusBadRequest}
	}

	if err := validatePostBody(req.Body); err != nil {
		return PostError{Message: gin.H{"error": err.Error()}, Status: http.StatusBadRequest}
	}

	post, getPostErr := user.GetPostByID(db, uint(postID))
	if getPostErr != nil {
		if errors.Is(getPostErr, gorm.ErrRecordNotFound) {
			return PostError{Message: gin.H{"error": constants.ErrNoPost}, Status: http.StatusNotFound}
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while fetching the post"})
		return PostError{Message: gin.H{"error": "An error occurred while fetching the post"}, Status: http.StatusInternalServerError}
	}

	if post.ParentID != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot edit the body of a repost"})
		return PostError{Message: gin.H{"error": "Cannot edit the body of a repost"}, Status: http.StatusForbidden}
	}

	post.Body = req.Body
	if db.Save(&post).Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return PostError{Message: gin.H{"error": "Failed to update post"}, Status: http.StatusInternalServerError}
	}
	return PostError{Message: gin.H{"message": "Post updated successfully"}, Status: http.StatusOK}
}

func createRepost(c *gin.Context, db *gorm.DB, parentID int) PostError {
	userIDVal, _ := c.Get("userID")
	currentUserID, ok := userIDVal.(uint)
	if !ok {
		return PostError{Message: gin.H{"error": "Invalid user ID type"}, Status: http.StatusInternalServerError}
	}

	var req struct {
		Quote string `json:"quote"`
	}
	if errJSON := c.ShouldBindJSON(&req); errJSON != nil {
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

	rawRepost := models.Post{
		UserID:   currentUserID,
		ParentID: &parentPost.ID,
		Body:     parentPost.Body,
	}
	if req.Quote != constants.Empty {
		rawRepost.Quote = &req.Quote
	}

	if errDB := db.Create(&rawRepost).Error; errDB != nil {
		return PostError{Message: gin.H{"error": "Failed to create repost"}, Status: http.StatusInternalServerError}
	}

	repost := user.ProcessPosts([]models.Post{rawRepost})

	return PostError{Message: gin.H{"message": "Repost created successfully", "post": repost}, Status: http.StatusCreated}
}

type PostError struct {
	Message gin.H
	Status  int
}
