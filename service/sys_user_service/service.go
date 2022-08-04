package sys_user_service

import (
	"github.com/unhejing/go-gin-api/dto/sys_user_dto"
	"github.com/unhejing/go-gin-api/model/sys_user_model"
	"github.com/unhejing/go-gin-api/utils/core"
	"github.com/unhejing/go-gin-api/utils/mysql"
	"github.com/unhejing/go-gin-api/utils/redis"

	"github.com/fatih/structs"
	"github.com/spf13/cast"
)

var _ Service = (*service)(nil)

type Service interface {
	i()

	PageList(ctx core.Context, searchData *sys_user_dto.PageReq) (res sys_user_dto.PageRes, err error)

	Create(ctx core.Context, user *sys_user_model.SysUser) (id int32, err error)

	Delete(c core.Context, id int32) (err error)

	Edit(ctx core.Context, user *sys_user_model.SysUser) (err error)

	Detail(c core.Context, data *sys_user_dto.SearchOneData) (info *sys_user_model.SysUser, err error)
}

type service struct {
	db    mysql.Repo
	cache redis.Repo
}

func (s *service) PageList(ctx core.Context, searchData *sys_user_dto.PageReq) (res sys_user_dto.PageRes, err error) {
	page := searchData.Page
	if page == 0 {
		page = 1
	}

	pageSize := searchData.Size
	if pageSize == 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	qb := sys_user_model.NewQueryBuilder()

	// 查询总数
	total, err := qb.Count(s.db.GetDbR().WithContext(ctx.RequestContext()))
	if err != nil {
		return res, err
	}
	res.Pagination.Total = cast.ToInt(total)
	res.Pagination.PerPageCount = pageSize
	res.Pagination.CurrentPage = page

	// 查询分页数据
	listData, err := qb.
		Limit(pageSize).
		Offset(offset).
		OrderById(false).
		QueryAll(s.db.GetDbR().WithContext(ctx.RequestContext()))
	if err != nil {
		return res, err
	}
	res.List = make([]sys_user_model.SysUser, 0, len(listData))
	for _, v := range listData {
		res.List = append(res.List, *v)
	}
	return
}

func (s *service) Create(ctx core.Context, model *sys_user_model.SysUser) (id int32, err error) {
	id, err = model.Create(s.db.GetDbW().WithContext(ctx.RequestContext()))
	if err != nil {
		return 0, err
	}
	return
}

func (s *service) Delete(ctx core.Context, id int32) (err error) {
	qb := sys_user_model.NewQueryBuilder()
	qb.WhereId(mysql.EqualPredicate, id)
	if err = qb.Delete(s.db.GetDbW().WithContext(ctx.RequestContext())); err != nil {
		return err
	}
	return
}

func (s *service) Edit(ctx core.Context, model *sys_user_model.SysUser) (err error) {
	id := model.Id
	qb := sys_user_model.NewQueryBuilder()
	qb.WhereId(mysql.EqualPredicate, id)
	err = qb.Updates(s.db.GetDbW().WithContext(ctx.RequestContext()), structs.Map(&model))
	if err != nil {
		return err
	}
	return
}

func (s *service) Detail(ctx core.Context, searchOneData *sys_user_dto.SearchOneData) (info *sys_user_model.SysUser, err error) {
	qb := sys_user_model.NewQueryBuilder()
	qb.WhereIsDeleted(mysql.EqualPredicate, -1)

	if searchOneData.Id != 0 {
		qb.WhereId(mysql.EqualPredicate, searchOneData.Id)
	}

	if searchOneData.Username != "" {
		qb.WhereUsername(mysql.EqualPredicate, searchOneData.Username)
	}

	if searchOneData.Nickname != "" {
		qb.WhereNickname(mysql.EqualPredicate, searchOneData.Nickname)
	}

	if searchOneData.Mobile != "" {
		qb.WhereMobile(mysql.EqualPredicate, searchOneData.Mobile)
	}

	if searchOneData.Password != "" {
		qb.WherePassword(mysql.EqualPredicate, searchOneData.Password)
	}

	if searchOneData.IsUsed != 0 {
		qb.WhereIsUsed(mysql.EqualPredicate, searchOneData.IsUsed)
	}

	info, err = qb.QueryOne(s.db.GetDbR().WithContext(ctx.RequestContext()))
	if err != nil {
		return nil, err
	}

	return
}

func New(db mysql.Repo, cache redis.Repo) Service {
	return &service{
		db:    db,
		cache: cache,
	}
}

func (s *service) i() {}
