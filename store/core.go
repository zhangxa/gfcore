package store

import (
	"github.com/gogf/gf/v2/frame/g"
)

type CurdDataItem any

// Context 请求上下文结构
type Context struct {
	Module          string       // 应用模块
	Ip              string       // 请求IP
	TakeUpTime      int64        // 请求耗时 ms
	AuthUser        *JWTAuthInfo // 上下文用户信息
	HandlerResponse interface{}  // 组件响应
	Data            g.Map        // 自定KV变量，业务模块根据需要设置，不固定
	ResMsg          string       // 自定义接口返回消息
	ResAct          string       // 自定义操作名称，如：新增，修改，删除等，自动生成返回消息，如新增成功、失败
	ResCode         any          // 自定义接口返回值
	ResData         any          // 自定义返回内容
}

// JWTAuthInfo JWT保存的授权用户信息
type JWTAuthInfo struct {
	AuthId    uint   `json:"authId"               dc:"登录用户id" v:"required"`
	Account   string `json:"account,omitempty"    dc:"登录账号"`
	LoginType string `json:"loginType,omitempty"  dc:"登录类型"`
	OpenId    string `json:"openId,omitempty"     dc:"第三方应用用户标识"`
	UnionId   string `json:"unionId,omitempty"    dc:"第三方应用联合用户标识"`
	Extra     g.Map  `json:"extra,omitempty"      dc:"附加信息"`
}

// JWTAuthResult JWT授权结果
type JWTAuthResult struct {
	AccessToken  string `json:"accessToken"            dc:"授权token"`
	RefreshToken string `json:"refreshToken,omitempty" dc:"刷新token"`
	ExpiresIn    int64  `json:"expiresIn"              dc:"授权到期时间"`
}

// DefaultHandlerResponse 默认api返回结构体
type DefaultHandlerResponse struct {
	Code interface{} `json:"code"                dc:"错误代码：0成功，!0其它错误" example:"0"`
	Msg  string      `json:"msg"                 dc:"错误信息" example:"操作成功"`
	Data interface{} `json:"data"                dc:"返回数据"`
}
