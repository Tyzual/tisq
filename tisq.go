package main

import (
	"fmt"
	"net/http"
	"tisq/conf"
	"tisq/server"
	"tisq/util"
)

func main() {
	conf.LoadConf()
	http.HandleFunc("/addComment", server.HandleAddComment)
	http.HandleFunc("/commentList", server.HandleCommentList)
	util.LogFatal(fmt.Sprintf("%v", http.ListenAndServe(fmt.Sprintf("localhost:%d", conf.GlobalConf().Port), nil)))
}
