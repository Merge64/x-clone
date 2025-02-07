package controllers

import (
	"main/constants"
	"main/models"
)

var GetAllPostsByUserIDEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "/profile/:userid/posts",
	HandlerFunction: GetPostsByUserIDHandler,
}

var GetAllPostsEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "posts/all",
	HandlerFunction: GetAllPostsHandler,
}

var GetSpecificPostEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "posts/:postid",
	HandlerFunction: GetSpecificPostHandler,
}

var DeletePostEndpoint = models.Endpoint{
	Method:          models.DELETE,
	Path:            constants.BASEURL + "posts/:postid/delete",
	HandlerFunction: DeletePostHandler,
}

var EditPostEndpoint = models.Endpoint{
	Method:          models.PUT,
	Path:            constants.BASEURL + "posts/:postid/edit",
	HandlerFunction: EditPostHandler,
}

var CreatePostEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "posts/create",
	HandlerFunction: CreatePostHandler,
}

var CreateRepostEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "/posts/:parentid/repost",
	HandlerFunction: CreateRepostHandler,
}

var GetFollowersProfileEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "profile/followers/user/:userid",
	HandlerFunction: GetFollowersProfileHandler,
}

var GetFollowingProfileEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "profile/following/user/:userid",
	HandlerFunction: GetFollowingProfileHandler,
}

var ViewUserProfileEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "profile/:userid",
	HandlerFunction: ViewUserProfileHandler,
}

var EditUserProfileEndpoint = models.Endpoint{
	Method:          models.PUT,
	Path:            constants.BASEURL + "profile/:userid/edit",
	HandlerFunction: EditUserProfileHandler,
}

var SearchUserEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "search/:username",
	HandlerFunction: SearchUserHandler,
}

var SearchPostEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "posts",
	HandlerFunction: SearchPostHandler,
}

var ToggleLikeEndPoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "togglelike/:postid",
	HandlerFunction: ToggleLikeHandler,
}

var FollowUserEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "follow/:userid",
	HandlerFunction: FollowUserHandler,
}

var UnfollowUserEndpoint = models.Endpoint{
	Method:          models.DELETE,
	Path:            constants.BASEURL + "unfollow/:userid",
	HandlerFunction: UnfollowUserHandler,
}

var UserSignUpEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "signup",
	HandlerFunction: SignUpHandler,
}

var UserLoginEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "login",
	HandlerFunction: LoginHandler,
}

var UserLogoutEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "logout",
	HandlerFunction: LogoutHandler,
}

var SendDirectMessageEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.BASEURL + "conversations/users/:userid/message",
	HandlerFunction: SendMessageHandler,
}

var ListConversationsEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "conversations/listall",
	HandlerFunction: ListConversationsHandler,
}

var GetConversationMessagesEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.BASEURL + "conversations/:conversationID/messages",
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
