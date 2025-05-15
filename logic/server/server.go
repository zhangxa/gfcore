package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gmode"
	"github.com/gogf/gf/v2/util/gvalid"
	"github.com/zhangxa/gfcore/core"
	"github.com/zhangxa/gfcore/plugins/swagger"
)

type sServer struct {
	servers map[string]*ghttp.Server
	modules map[string]func(group *ghttp.RouterGroup)
}

func init() {
	core.RegisterServer(NewServer())
}

// NewServer 服务器基础服务
func NewServer() core.IServer {
	return &sServer{
		servers: make(map[string]*ghttp.Server),
		modules: make(map[string]func(group *ghttp.RouterGroup)),
	}
}

// getServerName 获取服务名称
func (s *sServer) getServerName(name ...interface{}) string {
	instanceName := ghttp.DefaultServerName
	if len(name) > 0 && name[0] != "" {
		instanceName = gconv.String(name[0])
	}
	return instanceName
}

// InitServer 初始化服务
func (s *sServer) InitServer(name ...interface{}) *ghttp.Server {
	instanceName := s.getServerName(name...)
	if _, ok := s.servers[instanceName]; !ok {
		s.servers[instanceName] = g.Server(instanceName)
		if gmode.IsDevelop() {
			s.servers[instanceName].BindHookHandler("/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
				r.Response.Header().Set("Cache-Control", "no-store")
			})
		}
		gvalid.RegisterRule("unique", core.Validate().RuleUnique)
	}
	return s.servers[instanceName]
}

// AddModule 添加模块
func (s *sServer) AddModule(module string, groupFunc func(group *ghttp.RouterGroup)) {
	s.modules[module] = groupFunc
}

// Start 启动服务
func (s *sServer) Start(ctx context.Context, name ...interface{}) (err error) {
	instanceName := s.getServerName(name...)
	svr, ok := s.servers[instanceName]
	if !ok {
		return errors.New("server not init")
	}
	svr.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(
			core.Middleware().Base,
			core.Middleware().HandlerResponse,
		)
		for module, groupFunc := range s.modules {
			prefix := g.Config().MustGet(ctx, fmt.Sprintf("modules.%s.routePath", module)).String()
			if prefix == "" {
				err = fmt.Errorf("module %s routePath is empty", module)
				return
			}
			group.Group(fmt.Sprintf("/%s", prefix), func(group *ghttp.RouterGroup) {
				group.Hook("/*", ghttp.HookBeforeServe, func(r *ghttp.Request) {
					core.Context().SetModule(r.Context(), module)
				})
				group.Middleware(core.Middleware().VisitLimit)
				groupFunc(group)
			})
		}
	})
	// Custom enhance API document.
	svr.Plugin(&swagger.Swagger{})
	svr.Run()
	return nil
}
