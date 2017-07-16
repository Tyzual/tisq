package main

import (
	"fmt"
	"net/http"
	"tisq/conf"
	"tisq/db"
	"tisq/server"
	"tisq/util"
)

func main() {
	util.PrintTrace = true
	conf.LoadConf()
	var mysql db.Mysql
	mysql.Open()
	defer mysql.Close()
	// test(&mysql)

	http.HandleFunc("/addComment", server.HandleAddComment)
	http.HandleFunc("/commentList", server.HandleCommentList)
	util.LogFatal(fmt.Sprintf("%v", http.ListenAndServe(fmt.Sprintf("localhost:%d", conf.GlobalConf().Server.Port), nil)))
}

func test(mysql *db.Mysql) {
	user := db.NewUser("e.tyzual@gmail.com", "tyzual", "tyzual.com", "")
	mysql.InsertUser(user)

	user = mysql.GetUserByEmail("e.tyzual@gmail.com")
}
