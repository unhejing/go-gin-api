package sys_config_service

import (
	"github.com/unhejing/go-gin-api/dto/sys_config_dto"
	"github.com/unhejing/go-gin-api/model/sys_config_model"
	"github.com/unhejing/go-gin-api/utils/core"
	"github.com/unhejing/go-gin-api/utils/mysql"
	"github.com/unhejing/go-gin-api/utils/redis"

	"github.com/fatih/structs"
	"github.com/spf13/cast"
)

var _ Service = (*service)(nil)

type Service interface {
	i()

	PageList(ctx core.Context, searchData *sys_config_dto.PageReq) (res sys_config_dto.PageRes, err error)

	Create(ctx core.Context, config *sys_config_model.SysConfig) (id int32, err error)

	Delete(c core.Context, id int32) (err error)

	Edit(ctx core.Context, user *sys_config_model.SysConfig) (err error)
}

type service struct {
	db    mysql.Repo
	cache redis.Repo
}

func New(db mysql.Repo, cache redis.Repo) Service {
	return &service{
		db:    db,
		cache: cache,
	}
}

func (s *service) i() {}

func (s *service) PageList(ctx core.Context, searchData *sys_config_dto.PageReq) (res sys_config_dto.PageRes, err error) {
	page := searchData.Page
	if page == 0 {
		page = 1
	}

	pageSize := searchData.Size
	if pageSize == 0 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	qb := sys_config_model.NewQueryBuilder()

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
	res.List = make([]sys_config_model.SysConfig, 0, len(listData))
	for _, v := range listData {
		res.List = append(res.List, *v)
	}
	return
}

func (s *service) Create(ctx core.Context, model *sys_config_model.SysConfig) (id int32, err error) {
	id, err = model.Create(s.db.GetDbW().WithContext(ctx.RequestContext()))
	if err != nil {
		return 0, err
	}
	return
}

func (s *service) Delete(ctx core.Context, id int32) (err error) {
	qb := sys_config_model.NewQueryBuilder()
	qb.WhereId(mysql.EqualPredicate, id)
	if err = qb.Delete(s.db.GetDbW().WithContext(ctx.RequestContext())); err != nil {
		return err
	}
	return
}

func (s *service) Edit(ctx core.Context, model *sys_config_model.SysConfig) (err error) {
	id := model.Id
	qb := sys_config_model.NewQueryBuilder()
	qb.WhereId(mysql.EqualPredicate, id)
	err = qb.Updates(s.db.GetDbW().WithContext(ctx.RequestContext()), structs.Map(&model))
	if err != nil {
		return err
	}
	return
}
