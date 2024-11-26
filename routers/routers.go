package routers

import (
	"NothingBlog/controller"
	"NothingBlog/logger"
	"NothingBlog/middleware"
	"NothingBlog/settings"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetUp(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // 将gin设置为发布模式
	}

	eng := gin.New()

	eng.Use(logger.GinZapLogger(), logger.GinZapRecovery(true))

	// 注册路由 - 测试路由
	eng.GET("/test", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, fmt.Sprintf("this is webapp, version is:%s", settings.Confg.Version))
	})

	// 创建一个路由组
	v1 := eng.Group("/api/v1")
	// 注册请求验证码的路由
	v1.GET("/verifycode", controller.VerifyCodeHandler)
	// 注册 用户注册路由
	v1.POST("/signup", controller.SignUpHandler)

	// 注册 用户登录路由
	v1.POST("/login", controller.LoginHandler)

	/*v1.Use(middleware.JwtAuthorization())
	{
		// 获取用户组
		v1.GET("/community", controller.CommunityListHandler)
		// 路径参数
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		// 创建博客文章的post请求
		v1.POST("/newblog", controller.BlogCreaterHandler)

		// 查询具体的博文详情
		v1.GET("/getblog/:id", controller.BlogDetailHandler)

		// 查询博文列表
		// 完整路由：http://localhost:8080/api/v1/getblog/?size=1&page=1
		v1.GET("/getblog", controller.BlogListHandler)
		// 查询博文列表 - 并且支持指定排序方式 - 创建时间 或者 点赞数量
		// 完整路由：http://localhost:8080/api/v1/getblog/?size=1&page=1&order=time / score
		v1.GET("/getblog2", controller.BlogOrderListHandler)

		// 文章点赞或者投票
		v1.POST("/vote", controller.VotedHandler)
	}*/

	// 注册 认证测试路由
	eng.GET("/authtest", middleware.JwtAuthorization(), func(ctx *gin.Context) {
		uid, _ := ctx.Get("UserId")
		ctx.JSON(http.StatusOK, uid)
	})

	return eng
}
