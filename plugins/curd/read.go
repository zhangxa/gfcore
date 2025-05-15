package curd

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gpage"
	"github.com/zhangxa/gfcore/utils"
)

// ReqPageData 获取分页列表
func ReqPageData[T any](ctx context.Context, req *PageDataInput) (data *PageData[T], err error) {
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
	if req.Model == nil {
		req.Model = g.DB().Model(new(T)).Ctx(ctx)
	}
	db := req.Model
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
	data.Pager.CurrentPage = page.CurrentPage
	data.Pager.TotalSize = page.TotalSize
	err = db.Page(req.CurPage, req.Limit).Scan(&data.List)
	return
}

// ReqJoinPageData 获取ajax列表分页数据
func ReqJoinPageData[T any](ctx context.Context, req *JoinPageDataInput) (data *PageData[T], err error) {
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
	if req.Model == nil {
		req.Model = g.DB().Model(new(T)).Ctx(ctx)
	}
	db := req.Model.As(req.TableAlia)
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
	data.Pager.CurrentPage = page.CurrentPage
	data.Pager.TotalSize = page.TotalSize
	err = db.Page(req.CurPage, req.Limit).Scan(&data.List)
	return
}

// ReqListData 获取列表-不分页
func ReqListData[T any](ctx context.Context, req *ListDataInput) (data []T, err error) {
	data = make([]T, 0)
	if req.Model == nil {
		req.Model = g.DB().Model(new(T)).Ctx(ctx)
	}
	if req.Fields != "" {
		req.Model = req.Model.Fields(req.Fields)
	}
	if req.Where != nil {
		req.Model = req.Model.Where(req.Where)
	}
	if req.Order != "" {
		req.Model = req.Model.Order(req.Order)
	}
	err = req.Model.Scan(&data)
	return
}

// ReqJoinListData 获取联合表数据列表-不分页
func ReqJoinListData[T any](ctx context.Context, req *JoinPageDataInput) (data []T, err error) {
	data = make([]T, 0)
	if req.Model == nil {
		req.Model = g.DB().Model(new(T)).Ctx(ctx)
	}
	db := req.Model.As(req.TableAlia)
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
func ReqSingleData[T any](ctx context.Context, req *SingleDataInput) (data *T, err error) {
	if req.Model == nil {
		req.Model = g.DB().Model(new(T)).Ctx(ctx)
	}
	if req.Where != nil {
		req.Model = req.Model.Where(req.Where)
	}
	if req.Fields != "" {
		req.Model = req.Model.Fields(req.Fields)
	}
	if req.Order != "" {
		req.Model = req.Model.Order(req.Order)
	}
	err = req.Model.Scan(&data)
	return
}

// ReqJoinSingleData 获取联合表单条数据
func ReqJoinSingleData[T any](ctx context.Context, req *JoinSingleDataInput) (data *T, err error) {
	if req.Model == nil {
		req.Model = g.DB().Model(new(T)).Ctx(ctx)
	}
	db := req.Model.As(req.TableAlia)
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
