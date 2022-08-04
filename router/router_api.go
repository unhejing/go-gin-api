package router

import (
	"github.com/unhejing/go-gin-api/api/sys_config_handler"
	"github.com/unhejing/go-gin-api/api/sys_user_handler"
	"github.com/unhejing/go-gin-api/utils/core"
)

func setApiRouter(r *resource) {

	// 无需签名
	noLogin := r.mux.Group("/api")
	{
		// sys_user
		userHandler := sys_user_handler.New(r.logger, r.db, r.cache)
		noLogin.POST("/sys_user/register", userHandler.Create())   // 添加用户
		noLogin.POST("/sys_user/login", userHandler.Login())       // 用户登录
		noLogin.POST("/sys_user/pageList", userHandler.PageList()) // 分页请求信息

		// sys_config
		sysConfigHandler := sys_config_handler.New(r.logger, r.db, r.cache)
		noLogin.POST("/sys_config/pageList", sysConfigHandler.PageList())
		noLogin.POST("/sys_config/add", sysConfigHandler.Create())
		noLogin.POST("/sys_config/edit", sysConfigHandler.Edit())
		noLogin.POST("/sys_config/delete", sysConfigHandler.Delete())

	}

	// 需要签名验证、登录验证
	api := r.mux.Group("/api", core.WrapAuthHandler(r.interceptors.CheckLogin))
	{
		// sys_user_dto
		userHandler := sys_user_handler.New(r.logger, r.db, r.cache)
		api.POST("/sys_user_dto/info", userHandler.Create()) // 添加用户
	}
}
