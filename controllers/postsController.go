package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/constants"
	"main/models"
	"main/services/user"
	"net/http"
	"strconv"
)

func CreatePostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		if parseErr := parseFormData(c); parseErr != nil {
			sendErrorResponse(c, http.StatusBadRequest, parseErr.Error())
			return
		}

		userID, getUserIDerr := getUserIDFromContext(c)
		if getUserIDerr != nil {
			sendErrorResponse(c, http.StatusUnauthorized, getUserIDerr.Error())
			return
		}

		body := c.PostForm("body")
		if err := validatePostBody(body); err != nil {
			sendErrorResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		parentID, optionalIDErr := parseOptionalID(c, "parent")
		if optionalIDErr != nil {
			sendErrorResponse(c, http.StatusBadRequest, optionalIDErr.Error())
			return
		}

		quote := constants.EMPTY

		if createPostErr := user.CreatePost(db, userID, parentID, quote, body); createPostErr != nil {
			if createPostErr.Error() == constants.ERRNOUSER {
				sendErrorResponse(c, http.StatusBadRequest, constants.ERRNOUSER)
				return
			}
			sendErrorResponse(c, http.StatusInternalServerError, "Failed to create post")
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

func CreateRepostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		parentIDStr := c.Param("parentid")
		parentID, err := strconv.Atoi(parentIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid parent post id"})
			return
		}

		userIDVal, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		currentUserID := userIDVal.(uint)

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

		// 5. Flatten the repost chain:
		//    If the parent post is itself a repost (i.e. it has a ParentID), then
		//    fetch the original post from which it was reposted, as x.com does.
		if parentPost.ParentID != nil {
			// Logically, we’re now reposting the original content.
			if errDB := db.First(&parentPost, *parentPost.ParentID).Error; errDB != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Original post not found"})
				return
			}
		}

		repost := models.Post{
			UserID:   currentUserID,
			ParentID: &parentPost.ID,
			Body:     parentPost.Body,
		}
		if req.Quote != "" {
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
	if body == constants.EMPTY {
		return errors.New("body cannot be empty")
	}
	return nil
}

func parseOptionalID(c *gin.Context, paramName string) (*uint, error) {
	paramStr := c.PostForm(paramName)
	if paramStr == constants.EMPTY {
		return nil, errors.New(constants.ERRNOVALUE)
	}

	parsedID, err := strconv.ParseUint(paramStr, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid %s ID", paramName)
	}

	tempID := uint(parsedID)
	return &tempID, nil
}

func sendErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}
