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
	"strings"
)

func CreatePostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.Request.ParseForm(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
			return
		}

		userIDStr, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		parentIDStr := c.PostForm("parent")
		quoteIDStr := c.PostForm("quote")
		body := c.PostForm("body")

		if body == constants.EMPTY {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Body cannot be empty"})
			return
		}

		userID, ok := userIDStr.(uint)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var parentID, quoteID *uint
		if parentIDStr != constants.EMPTY {
			parsedParentID, parentErr := strconv.ParseUint(parentIDStr, 10, 32)
			if parentErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parent ID"})
				return
			}
			tempParentID := uint(parsedParentID)
			parentID = &tempParentID
		}

		if quoteIDStr != constants.EMPTY {
			parsedQuoteID, parsedErr := strconv.ParseUint(quoteIDStr, 10, 32)
			if parsedErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quote ID"})
				return
			}
			tempQuoteID := uint(parsedQuoteID)
			quoteID = &tempQuoteID
		}

		if err := user.CreatePost(db, userID, parentID, quoteID, body); err != nil {
			if err.Error() == constants.ERRNOUSER {
				c.JSON(http.StatusBadRequest, gin.H{"error": constants.ERRNOUSER})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Post created successfully"})
	}
}

func GetAllPostsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		posts, err := user.GetAllPosts(db)
		if err != nil {
			if errors.Is(err, errors.New("no posts found")) {
				c.JSON(http.StatusNotFound, gin.H{"error": "No posts found."})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"posts": posts})
	}
}

func GetPostsByUserIDHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseUint(c.Param("userid"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		posts, errDB := user.GetAllPostsByUserID(db, uint(userID))
		if errDB != nil {
			if errors.Is(errDB, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "No posts found with the given userID."})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"posts": posts})
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
				c.JSON(http.StatusNotFound, gin.H{"error": constants.ERRNOPOST})
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
		postID, err := strconv.Atoi(c.Param("postid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		body := c.PostForm("body")
		if body == constants.EMPTY {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Body cannot be empty"})
			return
		}

		post, getPostErr := user.GetPostByID(db, uint(postID))
		if getPostErr != nil {
			if errors.Is(getPostErr, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": constants.ERRNOPOST})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occurred while fetching the post"})
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
		postID, err := strconv.Atoi(c.Param("postid"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
			return
		}

		post, getPostErr := user.GetPostByID(db, uint(postID))
		if getPostErr != nil {
			if errors.Is(getPostErr, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": constants.ERRNOPOST})
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

//// AUX.

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
	Path:            constants.BASEURL + "/profile/:userid/posts",
	HandlerFunction: GetPostsByUserIDHandler,
}

var GetAllPostsEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "posts/all",
	HandlerFunction: GetAllPostsHandler,
}

var GetSpecificPostEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "posts/:postid",
	HandlerFunction: GetSpecificPostHandler,
}

var DeletePostEndpoint = models.Endpoint{
	Method:          models.DELETE,
	Path:            constants.BASEURL + "posts/:postid/delete",
	HandlerFunction: DeletePostHandler,
}

var EditPostEndpoint = models.Endpoint{
	Method:          models.PUT,
	Path:            constants.BASEURL + "posts/:postid/edit",
	HandlerFunction: EditPostHandler,
}

var CreatePostEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "posts/create",
	HandlerFunction: CreatePostHandler,
}
