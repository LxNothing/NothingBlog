package logger

import (
	"NothingBlog/settings"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger *zap.Logger

func Init(logcfg *settings.LogConfig, mode string) (err error) {
	// 使用 zap.New() 方法来自定义构造日志记录器
	writer := getLoggerWriter(
		logcfg.FileName,
		logcfg.MaxSize,
		logcfg.MaxBackup,
		logcfg.MaxAge)
	encoder := getLoggerEncoder()

	var level = new(zapcore.Level)
	if err = level.UnmarshalText([]byte(logcfg.Level)); err != nil {
		return
	}

	var core zapcore.Core
	// 开发者模式 - 日志向控制台和文件输出
	if mode == "dev" {
		consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
		core = zapcore.NewTee( // 将日志向两个方向写入
			zapcore.NewCore(encoder, writer, level),
			// zap.DebugLevel 设置日志等级为debug
			zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		// 非开发模式，只向文件输出日志
		core = zapcore.NewCore(encoder, writer, level)
	}

	// zap.AddCaller() - 添加调用者
	logger = zap.New(core, zap.AddCaller())

	// 替换全局的变量，就可以使用zap.L().Debug("Init success")来进行日志记录
	zap.ReplaceGlobals(logger)
	// 还可以将不同日志级别输出到不同的文件：参考 https://www.liwenzhou.com/posts/Go/zap/
	return
}

// 获取编码器 - 即以哪种格式写入日志，比如json格式、控制台打印格式
func getLoggerEncoder() zapcore.Encoder {
	// 这里使用json格式
	//return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	// 这里采用控制台格式
	// zap.NewDevelopmentEncoderConfig() 里面返回一个结构体，控制时间格式等
	return zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
}

// 指定将日志写到哪里去 - 这里写到固定的文件
func getLoggerWriter(filename string, max_size int, max_backups int, max_age int) zapcore.WriteSyncer {
	//f, _ := os.Create("/Users/mantall/Temp/GoLog/zaplogtest.log")
	// 利用io.MultiWriter支持文件和终端两个输出目标
	//ws := io.MultiWriter(f, os.Stdout)

	// 使用该库是为了实现日志切割
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    max_size,
		MaxBackups: max_backups,
		MaxAge:     max_age,
		Compress:   false,
	}

	return zapcore.AddSync(lumberJackLogger)
}

func GinZapLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery

		ctx.Next()

		cost := time.Since(start)

		zap.L().Info(path,
			zap.Int("status", ctx.Writer.Status()),
			zap.String("method", ctx.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", ctx.ClientIP()),
			zap.String("user-agent", ctx.Request.UserAgent()),
			zap.String("error", ctx.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("duration", cost))
	}
}

// GinRecovery recover掉项目可能出现的panic
func GinZapRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					// If the connection is dead, we can't write a status to it.
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
