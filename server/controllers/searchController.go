package controllers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"x-clone/server/mappers"
	"x-clone/server/models"
	"x-clone/server/services/user"
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

			c.JSON(http.StatusOK, posts)

		case "latest":
			posts, err := user.SearchPostsByKeywordsSortedByLatest(db, keywordProcessed)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, posts)

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

func PrivateSearchHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyword := strings.TrimSpace(c.Query("q"))
		filter := c.Query("f")
		switch filter {
		case "following":
			handleFollowingFilter(c, db, keyword)
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filter parameter"})
		}
	}
}

func fetchPostsFromFollowedUsers(db *gorm.DB, followedUsers []mappers.Response) ([]models.Post, error) {
	var allPosts []models.Post
	for _, followedUser := range followedUsers {
		posts, err := user.GetAllPostsByUsername(db, followedUser.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch posts for user %s: %w", followedUser.Username, err)
		}
		allPosts = append(allPosts, posts...)
	}
	return allPosts, nil
}

// _ should be replaced when keyword has any use.
func handleFollowingFilter(c *gin.Context, db *gorm.DB, _ string) {
	currentUser, err := getCurrentUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	followedUsers, err := user.GetFollowing(db, currentUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	allPosts, err := fetchPostsFromFollowedUsers(db, followedUsers)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(allPosts) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no posts found from followed users"})
		return
	}

	processedPosts := user.ProcessPosts(allPosts)
	c.JSON(http.StatusOK, processedPosts)
}

func getCurrentUser(c *gin.Context) (string, error) {
	currentUserAux, exists := c.Get("username")
	if !exists || currentUserAux == nil {
		return "", errors.New("authentication required for following filter")
	}

	currentUser, ok := currentUserAux.(string)
	if !ok {
		return "", errors.New("invalid username type in context")
	}

	if currentUser == "" {
		return "", errors.New("user parameter is required for following filter")
	}

	return currentUser, nil
}
