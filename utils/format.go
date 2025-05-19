package utils

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
)

type sFormat struct {
}

var Format = &sFormat{}

// FormFieldAttr 字符串配置参数值转换至数组
func (s *sFormat) FormFieldAttr(str string, sep ...string) g.Map {
	res := g.Map{}
	if str == "" {
		return res
	}
	sep1 := "<br>"
	sep2 := ":"
	if len(sep) == 1 {
		sep1 = sep[0]
	}
	if len(sep) > 1 {
		sep1 = sep[0]
		sep2 = sep[1]
	}
	arr1 := strings.Split(gstr.Nl2Br(str), sep1)
	if len(arr1) > 0 {
		for _, v := range arr1 {
			arr2 := strings.Split(v, sep2)
			if len(arr2) > 1 {
				res[arr2[0]] = arr2[1]
			}
		}
	}
	return res
}

// FormSelectOptionStrToList 将配置文本转为选择列表
func (s *sFormat) FormSelectOptionStrToList(str string) (list g.List) {
	if str != "" {
		sep1 := "<br>"
		sep2 := ":"
		arr1 := strings.Split(gstr.Nl2Br(str), sep1)
		for _, v := range arr1 {
			arr2 := strings.Split(v, sep2)
			if len(arr2) > 1 {
				list = append(list, g.Map{
					"value": arr2[0],
					"title": arr2[1],
				})
			}
		}
	}
	return
}

// FormatBytes 字节格式化为带单位
func (s *sFormat) FormatBytes(fileSize int64) string {
	if fileSize < 1024 {
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else { //if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

// Tree2Table 将树形列表转为table展示列表
func (s *sFormat) Tree2Table(trees interface{}, topPk interface{}, pkField, childField, PNode string) (list []g.Map) {
	data := gconv.Maps(trees)
	for _, v := range data {
		v["nodes"] = PNode + "-" + gconv.String(v[pkField])
		v["spl"] = strings.Repeat("&nbsp├&nbsp;", strings.Count(PNode, "-"))
		list = append(list, v)
		children := gconv.Maps(v[childField])
		delete(v, childField)
		if len(children) > 0 {
			suList := s.Tree2Table(children, topPk, pkField, childField, v["nodes"].(string))
			for _, vv := range suList {
				list = append(list, vv)
			}
		}
	}
	return
}

// ArrayToTree 数组转树形结构
func (s *sFormat) ArrayToTree(list g.List, topID interface{}, idField string, pidField string, childField string) g.List {
	nodes := make(g.List, 0)
	var topNode g.Map
	if len(list) > 0 {
		for _, v := range list {
			if !g.IsEmpty(topID) && gconv.String(v[idField]) == gconv.String(topID) {
				topNode = v
			}
			if gconv.String(topID) == gconv.String(v[pidField]) {
				sub := v
				s.MakeTree(list, sub, idField, pidField, childField)
				nodes = append(nodes, sub)
			}
		}
	}
	if topNode != nil {
		topNode[childField] = nodes
		return g.List{topNode}
	}
	return nodes
}

// GetTree 列表转为树形结构
func (s *sFormat) GetTree(data interface{}, topId interface{}, idKey, pidKey, childrenKey string, resPointer interface{}) {
	list := gconv.Maps(data)
	var top g.Map
	result := make(g.List, 0)
	if len(list) > 0 {
		for _, v := range list {
			if !g.IsEmpty(topId) && fmt.Sprintf("%v", v[idKey]) == fmt.Sprintf("%v", topId) {
				top = v
			}
			if fmt.Sprintf("%v", topId) == fmt.Sprintf("%v", v[pidKey]) {
				sub := v
				s.MakeTree(list, sub, idKey, pidKey, childrenKey)
				result = append(result, sub)
			}
		}
	}
	if top != nil {
		top[childrenKey] = result
		result = g.List{top}
	}
	_ = gconv.Scan(result, resPointer)
	return
}

// MakeTree 获取节点树
func (s *sFormat) MakeTree(list g.List, node g.Map, idKey, pidKey, childrenKey string) { //参数为父节点，添加父节点的子节点指针切片
	children, hav := s.haveChild(list, node, idKey, pidKey) //判断节点是否有子节点并返回
	if hav {
		node[childrenKey] = children //添加子节点
		for _, v := range children { //查询子节点的子节点，并添加到子节点
			s.MakeTree(list, v, idKey, pidKey, childrenKey) //递归添加节点
		}
	}
}

// haveChild 判断是否又子菜单，并返回
func (s *sFormat) haveChild(list g.List, node g.Map, idKey, pidKey string) (child []g.Map, yes bool) {
	for _, v := range list {
		if fmt.Sprintf("%v", v[pidKey]) == fmt.Sprintf("%v", node[idKey]) {
			child = append(child, v)
		}
	}
	yes = len(child) > 0
	return
}

// FormatMac Mac标准化格式
func (s *sFormat) FormatMac(macStr string) string {
	return Strings.FormatMac(macStr)
}

// FormatURL url标准化格式
func (s *sFormat) FormatURL(url string, module ...string) string {
	return Strings.FormatURL(url, module...)
}

// ListToMap 获取列表内指定字段Map
func (s *sFormat) ListToMap(data interface{}, keyField string, dataField ...string) g.Map {
	var res = g.Map{}
	if len(dataField) == 0 {
		return res
	}
	dataList := gconv.Maps(data)
	for _, v := range dataList {
		keyStr := gconv.String(v[keyField])
		if len(dataField) == 1 {
			res[keyStr] = v[dataField[0]]
		} else {
			mp := g.Map{}
			for _, f := range dataField {
				mp[f] = v[f]
			}
			res[keyStr] = mp
		}
	}
	return res
}
