// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package core

import (
	"context"
)

type (
	IIpLimit interface {
		// ClearCache 清除黑白名单缓存
		ClearCache(ctx context.Context, module ...string)
		// IsLimited 检测访问是否受限
		IsLimited(ctx context.Context, module ...string) bool
	}
)

var (
	localIpLimit IIpLimit
)

func IpLimit() IIpLimit {
	if localIpLimit == nil {
		panic("implement not found for interface IIpLimit, forgot register?")
	}
	return localIpLimit
}

func RegisterIpLimit(i IIpLimit) {
	localIpLimit = i
}
