package interceptor

import (
	"encoding/json"
	"net/http"

	"github.com/unhejing/go-gin-api/config"
	"github.com/unhejing/go-gin-api/session"
	"github.com/unhejing/go-gin-api/utils/code"
	"github.com/unhejing/go-gin-api/utils/core"
	"github.com/unhejing/go-gin-api/utils/errors"
	"github.com/unhejing/go-gin-api/utils/redis"
)

func (i *interceptor) CheckLogin(ctx core.Context) (sessionUserInfo session.SessionUserInfo, err core.BusinessError) {
	token := ctx.GetHeader(config.HeaderLoginToken)
	if token == "" {
		err = core.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			code.Text(code.AuthorizationError)).WithError(errors.New("Header 中缺少 Token 参数"))

		return
	}

	if !i.cache.Exists(config.RedisKeyPrefixLoginUser + token) {
		err = core.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			code.Text(code.AuthorizationError)).WithError(errors.New("请先登录"))

		return
	}

	cacheData, cacheErr := i.cache.Get(config.RedisKeyPrefixLoginUser+token, redis.WithTrace(ctx.Trace()))
	if cacheErr != nil {
		err = core.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			code.Text(code.AuthorizationError)).WithError(cacheErr)

		return
	}

	jsonErr := json.Unmarshal([]byte(cacheData), &sessionUserInfo)
	if jsonErr != nil {
		core.Error(
			http.StatusUnauthorized,
			code.AuthorizationError,
			code.Text(code.AuthorizationError)).WithError(jsonErr)

		return
	}

	return
}
