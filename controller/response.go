package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type baseResponse struct {
	Code ResponseCodeType `json:"code"`
	Msg  interface{}      `json:"msg"`
}

type ResponseData struct {
	baseResponse
	Data interface{} `json:"data,omitempty"` //omitempty 该字段为空时忽略
}

func ResponseError(ctx *gin.Context, code ResponseCodeType) {
	rsp := &ResponseData{
		baseResponse: baseResponse{
			Code: code,
			Msg:  code.Msg(),
		},
		Data: nil,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func ResponseSuccess(ctx *gin.Context, data interface{}) {
	rsp := &ResponseData{
		baseResponse: baseResponse{
			Code: CodeSuccess,
			Msg:  CodeSuccess.Msg(),
		},

		Data: data,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func ResponseErrorWithMsg(ctx *gin.Context, code ResponseCodeType, msg interface{}) {
	rsp := &ResponseData{
		baseResponse: baseResponse{
			Code: code,
			Msg:  msg,
		},
		Data: nil,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func ResponseErrorWithDataMsg(ctx *gin.Context, code ResponseCodeType, msg interface{}, data interface{}) {
	rsp := &ResponseData{
		baseResponse: baseResponse{
			Code: code,
			Msg:  msg,
		},
		Data: data,
	}
	ctx.JSON(http.StatusOK, rsp)
}

func ResponseSuccessWithMsg(ctx *gin.Context, msg interface{}, data interface{}) {
	rsp := &ResponseData{
		baseResponse: baseResponse{
			Code: CodeSuccess,
			Msg:  msg,
		},
		Data: data,
	}
	ctx.JSON(http.StatusOK, rsp)
}
