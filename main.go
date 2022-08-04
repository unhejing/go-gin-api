package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/unhejing/go-gin-api/config"
	"github.com/unhejing/go-gin-api/router"
	"github.com/unhejing/go-gin-api/utils/env"
	"github.com/unhejing/go-gin-api/utils/logger"
	"github.com/unhejing/go-gin-api/utils/shutdown"
	"github.com/unhejing/go-gin-api/utils/timeutil"

	"go.uber.org/zap"
	//"github.com/gin-gonic/gin"
)

// @title swagger 接口文档
// @version 2.0
// @description

// @securityDefinitions.apikey  LoginToken
// @in                          header
// @name                        token

func main() {
	// 初始化 access logger
	accessLogger, err := logger.NewJSONLogger(
		logger.WithDisableConsole(),
		logger.WithField("domain", fmt.Sprintf("%s[%s]", config.ProjectName, env.Active().Value())),
		logger.WithTimeLayout(timeutil.CSTLayout),
		logger.WithFileP(config.ProjectAccessLogFile),
	)
	if err != nil {
		panic(err)
	}

	// 初始化 cron logger
	cronLogger, err := logger.NewJSONLogger(
		logger.WithDisableConsole(),
		logger.WithField("domain", fmt.Sprintf("%s[%s]", config.ProjectName, env.Active().Value())),
		logger.WithTimeLayout(timeutil.CSTLayout),
		logger.WithFileP(config.ProjectCronLogFile),
	)

	if err != nil {
		panic(err)
	}

	defer func() {
		_ = accessLogger.Sync()
		_ = cronLogger.Sync()
	}()

	// 初始化 HTTP 服务
	s, err := router.NewHTTPServer(accessLogger, cronLogger)
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr:    config.ProjectPort,
		Handler: s.Mux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			accessLogger.Fatal("http server startup err", zap.Error(err))
		}
		accessLogger.Info("http server started: " + config.ProjectDomain + config.ProjectPort)
	}()

	// 优雅关闭
	shutdown.NewHook().Close(
		// 关闭 http server
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			if err := server.Shutdown(ctx); err != nil {
				accessLogger.Error("server shutdown err", zap.Error(err))
			}
		},

		// 关闭 db
		func() {
			if s.Db != nil {
				if err := s.Db.DbWClose(); err != nil {
					accessLogger.Error("dbw close err", zap.Error(err))
				}

				if err := s.Db.DbRClose(); err != nil {
					accessLogger.Error("dbr close err", zap.Error(err))
				}
			}
		},

		// 关闭 cache
		func() {
			if s.Cache != nil {
				if err := s.Cache.Close(); err != nil {
					accessLogger.Error("cache close err", zap.Error(err))
				}
			}
		},
	)
}
