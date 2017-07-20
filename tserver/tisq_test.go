package main

import (
	"fmt"
	"testing"

	"github.com/tyzual/tisq/tdb"
)

func TestDb(t *testing.T) {
	user := tdb.NewUser("echizen@foxmail.com", "echizen", "tyzual.com", "")
	tdb.GlobalSqlMgr().InsertUser(user)

	site := tdb.NewSite("tyzual.com")
	tdb.GlobalSqlMgr().InsertSite(site)

	comm := tdb.NewComment("tyzual.com", "abcdefg", "echizen@foxmail.com", "echizen content")
	if comm != nil {
		tdb.GlobalSqlMgr().InsertComment(comm)
	}

	comms, _ := tdb.GlobalSqlMgr().GetCommentByArticleKey("abcdefg")
	fmt.Println(len(comms))
}