// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package core

import (
	"context"
)

type (
	ISysConfig interface {
		// Dict 配置字典集合
		Dict(ctx context.Context) (dict map[string]string)
		// Get 根据key 获取对应配置值
		Get(key string, def ...string) string
		// ClearDbcCache 清除dbc cache
		ClearDbcCache(ctx context.Context)
	}
)

var (
	localSysConfig ISysConfig
)

func SysConfig() ISysConfig {
	if localSysConfig == nil {
		panic("implement not found for interface ISysConfig, forgot register?")
	}
	return localSysConfig
}

func RegisterSysConfig(i ISysConfig) {
	localSysConfig = i
}
