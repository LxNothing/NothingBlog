package controller

import (
	"NothingBlog/logic"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// FileUploadHandler 文件上传的处理函数
func FileUploadHandler(ctx *gin.Context) {
	fh, err := ctx.FormFile("file")
	if err != nil {
		zap.L().Debug("获取文件传输表单数据错误", zap.Error(err))
		ResponseError(ctx, CodeParameterInvalid)
		return
	}

	// 调用logic层的写文件的方法
	visitPath, err := logic.WriteFileToLocal(fh)
	if err != nil {
		zap.L().Debug("保存文件失败", zap.Error(err))
		ResponseError(ctx, CodeServerBusy)
		return
	}

	ResponseSuccessWithMsg(ctx, "保存文件成功", gin.H{
		"FilePath": visitPath,
	})
}
