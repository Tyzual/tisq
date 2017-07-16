package util

import "log"

/*
LogFatal 严重错误
*/
func LogFatal(msg string) {
	log.Fatal(msg)
}

/*
LogWarn 非严重错误
*/
func LogWarn(msg string) {
	log.Println(msg)
}

/*
Log 普通输出
*/
func Log(msg string) {
	log.Println(msg)
}
