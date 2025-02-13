package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/constants"
	"main/services/user"
	"net/http"
)

func SearchUserHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Param("username")
		if username == constants.Empty {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'username' parameter"})
			return
		}

		users, err := user.SearchUserByUsername(db, username)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": constants.ErrNoUser})
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}

func SearchPostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyword := c.Query("search")
		if keyword == constants.Empty {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'search' query parameter"})
			return
		}

		posts, err := user.SearchPostsByKeywords(db, keyword)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"posts": posts})
	}
}
