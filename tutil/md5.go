package tutil

import (
	"crypto/md5"
	"fmt"
)

/*
MD5 计算MD5
*/
func MD5(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}
