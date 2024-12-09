package controller

import (
	"NothingBlog/models"
)

// 本文件定义了所有handler的返回数据的格式，仅用于描述swagger文档

// 文章（article）相关

// 类别（class）相关
type _ResponseAllClassesList struct {
	Code ResponseCodeType             `json:"code"`
	Msg  interface{}                  `json:"msg"`
	Data []*models.ResponseClassBrief //`json:"data,omitempty"` //omitempty 该字段为空时忽略
}

type _ResponseClassDetailList struct {
	Code ResponseCodeType            `json:"code"`
	Msg  interface{}                 `json:"msg"`
	Data *models.ResponseClassDetail //`json:"data,omitempty"` //omitempty 该字段为空时忽略
}

type _ResponseCreateClass struct {
	Code ResponseCodeType           `json:"code"`
	Msg  interface{}                `json:"msg"`
	Data *models.ResponseClassBrief //`json:"data,omitempty"` //omitempty 该字段为空时忽略
}

// 返回的数据不包含data域
type _ResponseDeleteClass struct {
	Code ResponseCodeType `json:"code"`
	Msg  interface{}      `json:"msg"`
	Data []int64          `json:"data"`
}

type _ResponseNoDataArea struct {
	Code ResponseCodeType `json:"code"`
	Msg  interface{}      `json:"msg"`
}

// 标签（tag）相关
type _ResponseAllTagList struct {
	Code ResponseCodeType           `json:"code"`
	Msg  interface{}                `json:"msg"`
	Data []*models.ResponseTagBrief //`json:"data,omitempty"` //omitempty 该字段为空时忽略
}

type _ResponseTagDetailList struct {
	Code ResponseCodeType          `json:"code"`
	Msg  interface{}               `json:"msg"`
	Data *models.ResponseTagDetail //`json:"data,omitempty"` //omitempty 该字段为空时忽略
}

type _ResponseCreateTag struct {
	Code ResponseCodeType         `json:"code"`
	Msg  interface{}              `json:"msg"`
	Data *models.ResponseTagBrief //`json:"data,omitempty"` //omitempty 该字段为空时忽略
}
