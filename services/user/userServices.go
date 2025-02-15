package user

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"main/constants"
	"main/models"
	"regexp"
)

func CreateAccount(db *gorm.DB, username, password, mail string, location *string) error {
	if password == constants.Empty || username == constants.Empty {
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
		return errors.New("invalid ID: user cannot follow themselves")
	}

	if alreadyFollows(db, followingUserID, followedUserID) {
		return errors.New("user already follows this account")
	}

	follow := models.Follow{
		FollowingUserID: followingUserID,
		FollowedUserID:  followedUserID,
	}

	if err := db.Create(&follow).Error; err != nil {
		log.Printf("Error creating follow record: %v", err)
		return err
	}

	return nil
}

func UnfollowAccount(db *gorm.DB, followingUserID, followedUserID uint) error {
	if followingUserID == followedUserID {
		return errors.New("invalid ID: user cannot unfollow themselves")
	}

	result := db.Where("following_user_id = ? AND followed_user_id = ?", followingUserID, followedUserID).
		Delete(&models.Follow{})

	if result.Error != nil {
		log.Printf("Error deleting follow record: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("no follow relationship found to delete")
	}

	return nil
}

func ToggleLike(db *gorm.DB, userID uint, parentID uint) (ToggleInfo, error) {
	if !userExists(db, userID) {
		return ToggleInfo{}, errors.New(constants.ErrNoUser)
	}

	var toggleResult ToggleInfo
	var currentUser models.Like
	if isLiked(db, userID, parentID) {
		db.Model(models.Like{}).First(&currentUser, "user_id = ? AND parent_id = ?", userID, parentID)
		db.Model(models.Like{}).Delete(&currentUser)

		toggleResult = ToggleInfo{
			IsLiked:       false,
			MessageStatus: "unliked post successfully",
		}
	} else {
		currentUser = models.Like{
			Model:    gorm.Model{},
			ParentID: parentID,
			UserID:   userID,
		}
		db.Model(models.Like{}).Create(&currentUser)

		toggleResult = ToggleInfo{
			IsLiked:       true,
			MessageStatus: "liked post successfully",
		}
	}

	return toggleResult, nil
}

func SearchUserByUsername(db *gorm.DB, username string) ([]models.User, error) {
	var users []models.User
	result := db.Where("Username LIKE ?", username).First(&users)
	if result.RowsAffected == 0 {
		return nil, errors.New(constants.ErrNoUser)
	}
	return users, nil
}

func SearchPostsByKeywords(db *gorm.DB, keyword string) ([]models.Post, error) {
	var posts []models.Post
	result := db.Where("Body ILIKE ?", "%"+keyword+"%").Find(&posts)
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf(constants.ErrNoPost+" keyword used: %s", keyword)
	}

	return posts, nil
}

func GetAllPosts(db *gorm.DB) ([]models.Post, error) {
	var posts []models.Post
	// TODO: SCALE THIS
	result := db.Find(&posts)
	if result.Error != nil {
		return nil, fmt.Errorf("internal server error: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, errors.New(constants.ErrNoPost)
	}

	return posts, nil
}

func GetAllPostsByUsername(db *gorm.DB, username string) ([]models.Post, error) {
	var posts []models.Post
	var user models.User

	db.Where("username = ?", username).First(&user)
	result := db.Where("user_id = ?", user.ID).Find(&posts)
	if result.Error != nil {
		return nil, fmt.Errorf("internal server error: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, errors.New(constants.ErrNoPost)
	}

	return posts, nil
}

func CreatePost(db *gorm.DB, userID uint, parentID *uint, quote *string, body string) error {
	if !userExists(db, userID) {
		return errors.New(constants.ErrNoUser)
	}
	post := models.Post{
		UserID:   userID,
		ParentID: parentID,
		Quote:    quote,
		Body:     body,
	}

	return db.Create(&post).Error
}

// AUX.

type ToggleInfo struct {
	IsLiked       bool
	MessageStatus string
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

//	func ValidateCredentials(db *gorm.DB, inputUser, password string) bool {
//		var user models.User
//
//		field := "Mail"
//		if !IsEmail(inputUser) {
//			field = "Username"
//		}
//		err := queryUserByField(db, field, inputUser, password, &user)
//		if err != nil {
//			if errors.Is(err, gorm.ErrRecordNotFound) {
//				return false
//			}
//			log.Printf("Error querying user by %s: %v", field, err)
//			return false
//		}
//
//		return true
//	}
func IsEmail(email string) bool {
	re := regexp.MustCompile(constants.EmailRegexPatterns)
	return re.MatchString(email)
}

func GetUserByID(db *gorm.DB, userID uint) (models.User, error) {
	var user models.User
	err := db.First(&user, userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.New(constants.ErrNoUser)
		}
		return user, errors.New("failed to retrieve the user from the database")
	}

	return user, nil
}

func GetPostByID(db *gorm.DB, postID uint) (models.Post, error) {
	var post models.Post
	err := db.First(&post, postID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return post, errors.New(constants.ErrNoPost)
		}
		return post, errors.New("failed to retrieve the post from the database")
	}

	return post, nil
}

func UpdateProfile(db *gorm.DB, user *models.User) error {
	hashedPassword, hasedPasswordErr := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if hasedPasswordErr != nil {
		return errors.New("failed to hash password")
	}
	user.Password = string(hashedPassword)
	return db.Save(user).Error
}

func GetFollowers(db *gorm.DB, userID uint) ([]models.User, error) {
	var followers []models.User
	result := db.Table("users").
		Select("users.*").
		Joins("JOIN follows ON users.id = follows.following_user_id").
		Where("followed_user_id = ?", userID).
		Find(&followers)

	if result.Error != nil {
		return nil, fmt.Errorf("internal server error: %w", result.Error)
	}

	return followers, nil
}

func GetFollowing(db *gorm.DB, u uint) ([]models.User, error) {
	var following []models.User
	result := db.Table("users").
		Select("users.*").
		Joins("JOIN follows ON users.id = follows.followed_user_id").
		Where("following_user_id = ?", u).
		Find(&following)

	if result.Error != nil {
		return nil, fmt.Errorf("internal server error: %w", result.Error)
	}

	return following, nil
}

func IsPostOwner(db *gorm.DB, userID, postID uint) bool {
	var post models.Post
	err := db.Where("id = ? AND user_id = ?", postID, userID).First(&post).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	} else if err != nil {
		log.Printf("Error querying post by id: %v", err)
		return false
	}

	return true
}

//	func queryUserByField(db *gorm.DB, field, value, password string, user *models.User) error {
//		return db.Where(fmt.Sprintf("%s = ? AND Password = ?", field), value, password).First(user).Error
//	}

func alreadyFollows(db *gorm.DB, followingUserID, followedUserID uint) bool {
	var follow models.Follow
	result := db.Where("following_user_id = ? AND followed_user_id = ?", followingUserID, followedUserID).First(&follow)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false
		}
		log.Printf("Error querying database: %v", result.Error)
		return false
	}
	return true
}

func isLiked(db *gorm.DB, userID, parentID uint) bool {
	result := db.Model(models.Like{}).Where("user_id = ? AND parent_id = ?", userID, parentID).First(&models.Like{})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false
		}
		log.Printf("Error querying database: %v", result.Error)
		return false
	}
	if result.RowsAffected == 0 {
		return false
	}

	return true
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

func FindOrCreateConversation(db *gorm.DB, currentSenderID, currentReceiverID uint) (*models.Conversation, error) {
	var convo models.Conversation
	err := db.Where("sender_id = ? AND receiver_id = ?", currentSenderID, currentReceiverID).First(&convo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			convo = models.Conversation{
				SenderID:   currentSenderID,
				ReceiverID: currentReceiverID,
			}
			if createErr := db.Create(&convo).Error; createErr != nil {
				return nil, createErr
			}
			return &convo, nil
		}
		return nil, err
	}
	return &convo, nil
}

func SendMessage(db *gorm.DB, currentSenderID, currentReceiverID uint, content string) error {
	if currentSenderID == 0 || currentReceiverID == 0 {
		return errors.New("invalid sender or receiver ID")
	}
	if content == constants.Empty {
		return errors.New("message content cannot be empty")
	}
	convo, err := FindOrCreateConversation(db, currentSenderID, currentReceiverID)
	if err != nil {
		return err
	}
	message := models.Message{
		ConversationID: convo.ID,
		SenderID:       currentSenderID,
		Content:        content,
	}
	return db.Create(&message).Error
}
