package tserverlogic

import (
	"fmt"
	"net/http"

	"github.com/tyzual/tisq/tutil"
)

/*
HandleAddComment 添加评论Handler
*/
func HandleAddComment(w http.ResponseWriter, r *http.Request) {
	go func() {
		defer r.Body.Close()
		if r.Method != http.MethodPost {
			tutil.Log(fmt.Sprintf("request from:%v", r.Host))
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("ERROR: USE POST"))
			return
		}
		tutil.Log(fmt.Sprintf("%v", r.Method))
	}()
}

/*
HandleCommentList 获取评论列表Handler
*/
func HandleCommentList(w http.ResponseWriter, r *http.Request) {
	go func() {
		defer r.Body.Close()
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte("ERROR: USE POST"))
			return
		}
		tutil.Log(fmt.Sprintf("%v", r.Method))
	}()
}
