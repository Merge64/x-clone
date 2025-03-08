package mappers

import "main/models"

type Response struct {
	ID            uint   `json:"id"`
	Username      string `json:"username"`
	Nickname      string `json:"nickname"`
	FollowerCount uint   `json:"follower_count"`
}

// MapUserToResponse converts ONLY ONE models.User to a Response.
func MapUserToResponse(u models.User) Response {
	return Response{
		ID:            u.ID,
		Nickname:      u.Nickname,
		Username:      u.Username,
		FollowerCount: u.FollowerCount,
	}
}

func MapUsersToResponses(users []models.User) []Response {
	responses := make([]Response, len(users))
	for i, u := range users {
		responses[i] = MapUserToResponse(u)
	}
	return responses
}

func MapPostsToResponses(posts []models.Post) []PostResponse {
	responses := make([]PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = ProcessPost(post)
	}
	return responses
}

type PostResponse struct {
	ID           uint                `json:"id"`
	CreatedAt    string              `json:"created_at"`
	UserID       uint                `json:"userid"`
	Nickname     string              `json:"nickname"`
	Username     string              `json:"username"`
	ParentID     *uint               `json:"parent_id"`
	Quote        *string             `json:"quote"`
	Body         string              `json:"body"`
	RepostsCount uint                `json:"reposts_count"`
	LikesCount   uint                `json:"likes_count"`
	IsRepost     bool                `json:"is_repost"`
	ParentPost   *ParentPostResponse `json:"parent_post,omitempty"`
}

type ParentPostResponse struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"created_at"`
	Username  string `json:"username"`
	Nickname  string `json:"nickname"`
	Body      string `json:"body"`
}

func ProcessPost(post models.Post) PostResponse {
	var parentPost *ParentPostResponse
	if post.ParentPost != nil {
		parentPost = &ParentPostResponse{
			ID:        post.ParentPost.ID,
			CreatedAt: post.ParentPost.CreatedAt.Format("2006-01-02 15:04:05.999999999 -0700 MST"),
			Username:  post.ParentPost.Username,
			Nickname:  post.ParentPost.Nickname,
			Body:      post.ParentPost.Body,
		}
	}

	return PostResponse{
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
		ParentPost:   parentPost, // <- This ensures parent post is included
	}
}
