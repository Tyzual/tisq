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

	lastCommentID *string
}

type OutUser struct {
	Email       string
	DisplayName *string
	Site        *string
}

type OutComment struct {
	UserID     string
	Content    string
	CommentID  string
	CreateTime time.Time
}

type Result struct {
	//key:userid value:user
	User map[string]OutUser

	Comment []OutComment
}

func newResult() *Result {
	oResult := Result{}
	oResult.User = make(map[string]OutUser)
	oResult.Comment = make([]OutComment, 0)
	return &oResult
}
