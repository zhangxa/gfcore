package server

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gvalid"
	"github.com/zhangxa/gfcore/core"
	"strings"
	"time"
)

type sValidate struct {
}

func init() {
	core.RegisterValidate(NewValidate())
}

// NewValidate 验证规则自定义服务
func NewValidate() core.IValidate {
	return &sValidate{}
}

func (s *sValidate) RuleUnique(ctx context.Context, in gvalid.RuleFuncInput) error {
	arr := strings.Split(in.Rule, ":")
	var array = strings.Split(arr[1], ",")
	data := gconv.Map(in.Data)
	// g.Dump(array, in)
	// SELECT COUNT(*) FROM `user` WHERE `id` != xxx AND `name` = 'xxx'
	group := ""
	table := array[0]
	tblArr := strings.Split(array[0], ".")
	if len(tblArr) > 1 {
		group = tblArr[0]
		table = tblArr[1]
	}
	fmt.Println("group:", group, "table:", table)
	m := g.DB(group).Model(table).Safe().Ctx(ctx)
	pk := "id"
	al := len(array)
	if al > 1 {
		pk = array[1]
		if v, ok := data[pk]; ok {
			m = m.WhereNot(pk, v)
		}

		if al > 2 {
			for i := 2; i < al; i++ {
				f := array[i]
				if v, ok := data[f]; ok {
					m = m.Where(f, v)
				}
			}
		}
	} else {
		return gerror.New("unique校验配置有误")
	}

	count, err := m.Cache(gdb.CacheOption{
		Duration: 5 * time.Minute,
		Name:     "",
		Force:    false,
	}).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		if in.Message != "" {
			msg := gstr.ReplaceByMap(in.Message, map[string]string{
				"{value}": in.Value.String(), // Current validating value.
			})
			msg, _ = gregex.ReplaceString(`\s{2,}`, ` `, msg)
			return gerror.New(msg)
		}
		return gerror.Newf(`%s is not unique`, in.Value.String())
	}
	return nil
}
