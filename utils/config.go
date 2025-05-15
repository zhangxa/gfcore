package utils

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"strings"
)

type sConfig struct {
}

// Config 配置服务
var Config = new(sConfig)

// GetDeployDomain 获取部署域名
func (s *sConfig) GetDeployDomain() string {
	ctx := context.Background()
	return strings.TrimSuffix(g.Config().MustGet(ctx, "core.deployDomain").String(), "/")
}

// GetDeployPath 获取部署路径
func (s *sConfig) GetDeployPath() string {
	ctx := context.Background()
	return strings.TrimPrefix(strings.TrimSuffix(g.Config().MustGet(ctx, "core.deployPath").String(), "/"), "/")
}

// GetDeployUrl 获取应用部署链接
func (s *sConfig) GetDeployUrl() string {
	domain := s.GetDeployDomain()
	path := s.GetDeployPath()
	if path != "" {
		return fmt.Sprintf("%s/%s", domain, path)
	}
	return domain
}

// GetModuleRoutePath 获取应用路由路径
func (s *sConfig) GetModuleRoutePath(module string) string {
	ctx := context.Background()
	pattern := fmt.Sprintf("modules.%s.routePath", module)
	return strings.TrimPrefix(strings.TrimSuffix(g.Config().MustGet(ctx, pattern).String(), "/"), "/")
}

// GetModuleUrlPrefix 获取应用路由路径
func (s *sConfig) GetModuleUrlPrefix(module string, addDomain ...bool) string {
	prefix := ""
	if len(addDomain) > 0 && addDomain[0] {
		prefix = s.GetDeployUrl()
	} else {
		prefix = s.GetDeployPath()
	}
	routePath := s.GetModuleRoutePath(module)
	if prefix != "" {
		return fmt.Sprintf("%s/%s", prefix, routePath)
	}
	return routePath
}

// GetDefaultPageSize 获取默认分页大小
func (s *sConfig) GetDefaultPageSize(def ...int) int {
	ctx := context.Background()
	defPS := 20
	if len(def) > 0 && def[0] > 0 {
		defPS = def[0]
	}
	return g.Config().MustGet(ctx, "core.defaultPageSize", defPS).Int()
}
