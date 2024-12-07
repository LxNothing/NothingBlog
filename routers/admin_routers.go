package routers

import (
	"NothingBlog/controller"

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

	// 文章路由
	articleRouter := eng.Group(basePath)
	{
		//articleRouter.Use(middleware.JwtAuthorization())
		// 获取所有文章
		articleRouter.GET("/articles", controller.GetAllArticleHandler)
		// 根据文章ID获取文章
		articleRouter.GET("/article/:id", controller.GetArticleWithIdHandler)
		// 添加新的文章
		articleRouter.POST("/article", controller.CreateArticleHandler)

		// 删除指定ID的文章
		articleRouter.DELETE("/article/:id", controller.DeleteArticleHandler)
		// 删除多篇文章
		articleRouter.DELETE("/articles", controller.DeleteMultiArticleHandler)
		// 更新文章
		articleRouter.PUT("/article", controller.UpdateArticleHandler)
	}

	// Tag路由
	tagRouter := eng.Group(basePath)
	{
		//tagRouter.Use(middleware.JwtAuthorization())
		// 获取所有tag
		tagRouter.GET("/tags", controller.GetAllTagsHandler)
		// 根据tag id获取tag详细信息
		tagRouter.GET("/tag/:id", controller.GetTagByIdHandler)
		// 添加新的tag
		tagRouter.POST("/tag", controller.CreateTagHandler)
		// 删除单个tag
		tagRouter.DELETE("/tag/:id", controller.DeleteTagHandler)
		// 删除多个tag
		tagRouter.DELETE("/tags", controller.DeleteMultiTagHandler)
		// 更新tag
		tagRouter.PUT("/tag", controller.UpdateTagHandler)
	}

	// Class类别路由
	classRouter := eng.Group(basePath)
	{
		//classRouter.Use(middleware.JwtAuthorization())
		// 获取所有class
		classRouter.GET("/classes", controller.GetAllClassesHandler)
		// 根据class id获取tag详细信息
		classRouter.GET("/class/:id", controller.GetClassByIdHandler)
		// 添加新的class
		classRouter.POST("/class", controller.CreateClassHandler)

		// 删除单个类别
		classRouter.DELETE("/class/:id", controller.DeleteClassHandler)
		// 删除多个类别
		classRouter.DELETE("/classes", controller.DeleteMultiClassHandler)
		// 更新class
		tagRouter.PUT("/class", controller.UpdateClassHandler)
	}
}
