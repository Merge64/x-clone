package models

import (
	"gorm.io/gorm"
	"time"
)

type Post struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	CreatedAt  time.Time `json:"created_at"`
	UserID     uint      `json:"user_id"`
	Nickname   string    `json:"nickname"`
	Username   string    `json:"username"`
	ParentID   *uint     `json:"parent_id"`
	Quote      *string   `json:"quote"`
	Body       string    `json:"body"`
	LikesCount uint      `json:"likes_count"`
	IsRepost   bool      `json:"is_repost"`
	ParentPost *Post     `gorm:"foreignKey:ParentID"`
}

type Like struct {
	gorm.Model
	PostID uint
	UserID uint
}
