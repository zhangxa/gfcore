// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package core

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/zhangxa/gfcore/store"
)

type (
	IContext interface {
		// Init 初始化上下文对象指针到上下文对象中，以便后续的请求流程中可以修改。
		Init(ctx context.Context) *store.Context
		// Get 获得上下文变量，如果没有设置，那么返回nil
		Get(ctx context.Context) *store.Context
		// GetAuthId 获取授权用户ID
		GetAuthId(ctx context.Context) uint
		// SetUser 设置上下文用户信息
		SetUser(ctx context.Context, ctxUser *store.JWTAuthInfo)
		// GetUser 获取上下文用户信息
		GetUser(ctx context.Context) *store.JWTAuthInfo
		// SetModule 设置应用标识
		SetModule(ctx context.Context, module string)
		GetModule(ctx context.Context) string
		// SetResponse 设置响应结果
		SetResponse(ctx context.Context, response interface{})
		// SetData 将上下文信息设置到上下文请求中，注意是完整覆盖
		SetData(ctx context.Context, data g.Map)
		// AddData 追加额外参数
		AddData(ctx context.Context, key string, value interface{})
		// GetData 获取额外参数
		GetData(ctx context.Context, key string) interface{}
		// DelData 删除额外参数
		DelData(ctx context.Context, key string)
		// MergeData 合并额外参数
		MergeData(ctx context.Context, data g.Map)
		// SetResMsg 设置响应消息
		SetResMsg(ctx context.Context, msg string)
		// SetAct 设置操作名称
		SetAct(ctx context.Context, act string)
		// SetResCode 将上下文信息设置到上下文请求中，注意是完整覆盖
		SetResCode(ctx context.Context, code interface{})
		// SetResData 将上下文信息设置到上下文请求中，注意是完整覆盖
		SetResData(ctx context.Context, data interface{})
	}
)

var (
	localContext IContext
)

func Context() IContext {
	if localContext == nil {
		panic("implement not found for interface IContext, forgot register?")
	}
	return localContext
}

func RegisterContext(i IContext) {
	localContext = i
}
