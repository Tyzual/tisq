package server

import (
	"fmt"
	"net/http"
	"tisq/util"
)

/*
HandleAddComment 添加评论Handler
*/
func HandleAddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("ERROR: USE POST"))
		return
	}
	util.Log(fmt.Sprintf("%v", r.Method))
}

/*
HandleCommentList 获取评论列表Handler
*/
func HandleCommentList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("ERROR: USE POST"))
		return
	}
	util.Log(fmt.Sprintf("%v", r.Method))
}
