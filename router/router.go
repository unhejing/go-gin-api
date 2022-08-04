package router

import (
	"github.com/unhejing/go-gin-api/config"
	"github.com/unhejing/go-gin-api/interceptor"
	"github.com/unhejing/go-gin-api/utils/core"
	"github.com/unhejing/go-gin-api/utils/errors"
	"github.com/unhejing/go-gin-api/utils/file"
	"github.com/unhejing/go-gin-api/utils/mysql"
	"github.com/unhejing/go-gin-api/utils/redis"

	"go.uber.org/zap"
)

type resource struct {
	mux          core.Mux
	logger       *zap.Logger
	db           mysql.Repo
	cache        redis.Repo
	interceptors interceptor.Interceptor
}

type Server struct {
	Mux   core.Mux
	Db    mysql.Repo
	Cache redis.Repo
}

func NewHTTPServer(logger *zap.Logger, cronLogger *zap.Logger) (*Server, error) {
	if logger == nil {
		return nil, errors.New("logger required")
	}

	r := new(resource)
	r.logger = logger

	openBrowserUri := config.ProjectDomain + config.ProjectPort

	_, ok := file.IsExists(config.ProjectInstallMark)
	if !ok { // 未安装
		openBrowserUri += "/install"
	} else { // 已安装
		openBrowserUri += "/swagger/index.html"
		// 初始化 DB
		dbRepo, err := mysql.New()
		if err != nil {
			logger.Fatal("new db err", zap.Error(err))
		}
		r.db = dbRepo

		// 初始化 Cache
		cacheRepo, err := redis.New()
		if err != nil {
			logger.Fatal("new cache err", zap.Error(err))
		}
		r.cache = cacheRepo
	}

	mux, err := core.New(logger,
		core.WithEnableOpenBrowser(openBrowserUri),
		core.WithEnableCors(),
		core.WithEnableRate(),
	)

	if err != nil {
		panic(err)
	}

	r.mux = mux
	r.interceptors = interceptor.New(logger, r.cache, r.db)

	// 设置 API 路由
	setApiRouter(r)

	s := new(Server)
	s.Mux = mux
	s.Db = r.db
	s.Cache = r.cache

	return s, nil
}
