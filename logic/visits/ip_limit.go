package visits

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/zhangxa/gfcore/core"
	"strings"
)

type sIpLimit struct {
}

func init() {
	core.RegisterIpLimit(NewIpLimit())
}

// NewIpLimit ip访问限制服务
func NewIpLimit() core.IIpLimit {
	return &sIpLimit{}
}

const (
	blockListCache = "block_list"
	allowListCache = "allow_list"
)

// equalIP 对比IP是否相等
func (s *sIpLimit) equalIP(requestIP string, ruleIP string) bool {
	if strings.Contains(ruleIP, "*") {
		arr := strings.Split(ruleIP, "*")
		return strings.HasPrefix(requestIP, arr[0])
	}
	return requestIP == ruleIP
}

// getLimitParams 获取当前模块所属限制类型数据表名称及缓存名称
func (s *sIpLimit) getLimitParams(module ...string) (limitType int, resModule string, cacheName string) {
	resModule = core.Modules().FormatName(module...)
	limitType = core.Modules().GetConfig("limitType", 0, resModule).Int()
	switch limitType {
	case 1:
		cacheName = fmt.Sprintf("%s_%s", resModule, blockListCache)
	case 2:
		cacheName = fmt.Sprintf("%s_%s", resModule, allowListCache)
	default:
		cacheName = ""
	}
	return
}

// ClearCache 清除黑白名单缓存
func (s *sIpLimit) ClearCache(ctx context.Context, module ...string) {
	_, _, cacheName := s.getLimitParams(module...)
	if cacheName != "" {
		_, _ = gcache.Remove(ctx, cacheName)
	}
}

// IsLimited 检测访问是否受限
func (s *sIpLimit) IsLimited(ctx context.Context, module ...string) bool {
	limitType, resModule, cacheName := s.getLimitParams(module...)
	if limitType == 0 {
		return false
	}
	ip := g.RequestFromCtx(ctx).GetClientIp()
	inDatabase := false
	if len(ip) > 0 {
		ips := strings.Split(ip, ".")
		prefix := ips[0]
		//本地地址除外
		if prefix == "0" || prefix == "10" || prefix == "127" || prefix == "172" || prefix == "192" {
			return false
		}

		cacheLst, _ := gcache.Get(ctx, cacheName)
		lst := make([]string, 0)
		if cacheLst == nil {
			dbList, _ := g.Model("ip_limit").Ctx(ctx).
				Where("type", limitType).
				Where("module", resModule).
				Where("status", 1).
				All()
			for _, v := range dbList {
				lst = append(lst, v["ip"].String())
			}
			if len(lst) > 0 {
				_ = gcache.Set(ctx, cacheName, lst, 0)
			}
		} else {
			lst = gconv.Strings(cacheLst)
		}
		if len(lst) > 0 {
			for _, v := range lst {
				if s.equalIP(ip, v) {
					inDatabase = true
					break
				}
			}
			if limitType == 2 {
				//只允许白名单模式：不在数据库内，表示受限
				return !inDatabase
			}
			//不允许黑名单模式：在数据库内，表示受限
			return inDatabase
		}
	}
	return false
}
