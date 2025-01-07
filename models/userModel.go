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
