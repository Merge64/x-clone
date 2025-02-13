package controllers

import (
	"main/constants"
	"main/models"
)

var UserSignUpEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURLAuth + "signup",
	HandlerFunction: SignUpHandler,
}

var UserLoginEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURLAuth + "login",
	HandlerFunction: LoginHandler,
}

var UserLogoutEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURLAuth + "logout",
	HandlerFunction: LogoutHandler,
}

//TODO User Profile Endpoints

var ViewUserProfileEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURL + "profile/:userid",
	HandlerFunction: ViewUserProfileHandler,
}

var GetFollowingProfileEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURL + "profile/following/user/:userid",
	HandlerFunction: GetFollowingProfileHandler,
}

var GetFollowersProfileEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURL + "profile/followers/user/:userid",
	HandlerFunction: GetFollowersProfileHandler,
}

//TODO Post Endpoints

var GetAllPostsEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURL + "posts/all",
	HandlerFunction: GetAllPostsHandler,
}

var CreatePostEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURL + "posts/create",
	HandlerFunction: CreatePostHandler,
}

var GetAllPostsByUserIDEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURL + "/profile/:userid/posts",
	HandlerFunction: GetPostsByUserIDHandler,
}

var GetSpecificPostEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURL + "posts/:postid",
	HandlerFunction: GetSpecificPostHandler,
}

var EditPostEndpoint = models.Endpoint{
	Method:          models.PUT,
	Path:            constants.InitialURL + "posts/:postid/edit",
	HandlerFunction: EditPostHandler,
}

var DeletePostEndpoint = models.Endpoint{
	Method:          models.DELETE,
	Path:            constants.InitialURL + "posts/:postid/delete",
	HandlerFunction: DeletePostHandler,
}

var CreateRepostEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURL + "/posts/:parentid/repost",
	HandlerFunction: CreateRepostHandler,
}

var ToggleLikeEndPoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURL + "togglelike/:postid",
	HandlerFunction: ToggleLikeHandler,
}

var EditUserProfileEndpoint = models.Endpoint{
	Method:          models.PUT,
	Path:            constants.InitialURL + "profile/:userid/edit",
	HandlerFunction: EditUserProfileHandler,
}

//TODO Search Endpoints

var SearchPostEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURL + "posts",
	HandlerFunction: SearchPostHandler,
}

//TODO SearchLatestPost

var SearchUserEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURL + "search/:username",
	HandlerFunction: SearchUserHandler,
}

//TODO Direct Messaging

var ListConversationsEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURL + "conversations/listall",
	HandlerFunction: ListConversationsHandler,
}

var GetConversationMessagesEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURL + "conversations/:conversationID/messages",
	HandlerFunction: GetMessagesForConversationHandler,
}

var SendDirectMessageEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURL + "conversations/users/:userid/message",
	HandlerFunction: SendMessageHandler,
}

var FollowUserEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURL + "follow/:userid",
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
