package tserverlogic

type inQueryComment struct {
	domain        string
	articleKey    string
	lastCommentID *string
}

type inComment struct {
	inQueryComment

	email       string
	displayName *string
	site        *string

	content string
	replyID *uint32
}

/*
OutUser 服务器返回给客户端的User数据囧GB
*/
type OutUser struct {
	Email       string
	DisplayName *string `json:",omitempty"`
	Site        *string `json:",omitempty"`
}

/*
OutComment 服务器返回给客户端的Comment数据结构
*/
type OutComment struct {
	UserID         string
	Content        string
	CommentID      uint32
	CreateTime     int64
	ReplyCommentID *uint32 `json:",omitempty"`
}

/*
CommentResult 服务器返回给客户端的结果数据
*/
type CommentResult struct {
	//key:userid value:user
	User map[string]OutUser

	Comment []OutComment
}

func newResult() *CommentResult {
	oResult := CommentResult{}
	oResult.User = make(map[string]OutUser)
	oResult.Comment = make([]OutComment, 0)
	return &oResult
}
