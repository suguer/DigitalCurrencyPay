package util

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"strings"
)

func PadLeft(s string, pad string, totalLen int) string {
	if len(s) >= totalLen {
		return s
	}
	padLen := totalLen - len(s)
	// 使用strings.Repeat重复填充字符串，然后与原始字符串拼接
	return strings.Repeat(pad, padLen) + s
}

func Md5BySalt(s, salt string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(s+salt)))
}
func StringToInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}
