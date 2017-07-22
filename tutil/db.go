package tutil

import (
	"time"
)

/*
TimeToDateTimeString 将Time转换成MysqlDateTime认识的字符串
*/
func TimeToDateTimeString(t *time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
