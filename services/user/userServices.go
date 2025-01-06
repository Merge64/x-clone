package user

import (
	"errors"
	"gorm.io/gorm"
	"main/models"
)

const EMPTY = ""

func CreateAccount(db *gorm.DB, username, password string) error {
	if password == EMPTY || username == EMPTY {
		return errors.New("fields must not be empty")
	}

	db.Model(models.User{}).Create(models.User{
		Model:    gorm.Model{},
		Username: username,
		Password: password,
	})
	return nil
}

func FollowAccount(db *gorm.DB, followingUserId, followedUserId int) error {
	if followingUserId == followedUserId {
		return errors.New("invalid Id")
	}
	if alreadyFollows(db, followingUserId, followedUserId) {
		return errors.New("user already follows")
	}

	db.Model(models.Follows{}).Create(models.Follows{
		Model:             gorm.Model{},
		Following_user_id: followingUserId,
		Followed_user_id:  followedUserId,
	})
	return nil
}

func UnfollowAccount(db *gorm.DB, followingUserId, followedUserId int) error {
	if followingUserId == followedUserId {
		return errors.New("invalid Id")
	}
	var user models.Follows
	db.Model(models.Follows{}).First(&user, "Following_user_id = ? AND Followed_user_id = ?", followingUserId, followedUserId)
	db.Model(models.Follows{}).Delete(&user)
	return nil
}

// AUX
func alreadyFollows(db *gorm.DB, followingUserId, followedUserId int) bool {
	var user models.Follows
	return db.Model(models.Follows{}).First(&user, "Following_user_id = ? AND Followed_user_id = ?", followingUserId, followedUserId).Error != nil
}
