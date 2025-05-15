// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package core

import (
	"context"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gvalid"
)

type (
	ICaptcha interface {
		// NewAndStore 创建验证码，直接输出验证码图片内容到HTTP Response.
		NewAndStore() (id string, img string, err error)
		// VerifyAndClear 校验验证码，并清空缓存的验证码信息
		VerifyAndClear(id string, value string) bool
	}
	IModules interface {
		// IsDebug 是否为调试模式
		IsDebug(module ...string) bool
		// FormatName 格式化模块名称
		FormatName(module ...string) string
		// GetConfig 获取配置
		GetConfig(key string, def any, module ...string) *gvar.Var
	}
	IServer interface {
		// InitServer 初始化服务
		InitServer(svr *ghttp.Server)
		// AddModule 添加模块
		AddModule(module string, groupFunc func(group *ghttp.RouterGroup))
		// Start 启动服务
		Start(ctx context.Context, svr *ghttp.Server) (err error)
	}
	IValidate interface {
		RuleUnique(ctx context.Context, in gvalid.RuleFuncInput) error
	}
)

var (
	localCaptcha  ICaptcha
	localModules  IModules
	localServer   IServer
	localValidate IValidate
)

func Captcha() ICaptcha {
	if localCaptcha == nil {
		panic("implement not found for interface ICaptcha, forgot register?")
	}
	return localCaptcha
}

func RegisterCaptcha(i ICaptcha) {
	localCaptcha = i
}

func Modules() IModules {
	if localModules == nil {
		panic("implement not found for interface IModules, forgot register?")
	}
	return localModules
}

func RegisterModules(i IModules) {
	localModules = i
}

func Server() IServer {
	if localServer == nil {
		panic("implement not found for interface IServer, forgot register?")
	}
	return localServer
}

func RegisterServer(i IServer) {
	localServer = i
}

func Validate() IValidate {
	if localValidate == nil {
		panic("implement not found for interface IValidate, forgot register?")
	}
	return localValidate
}

func RegisterValidate(i IValidate) {
	localValidate = i
}
