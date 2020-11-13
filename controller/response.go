// Package controller provides ...
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
{
	"code": 1001, // 程序中的错误码
	"msg": xx, // 提示信息
	"data": {} // 数据
}
*/

type ResponseData struct {
	Code ResCode     `json:"code"`           // 业务状态码
	Msg  interface{} `json:"msg"`            // 提示信息
	Data interface{} `json:"data,omitempty"` // 数据
}

// ResponseError 错误时返回定义好的状态码和信息
func ResponseError(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  code.Msg(), // 通过状状态码获取对应信息
		Data: nil,
	})
}

// ResponseErrorWithMsg 返回自定义的错误状态码和信息
func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// ResponseSuccess 成功时返回定义好的状态码、信息
func ResponseSuccess(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}
