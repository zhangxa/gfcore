package utils

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/zhangxa/gfcore/core"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type sUtils struct {
	ctx context.Context
}

func init() {
	core.RegisterUtils(NewUtils())
}

// NewUtils 基础工具服务
func NewUtils() core.IUtils {
	return &sUtils{
		ctx: context.Background(),
	}
}

// InArray in_array()
// haystack supported types: slice, array or map
func (s *sUtils) InArray(needle interface{}, haystack interface{}) bool {
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
func (s *sUtils) StringToArray(str string, sep ...string) []string {
	sepT := ","
	if len(sep) > 0 {
		sepT = sep[0]
	}
	return strings.Split(str, sepT)
}

// IsValidDomain 合法域名判断
func (s *sUtils) IsValidDomain(domain string, withHTTP bool) bool {
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
func (s *sUtils) EncryptStrings(str ...string) string {
	if len(str) == 0 {
		return ""
	}
	w := strings.Join(str, "")
	return gmd5.MustEncryptString(w)
}

// EmojiDecode 表情解码
func (s *sUtils) EmojiDecode(str string) string {
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
func (s *sUtils) EmojiEncode(str string) string {
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
func (s *sUtils) PadLeft(str interface{}, padStr string, padLen int) string {
	return s.buildPadStr(gconv.String(str), padStr, padLen, true, false)
}

// PadRight pads right side of a string if size of string is less than indicated pad length
func (s *sUtils) PadRight(str interface{}, padStr string, padLen int) string {
	return s.buildPadStr(gconv.String(str), padStr, padLen, false, true)
}

// PadString either left, right or both sides.
// Note that padding string can be Unicode and more than one character
func (s *sUtils) buildPadStr(str string, padStr string, padLen int, padLeft bool, padRight bool) string {

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

// GetDeployDomain 获取部署域名 示例：http://www.example.com
func (s *sUtils) GetDeployDomain() string {
	return strings.TrimSuffix(g.Config().MustGet(s.ctx, "core.deployDomain").String(), "/")
}

// GetDeployPath 获取部署路径 示例：/app
func (s *sUtils) GetDeployPath() string {
	path := strings.TrimSuffix(g.Config().MustGet(s.ctx, "core.deployPath").String(), "/")
	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}

// GetModuleRoutePath 获取应用路由路径 示例：/api
func (s *sUtils) GetModuleRoutePath(module string) string {
	pattern := fmt.Sprintf("modules.%s.routePath", module)
	path := strings.TrimSuffix(g.Config().MustGet(s.ctx, pattern).String(), "/")
	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}

// UrlAddDomain 文字追加域名前缀
func (s *sUtils) UrlAddDomain(url string, domainStr ...string) string {
	var domain string
	if len(domainStr) == 0 {
		domain = s.GetDeployDomain()
	} else {
		domain = strings.TrimSuffix(domainStr[0], "/")
	}
	if !strings.HasPrefix(strings.ToLower(url), "http") && !strings.HasPrefix(url, "//") && strings.HasPrefix(url, "/") {
		return fmt.Sprintf("%s%s", domain, url)
	}
	return url
}

// CallClassFunc 执行指定类的指定名称函数
func (s *sUtils) CallClassFunc(myClass interface{}, funcName string, params ...interface{}) (out []reflect.Value, err error) {
	myClassValue := reflect.ValueOf(myClass)
	m := myClassValue.MethodByName(funcName)
	if !m.IsValid() {
		return make([]reflect.Value, 0), fmt.Errorf("method not found \"%s\"", funcName)
	}
	in := make([]reflect.Value, len(params))
	for i, param := range params {
		in[i] = reflect.ValueOf(param)
	}
	out = m.Call(in)
	return
}
