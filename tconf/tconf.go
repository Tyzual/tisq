package tconf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/tyzual/tisq/tutil"
)

var gConf = TConf{}

const confFile = "/etc/tisq.cnf"

/*
配置文件格式
{
"Server": {
	"Domain": "localhost",
	"Port": 34958
},
"Mysql": {
	"User": "root",
	"Password": "",
	"Host": "localhost",
	"Port": 3306,
	"DbName": "TISQ"
},
"Site": [
		"tyzual.moe",
		"tyzual.com"
	]
}
*/

/*
TConf 服务器的全局配置
*/
type TConf struct {
	Server   ServerConf
	Mysql    MysqlConf
	Site     []string
	siteDict map[string]struct{}
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

func init() {
	gConf.siteDict = make(map[string]struct{})
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		gConf.Server.Port = 34958
		gConf.Server.Domain = "localhost"
		gConf.Mysql.Host = "localhost"
		gConf.Mysql.DbName = "TISQ"
		gConf.Mysql.Password = ""
		gConf.Mysql.Port = 3306
		gConf.Mysql.User = "root"
		gConf.Site = make([]string, 0)
		jsonByte, err := json.MarshalIndent(gConf, "", "\t")
		if err != nil {
			tutil.LogWarn(fmt.Sprintf("创建配置文件出错\n错误原因:%v", err))
		} else {
			tutil.Log(fmt.Sprintf("生成默认配置文件\"%v\":\n%v", confFile, string(jsonByte)))
			if err = ioutil.WriteFile(confFile, jsonByte, 0644); err != nil {
				tutil.LogWarn(fmt.Sprintf("创建配置文件出错\n错误原因:%v", err))
			}
		}
	} else {
		if jsonByte, err := ioutil.ReadFile(confFile); err != nil {
			tutil.LogWarn(fmt.Sprintf("读取配置文件出错\n错误原因:%v", err))
		} else {
			if err := json.Unmarshal(jsonByte, &gConf); err != nil {
				tutil.LogWarn(fmt.Sprintf("读取配置文件出错\n错误原因:%v", err))
			} else {
				tutil.Log(fmt.Sprintf("域名:%v", gConf.Server.Domain))
				tutil.Log(fmt.Sprintf("端口:%v", gConf.Server.Port))
				tutil.Log("管理的博客域名:")
				for _, domain := range gConf.Site {
					tutil.Log(domain)
					gConf.siteDict[domain] = struct{}{}
				}
			}
		}
	}
}

func (cfg *TConf) IsSiteRegistered(site string) bool {
	_, ok := cfg.siteDict[site]
	return ok
}
