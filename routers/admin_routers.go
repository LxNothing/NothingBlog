package routers

import (
	"NothingBlog/controller"
	"NothingBlog/middleware"

	"github.com/gin-gonic/gin"
)

func AdminSetUp(basePath string, eng *gin.Engine) {
	// 创建与认证相关的路由组
	authRouter := eng.Group(basePath)
	{
		// 注册请求验证码的路由
		authRouter.GET("/auth/verifycode", controller.VerifyCodeHandler)
		// 注册 用户注册路由
		authRouter.POST("/auth/signup", controller.SignUpHandler)
		// 注册 用户登录路由
		authRouter.POST("/auth/login", controller.LoginHandler)
		// 修改密码 - 输入用户名，输入旧的密码，输入新的密码，输入验证码
		authRouter.POST("/auth/password/reset", controller.PwdResetHandler)
		// 重置密码
		authRouter.POST("/auth/password/modify", controller.PwdModifyHandler)
	}

	// 以下的内容都需要进行jwt认证才可以访问
	//v1.Use(middleware.JwtAuthorization())

	// 文章路由
	articleRouter := eng.Group(basePath)
	{
		articleRouter.Use(middleware.JwtAuthorization())
		// 获取所有文章
		articleRouter.GET("/article/all", controller.GetAllArticleHandler)
		// 根据文章ID获取文章
		articleRouter.GET("/article/:id", controller.GetArticleWithIdHandler)
		// 添加新的文章
		articleRouter.POST("/article", controller.CreateArticleHandler)

		// 删除指定ID的文章
		articleRouter.DELETE("/article/:id")
		// 删除多篇文章
	}

	// Tag路由
	tagRouter := eng.Group(basePath)
	{
		tagRouter.Use(middleware.JwtAuthorization())
		// 获取所有tag
		tagRouter.GET("/tag-all", controller.GetAllArticleHandler)
		// 根据tag id获取tag详细信息
		tagRouter.GET("/tag/:id", controller.GetArticleWithIdHandler)
		// 添加新的tag
		tagRouter.POST("/tag", controller.CreateTagHandler)
	}

	// Class类别路由
	classRouter := eng.Group(basePath)
	{
		classRouter.Use(middleware.JwtAuthorization())
		// 获取所有tag
		classRouter.GET("/class-all", controller.GetAllArticleHandler)
		// 根据tag id获取tag详细信息
		classRouter.GET("/class/:id", controller.GetArticleWithIdHandler)
		// 添加新的tag
		classRouter.POST("/class", controller.CreateClassHandler)
	}
}
