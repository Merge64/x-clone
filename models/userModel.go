package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Password string
}

type Follows struct {
	gorm.Model
	Following_user_id int
	Followed_user_id  int
}
