package user

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"main/constants"
	"main/mappers"
	"main/models"
	"regexp"
)

func FollowAccount(db *gorm.DB, followingUsername, followedUsername string) error {
	// 1. Check if the user is already following
	var existing models.Follow
	if err := db.Where("following_username = ? AND followed_username = ?",
		followingUsername, followedUsername).
		First(&existing).Error; err == nil {
		return errors.New("already following this user")
	}

	// 2. Create a new Follow record
	follow := models.Follow{
		FollowingUsername: followingUsername,
		FollowedUsername:  followedUsername,
	}
	if err := db.Create(&follow).Error; err != nil {
		return err // This error triggers "Failed to follow user"
	}

	return nil
}

func UnfollowAccount(db *gorm.DB, followingUsername, followedUsername string) error {
	if followingUsername == followedUsername {
		return errors.New("invalid ID: user cannot unfollow themselves")
	}

	result := db.Where("following_username = ? AND followed_username = ?", followingUsername, followedUsername).
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

// IsLiked Like-specific functions.
func IsLiked(db *gorm.DB, userID uint, postID uint) bool {
	var count int64
	db.Model(&models.Like{}).Where("user_id = ? AND post_id = ?", userID, postID).Count(&count)
	return count > 0
}

func ToggleLike(db *gorm.DB, userID uint, postID uint) (ToggleInfo, error) {
	if !userExists(db, userID) {
		return ToggleInfo{}, errors.New(constants.ErrNoUser)
	}

	var toggleResult ToggleInfo

	if IsLiked(db, userID, postID) {
		db.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&models.Like{})
		db.Model(&models.Post{}).Where("id = ?", postID).Update("likes_count", gorm.Expr("likes_count - 1"))

		toggleResult = ToggleInfo{
			IsLiked:       false,
			MessageStatus: "unliked successfully",
		}
	} else {
		newLike := models.Like{PostID: postID, UserID: userID}
		db.Create(&newLike)
		db.Model(&models.Post{}).Where("id = ?", postID).Update("likes_count", gorm.Expr("likes_count + 1"))

		toggleResult = ToggleInfo{
			IsLiked:       true,
			MessageStatus: "liked successfully",
		}
	}

	return toggleResult, nil
}

// searchPostsByKeywords is a helper.
func searchPostsByKeywords(db *gorm.DB, keyword, orderBy string) ([]mappers.PostResponse, error) {
	var rawPosts []models.Post
	var result *gorm.DB

	// Start building the query with Preload to fetch ParentPost
	query := db.Preload("ParentPost")

	switch {
	case keyword == constants.Empty:
		// If no keyword is provided, fetch all posts
		if orderBy != constants.Empty {
			query = query.Order(orderBy)
		}
		result = query.Find(&rawPosts)

	case len(keyword) < constants.SearchedWordLen:
		// If the keyword is too short, use regex search
		queryPattern := fmt.Sprintf("\\m%s\\M", keyword)
		q := query.Where("body ~* ?", queryPattern)
		if orderBy != constants.Empty {
			q = q.Order(orderBy)
		}
		result = q.Find(&rawPosts)

	default:
		// For longer keywords, use case-insensitive search
		queryPattern := "%" + keyword + "%"
		q := query.Where("body ILIKE ?", queryPattern)
		if orderBy != constants.Empty {
			q = q.Order(orderBy)
		}
		result = q.Find(&rawPosts)
	}

	// Handle errors
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf(constants.ErrNoPost+" keyword used: %s", keyword)
	}

	// Process the raw posts into the desired response format
	processedPosts := ProcessPosts(rawPosts)

	return processedPosts, nil
}

func SearchPostsByKeywords(db *gorm.DB, keyword string) ([]mappers.PostResponse, error) {
	return searchPostsByKeywords(db, keyword, "likes_count DESC")
}

func SearchPostsByKeywordsSortedByLatest(db *gorm.DB, keyword string) ([]mappers.PostResponse, error) {
	return searchPostsByKeywords(db, keyword, "created_at DESC")
}

func SearchUsersByUsername(db *gorm.DB, username string) ([]mappers.Response, error) {
	var users []models.User
	result := db.Table("users").
		Select(`
            users.id, 
			users.nickname,
            users.username,
            users.mail,
            users.password,
            users.location,
            COUNT(follows.id) AS follower_count,
            CASE WHEN LOWER(users.username) = LOWER(?) THEN 0 ELSE 1 END AS priority
        `, username).
		Joins("LEFT JOIN follows ON follows.followed_username = users.username").
		Where("users.username ILIKE ?", "%"+username+"%").
		Group("users.id, users.nickname, users.username, users.mail, users.password, users.location, priority").
		Order("priority ASC, follower_count DESC").
		Scan(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	if len(users) == 0 {
		return nil, errors.New("no users found")
	}

	return mappers.MapUsersToResponses(users), nil
}

func SearchUniqueMailUsername(db *gorm.DB, field string, value string) (bool, error) {
	var count int64
	err := db.Model(&models.User{}).Where("LOWER("+field+") = LOWER(?)", value).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func GetAllPosts(db *gorm.DB) ([]models.Post, error) {
	var posts []models.Post

	// Ensure ParentPost is loaded to support reposts
	result := db.Preload("ParentPost").Order("created_at desc").Find(&posts)
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	if result.Error != nil {
		return nil, fmt.Errorf("internal server error: %w", result.Error)
	}

	return posts, nil
}

func GetAllPostsByUsername(db *gorm.DB, username string) ([]models.Post, error) {
	var posts []models.Post
	var user models.User

	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with username %s not found", username)
		}
		return nil, fmt.Errorf("internal server error: %w", err)
	}

	result := db.Preload("ParentPost").Where("user_id = ?", user.ID).Order("created_at desc").Find(&posts)
	if result.Error != nil {
		return nil, fmt.Errorf("internal server error: %w", result.Error)
	}

	return posts, nil
}

func GetAllRepliesByUsername(db *gorm.DB, username string) ([]models.Post, error) {
	var posts []models.Post
	var user models.User

	// Fetch user by username
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with username %s not found", username)
		}
		return nil, fmt.Errorf("internal server error: %w", err)
	}

	// Fetch posts where IsRepost is false and ParentID is not nil
	result := db.Preload("ParentPost").
		Where("user_id = ?", user.ID).
		Where("is_repost = ?", false).  // Ensure it's not a repost
		Where("parent_id IS NOT NULL"). // Ensure ParentID is not nil
		Order("created_at desc").
		Find(&posts)

	if result.Error != nil {
		return nil, fmt.Errorf("internal server error: %w", result.Error)
	}

	return posts, nil
}

func PostsWLikesByUsername(db *gorm.DB, username string) ([]models.Post, error) {
	var posts []models.Post
	var user models.User

	// Fetch user by username
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with username %s not found", username)
		}
		return nil, fmt.Errorf("internal server error: %w", err)
	}

	// Fetch the posts liked by the user
	result := db.Joins("JOIN likes ON likes.post_id = posts.id").
		Preload("ParentPost").
		Where("likes.user_id = ?", user.ID).
		Order("posts.created_at desc").
		Find(&posts)

	if result.Error != nil {
		return nil, fmt.Errorf("internal server error: %w", result.Error)
	}

	return posts, nil
}

func CreatePost(db *gorm.DB,
	userID uint,
	nickname string,
	parentID *uint,
	username string,
	quote *string,
	body string,
	isRepost bool) (*models.Post, error) {
	if !userExists(db, userID) {
		return nil, errors.New(constants.ErrNoUser)
	}

	post := models.Post{
		UserID:   userID,
		Username: username,
		Nickname: nickname,
		ParentID: parentID,
		Quote:    quote,
		Body:     body,
		IsRepost: isRepost,
	}

	if err := db.Create(&post).Error; err != nil {
		return nil, err
	}
	// Ensure ID is assigned after creation
	if post.ID == 0 {
		return nil, errors.New("failed to create post: ID not assigned")
	}
	return &post, nil
}

// AUX.

type ToggleInfo struct {
	IsLiked       bool
	IsReposted    bool
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

func IsEmail(email string) bool {
	re := regexp.MustCompile(constants.EmailRegexPatterns)
	return re.MatchString(email)
}

func GetUserProfileByUsername(db *gorm.DB, username string) (*mappers.Response, error) {
	var user models.User

	result := db.Table("users").
		Select(`
            users.id, 
            users.nickname,
            users.username,
            users.created_at,
            COUNT(follows.id) AS follower_count
        `).
		Joins("LEFT JOIN follows ON follows.followed_username = users.username").
		Where("LOWER(users.username) = LOWER(?)", username).
		Group("users.id, users.nickname, users.username, users.created_at").
		First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	response := mappers.MapUserToResponse(user)
	return &response, nil
}

func GetSimplePostByID(db *gorm.DB, postID uint) (models.Post, error) {
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

func UpdateProfile(db *gorm.DB, username string, user *models.User) error {
	return db.Where("username = ?", username).Updates(user).Error
}

func GetFollowers(db *gorm.DB, username string) ([]mappers.Response, error) {
	var followers []models.User
	currentUser, getUserErr := GetUserProfileByUsername(db, username)
	if getUserErr != nil {
		return nil, getUserErr
	}

	result := db.Table("users").
		Select(`
            users.id, 
            users.nickname,
            users.username,
            users.created_at,
            COUNT(follows.id) AS follower_count
        `).
		Joins("JOIN follows ON users.username = follows.following_username").
		Where("followed_username = ?", currentUser.Username).
		Group("users.id, users.nickname, users.username, users.created_at").
		Find(&followers)

	if result.Error != nil {
		return nil, fmt.Errorf("internal server error: %w", result.Error)
	}

	return mappers.MapUsersToResponses(followers), nil
}

func GetFollowing(db *gorm.DB, username string) ([]mappers.Response, error) {
	var following []models.User
	currentUser, getUserErr := GetUserProfileByUsername(db, username)
	if getUserErr != nil {
		return nil, getUserErr
	}

	result := db.Table("users").
		Select(`
            users.id, 
            users.nickname,
            users.username,
            users.created_at,
            COUNT(follows.id) AS follower_count
        `).
		Joins("JOIN follows ON users.username = follows.followed_username").
		Where("following_username = ?", currentUser.Username).
		Group("users.id, users.nickname, users.username, users.created_at").
		Find(&following)

	if result.Error != nil {
		return nil, fmt.Errorf("internal server error: %w", result.Error)
	}

	return mappers.MapUsersToResponses(following), nil
}

func GetUsernameIDFromContext(c *gin.Context) (string, error) {
	username, exists := c.Get("username")
	if !exists {
		return "", errors.New("unauthorized")
	}
	usernameStr, ok := username.(string)
	if !ok {
		return "", errors.New("invalid username type")
	}

	return usernameStr, nil
}

func GetNicknameFromContext(c *gin.Context) (string, error) {
	nickname, exists := c.Get("nickname")
	if !exists {
		return "", errors.New("unauthorized")
	}
	nicknameStr, ok := nickname.(string)
	if !ok {
		return "", errors.New("invalid nickname type")
	}

	return nicknameStr, nil
}

func GetUserIDFromContext(c *gin.Context) (uint, error) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("unauthorized")
	}

	userID, ok := userIDStr.(uint)
	if !ok {
		return 0, errors.New("invalid user ID")
	}

	return userID, nil
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

func FindOrCreateConversation(db *gorm.DB, currentSenderUsername, currentReceiverUsername string) (*models.
	Conversation, error) {
	var convo models.Conversation

	// Try to find an existing conversation
	err := findExistingConversation(db, currentSenderUsername, currentReceiverUsername, &convo)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// If conversation not found, create a new one
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return createNewConversation(db, currentSenderUsername, currentReceiverUsername)
	}

	// Fetch nicknames for an existing conversation
	if errUpdateNickname := updateConversationNicknames(db, &convo); errUpdateNickname != nil {
		return nil, errUpdateNickname
	}

	return &convo, nil
}

func findExistingConversation(db *gorm.DB, sender, receiver string, convo *models.Conversation) error {
	return db.Where(
		"(sender_username = ? AND receiver_username = ?) OR (sender_username = ? AND receiver_username = ?)",
		sender, receiver, receiver, sender,
	).First(convo).Error
}

func createNewConversation(db *gorm.DB, sender, receiver string) (*models.Conversation, error) {
	senderNickname, err := getNickname(db, sender)
	if err != nil {
		return nil, err
	}

	receiverNickname, err := getNickname(db, receiver)
	if err != nil {
		return nil, err
	}

	convo := models.Conversation{
		SenderUsername:   sender,
		SenderNickname:   senderNickname,
		ReceiverUsername: receiver,
		ReceiverNickname: receiverNickname,
	}

	if errDB := db.Create(&convo).Error; errDB != nil {
		return nil, errDB
	}

	return &convo, nil
}

func getNickname(db *gorm.DB, username string) (string, error) {
	var nickname string
	err := db.Model(&models.User{}).Where("username = ?", username).Pluck("nickname", &nickname).Error
	return nickname, err
}

func updateConversationNicknames(db *gorm.DB, convo *models.Conversation) error {
	if err := db.Model(&models.User{}).
		Where("username = ?", convo.SenderUsername).
		Pluck("nickname", &convo.SenderNickname).Error; err != nil {
		return err
	}

	if err := db.Model(&models.User{}).
		Where("username = ?", convo.ReceiverUsername).
		Pluck("nickname", &convo.ReceiverNickname).Error; err != nil {
		return err
	}

	return nil
}

func SendMessage(db *gorm.DB, currentSenderUsername, currentReceiverUsername string, content string) error {
	if currentSenderUsername == constants.Empty || currentReceiverUsername == constants.Empty {
		return errors.New("invalid sender or receiver username")
	}
	if content == constants.Empty {
		return errors.New("message content cannot be empty")
	}
	convo, err := FindOrCreateConversation(db, currentSenderUsername, currentReceiverUsername)
	if err != nil {
		return err
	}
	message := models.Message{
		ConversationID: convo.ID,
		SenderUsername: currentSenderUsername,
		Content:        content,
	}
	return db.Create(&message).Error
}

func ProcessPosts(rawPosts []models.Post) []mappers.PostResponse {
	processedPosts := make([]mappers.PostResponse, len(rawPosts))

	for i, post := range rawPosts {
		processedPosts[i] = mappers.ProcessPost(post)
	}

	return processedPosts
}

func ProcessPost(post models.Post) mappers.PostResponse {
	var parentPost *mappers.ParentPostResponse
	if post.ParentPost != nil {
		parentPost = &mappers.ParentPostResponse{
			ID:        post.ParentPost.ID,
			CreatedAt: post.ParentPost.CreatedAt.Format("2006-01-02 15:04:05.999999999 -0700 MST"),
			Username:  post.ParentPost.Username,
			Nickname:  post.ParentPost.Nickname,
			Body:      post.ParentPost.Body,
		}
	}

	return mappers.PostResponse{
		ID:           post.ID,
		CreatedAt:    post.CreatedAt.Format("2006-01-02 15:04:05.999999999 -0700 MST"),
		UserID:       post.UserID,
		Nickname:     post.Nickname,
		Username:     post.Username,
		ParentID:     post.ParentID,
		Quote:        post.Quote,
		Body:         post.Body,
		RepostsCount: post.RepostsCount,
		LikesCount:   post.LikesCount,
		IsRepost:     post.IsRepost,
		ParentPost:   parentPost,
	}
}

func EnlistUsers(arrayOfUsers []mappers.Response) []string {
	var usersList []string

	for _, currentUser := range arrayOfUsers {
		usersList = append(usersList, currentUser.Username)
	}

	return usersList
}

func IsFollowing(db *gorm.DB, followedUsername, currentUsername string) (bool, error) {
	var follow models.Follow
	db.Where("following_username = ? AND followed_username = ?", currentUsername, followedUsername).First(&follow)

	if follow.ID == 0 {
		return false, errors.New("not following user")
	}

	return true, nil
}

func GetMissingUserFields(db *gorm.DB, username string, user *models.User) {
	var aux models.User
	db.Where("username = ?", username).First(&aux)
	user.Username = aux.Username
	user.Mail = aux.Mail
	user.FollowerCount = aux.FollowerCount
	user.Password = aux.Password
}

func UpdateNicknamePosts(db *gorm.DB, username, nickname string) error {
	rawPosts, errDB := GetAllPostsByUsername(db, username)
	if errDB != nil {
		return errDB
	}
	for _, post := range rawPosts {
		post.Nickname = nickname
		db.Updates(&post)
	}
	return nil
}
