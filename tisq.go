package main

import (
	"fmt"
	"net/http"

	"github.com/tyzual/tisq/db"
	"github.com/tyzual/tisq/server"
	"github.com/tyzual/tisq/tconf"
	"github.com/tyzual/tisq/util"
)

func main() {
	util.PrintTrace = true
	tconf.LoadConf()
	var mysql db.Mysql
	mysql.Open()
	defer mysql.Close()
	test(&mysql)

	http.HandleFunc("/addComment", server.HandleAddComment)
	http.HandleFunc("/commentList", server.HandleCommentList)
	util.LogFatal(fmt.Sprintf("%v", http.ListenAndServe(fmt.Sprintf("%v:%d", tconf.GlobalConf().Server.Domain, tconf.GlobalConf().Server.Port), nil)))
}

func test(mysql *db.Mysql) {
	// user := db.NewUser("echizen@foxmail.com", "echizen", "tyzual.com", "")
	// mysql.InsertUser(user)

	// comm := db.NewComment(mysql, "abcdefg", "echizen@foxmail.com", "echizen content")
	// if comm != nil {
	// 	mysql.InsertComment(comm)
	// }

	comms, _ := mysql.GetCommentByArticleKey("abcdefg")
	fmt.Println(len(comms))
}
