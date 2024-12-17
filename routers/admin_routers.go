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
		// 根据特定的参数分页获取文章 - 仅包含文章的大致信息，不包含文本的具体内容
		// /api/v1/articles?p=xx&s=xx&c=xx&t=xx&n=xx&st=xx&pr=&
		// p-page, s-size, c-class, t-tag, n-name(文章名称,模糊搜索), st-status, pr-Privilege
		articleRouter.GET("/articles", controller.GetAllArticleHandler)

		// 根据文章ID获取文章
		articleRouter.GET("/article/:id", controller.GetArticleWithIdHandler)
		// 添加新的文章
		articleRouter.POST("/article", controller.CreateArticleHandler)
		// 硬删除指定ID的文章
		articleRouter.DELETE("/article/:id", controller.DeleteArticleHandler)
		// 硬删除多篇文章
		articleRouter.DELETE("/articles", controller.DeleteMultiArticleHandler)
		// 软删除文章
		articleRouter.PUT("/soft-article", controller.UpdateArticleStatusHandler)
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

	// 文件相关路由
	fileRouter := eng.Group(basePath)
	{
		//classRouter.Use(middleware.JwtAuthorization())
		fileRouter.POST("/upload", controller.FileUploadHandler)
	}

	// 评论相关路由
	commentRouter := eng.Group(basePath)
	{
		commentRouter.GET("/comments", controller.GetCommentWithPageHandler) // 按页获取评论
		commentRouter.POST("/comment", controller.CreateCommentHandler)      // 创建新的评论
		commentRouter.PUT("/comment/vote")                                   // 评论投票 - 即支持 或者 不支持
		//以下的接口需要进行JWT认证
		//commentRouter.Use(middleware.JwtAuthorization())
		commentRouter.PUT("/comment/state", controller.UpdateCommentStateHandler) // 评论审核
		commentRouter.DELETE("/comment/:id", controller.DeleteCommentById)        // 根据评论ID删除评论
		commentRouter.DELETE("/comments", controller.DeleteCommentsByIds)         // 根据评论ID删除评论
	}
}
