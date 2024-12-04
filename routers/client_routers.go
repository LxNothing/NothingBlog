package routers

import (
	"NothingBlog/controller"

	"github.com/gin-gonic/gin"
)

// ClientSetUp 注册客户端访问时的路由 - 不需要进行jwt认证
func ClientSetUp(basePath string, eng *gin.Engine) {
	// 文章路由
	router := eng.Group(basePath)
	{
		// index 主页面
		router.GET("", controller.GetIndexHandler)

		router.GET("/page/:page", controller.GetIndexHandler)
		// 根据文章ID获取文章
		router.GET("/article/:id", controller.GetArticleWithIdHandler)
		// 获取所有tag
		router.GET("/tags", controller.GetAllTagsHandler)
		// 根据tag id获取tag详细信息
		router.GET("/tag/:id", controller.GetArticleWithIdHandler)
		// 获取指定tag下的文章列表
		router.GET("/tag/:id/atc/:page", controller.GetArticleWithIdHandler)

		// 获取所有class
		router.GET("/classes", controller.GetAllArticleHandler)
		// 根据class id获取tag详细信息
		router.GET("/classes/:id", controller.GetArticleWithIdHandler)
		router.GET("/classes/:id/atc/:page", controller.GetArticleWithIdHandler)
	}
}
