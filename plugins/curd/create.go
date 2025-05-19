package curd

import (
	"github.com/gogf/gf/v2/database/gdb"
)

type CreateModel struct {
	gdbModel *gdb.Model
}

func NewCreate(gdbModel *gdb.Model) *CreateModel {
	return &CreateModel{
		gdbModel: gdbModel,
	}
}

// Submit 创建提交请求
func (s *CreateModel) Submit(form any) (newId int64, err error) {
	newId, err = s.gdbModel.InsertAndGetId(form)
	return
}
