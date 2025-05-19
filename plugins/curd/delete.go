package curd

import (
	"errors"
	"github.com/gogf/gf/v2/database/gdb"
)

type DeleteModel struct {
	gdbModel *gdb.Model
}

func NewDelete(gdbModel *gdb.Model) *DeleteModel {
	return &DeleteModel{
		gdbModel: gdbModel,
	}
}

// Submit 删除提交请求
func (s *DeleteModel) Submit(where any) (err error) {
	_, err = s.gdbModel.Where(where).Delete()
	return
}

// SubmitStrict 删除提交请求
func (s *DeleteModel) SubmitStrict(where any) error {
	count, err := s.gdbModel.Where(where).Count()
	if err != nil {
		return err
	}
	if count <= 0 {
		return errors.New("data not exists")
	}
	_, err = s.gdbModel.Where(where).Delete()
	return err
}
