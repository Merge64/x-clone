package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Mail     string
	Location *string
	Password string
}

type Follows struct {
	gorm.Model
	FollowingUserID uint
	FollowedUserID  uint
}
