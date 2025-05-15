package utils

import (
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/mozillazg/go-pinyin"
	"regexp"
	"strings"
)

type sCharacter struct {
}

var insCharacter = sCharacter{}

// Character 汉字拼音及字符长度处理服务
var Character = &insCharacter

// ToPinyin 获取拼音
func (s *sCharacter) ToPinyin(text string, limit ...int) string {
	a := pinyin.NewArgs()
	arr := pinyin.Pinyin(text, a)
	lmt := 0
	if len(limit) > 0 {
		lmt = limit[0]
	}
	count := 0
	result := make([]string, 0)
	for _, v := range arr {
		if len(v) > 0 {
			count += len(v[0])
			if lmt > 0 && count > lmt {
				break
			}
			result = append(result, v[0])
		}
	}
	res := strings.Join(result, "")
	return res
}

// ToPinyinShort 获取拼音简写
func (s *sCharacter) ToPinyinShort(text string, limit ...int) string {
	a := pinyin.NewArgs()
	arr := pinyin.Pinyin(text, a)
	lmt := 0
	if len(limit) > 0 {
		lmt = limit[0]
	}
	count := 0
	result := make([]string, 0)
	for _, v := range arr {
		if len(v) > 0 {
			count += 1
			if lmt > 0 && count > lmt {
				break
			}
			result = append(result, gstr.StrLimit(v[0], 1, ""))
		}
	}
	res := strings.Join(result, "")
	return res
}

// TrimHtml 去除html标签
func (s *sCharacter) TrimHtml(src string, oneSentence ...bool) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile(`<[\S\s]+?>`)
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除空格
	re, _ = regexp.Compile(`&nbsp`)
	src = re.ReplaceAllString(src, "")
	re, _ = regexp.Compile(` `)
	src = re.ReplaceAllString(src, "")
	re, _ = regexp.Compile(`　`)
	src = re.ReplaceAllString(src, "")
	//去除STYLE
	re, _ = regexp.Compile(`<style[\S\s]+?</style>`)
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile(`<script[\S\s]+?</script>`)
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile(`<[\S\s]+?>`)
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile(`\s{2,}`)
	src = re.ReplaceAllString(src, "\n")
	if len(oneSentence) > 0 && oneSentence[0] {
		//去除换行符
		re, _ = regexp.Compile(`\r\n`)
		src = re.ReplaceAllString(src, "")
		re, _ = regexp.Compile(`\t`)
		src = re.ReplaceAllString(src, "")
		re, _ = regexp.Compile(`\n`)
		src = re.ReplaceAllString(src, "")
	}
	return strings.TrimSpace(src)
}

// TrimHtmlWithLimit 去除html标签并限制长度
func (s *sCharacter) TrimHtmlWithLimit(src string, limit int) string {
	str := s.TrimHtml(src, true)
	return s.TextLimit(str, limit)
}

// TextLimit 文本限制长度
func (s *sCharacter) TextLimit(src string, limit int) string {
	return gstr.StrLimitRune(src, limit)
}
