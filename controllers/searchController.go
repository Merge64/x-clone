package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/constants"
	"main/services/user"
	"net/http"
)

func SearchHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyword := c.Query("q")
		if keyword == constants.Empty {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'q' query parameter"})
			return
		}

		filter := c.Query("f")
		switch filter {
		case constants.Empty:
			posts, err := user.SearchPostsByKeywords(db, keyword)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"posts": posts})

		case "latest":
			posts, err := user.SearchPostsByKeywordsSortedByLatest(db, keyword)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"posts": posts})

		case "user":
			users, err := user.SearchUserByUsername(db, keyword)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"users": users})

		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter parameter."})
		}
	}
}
