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

type outUser struct {
	email       string
	displayName *string
	site        *string
}

type outComment struct {
	userID     string
	content    string
	commentID  string
	createTime time.Time
}

type result struct {
	//key:userid value:user
	user map[string]outUser

	comment []outComment
}

func newResult() *result {
	oResult := result{}
	oResult.user = make(map[string]outUser)
	oResult.comment = make([]outComment, 0)
	return &oResult
}
