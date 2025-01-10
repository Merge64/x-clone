package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UserID   uint   `json:"userid"`
	ParentID *uint  `json:"parentid"`
	Quote    *uint  `json:"quote"`
	Body     string `json:"body"`
}

type Like struct {
	gorm.Model
	ParentID uint
	UserID   uint
}
