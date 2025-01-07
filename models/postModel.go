package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UserID uint
	Title  string
	Body   string
	Likes  uint
}

type Like struct {
	gorm.Model
	ParentID uint
	UserID   uint
}

type Comment struct {
	gorm.Model
	ParentID uint
	UserID   uint
	Body     string
}
