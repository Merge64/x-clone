package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
			if errors.Is(err, errors.New("no posts found")) {
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
		if parseErr := parseFormData(c); parseErr != nil {
			sendErrorResponse(c, http.StatusBadRequest, parseErr.Error())
			return
		}

		userID, getUserIDErr := getUserIDFromContext(c)
		if getUserIDErr != nil {
			sendErrorResponse(c, http.StatusUnauthorized, getUserIDErr.Error())
			return
		}

		body := c.PostForm("body")
		if err := validatePostBody(body); err != nil {
			sendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		var quote *string // its empty when created
		var parentID *uint
		if createPostErr := user.CreatePost(db, userID, parentID, quote, body); createPostErr != nil {
			if createPostErr.Error() == constants.ErrNoUser {
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

		body := c.PostForm("body")
		if body == constants.Empty {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Body cannot be empty"})
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

		// If this post is a repost (has a ParentID), do not allow editing the body. forbidden
		if post.ParentID != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Cannot edit the body of a repost"})
			return
		}

		post.Body = body
		if db.Save(&post).Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Post updated successfully"})
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
		// Parse the parent post id from the URL parameter.
		parentIDStr := c.Param("parentid")
		parentID, err := strconv.Atoi(parentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parent post id"})
			return
		}

		// Get the current user's id from the context.
		userIDVal, _ := c.Get("userID")
		currentUserID, ok := userIDVal.(uint)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type"})
			return
		}

		var req struct {
			Quote string `json:"quote"`
		}
		if errJSON := c.ShouldBindJSON(&req); errJSON != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}

		var parentPost models.Post
		if errDB := db.First(&parentPost, parentID).Error; errDB != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Original post not found"})
			return
		}

		// Flatten the repost chain:
		// If the parent post is itself a repost (i.e. it has a ParentID),
		// then we want to fetch the original post.
		if parentPost.ParentID != nil && *parentPost.ParentID != 0 {
			var originalPost models.Post
			if errDB := db.First(&originalPost, *parentPost.ParentID).Error; errDB != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Original post not found"})
				return
			}
			parentPost = originalPost
		}

		repost := models.Post{
			UserID:   currentUserID,
			ParentID: &parentPost.ID,
			Body:     parentPost.Body,
		}
		if req.Quote != constants.Empty {
			repost.Quote = &req.Quote
		}

		if errDB := db.Create(&repost).Error; errDB != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create repost"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Repost created successfully",
			"post":    repost,
		})
	}
}

//// AUX.

func parseFormData(c *gin.Context) error {
	if err := c.Request.ParseForm(); err != nil {
		return errors.New("invalid form data")
	}
	return nil
}

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

func validatePostBody(body string) error {
	if body == constants.Empty {
		return errors.New("body cannot be empty")
	}
	return nil
}

func sendErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}
