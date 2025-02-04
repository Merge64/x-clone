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
	SenderID   uint
	ReceiverID uint
	Messages   []Message
}

type Message struct {
	gorm.Model
	ConversationID uint   `gorm:"index;not null"`
	SenderID       uint   `gorm:"index;not null"`
	Content        string `gorm:"type:text;not null"`
}
