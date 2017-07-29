package tserverlogic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tyzual/tisq/tutil"
)

const (
	keyDomain        = "domain"
	keyEmail         = "email"
	keyDisplayName   = "displayname"
	keySite          = "site"
	keyArticleKey    = "articlekey"
	keyContent       = "content"
	keyReplyID       = "replyid"
	keyLastCommentID = "lastcommentid"
)

/*
HandleAddComment 添加评论Handler
请求地址：
domain:port/addComment
body：
先将key和value进行URL编码 js调用encodeURIComponent()函数
然后再将 key和value 以
key1=value1&key2=value2
的形式编码发送

参数	含义
domain	博客的域名

email	评论者的email
displayname	评论者显示的名字(昵称)，可选
site		评论者的主页，可选

articlekey	评论文章的特征码(可以使用文章的URL地址)
content		评论内容
replyid		如果评论是回复某条评论，则在这里填写回复评论的评论ID，如果不是回复，不要设置这个字段

lastcommentid	客户端最新一条评论的评论id，若不传这个值，服务器会返回articlekey下的所有评论，否则返回lastcommentid以后的评论。若找不到lastcommentid所对应的评论。则返回错误。
*/
func HandleAddComment(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		tutil.Log(fmt.Sprintf("request from:%v", r.Host))
		w.Write([]byte("ERROR: USE POST"))
		return
	}
	err := r.ParseForm()
	if err != nil {
		tutil.LogWarn(fmt.Sprintf("%v处理请求错误", r.URL.String()))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("parse args error"))
		return
	}
	comment := inComment{}
	domains, ok := r.Form[keyDomain]
	if !ok || len(domains) == 0 || len(domains[0]) == 0 {
		tutil.LogWarn(fmt.Sprintf("%vdomain错误", r.URL.String()))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("parse domain error"))
		return
	}
	comment.domain = domains[0]

	emails, ok := r.Form[keyEmail]
	if !ok || len(emails) == 0 || len(emails[0]) == 0 {
		tutil.LogWarn("addComment email错误")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("parse email error"))
		return
	}
	comment.email = emails[0]

	displayNames, ok := r.Form[keyDisplayName]
	if ok && len(displayNames) != 0 && len(displayNames[0]) != 0 {
		displayName := displayNames[0]
		comment.displayName = &displayName
	}

	sites, ok := r.Form[keySite]
	if ok && len(sites) != 0 && len(sites[0]) != 0 {
		site := sites[0]
		comment.site = &site
	}

	articleKeys, ok := r.Form[keyArticleKey]
	if !ok || len(articleKeys) == 0 || len(articleKeys[0]) == 0 {
		tutil.LogWarn("addComment articlekey错误")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("parse articlekey error"))
		return
	}
	comment.articleKey = articleKeys[0]

	contents, ok := r.Form[keyContent]
	if !ok || len(contents) == 0 || len(contents[0]) == 0 {
		tutil.LogWarn("addComment content错误")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("parse content error"))
		return
	}
	comment.content = contents[0]

	replyIDs, ok := r.Form[keyReplyID]
	if ok && len(replyIDs) != 0 && len(replyIDs[0]) != 0 {
		replyID64, err := strconv.ParseUint(replyIDs[0], 10, 32)
		if err != nil {
			tutil.LogWarn("解析replyID错误")
		}
		replyID32 := uint32(replyID64)
		comment.replyID = &replyID32
	}

	lastCommentIDs, ok := r.Form[keyLastCommentID]
	if ok && len(lastCommentIDs) != 0 && len(lastCommentIDs[0]) != 0 {
		lastCommentID := lastCommentIDs[0]
		comment.lastCommentID = &lastCommentID
	}

	cmd := newCmd(cmdInsertComment, &comment)
	cmdQueue <- cmd
	res, ok := <-cmd.result
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	outRes, ok := res.(*CommentResult)
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
请求地址：
domain:port/commentList
body：
先将key和value进行URL编码 js调用encodeURIComponent()函数
然后再将 key和value 以
key1=value1&key2=value2
的形式编码发送

参数	含义
domain	博客的域名
articlekey	评论文章的特征码(可以使用文章的URL地址)
lastcommentid	客户端最新一条评论的评论id，若不传这个值，服务器会返回articlekey下的所有评论，否则返回lastcommentid以后的评论。若找不到lastcommentid所对应的评论。则返回错误。
*/
func HandleCommentList(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("ERROR: USE POST"))
		return
	}
	err := r.ParseForm()
	if err != nil {
		tutil.LogWarn(fmt.Sprintf("%v处理请求错误", r.URL.String()))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("parse args error"))
		return
	}

	var queryComment inQueryComment

	domains, ok := r.Form[keyDomain]
	if !ok || len(domains) == 0 || len(domains[0]) == 0 {
		tutil.LogWarn(fmt.Sprintf("%vdomain错误", r.URL.String()))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("parse domain error"))
		return
	}
	queryComment.domain = domains[0]

	articleKeys, ok := r.Form[keyArticleKey]
	if !ok || len(articleKeys) == 0 || len(articleKeys[0]) == 0 {
		tutil.LogWarn(fmt.Sprintf("%varticlekey错误", r.URL.String()))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("parse articlekey error"))
		return
	}
	queryComment.articleKey = articleKeys[0]

	lastCommentIDs, ok := r.Form[keyLastCommentID]
	if ok && len(lastCommentIDs) != 0 && len(lastCommentIDs[0]) != 0 {
		lastCommentID := lastCommentIDs[0]
		queryComment.lastCommentID = &lastCommentID
	}

	cmd := newCmd(cmdQueryComment, &queryComment)
	cmdQueue <- cmd
	res, ok := <-cmd.result
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	outRes, ok := res.(*CommentResult)
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
