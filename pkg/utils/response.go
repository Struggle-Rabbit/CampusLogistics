package utils

import (
	"net/http"

	"github.com/Struggle-Rabbit/CampusLogistics/pkg/constant"
	"github.com/gin-gonic/gin"
)

// SuccessResponse 成功响应结构体
type SuccessResponse struct {
	Code int         `json:"code" example:"200"`
	Msg  string      `json:"msg" example:"操作成功"`
	Data interface{} `json:"data,omitempty" swaggertype:"object"`
}

type ErrResponse struct {
	Code int         `json:"code" example:"400"`
	Msg  string      `json:"msg" example:"失败"`
	Err  interface{} `json:"err,omitempty" swaggertype:"object"`
}

// Success 成功响应
func Success(c *gin.Context, args ...interface{}) {

	msg := constant.MsgSuccess
	var data interface{} = nil

	switch len(args) {
	case 1:
		data = args[0]
	case 2:
		if m, ok := args[1].(string); ok {
			msg = m
		}
		data = args[0]
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Code: constant.CodeSuccess,
		Msg:  msg,
		Data: data,
	})
}

// Fail 业务失败响应
func Fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, ErrResponse{
		Code: constant.CodeFail,
		Msg:  msg,
	})
}

// NoPermission 无接口权限
func NoPermission(c *gin.Context) {
	c.JSON(http.StatusOK, ErrResponse{
		Code: constant.CodeForbid,
		Msg:  "无权限访问！",
	})
}

// Error 系统错误响应
func Error(c *gin.Context, msg string, err error) {
	c.JSON(http.StatusInternalServerError, ErrResponse{
		Code: constant.CodeError,
		Msg:  msg,
		Err:  err,
	})
}

// Unauth 未登录响应
func Unauth(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, ErrResponse{
		Code: constant.CodeUnauth,
		Msg:  msg,
	})
}
