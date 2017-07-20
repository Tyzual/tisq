package tserverlogic

import (
	"github.com/tyzual/tisq/tdb"
)

const (
	cmdInsertComment = iota
	cmdQueryComment
)

var cmdQueue = make(chan *dbCmd)

type dbCmd struct {
	cmd    uint16
	cmdArg interface{}
	result chan interface{}
}

func init() {
	go func() {
		for cmd := range cmdQueue {
			switch cmd.cmd {
			case cmdInsertComment:
				{
					insertComment(cmd)
				}
			case cmdQueryComment:
				{
					queryComment(cmd)
				}
			}
		}
	}()
}

func newCmd(cmd uint16, arg interface{}) *dbCmd {
	return &dbCmd{cmd: cmd, cmdArg: arg, result: make(chan interface{})}
}

func insertComment(cmd *dbCmd) {
	defer close(cmd.result)
	comm, ok := cmd.cmdArg.(inComment)
	if !ok {
		return
	}
	dbSite := tdb.GlobalSQLMgr().GetSiteByDomain(comm.domain)
	if dbSite == nil {
		return
	}

	dbUser := tdb.GlobalSQLMgr().GetUserByEmail(comm.email)
	if dbUser == nil {
		webSite := ""
		if comm.site != nil {
			webSite = *comm.site
		}
		displayName := ""
		if comm.displayName != nil {
			displayName = *comm.displayName
		}
		dbUser = tdb.NewUser(comm.email, displayName, webSite)
	}
	if dbUser == nil {
		return
	}

	dbComment := tdb.NewComment(dbSite.SiteID, comm.articleKey, dbUser.Email, comm.content)
	if dbComment == nil {
		return
	}
	if !tdb.GlobalSQLMgr().InsertComment(dbComment) {
		return
	}

	oUser := outUser{email: dbUser.Email}
	if dbUser.DisplayName.Valid {
		oUser.displayName = &dbUser.DisplayName.String
	}
	if dbUser.WebSite.Valid {
		oUser.site = &dbUser.WebSite.String
	}
	oComment := outComment{userID: dbUser.UserID, content: dbComment.Content, createTime: dbComment.TimeStamp, commentID: dbComment.CommentID}

	oResult := newResult()
	oResult.user[dbUser.UserID] = oUser
	oResult.comment = append(oResult.comment, oComment)

	cmd.result <- oComment
}

func queryComment(cmd *dbCmd) {
}
