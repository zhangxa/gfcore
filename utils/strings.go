package utils

import (
	"fmt"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/encoding/ghtml"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type sStrings struct {
}

var insStrings = sStrings{}

// Strings 字符串处理服务
var Strings = &insStrings

// IsInStrArray 检测值是否在,拼接的字符串中
func (s *sStrings) IsInStrArray(key interface{}, values string) bool {
	if values != "" {
		val := gconv.String(key)
		arr := strings.Split(values, ",")
		for _, v := range arr {
			if val == v {
				return true
			}
		}
	}
	return false
}

// InArray in_array()
// haystack supported types: slice, array or map
func (s *sStrings) InArray(needle interface{}, haystack interface{}) bool {
	val := reflect.ValueOf(haystack)
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if reflect.DeepEqual(needle, val.Index(i).Interface()) {
				return true
			}
		}
	case reflect.Map:
		for _, k := range val.MapKeys() {
			if reflect.DeepEqual(needle, val.MapIndex(k).Interface()) {
				return true
			}
		}
	default:
		panic("haystack: haystack type must be slice, array or map")
	}

	return false
}

// StringToArray 字符串分割成数组
func (s *sStrings) StringToArray(str string, sep ...string) []string {
	sepT := ","
	if len(sep) > 0 {
		sepT = sep[0]
	}
	return strings.Split(str, sepT)
}

// IsValidDomain 合法域名判断
func (s *sStrings) IsValidDomain(domain string, withHTTP bool) bool {
	if withHTTP {
		arr, err := gregex.MatchString(`http[s]{0,1}:\/\/([\w.]+\/?)\S*`, domain)
		if err != nil || len(arr) < 1 {
			return false
		}
		domain = arr[1]
	}
	return gregex.IsMatchString(`^([0-9a-zA-Z][0-9a-zA-Z\-]{0,62}\.)+([a-zA-Z]{0,62})$`, domain)
}

// EncryptStrings 字符串加盟
func (s *sStrings) EncryptStrings(str ...string) string {
	if len(str) == 0 {
		return ""
	}
	w := strings.Join(str, "")
	return gmd5.MustEncryptString(w)
}

// EmojiDecode 表情解码
func (s *sStrings) EmojiDecode(str string) string {
	//emoji表情的数据表达式
	re := regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+]")
	//提取emoji数据表达式
	reg := regexp.MustCompile("\\[\\\\u|]")
	src := re.FindAllString(str, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			str = strings.Replace(str, src[i], string(rune(p)), -1)
		}
	}
	return str
}

// EmojiEncode 表情转换
func (s *sStrings) EmojiEncode(str string) string {
	ret := ""
	rs := []rune(str)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u

		} else {
			ret += string(rs[i])
		}
	}
	return ret
}

// PadLeft pads left side of a string if size of string is less than indicated pad length
func (s *sStrings) PadLeft(str interface{}, padStr string, padLen int) string {
	return buildPadStr(gconv.String(str), padStr, padLen, true, false)
}

// PadRight pads right side of a string if size of string is less than indicated pad length
func (s *sStrings) PadRight(str interface{}, padStr string, padLen int) string {
	return buildPadStr(gconv.String(str), padStr, padLen, false, true)
}

// PadString either left, right or both sides.
// Note that padding string can be unicode and more than one character
func buildPadStr(str string, padStr string, padLen int, padLeft bool, padRight bool) string {

	// When padded length is less than the current string size
	if padLen < utf8.RuneCountInString(str) {
		return str
	}

	padLen -= utf8.RuneCountInString(str)

	targetLen := padLen

	targetLenLeft := targetLen
	targetLenRight := targetLen
	if padLeft && padRight {
		targetLenLeft = padLen / 2
		targetLenRight = padLen - targetLenLeft
	}

	strToRepeatLen := utf8.RuneCountInString(padStr)

	repeatTimes := int(math.Ceil(float64(targetLen) / float64(strToRepeatLen)))
	repeatedString := strings.Repeat(padStr, repeatTimes)

	leftSide := ""
	if padLeft {
		leftSide = repeatedString[0:targetLenLeft]
	}

	rightSide := ""
	if padRight {
		rightSide = repeatedString[0:targetLenRight]
	}

	return leftSide + str + rightSide
}

// ReplaceAllImgSrcWithDomain 将文字中所有图片连接追加域名前缀
func (s *sStrings) ReplaceAllImgSrcWithDomain(str string) string {
	url := Config.GetDeployUrl()
	reg := "<img[^>]*src\\s*=\\s*['\"]([\\w/\\-\\.]*)['\"][^>]*"
	res, err := gregex.MatchAllString(reg, str)
	if err != nil {
		return str
	}
	if len(res) > 0 {
		replaceMap := make(map[string]string)
		for _, v := range res {
			src := strings.ToLower(v[1])
			if !strings.HasPrefix(src, "http") && !strings.HasPrefix(src, "//") && strings.HasPrefix(src, "/") {
				replaceMap[v[1]] = fmt.Sprintf("%s%s", url, v[1])
			}
		}
		if len(replaceMap) > 0 {
			return gstr.ReplaceByMap(str, replaceMap)
		}
	}
	return str
}

// UrlAddParams URL参数添加
func (s *sStrings) UrlAddParams(oriUrl string, params g.Map) string {
	paramArr := make([]string, 0)
	for k, v := range params {
		paramArr = append(paramArr, fmt.Sprintf("%s=%v", k, v))
	}
	if strings.Contains(oriUrl, "?") {
		oriUrl += "&" + strings.Join(paramArr, "&")
	} else {
		oriUrl += "?" + strings.Join(paramArr, "&")
	}
	return oriUrl
}

// UrlAddDomain 文字追加域名前缀
func (s *sStrings) UrlAddDomain(url string, domainStr ...string) string {
	var domain string
	if len(domainStr) == 0 {
		domain = Config.GetDeployUrl()
	} else {
		domain = strings.TrimSuffix(domainStr[0], "/")
	}
	if !strings.HasPrefix(strings.ToLower(url), "http") && !strings.HasPrefix(url, "//") && strings.HasPrefix(url, "/") {
		return fmt.Sprintf("%s%s", domain, url)
	}
	return url
}

// GetStringLen 获取字符串长度，含中文
func (s *sStrings) GetStringLen(str string) int {
	return utf8.RuneCountInString(str)
}

// SubStrDecodeRuneInString 截取字符串长度，含中文
func (s *sStrings) SubStrDecodeRuneInString(str string, length int) string {
	var size, n int
	for i := 0; i < length && n < len(str); i++ {
		_, size = utf8.DecodeRuneInString(str[n:])
		n += size
	}

	return str[:n]
}

// GetOrderByFieldValues 获取mysql 根据field 数组 排序条件
func (s *sStrings) GetOrderByFieldValues(field string, valArr g.Slice) string {
	var orders []string
	orders = append(orders, fmt.Sprintf("field(%s", field))
	for _, v := range valArr {
		orders = append(orders, fmt.Sprintf("'%v'", v))
	}
	return fmt.Sprintf("%s)", strings.Join(orders, ","))
}

// Format2FolderName 获取文件夹名称
func (s *sStrings) Format2FolderName(name string) string {
	// 去除前后空格
	name = strings.TrimSpace(name)
	// 替换Windows/Linux下不允许的字符为下划线
	reg := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)
	name = reg.ReplaceAllString(name, "_")
	// 可选：限制长度，例如最大255个字符
	const maxLen = 255
	if len(name) > maxLen {
		name = name[:maxLen]
	}
	return name
}

// StrLimit returns a portion of string `str` specified by `length` parameters, if the length
// of `str` is greater than `length`, then the `suffix` will be appended to the result string.
// StrLimitRune considers parameter `str` as unicode string.
func (s *sStrings) StrLimit(str string, length int, suffix ...string) string {
	runes := []rune(str)
	if len(runes) < length {
		return str
	}
	suffixStr := "..."
	if len(suffix) > 0 {
		suffixStr = suffix[0]
	}
	length -= len(suffixStr)
	return string(runes[0:length]) + suffixStr
}

// StrLimitFromHtml 过滤掉html标签，并截取指定长度的字符串
func (s *sStrings) StrLimitFromHtml(html string, length int, suffix ...string) string {
	html = strings.Replace(html, "\n", "", -1)
	html = strings.Replace(html, "\r", "", -1)
	html = strings.Replace(html, "\t", "", -1)
	html = strings.Replace(html, " ", "", -1)
	html = strings.Replace(html, "&nbsp;", "", -1)
	html = strings.Replace(html, "&emsp;", "", -1)
	html = strings.Replace(html, "&ensp;", "", -1)
	html = strings.Replace(html, "&ldquo;", "", -1)
	html = strings.Replace(html, "&rdquo;", "", -1)
	html = strings.Replace(html, "&mdash;", "", -1)
	html = strings.Replace(html, "&hellip;", "", -1)
	html = strings.Replace(html, "&lsquo;", "", -1)
	html = strings.Replace(html, "&rsquo;", "", -1)
	html = strings.Replace(html, "&bull;", "", -1)
	html = strings.Replace(html, "&hellip;", "", -1)
	return s.StrLimit(ghtml.StripTags(html), length, suffix...)
}

// IsChinese 是否中文
func (s *sStrings) IsChinese(str string) bool {
	matched, _ := regexp.MatchString("[\\u4e00-\\u9fa5]", str)
	return matched
}

// FormatURL 格式化链接
func (s *sStrings) FormatURL(url string, module ...string) string {
	if url == "" {
		return "javascript:void(0);"
	}
	str := strings.ToLower(url)
	if strings.HasPrefix(str, "javascript") || strings.HasPrefix(str, "http") || strings.HasPrefix(str, "?") || strings.HasPrefix(str, "/") {
		return url
	}
	prefix := ""
	if len(module) > 0 {
		prefix = Config.GetModuleUrlPrefix(module[0])
	}
	if prefix == "" {
		return fmt.Sprintf("/%s", url)
	}
	return fmt.Sprintf("/%s/%s", prefix, url)
}

// FormatURLWithDomain 格式化链接
func (s *sStrings) FormatURLWithDomain(url string, module ...string) string {
	if url == "" {
		return "javascript:void(0);"
	}
	str := strings.ToLower(url)
	if strings.HasPrefix(str, "javascript") || strings.HasPrefix(str, "http") {
		return url
	}
	prefix := ""
	if len(module) > 0 {
		prefix = Config.GetModuleUrlPrefix(module[0], true)
	} else {
		prefix = Config.GetDeployUrl()
	}
	return fmt.Sprintf("%s/%s", prefix, url)
}

// FormatMac Mac标准化格式
func (s *sStrings) FormatMac(macStr string) string {
	if strings.Count(macStr, ":") == 5 {
		return macStr
	}
	mac := ""
	macArr := make([]string, 0)
	arr := strings.Split(macStr, "")
	for i, v := range arr {
		mac += v
		if (i+1)%2 == 0 {
			macArr = append(macArr, mac)
			mac = ""
		}
	}
	return strings.Join(macArr, ":")
}
