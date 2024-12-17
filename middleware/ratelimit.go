package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
	"go.uber.org/zap"
)

// 限流中间件
func RateLimit(interval time.Duration, cap int64) func(ctx *gin.Context) {
	bucket := ratelimit.NewBucket(interval, cap)
	return func(ctx *gin.Context) { // 闭包
		if bucket.TakeAvailable(1) <= 0 {
			zap.L().Debug("无可用令牌，流量过大")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
