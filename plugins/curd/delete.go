package curd

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// Delete 删除提交请求
func Delete[T any](ctx context.Context, where any, gdbModel ...*gdb.Model) (err error) {
	var db *gdb.Model
	if len(gdbModel) > 0 {
		db = gdbModel[0]
	} else {
		db = g.DB().Model(new(T)).Ctx(ctx)
	}
	_, err = db.Where(where).Delete()
	return
}
