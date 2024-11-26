package middleware

import (
	"NothingBlog/controller"
	"NothingBlog/dao/redis"
	"NothingBlog/package/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

// 请求头中的字段名称
const authKey = "Authorization"

// 认证中间件 - 从请求头中获取token
func JwtAuthorization() func(ctx *gin.Context) {
	// 代码实现的是：Authorization: Bearer xxx.xxx.xxx这种格式的token
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get(authKey)
		if authHeader == "" {
			controller.ResponseError(ctx, controller.CodeTokenEmpty)
			ctx.Abort()
			return
		}

		ps := strings.SplitN(authHeader, " ", 2)
		if len(ps) != 2 || ps[0] != "Bearer" {
			controller.ResponseError(ctx, controller.CodeTokenInvaild)
			ctx.Abort()
			return
		}

		sc, err := jwt.ParseJwtToken(ps[1])
		if err != nil {
			controller.ResponseError(ctx, controller.CodeTokenInvaild)
			ctx.Abort()
			return
		}

		// 获取redis中的token和认证信息中的token进行对比，实现单用户单点登录
		tk, err := redis.QueryTokenByUserId(sc.UserId)
		if err != nil || tk != ps[1] {
			controller.ResponseError(ctx, controller.CodeNeedReLogin)
			ctx.Abort()
			return
		}

		// 将用户id存到上下文中
		ctx.Set(controller.ContextUserIdKey, sc.UserId)
		ctx.Next()
	}
}
