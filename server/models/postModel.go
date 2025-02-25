package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UserID     uint    `json:"user_id"`
	ParentID   *uint   `json:"parent_id"`
	Quote      *string `json:"quote"`
	Body       string  `json:"body"`
	LikesCount uint    `json:"likes_count"` // <-- new field
}

type Like struct {
	gorm.Model
	PostID uint
	UserID uint
}
