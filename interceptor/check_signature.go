package interceptor

import (
	"net/http"
	"strings"

	"github.com/unhejing/go-gin-api/config"
	"github.com/unhejing/go-gin-api/utils/code"
	"github.com/unhejing/go-gin-api/utils/core"
	"github.com/unhejing/go-gin-api/utils/env"
	"github.com/unhejing/go-gin-api/utils/errors"
	"github.com/unhejing/go-gin-api/utils/signature"
)

var whiteListPath = map[string]bool{
	"/login/web": true,
}

func (i *interceptor) CheckSignature() core.HandlerFunc {
	return func(c core.Context) {
		if !env.Active().IsPro() {
			return
		}

		// 签名信息
		authorization := c.GetHeader(config.HeaderSignToken)
		if authorization == "" {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AuthorizationError,
				code.Text(code.AuthorizationError)).WithError(errors.New("Header 中缺少 Authorization 参数")),
			)
			return
		}

		// 时间信息
		date := c.GetHeader(config.HeaderSignTokenDate)
		if date == "" {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AuthorizationError,
				code.Text(code.AuthorizationError)).WithError(errors.New("Header 中缺少 Date 参数")),
			)
			return
		}

		// 通过签名信息获取 key
		authorizationSplit := strings.Split(authorization, " ")
		if len(authorizationSplit) < 2 {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AuthorizationError,
				code.Text(code.AuthorizationError)).WithError(errors.New("Header 中 Authorization 格式错误")),
			)
			return
		}

		key := authorizationSplit[0]

		// 查询数据库，权限判断
		Secret := ""

		if !whiteListPath[c.Path()] {
			// 验证 c.Method() + c.Path() 是否授权

		}

		ok, err := signature.New(key, Secret, config.HeaderSignTokenTimeout).Verify(authorization, date, c.Path(), c.Method(), c.RequestInputParams())
		if err != nil {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AuthorizationError,
				code.Text(code.AuthorizationError)).WithError(err),
			)
			return
		}

		if !ok {
			c.AbortWithError(core.Error(
				http.StatusBadRequest,
				code.AuthorizationError,
				code.Text(code.AuthorizationError)).WithError(errors.New("Header 中 Authorization 信息错误")),
			)
			return
		}
	}
}
