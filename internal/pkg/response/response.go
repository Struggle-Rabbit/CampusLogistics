package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构体
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// 状态码定义
const (
	CodeSuccess = 200
	CodeFail    = 400
	CodeUnauth  = 401
	CodeForbid  = 403
	CodeError   = 500
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  "操作成功",
		Data: data,
	})
}

// Fail 业务失败响应
func Fail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: CodeFail,
		Msg:  msg,
	})
}

// Error 系统错误响应
func Error(c *gin.Context, msg string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code: CodeError,
		Msg:  msg,
	})
}

// Unauth 未登录响应
func Unauth(c *gin.Context, msg string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code: CodeUnauth,
		Msg:  msg,
	})
}
