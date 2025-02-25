package mappers

import "main/models"

type Response struct {
	ID            uint   `json:"id"`
	Username      string `json:"username"`
	FollowerCount uint   `json:"follower_count"`
}

// MapUserToResponse converts ONLY ONE models.User to a Response.
func MapUserToResponse(u models.User) Response {
	return Response{
		ID:            u.ID,
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
