package controllers

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/services/user"
	"net/http"
	"strings"
)

// TODO: Make this functional when searching with more spaces in between
func SearchHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyword := c.Query("q")

		keyword = strings.TrimSpace(keyword)
		keywords := strings.Fields(keyword)
		keywordProcessed := strings.Join(keywords, " ")

		filter := c.Query("f")
		switch filter {
		case "":
			posts, err := user.SearchPostsByKeywords(db, keyword)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"posts": posts})

		case "latest":
			posts, err := user.SearchPostsByKeywordsSortedByLatest(db, keywordProcessed)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"posts": posts})

		case "user":
			users, err := user.SearchUsersByUsername(db, keywordProcessed)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"users": users})

		case "unique-user":
			exists, err := user.SearchUniqueMailUsername(db, "username", keywordProcessed)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"exists": exists})

		case "unique-mail":
			exists, err := user.SearchUniqueMailUsername(db, "mail", keywordProcessed)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"exists": exists})

		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter parameter"})
		}
	}
}
