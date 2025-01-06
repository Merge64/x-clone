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
		Model:           gorm.Model{},
		FollowingUserID: followingUserId,
		FollowedUserID:  followedUserId,
	})
	return nil
}

func UnfollowAccount(db *gorm.DB, followingUserId, followedUserId int) error {
	if followingUserId == followedUserId {
		return errors.New("invalid Id")
	}
	var user models.Follows
	db.Model(models.Follows{}).First(&user, "FollowingUserID = ? AND FollowedUserID = ?", followingUserId, followedUserId)
	db.Model(models.Follows{}).Delete(&user)
	return nil
}

// currentModel can only be either 'Post_liked' or 'Comment_liked'.
func ToggleLike(db *gorm.DB, userID, postID int, currentModel any) error {
	var currentUser models.PostLikes

	if isLiked(db, userID, postID, currentUser) {
		db.Model(currentModel).Delete(&currentUser)
	} else {
		db.Model(currentModel).Create(models.PostLikes{
			Model:  gorm.Model{},
			PostID: postID,
			UserID: userID,
		})
	}

	return nil
}

// AUX

func alreadyFollows(db *gorm.DB, followingUserId, followedUserId int) bool {
	var user models.Follows
	return db.Model(models.Follows{}).First(&user, "FollowingUserID = ? AND FollowedUserID = ?", followingUserId, followedUserId).Error != nil
}

func isLiked(db *gorm.DB, userID, postID int, currentUser models.PostLikes) bool {
	return db.Model(models.PostLikes{}).Where("UserID = ? AND PostID = ?", userID, postID).First(&currentUser).Error != nil
}
