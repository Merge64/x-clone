package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UserID uint
	Title  string
	Body   string
	Likes  uint
}

type PostLikes struct {
	gorm.Model
	PostID uint // id del post (gorm)
	UserID uint
}

type PostComment struct {
	gorm.Model
	PostID      uint // id del post (gorm)
	UserID      uint
	CommentBody string
}
type CommentComments struct {
	gorm.Model
	CommentID uint
	UserID    uint
	Body      string
}
type CommentLikes struct {
	gorm.Model
	CommentID uint
	UserID    uint
}
