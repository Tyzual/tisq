package tdb

import (
	"fmt"
	"time"

	"database/sql"

	"github.com/tyzual/tisq/tutil"
)

/*
Comment comment表对应的数据结构
*/
type Comment struct {
	CommentID  string
	ArticleID  string
	ArticleKey string
	UserID     string
	SiteID     string
	Content    string
	TimeStamp  time.Time
	Deleted    bool
}

/*
User user表对应的数据结构
*/
type User struct {
	UserID      string
	Email       string
	DisplayName sql.NullString
	WebSite     sql.NullString
}

/*
Site site表对应的数据结构
*/
type Site struct {
	SiteID     string
	SiteDomain string
	CreateTime time.Time
}

/*
NewUser 创建一个新用户的数据结构
*/
func NewUser(email, displayName, webSite string) *User {
	user := User{}
	id := tutil.MD5([]byte(email))
	user.UserID = id
	user.Email = email
	if len(displayName) != 0 {
		user.DisplayName.Valid = true
		user.DisplayName.String = displayName
	}
	if len(webSite) != 0 {
		user.WebSite.Valid = true
		user.WebSite.String = webSite
	}
	return &user
}

/*
NewSite 创建一个新站点的数据结构
*/
func NewSite(domain string) *Site {
	if len(domain) == 0 {
		return nil
	}
	site := Site{}
	site.SiteDomain = domain
	site.SiteID = tutil.MD5([]byte(domain))
	site.CreateTime = time.Now()
	return &site
}

/*
NewComment 创建一个新评论的数据结构
*/
func NewComment(siteID, articleKey, userEmail, content string) *Comment {
	m := GlobalSQLMgr()
	user := m.GetUserByEmail(userEmail)
	if user == nil {
		tutil.LogWarn(fmt.Sprintf("没找到用户%v", userEmail))
		return nil
	}

	comm := Comment{}
	comm.UserID = user.UserID
	comm.SiteID = siteID
	comm.TimeStamp = time.Now()
	comm.ArticleID = tutil.MD5([]byte(articleKey))
	comm.ArticleKey = articleKey
	comm.Content = content
	comm.Deleted = false
	return &comm
}
