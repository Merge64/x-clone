package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UserID int
	Title  string
	Body   string
	Likes  int
}

type PostLikes struct {
	gorm.Model
	PostID int // id del post (gorm)
	UserID int
}

type PostComment struct {
	gorm.Model
	PostID      int // id del post (gorm)
	UserID      int
	CommentBody string
}
type CommentComments struct {
	gorm.Model
	CommentID int
	UserID    int
	Body      string
}
type CommentLikes struct {
	gorm.Model
	CommentID int
	UserID    int
}
