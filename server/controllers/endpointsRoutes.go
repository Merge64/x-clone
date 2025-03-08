package controllers

import (
	"main/authentication"
	"main/constants"
	"main/models"
)

// TODO: User Endpoints - Add Nickname for user

var UserSignUpEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURLi + "/signup",
	HandlerFunction: SignUpHandler,
}

var UserLoginEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURLi + "/login",
	HandlerFunction: LoginHandler,
}

var UserLogoutEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURLi + "/logout",
	HandlerFunction: LogoutHandler,
}

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
	Path:            constants.InitialURLPosts + "/user/:username",
	HandlerFunction: GetPostsByUsernameHandler,
}

var GetSpecificPostEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURLPosts + "/:postid",
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

// Querystring parameters
// SearchEndpoint GET /search?q=keyword
// SearchEndpoint GET /search?q=keyword&f=user
// SearchEndpoint GET /search?q=keyword&f=unique-user
// SearchEndpoint GET /search?q=keyword&f=latest

var SearchEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURLSearch,
	HandlerFunction: SearchHandler,
}

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
	Path:            constants.InitialURLProfile + "/follow/:userid",
	HandlerFunction: FollowUserHandler,
}

var UnfollowUserEndpoint = models.Endpoint{
	Method:          models.DELETE,
	Path:            constants.InitialURLProfile + "/unfollow/:userid",
	HandlerFunction: UnfollowUserHandler,
}

var IsAlreadyFollowingEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURLProfile + "/is-following/:userid",
	HandlerFunction: IsAlreadyFollowingHandler,
}

var ValidateTokenEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURLAuth + "/validate",
	HandlerFunction: authentication.ValidateHandler,
}

var ExpireTokenEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURLAuth + "/logout",
	HandlerFunction: LogoutHandler,
}

var GetUserInfoEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            "/user/info",
	HandlerFunction: GetUserInfoHandler,
}

var UpdateUsernameEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            "/user/update-username",
	HandlerFunction: UpdateUsernameHandler,
}

var CheckIfReposted = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURLPosts + "/check/:postid/reposted",
	HandlerFunction: CheckRepostedHandler,
}

var CheckIfLiked = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURLPosts + "/check/:postid/liked",
	HandlerFunction: CheckIfLikedHandler,
}

var GetCommentsEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURLPosts + "/comments/:postid",
	HandlerFunction: GetCommentsHandler,
}

var CreateCommentEndpoint = models.Endpoint{
	Method:          models.POST,
	Path:            constants.InitialURLPosts + "/comments/:postid",
	HandlerFunction: CreateCommentHandler,
}

var CountRepostsEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURLPosts + "/count/:postid",
	HandlerFunction: CountRepostsHandler,
}
var CountLikesEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURLPosts + "/count/:postid/likes",
	HandlerFunction: CountLikesHandler,
}
var CountCommentsEndpoint = models.Endpoint{
	Method:          models.GET,
	Path:            constants.InitialURLPosts + "/count/:postid/comments",
	HandlerFunction: CountCommentsHandler,
}

var PublicEndpoints = []models.Endpoint{
	UserSignUpEndpoint,
	UserLoginEndpoint,
	ViewUserProfileEndpoint,
	SearchEndpoint,
	GetSpecificPostEndpoint,
	GetAllPostsByUserIDEndpoint,
	GetAllPostsEndpoint,
	ValidateTokenEndpoint,
	ExpireTokenEndpoint,
	GetCommentsEndpoint,
	CountRepostsEndpoint,
	CountLikesEndpoint,
	CountCommentsEndpoint,
}

var PrivateEndpoints = []models.Endpoint{
	FollowUserEndpoint,
	UnfollowUserEndpoint,
	IsAlreadyFollowingEndpoint,
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
	GetUserInfoEndpoint,
	UpdateUsernameEndpoint,
	CheckIfReposted,
	CheckIfLiked,
	CreateCommentEndpoint,
}
