package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"tisq/util"
)

var gConf = TConf{Domain: ""}

const confFile = "/etc/tisq.cnf"

/*
TConf 储存服务器的全局配置
*/
type TConf struct {
	Domain string
	Port   uint16
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
		gConf.Port = 34958
		gConf.Domain = "localhost"
		jsonByte, err := json.Marshal(gConf)
		if err != nil {
			util.LogWarn(fmt.Sprintf("创建配置文件出错\n错误原因:%v", err))
		} else {
			if err = ioutil.WriteFile(confFile, jsonByte, 0640); err != nil {
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
				util.Log(fmt.Sprintf("域名:%v", gConf.Domain))
				util.Log(fmt.Sprintf("端口:%v", gConf.Port))
			}
		}
	}
}
