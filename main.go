package main

import (
	"NothingBlog/dao/mysql"
	"NothingBlog/dao/redis"
	"NothingBlog/dao/table"
	"NothingBlog/logger"
	"NothingBlog/package/snowflake"
	"NothingBlog/package/verifycode"
	"NothingBlog/routers"
	"NothingBlog/settings"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	docs "NothingBlog/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// 参考博客 ： https://cangmang.xyz

// @title NgBlog
// @version 1.0
// @description NgBlog Go博客项目 API 接口文档
// @host localhost:8080
func main() {
	// 1.配置初始化
	if err := settings.Init(); err != nil {
		fmt.Println("Config init failed...")
		return
	}

	// 2.日志初始化
	if err := logger.Init(settings.Confg.LogConfig, settings.Confg.Mode); err != nil {
		fmt.Println("Logging init failed...")
		return
	}
	zap.L().Debug("Init success")
	defer zap.L().Sync()

	// 3.初始化Mysql
	if err := mysql.Init(settings.Confg.MysqlConfig, settings.Confg.Mode); err != nil {
		fmt.Println("Mysql init failed...")
		return
	}
	defer mysql.Close()
	// 4. 初始化项目所用到的mysql 表结构
	if err := table.DbTableInit(); err != nil {
		zap.L().Error("初始数据库表失败", zap.Error(err))
		return
	}
	// 4.初始化redis
	if err := redis.Init(settings.Confg.RedisConfig); err != nil {
		fmt.Println("Redis init failed...")
		fmt.Println(err)
		return
	}
	defer redis.Close()

	// 初始化雪花算法-用于生成用户ID
	if err := snowflake.Init(settings.Confg.StartTime, settings.Confg.MachineId); err != nil {
		fmt.Println("snowflake init failed...")
		fmt.Println(err)
		return
	}

	// 初始化验证码生成器
	verifycode.Init(settings.Confg.AuthConfig)

	// 初始化文章数据 - 主要是从mysql中读出文章的点赞数，以及点赞的人（对于博客的点赞应该不需要）

	// 5.注册路由
	r := gin.New()
	if settings.Confg.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) // 将gin设置为发布模式
	}

	// 设置日志记录中间件
	r.Use(logger.GinZapLogger(), logger.GinZapRecovery(true))
	// 这个路由用于配置接口文档的访问路由 - 仅在开发阶段使用
	docs.SwaggerInfo.BasePath = settings.Confg.AdminBasePath
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// 配置客户端访问和admin访问的路由
	routers.ClientSetUp(settings.Confg.ClientBasePath, r)
	routers.AdminSetUp(settings.Confg.AdminBasePath, r)

	// 6.启动服务：优雅关机与重启
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", settings.Confg.AppConfig.Port),
		Handler: r,
	}

	// 开启协程来处理请求
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Error("Server listen error...", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	// 捕获信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Server shuttdown...")

	// 创建一个5秒延时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Error("Shutdown error...", zap.Error(err))
	}
	zap.L().Info("Server has closed...")
}
