package utils

import (
	"crypto/md5"
	"fmt"
)

func GetMd5(s string) string {
	// %x转化成16进制表示
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}
