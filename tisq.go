package main

import (
	"fmt"
	"net/http"

	"github.com/tyzual/tisq/tconf"
	"github.com/tyzual/tisq/tserver"
	"github.com/tyzual/tisq/tutil"
)

func main() {
	util.PrintTrace = true

	http.HandleFunc("/addComment", tserver.HandleAddComment)
	http.HandleFunc("/commentList", tserver.HandleCommentList)
	util.LogFatal(fmt.Sprintf("%v", http.ListenAndServe(fmt.Sprintf("%v:%d", tconf.GlobalConf().Server.Domain, tconf.GlobalConf().Server.Port), nil)))
}
