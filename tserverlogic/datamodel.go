package tserverlogic

import (
	"time"
)

type inComment struct {
	domain string

	email       string
	displayName *string
	site        *string

	articleKey string
	content    string
	replyID    *uint32

	lastCommentID *string
}

/*
OutUser 服务器返回给客户端的User数据囧GB
*/
type OutUser struct {
	Email       string
	DisplayName *string
	Site        *string
}

/*
OutComment 服务器返回给客户端的Comment数据结构
*/
type OutComment struct {
	UserID     string
	Content    string
	CommentID  uint32
	CreateTime time.Time
}

/*
AddCommentResult 服务器返回给客户端的结果数据
*/
type AddCommentResult struct {
	//key:userid value:user
	User map[string]OutUser

	Comment []OutComment
}

func newResult() *AddCommentResult {
	oResult := AddCommentResult{}
	oResult.User = make(map[string]OutUser)
	oResult.Comment = make([]OutComment, 0)
	return &oResult
}
