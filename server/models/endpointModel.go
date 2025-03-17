package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

type Endpoint struct {
	Method          string
	Path            string
	HandlerFunction func(db *gorm.DB) gin.HandlerFunc
}
