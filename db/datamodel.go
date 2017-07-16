package db

import (
	"crypto/md5"
	"fmt"
	"time"
)

/*
Comment comment表对应的数据结构
*/
type Comment struct {
	CommentID string
	ArticleID string
	UserID    string
	Content   string
	TimeStamp time.Time
}

/*
User user表对应的数据结构
*/
type User struct {
	UserID      string
	Email       string
	DisplayName string
	WebSite     string
	Avatar      string
}

/*
NewUser 创建一个新用户的数据结构
*/
func NewUser(email, displayName, webSite, avatar string) *User {
	user := User{}
	id := fmt.Sprintf("%x", md5.Sum([]byte(email)))
	user.UserID = id
	user.Email = email
	user.DisplayName = displayName
	user.WebSite = webSite
	user.Avatar = avatar
	return &user
}
