package servicepkg

// Make sure that the template compiles during package initialization

var HandlerOutputTemplate = ParseTemplateOrPanic(`
package {{.TableName}}_handler

import (
	"net/http"

	"github.com/unhejing/go-gin-api/config"
	"github.com/unhejing/go-gin-api/dto/common/request"
	"github.com/unhejing/go-gin-api/dto/common/response"
	"github.com/unhejing/go-gin-api/dto/{{.TableName}}_dto"
	"github.com/unhejing/go-gin-api/model/{{.TableName}}_model"
	"github.com/unhejing/go-gin-api/service/{{.TableName}}_service"
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
	{{.HumpName}}Service {{.TableName}}_service.Service
}

func New(logger *zap.Logger, db mysql.Repo, cache redis.Repo) Handler {
	return &handler{
		logger:           logger,
		cache:            cache,
		hashids:          hash.New(config.Get().HashIds.Secret, config.Get().HashIds.Length),
		{{.HumpName}}Service: {{.TableName}}_service.New(db, cache),
	}
}

func (h *handler) i() {}

// PageList 分页查询
// @Summary 分页查询
// @Description 分页查询
// @Tags API.{{.TableName}}
// @Accept application/json
// @Produce json
// @Param data body {{.TableName}}_dto.PageReq true "请求信息"
// @Success 200 {object} response.Result
// @Failure 400 {object} code.Failure
// @Router /api/{{.TableName}}/pageList [post]
// @Security LoginToken
func (h *handler) PageList() core.HandlerFunc {
	return func(c core.Context) {
		req := new({{.TableName}}_dto.PageReq)
		if err := c.ShouldBindForm(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}
		res, err := h.{{.HumpName}}Service.PageList(c, req)
		if err != nil {
			c.Payload(response.OkWithData(new({{.TableName}}_dto.PageRes)))
			return
		}
		c.Payload(response.OkWithData(res))
	}
}

// Login 编辑
// @Summary 编辑
// @Description 编辑
// @Tags API.{{.TableName}}
// @Accept application/json
// @Produce json
// @Param data body {{.TableName}}_model.{{.StructName}} true "请求信息"
// @Success 200 {object} response.Result
// @Failure 400 {object} code.Failure
// @Router /api/{{.TableName}}/edit [post]
// @Security LoginToken
func (h *handler) Edit() core.HandlerFunc {
	return func(c core.Context) {
		req := new({{.TableName}}_model.{{.StructName}})
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
		err := h.{{.HumpName}}Service.Edit(c, req)
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
// @Tags API.{{.TableName}}
// @Accept  json
// @Produce json
// @Param data body request.IdRequest true "请求实体"
// @Success 200 {object} response.Result
// @Failure 400 {object} code.Failure
// @Router /api/{{.TableName}}/delete [post]
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
		err := h.{{.HumpName}}Service.Delete(c, id)
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
// @Tags API.{{.TableName}}
// @Accept  json
// @Produce json
// @Param data body {{.TableName}}_model.{{.StructName}} true "请求信息"
// @Success 200 {object} response.Result
// @Failure 400 {object} code.Failure
// @Router /api/{{.TableName}}/add [post]
// @Security LoginToken
func (h *handler) Create() core.HandlerFunc {
	return func(c core.Context) {
		req := new({{.TableName}}_model.{{.StructName}})
		res := new(response.IdResponse)
		if err := c.ShouldBindJSON(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				validation.Error(err)).WithError(err),
			)
			return
		}
		id, err := h.{{.HumpName}}Service.Create(c, req)
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
`)
