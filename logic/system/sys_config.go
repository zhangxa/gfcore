package system

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/zhangxa/gfcore/core"
	"github.com/zhangxa/gfcore/fixed"
	"time"
)

type sSysConfig struct {
}

func init() {
	core.RegisterSysConfig(NewSysConfig())
}

// NewSysConfig 系统配置服务
func NewSysConfig() core.ISysConfig {
	return &sSysConfig{}
}

// Dict 配置字典集合
func (s *sSysConfig) Dict(ctx context.Context) (dict map[string]string) {
	cached, _ := gcache.Get(ctx, fixed.CacheSysConfig)
	if cached.IsEmpty() {
		data, _ := g.Model("sys_config").Ctx(ctx).
			Where("status", 1).All()
		if !data.IsEmpty() {
			for _, v := range data {
				dict[v["name"].String()] = v["value"].String()
			}
			_ = gcache.Set(ctx, fixed.CacheSysConfig, dict, 1*time.Hour)
		}
		return dict
	}
	return cached.MapStrStr()
}

// Get 根据key 获取对应配置值
func (s *sSysConfig) Get(key string, def ...string) string {
	res := ""
	if len(def) > 0 {
		res = def[0]
	}
	if key == "" {
		return res
	}
	dict := s.Dict(context.Background())
	if v, ok := dict[key]; ok {
		return v
	}
	return res
}

// ClearDbcCache 清除dbc cache
func (s *sSysConfig) ClearDbcCache(ctx context.Context) {
	_, _ = gcache.Remove(ctx, fixed.CacheSysConfig)
}
