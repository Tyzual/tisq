package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"tisq/conf"
	"tisq/util"

	//Mysql 驱动
	_ "github.com/go-sql-driver/mysql"
)

const (
	createUserTableStmt = `
	CREATE TABLE IF NOT EXISTS user
	(UserId CHAR(32) PRIMARY KEY,
	Email VARCHAR(32) NOT NULL,
	DisplayName VARCHAR(16),
	WebSite VARCHAR(64),
	Avatar VARCHAR(256)
	)
	`
	createCommentTableStmt = `CREATE TABLE IF NOT EXISTS comment
	(CommentID INT UNSIGNED PRIMARY KEY AUTO_INCREMENT, 
	ArticleID CHAR(32) NOT NULL,
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
	gConf := conf.GlobalConf()
	var strBuff bytes.Buffer
	strBuff.WriteString(gConf.Mysql.User)
	if len(gConf.Mysql.Password) != 0 {
		strBuff.WriteRune(':')
		strBuff.WriteString(gConf.Mysql.User)
	}
	strBuff.WriteString(fmt.Sprintf("@tcp(%v:%d)/", gConf.Mysql.Host, gConf.Mysql.Port))
	var err error
	m.dbconn, err = sql.Open("mysql", strBuff.String())
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
