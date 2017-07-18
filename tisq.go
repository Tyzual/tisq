package main

import (
	"fmt"
	"net/http"

	"github.com/tyzual/tisq/tconf"
	"github.com/tyzual/tisq/tdb"
	"github.com/tyzual/tisq/tserver"
	"github.com/tyzual/tisq/tutil"
)

func main() {
	util.PrintTrace = true

	http.HandleFunc("/addComment", tserver.HandleAddComment)
	http.HandleFunc("/commentList", tserver.HandleCommentList)
	util.LogFatal(fmt.Sprintf("%v", http.ListenAndServe(fmt.Sprintf("%v:%d", tconf.GlobalConf().Server.Domain, tconf.GlobalConf().Server.Port), nil)))
}

func test() {
	user := tdb.NewUser("echizen@foxmail.com", "echizen", "tyzual.com", "")
	tdb.GlobalSqlMgr().InsertUser(user)

	// comm := db.NewComment("abcdefg", "echizen@foxmail.com", "echizen content")
	// if comm != nil {
	// 	mysql.InsertComment(comm)
	// }

	comms, _ := tdb.GlobalSqlMgr().GetCommentByArticleKey("abcdefg")
	fmt.Println(len(comms))
}
