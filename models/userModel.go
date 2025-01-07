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

type Follow struct {
	gorm.Model
	FollowingUserID uint
	FollowedUserID  uint
}
