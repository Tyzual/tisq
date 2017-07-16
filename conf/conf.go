package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"tisq/util"
)

var gConf = TConf{}

const confFile = "/etc/tisq.cnf"

/*
TConf 服务器的全局配置
*/
type TConf struct {
	Server ServerConf
	Mysql  MysqlConf
}

/*
ServerConf 服务器相关配置
*/
type ServerConf struct {
	Domain string
	Port   uint16
}

/*
MysqlConf Mysql配置
*/
type MysqlConf struct {
	User     string
	Password string
	Host     string
	Port     uint16
	DbName   string
}

/*
GlobalConf 获取全局配置
*/
func GlobalConf() *TConf {
	return &gConf
}

/*
LoadConf 加载配置
*/
func LoadConf() {
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		gConf.Server.Port = 34958
		gConf.Server.Domain = "localhost"
		gConf.Mysql.Host = "localhost"
		gConf.Mysql.DbName = "TISQ"
		gConf.Mysql.Password = ""
		gConf.Mysql.Port = 3306
		gConf.Mysql.User = "root"
		jsonByte, err := json.MarshalIndent(gConf, "", "\t")
		if err != nil {
			util.LogWarn(fmt.Sprintf("创建配置文件出错\n错误原因:%v", err))
		} else {
			util.Log(fmt.Sprintf("生成默认配置文件\"%v\":\n%v", confFile, string(jsonByte)))
			if err = ioutil.WriteFile(confFile, jsonByte, 0644); err != nil {
				util.LogWarn(fmt.Sprintf("创建配置文件出错\n错误原因:%v", err))
			}
		}
	} else {
		if jsonByte, err := ioutil.ReadFile(confFile); err != nil {
			util.LogWarn(fmt.Sprintf("读取配置文件出错\n错误原因:%v", err))
		} else {
			if err := json.Unmarshal(jsonByte, &gConf); err != nil {
				util.LogWarn(fmt.Sprintf("读取配置文件出错\n错误原因:%v", err))
			} else {
				util.Log(fmt.Sprintf("域名:%v", gConf.Server.Domain))
				util.Log(fmt.Sprintf("端口:%v", gConf.Server.Port))
			}
		}
	}
}
