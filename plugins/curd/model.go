package curd

import "github.com/gogf/gf/v2/database/gdb"

// PageData 列表数据结构体
type PageData[T any] struct {
	List  []T    `json:"list"  dc:"数据列表"`
	Pager *Pager `json:"pager" dc:"分页信息"`
}

// Pager 列表分页结构体
type Pager struct {
	TotalSize   int `json:"total_size"   dc:"数据总量"`
	PageSize    int `json:"page_size"    dc:"单页数据量"`
	TotalPage   int `json:"total_page"   dc:"总页数"`
	CurrentPage int `json:"current_page" dc:"当前页码"`
}

// PageDataInput 分页请求
type PageDataInput struct {
	Model   *gdb.Model  `dc:"数据模型，如：dao.About.Ctx(ctx)"`
	Where   interface{} `dc:"查询条件，按需填写"`
	Fields  string      `dc:"查询字段，按需填写"`
	Order   string      `dc:"排序规则，按需填写"`
	Limit   int         `dc:"单页数据量，默认20"`
	CurPage int         `dc:"当前页码，默认1"`
}

// JoinPageDataInput 联合查询分页请求
type JoinPageDataInput struct {
	*PageDataInput
	TableAlia     string
	JoinType      JoinType // 联合类型，left/right/inner...
	JoinTable     string
	JoinTableAlia string
	JoinOn        string
}

// ListDataInput 列表请求
type ListDataInput struct {
	Model  *gdb.Model  `dc:"数据模型，如：dao.About.Ctx(ctx)"`
	Where  interface{} `dc:"查询条件，按需填写"`
	Fields string      `dc:"查询字段，按需填写"`
	Order  string      `dc:"排序规则，按需填写"`
}

// JoinListDataInput 联合查询列表请求
type JoinListDataInput struct {
	*ListDataInput
	TableAlia     string
	JoinType      JoinType // 联合类型，left/right/inner...
	JoinTable     string
	JoinTableAlia string
	JoinOn        string
}

// SingleDataInput 单数据请求
type SingleDataInput struct {
	Model  *gdb.Model  `dc:"数据模型，如：dao.About.Ctx(ctx)"`
	Where  interface{} `dc:"查询条件，按需填写"`
	Fields string      `dc:"查询字段，按需填写"`
	Order  string      `dc:"排序规则，按需填写"`
}

// JoinSingleDataInput 联合单数据请求
type JoinSingleDataInput struct {
	*SingleDataInput
	TableAlia     string
	JoinType      JoinType // 联合类型，left/right/inner...
	JoinTable     string
	JoinTableAlia string
	JoinOn        string
}
