package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID            uint `gorm:"primaryKey"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Nickname      string         `json:"nickname"`
	Username      string         `json:"username"`
	Mail          string         `json:"mail"`
	Password      string         `json:"password"`
	Location      *string        `json:"location"`
	Bio           *string        `json:"bio"`
	FollowerCount uint           `json:"follower_count"`
}

type Follow struct {
	gorm.Model
	FollowingUsername string
	FollowedUsername  string
}

type Conversation struct {
	gorm.Model
	SenderUsername   string    `json:"sender_username"`
	SenderNickname   string    `json:"sender_nickname"`
	ReceiverUsername string    `json:"receiver_username"`
	ReceiverNickname string    `json:"receiver_nickname"`
	Messages         []Message `json:"messages"`
}

type Message struct {
	gorm.Model
	ConversationID uint   `json:"conversation_id" gorm:"index;not null"`
	SenderUsername string `json:"sender_username" gorm:"index;not null"`
	Content        string `json:"content" gorm:"type:text;not null"`
}
