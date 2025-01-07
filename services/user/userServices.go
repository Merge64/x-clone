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

// currentModel can only be either 'Post' or 'Comment'.
func ToggleLike(db *gorm.DB, userID uint, itemID uint, isPost bool) error {
	if !userExists(db, userID) {
		return errors.New("user does not exist")
	}

	if isPost {
		likePost(db, userID, itemID)
	} else {
		likeComment(db, userID, itemID)
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

func CreateComment(db *gorm.DB, userID uint, itemID uint, body string, isPost bool) error {
	if !userExists(db, userID) {
		return errors.New("user does not exist")
	}

	if isPost {
		db.Model(models.PostComment{}).Create(models.PostComment{
			Model:       gorm.Model{},
			PostID:      itemID,
			UserID:      userID,
			CommentBody: body,
		})
	} else {
		db.Model(models.CommentComments{}).Create(models.CommentComments{
			Model:     gorm.Model{},
			CommentID: itemID,
			UserID:    userID,
			Body:      body,
		})
	}

	return nil
}

// AUX
func alreadyFollows(db *gorm.DB, followingUserID, followedUserID uint) bool {
	var user models.Follows
	return db.Model(models.Follows{}).
		First(&user, "FollowingUserID = ? AND FollowedUserID = ?", followingUserID, followedUserID).Error == nil
}

func isLiked(db *gorm.DB, userID, itemID uint, isPost bool) bool {
	if isPost {
		return db.Model(models.PostLikes{}).Where("UserID = ? AND PostID = ?", userID, itemID).First(&models.PostLikes{}).Error == nil
	}
	return db.Model(models.CommentLikes{}).Where("UserID = ? AND CommentID = ?", userID, itemID).First(&models.CommentLikes{}).Error == nil
}

func userExists(db *gorm.DB, userID uint) bool {
	return db.Model(models.User{}).Where("UserID = ?", userID).First(models.User{}).Error == nil
}

func likeComment(db *gorm.DB, userID, commentID uint) {
	var currentUser models.CommentLikes
	if isLiked(db, userID, commentID, false) {
		db.Model(models.CommentLikes{}).Delete(&currentUser)
	}

	db.Model(models.CommentLikes{}).Create(models.CommentLikes{
		Model:     gorm.Model{},
		CommentID: commentID,
		UserID:    userID,
	})
}

func likePost(db *gorm.DB, userID, postID uint) {
	var currentUser models.CommentLikes
	if isLiked(db, userID, postID, true) {
		db.Model(models.PostLikes{}).Delete(&currentUser)

	} else {

		db.Model(models.PostLikes{}).Create(models.PostLikes{
			Model:  gorm.Model{},
			PostID: postID,
			UserID: userID,
		})
	}
}
