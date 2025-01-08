package models

import (
	"gorm.io/gorm"
	"net/http"
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
	HandlerFunction func(w http.ResponseWriter, r *http.Request, db *gorm.DB)
}
