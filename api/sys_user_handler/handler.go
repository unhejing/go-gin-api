package sys_user_handler

import (
	"net/http"

	"github.com/unhejing/go-gin-api/service/sys_user_service"

	"github.com/unhejing/go-gin-api/config"
	"github.com/unhejing/go-gin-api/dto/common/request"
	"github.com/unhejing/go-gin-api/dto/common/response"
	"github.com/unhejing/go-gin-api/dto/sys_user_dto"
	"github.com/unhejing/go-gin-api/model/sys_user_model"
	"github.com/unhejing/go-gin-api/session"
	"github.com/unhejing/go-gin-api/utils/code"
	"github.com/unhejing/go-gin-api/utils/core"
	"github.com/unhejing/go-gin-api/utils/errors"
	"github.com/unhejing/go-gin-api/utils/hash"
	"github.com/unhejing/go-gin-api/utils/mysql"
	"github.com/unhejing/go-gin-api/utils/password"
	"github.com/unhejing/go-gin-api/utils/redis"
	"github.com/unhejing/go-gin-api/utils/validation"

	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	PageList() core.HandlerFunc

	Login() core.HandlerFunc

	Create() core.HandlerFunc

	Delete() core.HandlerFunc
}

type handler struct {
	logger      *zap.Logger
	cache       redis.Repo
	hashids     hash.Hash
	userService sys_user_service.Service
}

func New(logger *zap.Logger, db mysql.Repo, cache redis.Repo) Handler {
	return &handler{
		logger:      logger,
		cache:       cache,
		hashids:     hash.New(config.Get().HashIds.Secret, config.Get().HashIds.Length),
		userService: sys_user_service.New(db, cache),
	}
}

func (h *handler) i() {}

// PageList 分页查询
// @Summary 分页查询
// @Description 分页查询
// @Tags API.sys_user
// @Accept application/json
// @Produce json
// @Param data body sys_user_dto.PageReq true "请求信息"
// @Success 200 {object} response.Result
// @Failure 400 {object} code.Failure
// @Router /api/sys_user/pageList [post]
// @Security LoginToken
func (h *handler) PageList() core.HandlerFunc {
	return func(c core.Context) {
		req := new(sys_user_dto.PageReq)
		if err := c.ShouldBindForm(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}
		res, err := h.userService.PageList(c, req)
		if err != nil {
			c.Payload(response.OkWithData(new(sys_user_dto.PageRes)))
			return
		}
		c.Payload(response.OkWithData(res))
	}
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录
// @Tags API.sys_user
// @Accept application/json
// @Produce json
// @Param data body sys_user_dto.LoginReq true "登录信息"
// @Success 200 {object} response.Result
// @Failure 400 {object} code.Failure
// @Router /api/sys_user/login [post]
// @Security LoginToken
func (h *handler) Login() core.HandlerFunc {
	return func(c core.Context) {
		req := new(sys_user_dto.LoginReq)
		res := new(sys_user_dto.LoginRes)
		if err := c.ShouldBindJSON(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				code.Text(code.ParamBindError)).WithError(err),
			)
			return
		}

		searchOneData := new(sys_user_dto.SearchOneData)
		searchOneData.Username = req.Username
		searchOneData.Password = password.GeneratePassword(req.Password)
		searchOneData.IsUsed = 1

		info, err := h.userService.Detail(c, searchOneData)
		if err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminLoginError,
				code.Text(code.AdminLoginError)).WithError(err),
			)
			return
		}

		if info == nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminLoginError,
				code.Text(code.AdminLoginError)).WithError(errors.New("未查询出符合条件的用户")),
			)
			return
		}

		token := password.GenerateLoginToken(info.Id)

		// 用户信息
		sessionUserInfo := &session.SessionUserInfo{
			UserID:   info.Id,
			UserName: info.Username,
		}

		// 将用户信息记录到 Redis 中
		err = h.cache.Set(config.RedisKeyPrefixLoginUser+token, string(sessionUserInfo.Marshal()), config.LoginSessionTTL, redis.WithTrace(c.Trace()))
		if err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AdminLoginError,
				code.Text(code.AdminLoginError)).WithError(err),
			)
			return
		}

		res.Token = token
		c.Payload(res)
	}
}

// Create 删除
// @Summary 删除
// @Description 删除
// @Tags API.sys_user
// @Accept  json
// @Produce json
// @Param data body request.IdRequest true "请求实体"
// @Success 200 {object} response.Result
// @Failure 400 {object} code.Failure
// @Router /api/sys_user/delete [post]
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
		err := h.userService.Delete(c, id)
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

// Create 新增用户
// @Summary 新增用户
// @Description 新增用户
// @Tags API.sys_user
// @Accept  json
// @Produce json
// @Param data body sys_user_dto.RegisterReq true "用户信息"
// @Success 200 {object} response.Result
// @Failure 400 {object} code.Failure
// @Router /api/sys_user/register [post]
// @Security LoginToken
func (h *handler) Create() core.HandlerFunc {
	return func(c core.Context) {
		req := new(sys_user_dto.RegisterReq)
		res := new(response.IdResponse)
		if err := c.ShouldBindJSON(req); err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.ParamBindError,
				validation.Error(err)).WithError(err),
			)
			return
		}

		createData := new(sys_user_model.SysUser)
		createData.Nickname = req.Nickname
		createData.Username = req.Username
		createData.Mobile = req.Mobile
		createData.Password = password.GeneratePassword(req.Password)

		id, err := h.userService.Create(c, createData)
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
