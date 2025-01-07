package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UserID   uint
	ParentID *uint
	Quote    *uint
	Body     string
}

type Like struct {
	gorm.Model
	ParentID uint
	UserID   uint
}
