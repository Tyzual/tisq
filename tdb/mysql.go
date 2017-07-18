package tdb

import (
	"bytes"
	"database/sql"
	"fmt"

	"github.com/tyzual/tisq/tconf"
	"github.com/tyzual/tisq/tutil"

	//Mysql驱动
	_ "github.com/go-sql-driver/mysql"
)

const (
	createUserTableStmt = `
	CREATE TABLE IF NOT EXISTS user (
		UserId CHAR(32) PRIMARY KEY,
		Email VARCHAR(32) NOT NULL,
		DisplayName VARCHAR(16) NULL,
		WebSite VARCHAR(64) NULL,
		Avatar VARCHAR(256) NULL
	)
	`
	createSiteTableStmt = `CREATE TABLE IF NOT EXISTS site (
		SiteId CHAR(32) PRIMARY KEY,
		SiteDomain VARCHAR(64) NOT NULL,
		PrivateKey CHAR(16) NOT NULL
	)`

	createCommentTableStmt = `CREATE TABLE IF NOT EXISTS comment (
		CommentID INT UNSIGNED PRIMARY KEY AUTO_INCREMENT, 
		ArticleID CHAR(32) NOT NULL,
		ArticleKey VARCHAR(512) NOT NULL,
		UserId CHAR(32) NOT NULL,
		SiteId CHAR(32) NOT NULL,
		Content VARCHAR(512) NOT NULL,
		TimeStamp DATETIME NOT NULL,
		Deleted BOOL DEFAULT false,
			FOREIGN KEY (UserId)
			REFERENCES user(UserId)
			ON DELETE CASCADE,
			FOREIGN KEY (SiteId)
			REFERENCES site(SiteId)
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

var gMysql Mysql

func init() {
	gMysql.open()
}

/*
GlobalSqlMgr 全局数据库对象
*/
func GlobalSqlMgr() *Mysql {
	return &gMysql
}

func (m *Mysql) open() {
	gConf := tconf.GlobalConf()
	var strBuff bytes.Buffer
	strBuff.WriteString(gConf.Mysql.User)
	if len(gConf.Mysql.Password) != 0 {
		strBuff.WriteRune(':')
		strBuff.WriteString(gConf.Mysql.Password)
	}
	strBuff.WriteString(fmt.Sprintf("@tcp(%v:%d)/", gConf.Mysql.Host, gConf.Mysql.Port))
	strBuff.WriteString("?parseTime=true")
	str := strBuff.String()
	tutil.LogTrace(fmt.Sprintf("连接数据库字符串：%v", str))
	var err error
	m.dbconn, err = sql.Open("mysql", str)
	if err != nil {
		tutil.LogFatal(fmt.Sprintf("创建数据库出错，原因:%v", err))
	} else {
		tutil.Log("连接数据库成功")
	}

	if len(gConf.Mysql.DbName) == 0 {
		tutil.LogFatal("数据库名为空")
	}
	_, err = m.dbconn.Exec("CREATE DATABASE IF NOT EXISTS " + gConf.Mysql.DbName + " DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_general_ci")
	if err != nil {
		tutil.LogFatal(fmt.Sprintf("创建数据库出错，原因:%v", err))
	} else {
		tutil.Log(fmt.Sprintf("创建数据库%v成功", gConf.Mysql.DbName))
	}
	_, err = m.dbconn.Exec("USE " + gConf.Mysql.DbName)
	if err != nil {
		tutil.LogFatal(fmt.Sprintf("切换数据库出错，原因:%v", err))
	} else {
		tutil.Log(fmt.Sprintf("切换到数据库%v", gConf.Mysql.DbName))
	}

	_, err = m.dbconn.Exec(createUserTableStmt)
	if err != nil {
		tutil.LogFatal(fmt.Sprintf("创建User表出错，原因:%v", err))
	}

	_, err = m.dbconn.Exec(createSiteTableStmt)
	if err != nil {
		tutil.LogFatal(fmt.Sprintf("创建Site表出错，原因:%v", err))
	}

	_, err = m.dbconn.Exec(createCommentTableStmt)
	if err != nil {
		tutil.LogFatal(fmt.Sprintf("创建Comment表出错，原因:%v", err))
	}

	tutil.Log("初始化数据库成功")
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
	if user.DisplayName.Valid && len(user.DisplayName.String) >= 16 {
		return false
	}
	if user.WebSite.Valid && len(user.WebSite.String) >= 64 {
		return false
	}
	if user.Avatar.Valid && len(user.Avatar.String) >= 256 {
		return false
	}
	return true
}

/*
InsertUser 向数据库中插入新用户
*/
func (m *Mysql) InsertUser(user *User) bool {
	if !checkUser(user) {
		tutil.LogWarn(fmt.Sprintf("User数据错误:%#v", *user))
		return false
	}
	var strBuff bytes.Buffer
	var args = make([]interface{}, 0)
	strBuff.WriteString("INSERT user SET UserId=?, Email=?")
	args = append(args, user.UserID, user.Email)
	if user.DisplayName.Valid {
		strBuff.WriteString(",DisplayName=?")
		args = append(args, user.DisplayName.String)
	}
	if user.WebSite.Valid {
		strBuff.WriteString(",WebSite=?")
		args = append(args, user.WebSite.String)
	}
	if user.Avatar.Valid {
		strBuff.WriteString(",Avatar=?")
		args = append(args, user.Avatar.String)
	}
	strBuff.WriteString(" ON DUPLICATE KEY UPDATE ")
	if user.DisplayName.Valid {
		strBuff.WriteString("DisplayName=?")
		args = append(args, user.DisplayName.String)
	}
	if user.WebSite.Valid {
		strBuff.WriteString(",WebSite=?")
		args = append(args, user.WebSite.String)
	}
	if user.Avatar.Valid {
		strBuff.WriteString(",Avatar=?")
		args = append(args, user.Avatar.String)
	}
	str := strBuff.String()
	tutil.LogTrace(fmt.Sprintf("数据库插入语句:%v", str))
	tutil.LogTrace(fmt.Sprintf("参数:%v", args))
	stmt, err := m.dbconn.Prepare(str)
	if err != nil {
		tutil.LogWarn(fmt.Sprintf("插入数据库失败，原因:%v", err))
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(args...)
	if err != nil {
		tutil.LogWarn(fmt.Sprintf("插入数据库失败，原因:%v", err))
		return false
	}
	tutil.LogTrace("插入数据库成功")
	return true
}

/*
GetUserByEmail 通过email获取用户信息
*/
func (m *Mysql) GetUserByEmail(email string) *User {
	rows, err := m.dbconn.Query("select * from user where Email=?", email)
	if err != nil {
		tutil.LogWarn(fmt.Sprintf("查询数据库失败，原因:%v", err))
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
		tutil.LogWarn(fmt.Sprintf("查询数据库失败，原因:%v", err))
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

	if len(comm.SiteID) != 32 {
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
		tutil.LogWarn(fmt.Sprintf("Comment数据错误:%#v", *comm))
		return false
	}

	if m.GetUserByID(comm.UserID) == nil {
		tutil.LogWarn("用户ID不存在")
		return false
	}

	var args = make([]interface{}, 0, 4)
	str := "INSERT comment SET ArticleId=?, ArticleKey=?, UserId=?,Content=?,TimeStamp=?,SiteId=?"
	args = append(args, comm.ArticleID, comm.ArticleKey, comm.UserID, comm.Content, comm.TimeStamp, comm.SiteID)
	tutil.LogTrace(fmt.Sprintf("数据库插入语句:%v", str))
	tutil.LogTrace(fmt.Sprintf("参数:%v", args))

	stmt, err := m.dbconn.Prepare(str)
	if err != nil {
		tutil.LogWarn(fmt.Sprintf("插入数据库失败，原因:%v", err))
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(args...)
	if err != nil {
		tutil.LogWarn(fmt.Sprintf("插入数据库失败，原因:%v", err))
		return false
	}
	tutil.LogTrace("插入数据库成功")

	return true
}

/*
GetCommentByArticleKey 通过article key来查询评论
*/
func (m *Mysql) GetCommentByArticleKey(key string) ([]Comment, []User) {
	if len(key) == 0 {
		return nil, nil
	}
	comments, err := m.dbconn.Query("SELECT * FROM comment WHERE ArticleKey=? AND Deleted=false", key)
	if err != nil {
		tutil.LogWarn(fmt.Sprintf("查询数据库失败，原因:%v", err))
		return nil, nil
	}
	defer comments.Close()
	cols, err := comments.Columns()
	if err != nil {
		tutil.LogWarn(fmt.Sprintf("查询数据库失败，原因:%v", err))
		return nil, nil
	}
	if len(cols) != 7 {
		tutil.LogTrace("数据库格式不匹配")
		return nil, nil
	}
	comms := make([]Comment, 0)
	userIds := make([]interface{}, 0) // []string
	for comments.Next() {
		comm := Comment{}
		if err := comments.Scan(&comm.CommentID, &comm.ArticleID, &comm.ArticleKey, &comm.UserID, &comm.Content, &comm.TimeStamp, &comm.Deleted); err != nil {
			tutil.LogTrace(fmt.Sprintf("解析Comment数据结构错误，原因:%v", err))
			return nil, nil
		}
		comms = append(comms, comm)
		userIds = append(userIds, comm.UserID)
	}
	if len(comms) == 0 {
		tutil.Log(fmt.Sprintf("没找到key为%v的评论", key))
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
	tutil.LogTrace(fmt.Sprintf("数据库查询语句:%v", str))
	tutil.LogTrace(fmt.Sprintf("参数:%v", userIds))
	dbUsers, err := m.dbconn.Query(str, userIds...)
	if err != nil {
		tutil.LogWarn(fmt.Sprintf("查询数据库失败，原因:%v", err))
		return nil, nil
	}
	defer dbUsers.Close()
	users := make([]User, 0)
	for dbUsers.Next() {
		user := User{}
		if err := dbUsers.Scan(&user.UserID, &user.Email, &user.DisplayName, &user.WebSite, &user.Avatar); err != nil {
			tutil.LogTrace(fmt.Sprintf("解析User数据结构错误，原因:%v", err))
		}
		users = append(users, user)
	}

	return comms, users
}

func checkSite(site *Site) bool {
	if site == nil {
		return false
	}
	if len(site.SiteID) != 32 {
		return false
	}
	if len(site.SiteDomain) == 0 || len(site.SiteDomain) >= 64 {
		return false
	}
	if len(site.PrivateKey) != 16 {
		return false
	}
	return true
}

/*
InsertSite 插入site数据
*/
func (m *Mysql) InsertSite(site *Site) bool {
	if !checkSite(site) {
		tutil.LogWarn(fmt.Sprintf("Site数据错误:%#v", *site))
		return false
	}
	stmt, err := m.dbconn.Prepare("INSERT site SET SiteId=?, SiteDomain=?, PrivateKey=?")
	if err != nil {
		tutil.LogWarn(fmt.Sprintf("插入数据库失败，原因:%v", err))
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(site.SiteID, site.SiteDomain, site.PrivateKey)
	if err != nil {
		tutil.LogWarn(fmt.Sprintf("插入数据库失败，原因:%v", err))
		return false
	}
	tutil.LogTrace("插入数据库成功")
	return true
}

/*
GetSiteByDomain 通过域名返回站点
*/
func (m *Mysql) GetSiteByDomain(domain string) *Site {
	sites, err := m.dbconn.Query("SELECT * FROM site WHERE SiteDomain=?", domain)
	if err != nil {
		tutil.LogWarn(fmt.Sprintf("查询数据库失败，原因:%v", err))
		return nil
	}
	defer sites.Close()
	site := Site{}
	for sites.Next() {
		sites.Scan(&site.SiteID, &site.SiteDomain, &site.PrivateKey)
	}
	return &site
}
