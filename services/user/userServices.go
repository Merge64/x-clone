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

func FollowAccount(db *gorm.DB, followingUserID, followedUserID uint) error {
	if followingUserID == followedUserID {
		return errors.New("invalid Id")
	}
	if alreadyFollows(db, followingUserID, followedUserID) {
		return errors.New("user already follows")
	}

	db.Model(models.Follows{}).Create(models.Follows{
		Model:           gorm.Model{},
		FollowingUserID: followingUserID,
		FollowedUserID:  followedUserID,
	})
	return nil
}

func UnfollowAccount(db *gorm.DB, followingUserID, followedUserID uint) error {
	if followingUserID == followedUserID {
		return errors.New("invalid Id")
	}
	var user models.Follows
	db.Model(models.Follows{}).First(&user, "FollowingUserID = ? AND FollowedUserID = ?", followingUserID, followedUserID)
	db.Model(models.Follows{}).Delete(&user)
	return nil
}

// currentModel can only be either 'PostLiked' or 'CommentLiked'.
func ToggleLike(db *gorm.DB, userID uint, postID uint, currentModel any) error {
	if !userExists(db, userID) {
		return errors.New("user does not exist")
	}
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

func CreatePost(db *gorm.DB, userID uint, title string, body string) error {
	if !userExists(db, userID) {
		return errors.New("user does not exist")
	}

	var currentPost models.Post
	db.Model(currentPost).Create(models.Post{
		Model:  gorm.Model{},
		UserID: userID,
		Title:  title,
		Body:   body,
		Likes:  0,
	})
	return nil
}

// AUX.
func alreadyFollows(db *gorm.DB, followingUserID, followedUserID uint) bool {
	var user models.Follows
	return db.Model(models.Follows{}).
		First(&user, "FollowingUserID = ? AND FollowedUserID = ?", followingUserID, followedUserID).Error == nil
}

func isLiked(db *gorm.DB, userID, postID uint, currentUser models.PostLikes) bool {
	return db.Model(models.PostLikes{}).Where("UserID = ? AND PostID = ?", userID, postID).First(&currentUser).Error == nil
}

func userExists(db *gorm.DB, userID uint) bool {
	return db.Model(models.User{}).Where("UserID = ?", userID).First(models.User{}).Error == nil
}
