package db

import (
	"fmt"
	"time"

	"github.com/tyzual/tisq/util"
)

/*
Comment comment表对应的数据结构
*/
type Comment struct {
	CommentID  string
	ArticleID  string
	ArticleKey string
	UserID     string
	Content    string
	TimeStamp  time.Time
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
	id := util.MD5([]byte(email))
	user.UserID = id
	user.Email = email
	user.DisplayName = displayName
	user.WebSite = webSite
	user.Avatar = avatar
	return &user
}

/*
NewComment 创建一个新评论的数据结构
*/
func NewComment(m *Mysql, articleKey, userEmail, content string) *Comment {
	user := m.GetUserByEmail(userEmail)
	if user == nil {
		util.LogWarn(fmt.Sprintf("没找到用户%v", userEmail))
		return nil
	}

	comm := Comment{}
	comm.UserID = user.UserID
	comm.TimeStamp = time.Now()
	comm.ArticleID = util.MD5([]byte(articleKey))
	comm.ArticleKey = articleKey
	comm.Content = content
	return &comm
}
