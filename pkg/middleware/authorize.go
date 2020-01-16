package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/haleyrom/wallet/core"
	"github.com/haleyrom/wallet/internal/controllers/base"
	"github.com/haleyrom/wallet/internal/params"
	"github.com/haleyrom/wallet/internal/resp"
	"github.com/haleyrom/wallet/pkg/jwt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	// HttpHeadToken http请求包头的token数据
	HttpHeadToken string = "Authorization"

	// emptyData 空数据
	emptyData = struct{}{}

	// NotHandlerToken 不拦截token
	NotHandlerToken = map[string]struct{}{
		"/api/v1/wallet/check":                   emptyData,
		"/api/v1/wallet/deposit/top_up":          emptyData,
		"/api/v1/wallet/test_data/account/write": emptyData,
		"/api/v1/wallet/withdrawal/callback":     emptyData,
	}

	// Jwt
	Jwt *jwt.JWT

	// AdminJwt
	AdminJwt *jwt.JWT
)

// HttpInterceptor 拦截器
func HttpInterceptor(j *jwt.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err error
		)

		token := c.Request.Header.Get(HttpHeadToken)
		logrus.Infof("request url: %s, ContentType: %s, token: %s, body: %v", c.Request.URL, c.Request.Method, token, c.Request.Body)
		// NotHandlerToken 不拦截token
		if _, ok := NotHandlerToken[c.Request.RequestURI]; ok == true {
			c.Next()
			return
		}

		if len(token) > core.DefaultNilNum {
			claims := &jwt.CustomClaims{}
			// parseToken 解析token包含的信息
			if claims, err = j.ParseToken(token); err == nil {
				fmt.Println(claims, "=======================")
				info := core.UserInfoPool.Get().(*params.BaseParam)
				info.Claims = *claims
				// 判断用户id是否未空
				if len(info.Claims.UserID) == core.DefaultNilNum {
					err = errors.Errorf("%d", resp.CodeIllegalToken)
				} else if err = base.CreateUser(c, info); err == nil {
					core.UserInfoPool.Put(info)
				}
			}
		} else {
			err = errors.Errorf("%d", resp.CodeNoToken)
		}

		switch err {
		case nil:
			c.Next()
		case jwt.TokenExpired:
			fallthrough
		case jwt.TokenNotValidYet:
			fallthrough
		case jwt.TokenMalformed:
			fallthrough
		case jwt.TokenInvalid:
			err = errors.Errorf("%d", resp.CodeIllegalToken)
			fallthrough
		default:
			core.GResp.Failure(c, err)
			c.Abort()
			return
		}
	}
}

// HttpBindGResp HttpBindGResp
func HttpBindGResp() gin.HandlerFunc {
	return func(c *gin.Context) {
		core.GResp = &resp.Resp{}
	}
}

// HttpCors 跨域名
func HttpCors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}
