package utils

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

func MustNil(v interface{}) {
	if v != nil {
		panic(v)
	}
}

func StructToJson(v interface{}) string {
	if v == nil {
		return ""
	}
	bytes, _ := json.Marshal(v)
	return string(bytes)
}

// FirstUpper 字符串首字母大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// FirstLower 字符串首字母小写
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func GenConvId(userId1, userId2 int32) string {
	if userId1 < userId2 {
		return fmt.Sprintf("%d_%d", userId1, userId2)
	}
	return fmt.Sprintf("%d_%d", userId2, userId1)
}

func Md5(content string) (md string) {
	h := md5.New()
	_, _ = io.WriteString(h, content)
	md = fmt.Sprintf("%x", h.Sum(nil))
	return
}