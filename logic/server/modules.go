package server

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/zhangxa/gfcore/core"
	"github.com/zhangxa/gfcore/fixed"
)

type sModules struct {
}

func init() {
	core.RegisterModules(NewModules())
}

// NewModules 应用模块服务
func NewModules() core.IModules {
	return &sModules{}
}

// IsDebug 是否为调试模式
func (s *sModules) IsDebug(module ...string) bool {
	md := s.FormatName(module...)
	pattern := fmt.Sprintf("modules.%s.isDebug", md)
	return g.Config().MustGet(context.Background(), pattern, false).Bool()
}

// FormatName 格式化模块名称
func (s *sModules) FormatName(module ...string) string {
	md := fixed.ModuleDefault
	if len(module) > 0 {
		md = module[0]
	}
	return md
}

// GetConfig 获取配置
func (s *sModules) GetConfig(key string, def any, module ...string) *gvar.Var {
	md := s.FormatName(module...)
	pattern := fmt.Sprintf("modules.%s.%s", md, key)
	return g.Config().MustGet(context.Background(), pattern, def)
}
