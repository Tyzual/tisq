package util

import (
	"log"
	"runtime"
)

/*
PrintTrace 是否打印内部日志
*/
var PrintTrace = false

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

/*
LogTrace 打印内部日志
*/
func LogTrace(msg string) {
	if !PrintTrace {
		return
	}

	_, fileName, fileLine, ok := runtime.Caller(1)
	if ok {
		log.Println("[", fileName, ":", fileLine, "]", msg)
	} else {
		log.Println(msg)
	}
}
