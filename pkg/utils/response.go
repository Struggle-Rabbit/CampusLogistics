package utils

import (
	"net/http"

	"github.com/Struggle-Rabbit/CampusLogistics/pkg/constant"
	"github.com/gin-gonic/gin"
)

// Response 统一响应结构体
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
	Err  interface{} `json:"err,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: constant.CodeSuccess,
		Msg:  constant.MsgSuccess,
		Data: data,
	})
}

// Fail 业务失败响应
func Fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: constant.CodeFail,
		Msg:  msg,
	})
}

// Error 系统错误响应
func Error(c *gin.Context, msg string, err error) {
	c.JSON(http.StatusInternalServerError, Response{
		Code: constant.CodeError,
		Msg:  msg,
		Err:  err,
	})
}

// Unauth 未登录响应
func Unauth(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code: constant.CodeUnauth,
		Msg:  msg,
	})
}
