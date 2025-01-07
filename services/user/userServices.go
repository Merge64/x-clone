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

func ToggleLike(db *gorm.DB, userID uint, parentID uint) error {
	if !userExists(db, userID) {
		return errors.New("user does not exist")
	}

	var currentUser models.Like
	if isLiked(db, userID, parentID) {
		db.Model(models.Like{}).Delete(&currentUser)
	}

	db.Model(models.Like{}).Create(models.Like{
		Model:    gorm.Model{},
		ParentID: parentID,
		UserID:   userID,
	})

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

func CreateComment(db *gorm.DB, userID uint, parentID uint, body string) error {
	if !userExists(db, userID) {
		return errors.New("user does not exist")
	}

	db.Model(models.Comment{}).Create(models.Comment{
		Model:    gorm.Model{},
		ParentID: parentID,
		UserID:   userID,
		Body:     body,
	})

	return nil
}

// AUX
func alreadyFollows(db *gorm.DB, followingUserID, followedUserID uint) bool {
	var user models.Follows
	return db.Model(models.Follows{}).
		First(&user, "FollowingUserID = ? AND FollowedUserID = ?", followingUserID, followedUserID).Error == nil
}

func isLiked(db *gorm.DB, userID, parentID uint) bool {
	return db.Model(models.Like{}).Where("UserID = ? AND ParentID = ?", userID, parentID).First(&models.Like{}).Error == nil
}

func userExists(db *gorm.DB, userID uint) bool {
	return db.Model(models.User{}).Where("UserID = ?", userID).First(models.User{}).Error == nil
}
