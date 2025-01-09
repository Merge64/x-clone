package user

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"main/constants"
	"main/models"
	"regexp"
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

func MailAlreadyUsed(db *gorm.DB, mail string) bool {
	var user models.User
	err := db.Where("Mail = ?", mail).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false

	} else if err != nil {
		log.Printf("Error querying user by email: %v", err)
		return false

	}
	return true
}
func UsernameAlreadyUsed(db *gorm.DB, username string) bool {
	var user models.User
	err := db.Where("Username = ?", username).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false

	} else if err != nil {
		log.Printf("Error querying user by username: %v", err)
		return false

	}
	return true
}

func ValidateCredentials(db *gorm.DB, inputUser, password string) bool {
	var user models.User

	if IsEmail(inputUser) {
		err := db.Where("Mail = ? AND Password = ?", inputUser, password).First(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false

		} else if err != nil {
			log.Printf("Error querying user by email: %v", err)
			return false

		}
	} else {
		err := db.Where("Username = ? AND Password = ?", inputUser, password).First(&user).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false

		} else if err != nil {
			log.Printf("Error querying user by username: %v", err)
			return false

		}
	}

	return true
}

func IsEmail(email string) bool {
	re := regexp.MustCompile(constants.EMAILREGEXPATTERNS)
	return re.MatchString(email)
}

func SearchUsersByUsername(db *gorm.DB, username string) ([]models.User, error) {
	var users []models.User
	result := db.Where("Username LIKE ?", username).First(&users)
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("no users found")
	}
	return users, nil
}

func SearchPostsByKeywords(db *gorm.DB, keyword string) ([]models.Post, error) {
	var posts []models.Post
	result := db.Where("Body ILIKE ?", "%"+keyword+"%").Find(&posts)
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("no posts found containing the keyword: %s", keyword)
	}
	return posts, nil
}

func CreatePost(db *gorm.DB, userID uint, parentID *uint, quoteID *uint, body string) error {
	if !userExists(db, userID) {
		return errors.New("user does not exist")
	}
	post := models.Post{
		UserID:   userID,
		ParentID: parentID,
		Quote:    quoteID,
		Body:     body,
	}

	if err := db.Create(&post).Error; err != nil {
		return err
	}
	return nil
}

//func CreatePostComment(db *gorm.DB, userID uint, parentID, quote *uint, body string) error {
//	if userID == 0 || body == constants.EMPTY {
//		return errors.New("required fields must not be empty")
//	}
//	if !userExists(db, userID) {
//		return errors.New("user does not exist")
//	}
//	var currentPost = models.Post{
//		Model:    gorm.Model{},
//		UserID:   userID,
//		ParentID: parentID,
//		Quote:    quote,
//		Body:     body,
//	}
//	err := db.Model(models.Post{}).Create(&currentPost).Error
//	if err != nil {
//		return fmt.Errorf("failed to create post comment: %v", err)
//	}
//	return nil
//}

// AUX.

func alreadyFollows(db *gorm.DB, followingUserID, followedUserID uint) bool {
	return db.Model(models.Follow{}).
		Where(models.Follow{}, "FollowingUserID = ? AND FollowedUserID = ?", followingUserID, followedUserID).Error == nil
}

func isLiked(db *gorm.DB, userID, parentID uint) bool {
	return db.Model(models.Like{}).Where("UserID = ? AND ParentID = ?", userID, parentID).Error == nil
}

func userExists(db *gorm.DB, userID uint) bool {
	var user models.User
	err := db.Where("id = ?", userID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else if err != nil {
		log.Printf("Error querying user by id: %v", err)
		return false
	}
	return true
}
