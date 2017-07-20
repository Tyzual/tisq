package main

import (
	"fmt"
	"net/http"

	"github.com/tyzual/tisq/tconf"
	"github.com/tyzual/tisq/tserverlogic"
	"github.com/tyzual/tisq/tutil"
)

func main() {
	tutil.PrintTrace = true

	http.HandleFunc("/addComment", tserverlogic.HandleAddComment)
	http.HandleFunc("/commentList", tserverlogic.HandleCommentList)
	tutil.LogFatal(fmt.Sprintf("%v", http.ListenAndServe(fmt.Sprintf("%v:%d", tconf.GlobalConf().Server.Domain, tconf.GlobalConf().Server.Port), nil)))
}
