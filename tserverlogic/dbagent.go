package tserverlogic

import (
	"fmt"
	"sync"

	"github.com/tyzual/tisq/tconf"
	"github.com/tyzual/tisq/tdb"
	"github.com/tyzual/tisq/tutil"
)

const (
	cmdInsertComment = iota
	cmdQueryComment
)

var cmdQueue = make(chan *dbCmd)
var rwLock = sync.RWMutex{}

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
					go insertComment(cmd)
				}
			case cmdQueryComment:
				{
					go queryComment(cmd)
				}
			}
		}
	}()
}

func newCmd(cmd uint16, arg interface{}) *dbCmd {
	return &dbCmd{cmd: cmd, cmdArg: arg, result: make(chan interface{})}
}

func insertComment(cmd *dbCmd) {
	rwLock.Lock()
	defer rwLock.Unlock()
	defer close(cmd.result)
	comm, ok := cmd.cmdArg.(*inComment)
	if !ok {
		return
	}
	if !tconf.GlobalConf().IsSiteRegistered(comm.domain) {
		return
	}
	dbSite := tdb.GlobalSQLMgr().GetSiteByDomain(comm.domain)
	if dbSite == nil {
		dbSite = tdb.NewSite(comm.domain)
		if dbSite == nil {
			return
		}
		if !tdb.GlobalSQLMgr().InsertSite(dbSite) {
			return
		}
	}

	var dbUser *tdb.User
	if (comm.displayName != nil && len(*comm.displayName) > 0) ||
		(comm.site != nil && len(*comm.site) > 0) {
		displayName := ""
		if comm.displayName != nil {
			displayName = *comm.displayName
		}
		site := ""
		if comm.site != nil {
			site = *comm.site
		}
		dbUser = tdb.NewUser(comm.email, displayName, site)
		if !tdb.GlobalSQLMgr().InsertUser(dbUser) {
			tutil.LogWarn("插入User失败")
			dbUser = nil
		}
	} else {
		dbUser = tdb.GlobalSQLMgr().GetUserByEmail(comm.email)
	}

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
	oComment := OutComment{UserID: dbUser.UserID, Content: dbComment.Content, CreateTime: dbComment.TimeStamp.Unix(), CommentID: dbComment.CommentID}
	if dbComment.ReplyID.Valid {
		replyID := uint32(dbComment.ReplyID.Int64)
		oComment.ReplyCommentID = &replyID
	}

	oResult := newResult()
	oResult.User[dbUser.UserID] = oUser
	oResult.Comment = append(oResult.Comment, oComment)

	cmd.result <- oResult
}

func queryComment(cmd *dbCmd) {
	rwLock.RLock()
	defer rwLock.RUnlock()
}
