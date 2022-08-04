package interceptor

import (
	"github.com/unhejing/go-gin-api/session"
	"github.com/unhejing/go-gin-api/utils/core"
	"github.com/unhejing/go-gin-api/utils/mysql"
	"github.com/unhejing/go-gin-api/utils/redis"

	"go.uber.org/zap"
)

var _ Interceptor = (*interceptor)(nil)

type Interceptor interface {

	// CheckSignature 验证签名是否合法，对用签名算法 pkg/signature
	CheckSignature() core.HandlerFunc

	// CheckLogin 检查用户是否登录
	CheckLogin(ctx core.Context) (sessionUserInfo session.SessionUserInfo, err core.BusinessError)

	// i 为了避免被其他包实现
	i()
}

type interceptor struct {
	logger *zap.Logger
	cache  redis.Repo
	db     mysql.Repo
}

func New(logger *zap.Logger, cache redis.Repo, db mysql.Repo) Interceptor {
	return &interceptor{
		logger: logger,
		cache:  cache,
		db:     db,
	}
}

func (i *interceptor) i() {}
