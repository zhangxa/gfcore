package curd

import (
	"github.com/gogf/gf/v2/database/gdb"
)

type UpdateModel struct {
	gdbModel *gdb.Model
}

func NewUpdate(gdbModel *gdb.Model) *UpdateModel {
	return &UpdateModel{
		gdbModel: gdbModel,
	}
}

func (m *UpdateModel) Submit(where any, form any) (err error) {
	_, err = m.gdbModel.Where(where).Update(form)
	return
}
