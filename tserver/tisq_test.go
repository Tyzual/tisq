package main

import (
	"fmt"
	"testing"

	"github.com/tyzual/tisq/tdb"
	"github.com/tyzual/tisq/tutil"
)

func TestDb(t *testing.T) {
	tutil.PrintTrace = true

	user := tdb.NewUser("tyzual@gmail.com", "tyzual", "tyzual.com")
	tdb.GlobalSQLMgr().InsertUser(user)

	site := tdb.GlobalSQLMgr().GetSiteByDomain("tyzual.moe")
	if site == nil {
		site = tdb.NewSite("tyzual.moe")
		if !tdb.GlobalSQLMgr().InsertSite(site) {
			t.Fatal("插入site失败")
		}
	}
	if site == nil {
		t.Fatal("site 为空")
		return
	}

	tutil.LogTrace(fmt.Sprintf("%#v", site))
	comm := tdb.NewComment(site.SiteID, "abcdefg", "tyzual@gmail.com", "tyzual content250", nil)
	if comm != nil {
		tdb.GlobalSQLMgr().InsertComment(comm)
		fmt.Println("commID: ", comm.CommentID)
		comms, _ := tdb.GlobalSQLMgr().GetComment(comm.ArticleID, site.SiteID)
		fmt.Println(len(comms))
		for _, comm := range comms {
			fmt.Printf("%#v\n", comm)
		}
	}

}
