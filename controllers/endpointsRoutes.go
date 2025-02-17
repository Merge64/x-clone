package controllers

import (
	"main/constants"
	"main/models"
)

// TODO: User Endpoints - Add Nickname for user

var UserSignUpEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURLAuth + "/signup",
	HandlerFunction: SignUpHandler,
}

var UserLoginEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURLAuth + "/login",
	HandlerFunction: LoginHandler,
}

var UserLogoutEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURLAuth + "/logout",
	HandlerFunction: LogoutHandler,
}

// TODO User Profile Endpoints

var ViewUserProfileEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURLProfile + "/:username",
	HandlerFunction: ViewUserProfileHandler,
}

var EditUserProfileEndpoint = models.Endpoint{
	Method:          models.PUT,
	Path:            constants.InitialURLProfile + "/edit",
	HandlerFunction: EditUserProfileHandler,
}

var GetFollowingProfileEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURLProfile + "/:username/following",
	HandlerFunction: GetFollowingProfileHandler,
}

var GetFollowersProfileEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURLProfile + "/:username/followers",
	HandlerFunction: GetFollowersProfileHandler,
}

// TODO Post Endpoints

var GetAllPostsEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURLPosts,
	HandlerFunction: GetAllPostsHandler,
}

var CreatePostEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURLPosts + "/create",
	HandlerFunction: CreatePostHandler,
}

var CreateRepostEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURLPosts + "/:postid/repost",
	HandlerFunction: CreateRepostHandler,
}

var GetAllPostsByUserIDEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURLPosts + "/:username",
	HandlerFunction: GetPostsByUsernameHandler,
}

var GetSpecificPostEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURLPosts + "/:username/:postid",
	HandlerFunction: GetSpecificPostHandler,
}

var EditPostEndpoint = models.Endpoint{
	Method:          models.PUT,
	Path:            constants.InitialURLPosts + "/:postid/edit",
	HandlerFunction: EditPostHandler,
}

var DeletePostEndpoint = models.Endpoint{
	Method:          models.DELETE,
	Path:            constants.InitialURLPosts + "/:postid/delete",
	HandlerFunction: DeletePostHandler,
}

var ToggleLikeEndPoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURLPosts + "/:postid/like",
	HandlerFunction: ToggleLikeHandler,
}

// TODO Search Endpoints

// SearchEndpoint GET /search?q=keyword
// SearchEndpoint GET /search?q=keyword&f=user
// SearchEndpoint GET /search?q=keyword&f=latest

var SearchEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURLSearch,
	HandlerFunction: SearchHandler,
}

// TODO Direct Messaging

var ListConversationsEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURLDms,
	HandlerFunction: ListConversationsHandler,
}

var GetConversationMessagesEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURLDms + "/:receiverID/:senderID",
	HandlerFunction: GetMessagesForConversationHandler,
}

var SendDirectMessageEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURLDms + "/dm",
	HandlerFunction: SendMessageHandler,
}

var FollowUserEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURLDms + "follow/:userid",
	HandlerFunction: FollowUserHandler,
}

var UnfollowUserEndpoint = models.Endpoint{
	Method:          models.DELETE,
	Path:            constants.InitialURL + "unfollow/:userid",
	HandlerFunction: UnfollowUserHandler,
}

var PublicEndpoints = []models.Endpoint{
	UserSignUpEndpoint,
	UserLoginEndpoint,
	ViewUserProfileEndpoint,
	SearchEndpoint,
	GetSpecificPostEndpoint,
	GetAllPostsByUserIDEndpoint,
	GetAllPostsEndpoint,
}

var PrivateEndpoints = []models.Endpoint{
	FollowUserEndpoint,
	UnfollowUserEndpoint,
	GetFollowersProfileEndpoint,
	GetFollowingProfileEndpoint,
	EditUserProfileEndpoint,
	CreatePostEndpoint,
	EditPostEndpoint,
	DeletePostEndpoint,
	ToggleLikeEndPoint,
	SendDirectMessageEndpoint,
	ListConversationsEndpoint,
	GetConversationMessagesEndpoint,
	UserLogoutEndpoint,
	CreateRepostEndpoint,
}
