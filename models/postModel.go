package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	User_id int
	Title   string
	Body    string
	Likes   int
}

type Post_likes struct {
	gorm.Model
	Post_id int // id del post (gorm)
	User_id int
}

type Post_comment struct {
	gorm.Model
	Post_id      int // id del post (gorm)
	User_id      int
	Comment_body string
}
type Comment_comments struct {
	gorm.Model
	Comment_id int
	User_id    int
	Body       string
}
type Comment_likes struct {
	gorm.Model
	Comment_id int
	User_id    int
}
