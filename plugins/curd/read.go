package curd

import (
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/util/gpage"
	"github.com/zhangxa/gfcore/utils"
)

type ReadModel[T any] struct {
	gdbModel *gdb.Model
}

func NewRead[T any](gdbModel *gdb.Model) *ReadModel[T] {
	return &ReadModel[T]{
		gdbModel: gdbModel,
	}
}

// ReqPageData 获取分页列表
func (s *ReadModel[T]) ReqPageData(req *PageDataInput) (data *PageData[T], err error) {
	if req.Limit <= 0 {
		req.Limit = utils.Config.GetDefaultPageSize(20)
	}
	if req.CurPage == 0 {
		req.CurPage = 1
	}
	data = &PageData[T]{
		List: make([]T, 0),
		Pager: &Pager{
			TotalPage:   0,
			CurrentPage: req.CurPage,
			TotalSize:   0,
			PageSize:    req.Limit,
		},
	}
	db := s.gdbModel.Clone()
	if req.Where != nil {
		db = db.Where(req.Where)
	}
	var total int
	total, err = db.Count()
	if err != nil {
		return
	}
	if total == 0 {
		return
	}
	if req.Fields != "" {
		db = db.Fields(req.Fields)
	}
	if req.Order != "" {
		db = db.Order(req.Order)
	}
	page := gpage.New(total, req.Limit, req.CurPage, "")
	data.Pager.TotalPage = page.TotalPage
	data.Pager.TotalSize = page.TotalSize
	err = db.Page(req.CurPage, req.Limit).Scan(&data.List)
	return
}

// ReqJoinPageData 获取ajax列表分页数据
func (s *ReadModel[T]) ReqJoinPageData(req *JoinPageDataInput) (data *PageData[T], err error) {
	if req.Limit <= 0 {
		req.Limit = utils.Config.GetDefaultPageSize(20)
	}
	if req.CurPage == 0 {
		req.CurPage = 1
	}
	data = &PageData[T]{
		List: make([]T, 0),
		Pager: &Pager{
			TotalPage:   0,
			CurrentPage: req.CurPage,
			TotalSize:   0,
			PageSize:    req.Limit,
		},
	}
	db := s.gdbModel.Clone().As(req.TableAlia)
	switch req.JoinType {
	case JoinLeft:
		db = db.LeftJoin(req.JoinTable, req.JoinTableAlia, req.JoinOn)
	case JoinRight:
		db = db.RightJoin(req.JoinTable, req.JoinTableAlia, req.JoinOn)
	default:
		db = db.InnerJoin(req.JoinTable, req.JoinTableAlia, req.JoinOn)
	}
	if req.Where != nil {
		db = db.Where(req.Where)
	}
	var total int
	total, err = db.Count()
	if err != nil {
		return
	}
	if total == 0 {
		return
	}
	if req.Fields != "" {
		db = db.Fields(req.Fields)
	}
	if req.Order != "" {
		db = db.Order(req.Order)
	}
	page := gpage.New(total, req.Limit, req.CurPage, "")
	data.Pager.TotalPage = page.TotalPage
	data.Pager.TotalSize = page.TotalSize
	err = db.Page(req.CurPage, req.Limit).Scan(&data.List)
	return
}

// ReqListData 获取列表-不分页
func (s *ReadModel[T]) ReqListData(req *ListDataInput) (data []T, err error) {
	data = make([]T, 0)
	db := s.gdbModel.Clone()
	if req.Fields != "" {
		db = db.Fields(req.Fields)
	}
	if req.Where != nil {
		db = db.Where(req.Where)
	}
	if req.Order != "" {
		db = db.Order(req.Order)
	}
	err = db.Scan(&data)
	return
}

// ReqJoinListData 获取联合表数据列表-不分页
func (s *ReadModel[T]) ReqJoinListData(req *JoinPageDataInput) (data []T, err error) {
	data = make([]T, 0)
	db := s.gdbModel.Clone().As(req.TableAlia)
	switch req.JoinType {
	case JoinLeft:
		db = db.LeftJoin(req.JoinTable, req.JoinTableAlia, req.JoinOn)
	case JoinRight:
		db = db.RightJoin(req.JoinTable, req.JoinTableAlia, req.JoinOn)
	default:
		db = db.InnerJoin(req.JoinTable, req.JoinTableAlia, req.JoinOn)
	}
	if req.Where != nil {
		db = db.Where(req.Where)
	}
	if req.Fields != "" {
		db = db.Where(req.Fields)
	}
	if req.Order != "" {
		db = db.Order(req.Order)
	}
	err = db.Scan(&data)
	return
}

// ReqSingleData 获取单条数据
func (s *ReadModel[T]) ReqSingleData(req *SingleDataInput) (data *T, err error) {
	db := s.gdbModel.Clone()
	if req.Where != nil {
		db = db.Where(req.Where)
	}
	if req.Fields != "" {
		db = db.Fields(req.Fields)
	}
	if req.Order != "" {
		db = db.Order(req.Order)
	}
	err = db.Scan(&data)
	return
}

// ReqJoinSingleData 获取联合表单条数据
func (s *ReadModel[T]) ReqJoinSingleData(req *JoinSingleDataInput) (data *T, err error) {
	db := s.gdbModel.Clone().As(req.TableAlia)
	switch req.JoinType {
	case JoinLeft:
		db = db.LeftJoin(req.JoinTable, req.JoinTableAlia, req.JoinOn)
	case JoinRight:
		db = db.RightJoin(req.JoinTable, req.JoinTableAlia, req.JoinOn)
	default:
		db = db.InnerJoin(req.JoinTable, req.JoinTableAlia, req.JoinOn)
	}
	if req.Where != nil {
		db = db.Where(req.Where)
	}
	if req.Fields != "" {
		db = db.Where(req.Fields)
	}
	if req.Order != "" {
		db = db.Order(req.Order)
	}
	err = db.Scan(&data)
	return
}
