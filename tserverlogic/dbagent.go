package tserverlogic

import (
	"fmt"

	"github.com/tyzual/tisq/tdb"
	"github.com/tyzual/tisq/tutil"
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
	comm, ok := cmd.cmdArg.(*inComment)
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
		if !tdb.GlobalSQLMgr().InsertUser(dbUser) {
			tutil.LogWarn(fmt.Sprintf("插入用户失败，用户数据%#v", dbUser))
			return
		}
	}
	if dbUser == nil {
		return
	}

	dbComment := tdb.NewComment(dbSite.SiteID, comm.articleKey, dbUser.Email, comm.content, comm.replyID)
	if dbComment == nil {
		return
	}
	if !tdb.GlobalSQLMgr().InsertComment(dbComment) {
		return
	}

	oUser := OutUser{Email: dbUser.Email}
	if dbUser.DisplayName.Valid {
		oUser.DisplayName = &dbUser.DisplayName.String
	}
	if dbUser.WebSite.Valid {
		oUser.Site = &dbUser.WebSite.String
	}
	oComment := OutComment{UserID: dbUser.UserID, Content: dbComment.Content, CreateTime: dbComment.TimeStamp, CommentID: dbComment.CommentID}

	oResult := newResult()
	oResult.User[dbUser.UserID] = oUser
	oResult.Comment = append(oResult.Comment, oComment)

	cmd.result <- oResult
}

func queryComment(cmd *dbCmd) {
}
