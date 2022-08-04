package sys_config_handler

import (
	"net/http"

	"github.com/unhejing/go-gin-api/config"
	"github.com/unhejing/go-gin-api/dto/common/request"
	"github.com/unhejing/go-gin-api/dto/common/response"
	"github.com/unhejing/go-gin-api/dto/sys_config_dto"
	"github.com/unhejing/go-gin-api/model/sys_config_model"
	"github.com/unhejing/go-gin-api/service/sys_config_service"
	"github.com/unhejing/go-gin-api/utils/code"
	"github.com/unhejing/go-gin-api/utils/core"
	"github.com/unhejing/go-gin-api/utils/hash"
	"github.com/unhejing/go-gin-api/utils/mysql"
	"github.com/unhejing/go-gin-api/utils/redis"
	"github.com/unhejing/go-gin-api/utils/validation"

	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	PageList() core.HandlerFunc

	Edit() core.HandlerFunc

	Create() core.HandlerFunc

	Delete() core.HandlerFunc
}

type handler struct {
	logger           *zap.Logger
	cache            redis.Repo
	hashids          hash.Hash
	sysConfigService sys_config_service.Service
}

func New(logger *zap.Logger, db mysql.Repo, cache redis.Repo) Handler {
	return &handler{
		logger:           logger,
		cache:            cache,
		hashids:          hash.New(config.Get().HashIds.Secret, config.Get().HashIds.Length),
		sysConfigService: sys_config_service.New(db, cache),
	}
}

func (h *handler) i() {}

// PageList 分页查询
// @Summary 分页查询
// @Description 分页查询
// @Tags API.sys_config
// @Accept application/json
// @Produce json
// @Param data body sys_config_dto.PageReq true "请求信息"
// @Success 200 {object} response.Result
// @Failure 400 {object} code.Failure
// @Router /api/sys_config/pageList [post]
// @Security LoginToken
func (h *handler) PageList() core.HandlerFunc {
	return func(c core.Context) {
		req := new(sys_config_dto.PageReq)
		if err := c.ShouldBindForm(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}
		res, err := h.sysConfigService.PageList(c, req)
		if err != nil {
			c.Payload(response.OkWithData(new(sys_config_dto.PageRes)))
			return
		}
		c.Payload(response.OkWithData(res))
	}
}

// Login 编辑
// @Summary 编辑
// @Description 编辑
// @Tags API.sys_config
// @Accept application/json
// @Produce json
// @Param data body sys_config_model.SysConfig true "请求信息"
// @Success 200 {object} response.Result
// @Failure 400 {object} code.Failure
// @Router /api/sys_config/edit [post]
// @Security LoginToken
func (h *handler) Edit() core.HandlerFunc {
	return func(c core.Context) {
		req := new(sys_config_model.SysConfig)
		res := new(response.IdResponse)
		if err := c.ShouldBindJSON(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}
		id := req.Id
		err := h.sysConfigService.Edit(c, req)
		if err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.MenuDeleteError,
				code.Text(code.MenuDeleteError)).WithError(err),
			)
			return
		}

		res.Id = id
		c.Payload(res)
	}
}

// Create 删除
// @Summary 删除
// @Description 删除
// @Tags API.sys_config
// @Accept  json
// @Produce json
// @Param data body request.IdRequest true "请求实体"
// @Success 200 {object} response.Result
// @Failure 400 {object} code.Failure
// @Router /api/sys_config/delete [post]
// @Security LoginToken
func (h *handler) Delete() core.HandlerFunc {
	return func(c core.Context) {
		req := new(request.IdRequest)
		res := new(response.IdResponse)
		if err := c.ShouldBindURI(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}
		id := req.Id
		err := h.sysConfigService.Delete(c, id)
		if err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.MenuDeleteError,
				code.Text(code.MenuDeleteError)).WithError(err),
			)
			return
		}
		res.Id = id
		c.Payload(res)
	}
}

// Create 新增
// @Summary 新增
// @Description 新增
// @Tags API.sys_config
// @Accept  json
// @Produce json
// @Param data body sys_config_model.SysConfig true "请求信息"
// @Success 200 {object} response.Result
// @Failure 400 {object} code.Failure
// @Router /api/sys_config/add [post]
// @Security LoginToken
func (h *handler) Create() core.HandlerFunc {
	return func(c core.Context) {
		req := new(sys_config_model.SysConfig)
		res := new(response.IdResponse)
		if err := c.ShouldBindJSON(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				validation.Error(err)).WithError(err),
			)
			return
		}
		id, err := h.sysConfigService.Create(c, req)
		if err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminCreateError,
				code.Text(code.AdminCreateError)).WithError(err),
			)
			return
		}

		res.Id = id
		c.Payload(response.OkWithData(res))
	}
}
