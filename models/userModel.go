package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string  `json:"username"`
	Mail     string  `json:"mail"`
	Location *string `json:"location"`
	Password string  `json:"password"`
}

type Follow struct {
	gorm.Model
	FollowingUserID uint
	FollowedUserID  uint
}

type Conversation struct {
	gorm.Model
	SenderID   uint      `json:"sender_id"`
	ReceiverID uint      `json:"receiver_id"`
	Messages   []Message `json:"messages"`
}

type Message struct {
	gorm.Model
	ConversationID uint   `json:"conversation_id" gorm:"index;not null"`
	SenderID       uint   `json:"sender_id" gorm:"index;not null"`
	Content        string `json:"content" gorm:"type:text;not null"`
}
