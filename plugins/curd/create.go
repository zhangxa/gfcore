package curd

import (
	"context"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// Create 创建提交请求
func Create[T any](ctx context.Context, form any, gdbModel ...*gdb.Model) (newId int64, err error) {
	var db *gdb.Model
	if len(gdbModel) > 0 {
		db = gdbModel[0]
	} else {
		db = g.DB().Model(new(T)).Ctx(ctx)
	}
	newId, err = db.InsertAndGetId(form)
	return
}
