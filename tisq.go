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
	conf.LoadConf()
	var mysql db.Mysql
	mysql.Open()
	defer mysql.Close()
	http.HandleFunc("/addComment", server.HandleAddComment)
	http.HandleFunc("/commentList", server.HandleCommentList)
	util.LogFatal(fmt.Sprintf("%v", http.ListenAndServe(fmt.Sprintf("localhost:%d", conf.GlobalConf().Server.Port), nil)))
}
