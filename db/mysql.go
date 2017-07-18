package db

import (
	"bytes"
	"database/sql"
	"fmt"

	"github.com/tyzual/tisq/tconf"
	"github.com/tyzual/tisq/util"

	//Mysql驱动
	_ "github.com/go-sql-driver/mysql"
)

const (
	createUserTableStmt = `
	CREATE TABLE IF NOT EXISTS user
	(UserId CHAR(32) PRIMARY KEY,
	Email VARCHAR(32) NOT NULL,
	DisplayName VARCHAR(16) NULL,
	WebSite VARCHAR(64) NULL,
	Avatar VARCHAR(256) NULL
	)
	`

	createCommentTableStmt = `CREATE TABLE IF NOT EXISTS comment
	(CommentID INT UNSIGNED PRIMARY KEY AUTO_INCREMENT, 
	ArticleID CHAR(32) NOT NULL,
	ArticleKey VARCHAR(512) NOT NULL,
	UserId CHAR(32) NOT NULL,
	Content VARCHAR(512) NOT NULL,
	TimeStamp DATETIME NOT NULL,
		FOREIGN KEY (UserId)
		REFERENCES user(UserId)
		ON DELETE CASCADE
	)
	`
)

/*
Mysql Mysql连接类
*/
type Mysql struct {
	dbconn *sql.DB
}

/*
Open 打开数据库连接
*/
func (m *Mysql) Open() {
	gConf := tconf.GlobalConf()
	var strBuff bytes.Buffer
	strBuff.WriteString(gConf.Mysql.User)
	if len(gConf.Mysql.Password) != 0 {
		strBuff.WriteRune(':')
		strBuff.WriteString(gConf.Mysql.User)
	}
	strBuff.WriteString(fmt.Sprintf("@tcp(%v:%d)/", gConf.Mysql.Host, gConf.Mysql.Port))
	strBuff.WriteString("?parseTime=true")
	str := strBuff.String()
	util.LogTrace(fmt.Sprintf("连接数据库字符串：%v", str))
	var err error
	m.dbconn, err = sql.Open("mysql", str)
	if err != nil {
		util.LogFatal(fmt.Sprintf("创建数据库出错，原因:%v", err))
	} else {
		util.Log("连接数据库成功")
	}

	if len(gConf.Mysql.DbName) == 0 {
		util.LogFatal("数据库名为空")
	}
	_, err = m.dbconn.Exec("CREATE DATABASE IF NOT EXISTS " + gConf.Mysql.DbName + " DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_general_ci")
	if err != nil {
		util.LogFatal(fmt.Sprintf("创建数据库出错，原因:%v", err))
	} else {
		util.Log(fmt.Sprintf("创建数据库%v成功", gConf.Mysql.DbName))
	}
	_, err = m.dbconn.Exec("USE " + gConf.Mysql.DbName)
	if err != nil {
		util.LogFatal(fmt.Sprintf("切换数据库出错，原因:%v", err))
	} else {
		util.Log(fmt.Sprintf("切换到数据库%v", gConf.Mysql.DbName))
	}

	_, err = m.dbconn.Exec(createUserTableStmt)
	if err != nil {
		util.LogFatal(fmt.Sprintf("创建User表出错，原因:%v", err))
	}

	_, err = m.dbconn.Exec(createCommentTableStmt)
	if err != nil {
		util.LogFatal(fmt.Sprintf("创建Comment表出错，原因:%v", err))
	}

	util.Log("初始化数据库成功")
}

/*
Close 关闭数据库
*/
func (m *Mysql) Close() {
	m.dbconn.Close()
}

func checkUser(user *User) bool {
	if user == nil {
		return false
	}
	if len(user.UserID) != 32 {
		return false
	}
	if len(user.Email) == 0 || len(user.Email) >= 32 {
		return false
	}
	if len(user.DisplayName) >= 16 {
		return false
	}
	if len(user.WebSite) >= 64 {
		return false
	}
	if len(user.Avatar) >= 256 {
		return false
	}
	return true
}

/*
InsertUser 向数据库中插入新用户
*/
func (m *Mysql) InsertUser(user *User) bool {
	if !checkUser(user) {
		util.LogWarn(fmt.Sprintf("User数据错误:%#v", *user))
		return false
	}
	var strBuff bytes.Buffer
	var args = make([]interface{}, 0)
	strBuff.WriteString("INSERT user SET UserId=?, Email=?")
	args = append(args, user.UserID, user.Email)
	if len(user.DisplayName) > 0 {
		strBuff.WriteString(",DisplayName=?")
		args = append(args, user.DisplayName)
	}
	if len(user.WebSite) > 0 {
		strBuff.WriteString(",WebSite=?")
		args = append(args, user.WebSite)
	}
	if len(user.Avatar) > 0 {
		strBuff.WriteString(",Avatar=?")
		args = append(args, user.Avatar)
	}
	strBuff.WriteString(" ON DUPLICATE KEY UPDATE ")
	if len(user.DisplayName) > 0 {
		strBuff.WriteString("DisplayName=?")
		args = append(args, user.DisplayName)
	}
	if len(user.WebSite) > 0 {
		strBuff.WriteString(",WebSite=?")
		args = append(args, user.WebSite)
	}
	if len(user.Avatar) > 0 {
		strBuff.WriteString(",Avatar=?")
		args = append(args, user.Avatar)
	}
	str := strBuff.String()
	util.LogTrace(fmt.Sprintf("数据库插入语句:%v", str))
	util.LogTrace(fmt.Sprintf("参数:%v", args))
	stmt, err := m.dbconn.Prepare(str)
	if err != nil {
		util.LogWarn(fmt.Sprintf("插入数据库失败，原因:%v", err))
		return false
	}
	_, err = stmt.Exec(args...)
	if err != nil {
		util.LogWarn(fmt.Sprintf("插入数据库失败，原因:%v", err))
		return false
	}
	util.LogTrace("插入数据库成功")
	return true
}

/*
GetUserByEmail 通过email获取用户信息
*/
func (m *Mysql) GetUserByEmail(email string) *User {
	rows, err := m.dbconn.Query("select * from user where Email=?", email)
	if err != nil {
		util.LogWarn(fmt.Sprintf("查询数据库失败，原因:%v", err))
		return nil
	}
	defer rows.Close()
	var user *User
	for rows.Next() {
		user = new(User)
		rows.Scan(&user.UserID,
			&user.Email,
			&user.DisplayName,
			&user.WebSite,
			&user.Avatar)
	}
	return user
}

/*
GetUserByID 通过id获取用户信息
*/
func (m *Mysql) GetUserByID(id string) *User {
	rows, err := m.dbconn.Query("select * from user WHERE UserId=?", id)
	if err != nil {
		util.LogWarn(fmt.Sprintf("查询数据库失败，原因:%v", err))
		return nil
	}
	defer rows.Close()
	var user *User
	for rows.Next() {
		user = new(User)
		rows.Scan(&user.UserID,
			&user.Email,
			&user.DisplayName,
			&user.WebSite,
			&user.Avatar)
	}
	return user
}

func checkComment(comm *Comment) bool {
	if comm == nil {
		return false
	}

	if len(comm.ArticleID) != 32 {
		return false
	}

	if len(comm.ArticleKey) >= 128 {
		return false
	}

	if len(comm.UserID) != 32 {
		return false
	}

	if len(comm.Content) == 0 || len(comm.Content) >= 512 {
		return false
	}
	return true
}

/*
InsertComment 向数据库中插入新用户
*/
func (m *Mysql) InsertComment(comm *Comment) bool {
	if !checkComment(comm) {
		util.LogWarn(fmt.Sprintf("Comment数据错误:%#v", *comm))
		return false
	}

	if m.GetUserByID(comm.UserID) == nil {
		util.LogWarn("用户ID不存在")
		return false
	}

	var args = make([]interface{}, 0, 4)
	str := "INSERT comment SET ArticleId=?, ArticleKey=?, UserId=?,Content=?,TimeStamp=?"
	args = append(args, comm.ArticleID, comm.ArticleKey, comm.UserID, comm.Content, comm.TimeStamp)
	util.LogTrace(fmt.Sprintf("数据库插入语句:%v", str))
	util.LogTrace(fmt.Sprintf("参数:%v", args))

	stmt, err := m.dbconn.Prepare(str)
	if err != nil {
		util.LogWarn(fmt.Sprintf("插入数据库失败，原因:%v", err))
		return false
	}
	_, err = stmt.Exec(args...)
	if err != nil {
		util.LogWarn(fmt.Sprintf("插入数据库失败，原因:%v", err))
		return false
	}
	util.LogTrace("插入数据库成功")

	return true
}

/*
GetCommentByArticleKey 通过article key来查询评论
*/
func (m *Mysql) GetCommentByArticleKey(key string) ([]Comment, []User) {
	if len(key) == 0 {
		return nil, nil
	}
	rows, err := m.dbconn.Query("SELECT * FROM comment WHERE ArticleKey=?", key)
	if err != nil {
		util.LogWarn(fmt.Sprintf("查询数据库失败，原因:%v", err))
		return nil, nil
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		util.LogWarn(fmt.Sprintf("查询数据库失败，原因:%v", err))
		return nil, nil
	}
	if len(cols) != 6 {
		util.LogWarn("数据库格式不匹配")
		return nil, nil
	}
	comms := make([]Comment, 0)
	userIds := make([]interface{}, 0) // []string
	for rows.Next() {
		comm := Comment{}
		rows.Scan(&comm.CommentID, &comm.ArticleID, &comm.ArticleKey, &comm.UserID, &comm.Content, &comm.TimeStamp)
		comms = append(comms, comm)
		userIds = append(userIds, comm.UserID)
	}
	if len(comms) == 0 {
		util.Log(fmt.Sprintf("没找到key为%v的评论", key))
		return nil, nil
	}

	var strBuff bytes.Buffer
	strBuff.WriteString("SELECT * from user WHERE ")
	for index := range userIds {
		if index > 0 {
			strBuff.WriteString("OR ")
		}
		strBuff.WriteString("UserId=? ")
	}
	str := strBuff.String()
	util.LogTrace(fmt.Sprintf("数据库查询语句:%v", str))
	util.LogTrace(fmt.Sprintf("参数:%v", userIds))
	rows, err = m.dbconn.Query(str, userIds...)
	if err != nil {
		util.LogWarn(fmt.Sprintf("查询数据库失败，原因:%v", err))
		return nil, nil
	}
	users := make([]User, 0)
	for rows.Next() {
		user := User{}
		rows.Scan(&user.UserID, &user.Email, &user.DisplayName, &user.WebSite, &user.Avatar)
		users = append(users, user)
	}

	return comms, users
}
