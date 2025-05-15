package utils

import (
	"github.com/gogf/gf/v2/frame/g"
)

type sView struct {
}

// View 视图公共函数
var View = &sView{}

// DefaultViewValue 判断是否为空，为空则追加默认值
func (s *sView) DefaultViewValue(def interface{}, val interface{}) interface{} {
	if g.IsEmpty(val) {
		return def
	}
	return val
}
