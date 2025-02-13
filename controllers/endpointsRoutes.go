package controllers

import (
	"main/constants"
	"main/models"
)

var UserSignUpEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialUrlAuth + "signup",
	HandlerFunction: SignUpHandler,
}

var UserLoginEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialUrlAuth + "login",
	HandlerFunction: LoginHandler,
}

var UserLogoutEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialUrlAuth + "logout",
	HandlerFunction: LogoutHandler,
}

var GetAllPostsByUserIDEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialUrl + "/profile/:userid/posts",
	HandlerFunction: GetPostsByUserIDHandler,
}

var GetAllPostsEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialUrl + "posts/all",
	HandlerFunction: GetAllPostsHandler,
}

var GetSpecificPostEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialUrl + "posts/:postid",
	HandlerFunction: GetSpecificPostHandler,
}

var DeletePostEndpoint = models.Endpoint{
	Method:          models.DELETE,
	Path:            constants.InitialUrl + "posts/:postid/delete",
	HandlerFunction: DeletePostHandler,
}

var EditPostEndpoint = models.Endpoint{
	Method:          models.PUT,
	Path:            constants.InitialUrl + "posts/:postid/edit",
	HandlerFunction: EditPostHandler,
}

var CreatePostEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialUrl + "posts/create",
	HandlerFunction: CreatePostHandler,
}

var CreateRepostEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialUrl + "/posts/:parentid/repost",
	HandlerFunction: CreateRepostHandler,
}

var GetFollowersProfileEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialUrl + "profile/followers/user/:userid",
	HandlerFunction: GetFollowersProfileHandler,
}

var GetFollowingProfileEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialUrl + "profile/following/user/:userid",
	HandlerFunction: GetFollowingProfileHandler,
}

var ViewUserProfileEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialUrl + "profile/:userid",
	HandlerFunction: ViewUserProfileHandler,
}

var EditUserProfileEndpoint = models.Endpoint{
	Method:          models.PUT,
	Path:            constants.InitialUrl + "profile/:userid/edit",
	HandlerFunction: EditUserProfileHandler,
}

var SearchUserEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialUrl + "search/:username",
	HandlerFunction: SearchUserHandler,
}

var SearchPostEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialUrl + "posts",
	HandlerFunction: SearchPostHandler,
}

var ToggleLikeEndPoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialUrl + "togglelike/:postid",
	HandlerFunction: ToggleLikeHandler,
}

var FollowUserEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialUrl + "follow/:userid",
	HandlerFunction: FollowUserHandler,
}

var UnfollowUserEndpoint = models.Endpoint{
	Method:          models.DELETE,
	Path:            constants.InitialUrl + "unfollow/:userid",
	HandlerFunction: UnfollowUserHandler,
}

var SendDirectMessageEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialUrl + "conversations/users/:userid/message",
	HandlerFunction: SendMessageHandler,
}

var ListConversationsEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialUrl + "conversations/listall",
	HandlerFunction: ListConversationsHandler,
}

var GetConversationMessagesEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialUrl + "conversations/:conversationID/messages",
	HandlerFunction: GetMessagesForConversationHandler,
}

var PublicEndpoints = []models.Endpoint{
	UserSignUpEndpoint,
	UserLoginEndpoint,
	ViewUserProfileEndpoint,
	SearchUserEndpoint,
	SearchPostEndpoint,
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
