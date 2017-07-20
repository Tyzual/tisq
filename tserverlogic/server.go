package tserverlogic

import (
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/tyzual/tisq/tutil"
)

/*
HandleAddComment 添加评论Handler
*/
func HandleAddComment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		tutil.Log(fmt.Sprintf("request from:%v", r.Host))
		w.Write([]byte("ERROR: USE POST"))
		return
	}
	// TODO:Parse body
	comment := inComment{}
	cmd := newCmd(cmdInsertComment, &comment)
	cmdQueue <- cmd
	res, ok := <-cmd.result
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	outRes, ok := res.(*Result)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(outRes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tutil.LogWarn(fmt.Sprintf("创建JSON字符串出错，原因%v", err))
		return
	}
	tutil.Log(string(jsonData))
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

/*
HandleCommentList 获取评论列表Handler
*/
func HandleCommentList(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("ERROR: USE POST"))
		return
	}
	tutil.Log(fmt.Sprintf("%v", r.Method))
}
