// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package core

import (
	"reflect"
)

type (
	IUtils interface {
		// InArray in_array()
		// haystack supported types: slice, array or map
		InArray(needle interface{}, haystack interface{}) bool
		// StringToArray 字符串分割成数组
		StringToArray(str string, sep ...string) []string
		// IsValidDomain 合法域名判断
		IsValidDomain(domain string, withHTTP bool) bool
		// EncryptStrings 字符串加盟
		EncryptStrings(str ...string) string
		// EmojiDecode 表情解码
		EmojiDecode(str string) string
		// EmojiEncode 表情转换
		EmojiEncode(str string) string
		// PadLeft pads left side of a string if size of string is less than indicated pad length
		PadLeft(str interface{}, padStr string, padLen int) string
		// PadRight pads right side of a string if size of string is less than indicated pad length
		PadRight(str interface{}, padStr string, padLen int) string
		// GetDeployDomain 获取部署域名 示例：http://www.example.com
		GetDeployDomain() string
		// GetDeployPath 获取部署路径 示例：/app
		GetDeployPath() string
		// GetModuleRoutePath 获取应用路由路径 示例：/api
		GetModuleRoutePath(module string) string
		// UrlAddDomain 文字追加域名前缀
		UrlAddDomain(url string, domainStr ...string) string
		// CallClassFunc 执行指定类的指定名称函数
		CallClassFunc(myClass interface{}, funcName string, params ...interface{}) (out []reflect.Value, err error)
	}
)

var (
	localUtils IUtils
)

func Utils() IUtils {
	if localUtils == nil {
		panic("implement not found for interface IUtils, forgot register?")
	}
	return localUtils
}

func RegisterUtils(i IUtils) {
	localUtils = i
}
