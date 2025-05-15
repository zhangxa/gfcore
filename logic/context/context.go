package context

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gutil"
	"github.com/zhangxa/gfcore/core"
	"github.com/zhangxa/gfcore/fixed"
	"github.com/zhangxa/gfcore/store"
)

type sContext struct {
}

func init() {
	core.RegisterContext(NewContext())
}

// NewContext 上下文管理服务
func NewContext() core.IContext {
	return &sContext{}
}

// Init 初始化上下文对象指针到上下文对象中，以便后续的请求流程中可以修改。
func (s *sContext) Init(ctx context.Context) *store.Context {
	r := g.RequestFromCtx(ctx)
	customCtx := &store.Context{
		Ip:     r.GetClientIp(),
		Module: fixed.ModuleDefault,
		Data:   make(g.Map),
	}
	r.SetCtxVar(fixed.ContextKey, customCtx)
	return customCtx
}

// Get 获得上下文变量，如果没有设置，那么返回nil
func (s *sContext) Get(ctx context.Context) *store.Context {
	value := ctx.Value(fixed.ContextKey)
	if value != nil {
		if localCtx, ok := value.(*store.Context); ok {
			return localCtx
		}
	}
	return s.Init(ctx)
}

// GetAuthId 获取授权用户ID
func (s *sContext) GetAuthId(ctx context.Context) uint {
	return s.Get(ctx).AuthUser.AuthId
}

// SetUser 设置上下文用户信息
func (s *sContext) SetUser(ctx context.Context, ctxUser *store.JWTAuthInfo) {
	s.Get(ctx).AuthUser = ctxUser
}

// GetUser 获取上下文用户信息
func (s *sContext) GetUser(ctx context.Context) *store.JWTAuthInfo {
	return s.Get(ctx).AuthUser
}

// SetModule 设置应用标识
func (s *sContext) SetModule(ctx context.Context, module string) {
	s.Get(ctx).Module = module
}

func (s *sContext) GetModule(ctx context.Context) string {
	return s.Get(ctx).Module
}

// SetResponse 设置响应结果
func (s *sContext) SetResponse(ctx context.Context, response interface{}) {
	s.Get(ctx).HandlerResponse = response
}

// SetData 将上下文信息设置到上下文请求中，注意是完整覆盖
func (s *sContext) SetData(ctx context.Context, data g.Map) {
	s.Get(ctx).Data = data
}

// AddData 追加额外参数
func (s *sContext) AddData(ctx context.Context, key string, value interface{}) {
	s.Get(ctx).Data[key] = value
}

// GetData 获取额外参数
func (s *sContext) GetData(ctx context.Context, key string) interface{} {
	if value, ok := s.Get(ctx).Data[key]; ok {
		return value
	}
	return nil
}

// DelData 删除额外参数
func (s *sContext) DelData(ctx context.Context, key string) {
	delete(s.Get(ctx).Data, key)
}

// MergeData 合并额外参数
func (s *sContext) MergeData(ctx context.Context, data g.Map) {
	gutil.MapMerge(s.Get(ctx).Data, data)
}

// SetResMsg 设置响应消息
func (s *sContext) SetResMsg(ctx context.Context, msg string) {
	s.Get(ctx).ResMsg = msg
}

// SetAct 设置操作名称
func (s *sContext) SetAct(ctx context.Context, act string) {
	s.Get(ctx).ResAct = act
}

// SetResCode 将上下文信息设置到上下文请求中，注意是完整覆盖
func (s *sContext) SetResCode(ctx context.Context, code interface{}) {
	s.Get(ctx).ResCode = code
}

// SetResUrl 将上下文信息设置到上下文请求中，注意是完整覆盖
func (s *sContext) SetResUrl(ctx context.Context, url string) {
	s.Get(ctx).ResUrl = url
}

// SetResData 将上下文信息设置到上下文请求中，注意是完整覆盖
func (s *sContext) SetResData(ctx context.Context, data interface{}) {
	s.Get(ctx).ResData = data
}
