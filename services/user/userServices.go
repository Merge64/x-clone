package user

import (
	"errors"
	"gorm.io/gorm"
	"main/constants"
	"main/models"
)

func CreateAccount(db *gorm.DB, username, password, mail string, location *string) error {
	if password == constants.EMPTY || username == constants.EMPTY {
		return errors.New("fields must not be empty")
	}

	var currentUser = models.User{
		Model:    gorm.Model{},
		Username: username,
		Mail:     mail,
		Location: location,
		Password: password,
	}

	db.Model(models.User{}).Create(&currentUser)

	return nil
}

func FollowAccount(db *gorm.DB, followingUserID, followedUserID uint) error {
	if followingUserID == followedUserID {
		return errors.New("invalid Id")
	}
	if alreadyFollows(db, followingUserID, followedUserID) {
		return errors.New("user already follows")
	}

	db.Model(models.Follow{}).Create(models.Follow{
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

	var user models.Follow
	db.Model(models.Follow{}).First(&user, "FollowingUserID = ? AND FollowedUserID = ?", followingUserID, followedUserID)
	db.Model(models.Follow{}).Delete(&user)

	return nil
}

func ToggleLike(db *gorm.DB, userID uint, parentID uint) error {
	if !userExists(db, userID) {
		return errors.New("user does not exist")
	}

	var currentUser models.Like
	if isLiked(db, userID, parentID) {
		db.Model(models.Like{}).First(&currentUser, "UserID = ? AND ParentID = ?", userID, parentID)
		db.Model(models.Like{}).Delete(&currentUser)
	} else {
		db.Model(models.Like{}).Create(models.Like{
			Model:    gorm.Model{},
			ParentID: parentID,
			UserID:   userID,
		})
	}

	return nil
}

func CreatePost(db *gorm.DB, userID uint, parentID, quote *uint, body string) error {
	if !userExists(db, userID) {
		return errors.New("user does not exist")
	}

	db.Model(models.Post{}).Create(models.Post{
		Model:    gorm.Model{},
		UserID:   userID,
		ParentID: parentID,
		Quote:    quote,
		Body:     body,
	})

	return nil
}

func MailAlreadyUsed(db *gorm.DB, mail string) bool {
	return db.Model(models.User{}).Where("Mail = ?", mail).Error == nil
}

// AUX.
func alreadyFollows(db *gorm.DB, followingUserID, followedUserID uint) bool {
	return db.Model(models.Follow{}).
		Where(models.Follow{}, "FollowingUserID = ? AND FollowedUserID = ?", followingUserID, followedUserID).Error == nil
}

func isLiked(db *gorm.DB, userID, parentID uint) bool {
	return db.Model(models.Like{}).Where("UserID = ? AND ParentID = ?", userID, parentID).Error == nil
}

func userExists(db *gorm.DB, userID uint) bool {
	return db.Model(models.User{}).Where("UserID = ?", userID).First(models.User{}).Error == nil
}
