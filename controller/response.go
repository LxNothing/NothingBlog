package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Code ResponseCodeType `json:"code"`
	Msg  interface{}      `json:"msg"`
	Data interface{}      `json:"data,omitempty"` //omitempty 该字段为空时忽略
}

func ResponseError(ctx *gin.Context, code ResponseCodeType) {
	rsp := &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func ResponseSuccess(ctx *gin.Context, data interface{}) {
	rsp := &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func ResponseErrorWithMsg(ctx *gin.Context, code ResponseCodeType, msg interface{}) {
	rsp := &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func ResponseSuccessWithMsg(ctx *gin.Context, msg interface{}, data interface{}) {
	rsp := &ResponseData{
		Code: CodeSuccess,
		Msg:  msg,
		Data: data,
	}
	ctx.JSON(http.StatusOK, rsp)
}
