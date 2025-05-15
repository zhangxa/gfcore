package store

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Context 请求上下文结构
type Context struct {
	Module          string        // 应用模块
	Ip              string        // 请求IP
	TakeUpTime      int64         // 请求耗时 ms
	AuthUser        *AuthUserInfo // 上下文用户信息
	HandlerResponse interface{}   // 组件响应
	Data            g.Map         // 自定KV变量，业务模块根据需要设置，不固定
	ResMsg          string        // 自定义接口返回消息
	ResAct          string        // 自定义操作名称，如：新增，修改，删除等，自动生成返回消息，如新增成功、失败
	ResCode         any           // 自定义接口返回值
	ResUrl          string        // 自定义接口跳转链接
	ResData         any           // 自定义返回内容
}

// AuthUserInfo 授权用户信息
type AuthUserInfo struct {
	AuthId    uint   `json:"auth_id" dc:"登录用户id"`
	UserName  string `json:"user_name,omitempty" dc:"用户名称"`
	LoginType string `json:"login_type,omitempty" dc:"登录类型"`
	RoleId    uint   `json:"role_id,omitempty" dc:"所属角色ID"`
	RoleTitle string `json:"role_title,omitempty" dc:"所属角色描述"`
	OpenId    string `json:"open_id,omitempty" dc:"微信用户标识"`
	UnionId   string `json:"union_id,omitempty" dc:"微信联合用户标识"`
	SystemId  uint   `json:"system_id,omitempty" dc:"独立系统ID，如单灯独立系统"`
}

// JWTAuthResult JWT授权结果
type JWTAuthResult struct {
	AccessToken  string `json:"access_token" dc:"授权token"`
	RefreshToken string `json:"refresh_token,omitempty" dc:"刷新token"`
	ExpiresIn    int64  `json:"expires_in" dc:"授权到期时间"`
}

// DefaultHandlerResponse 默认api返回结构体
type DefaultHandlerResponse struct {
	Code      interface{} `json:"code" example:"0"    dc:"错误代码：0成功，!0其它错误"`
	Msg       string      `json:"msg" example:"操作成功" dc:"错误信息"`
	Data      interface{} `json:"data"    dc:"返回数据"`
	Url       string      `json:"url,omitempty" example:"/index" dc:"跳转链接"`
	Success   bool        `json:"success" example:"true"  dc:"是否成功"`
	Timestamp int64       `json:"timestamp,omitempty" example:"1640966400" dc:"服务器时间戳"`
}
