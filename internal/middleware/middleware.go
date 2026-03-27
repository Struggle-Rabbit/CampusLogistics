package middleware

import (
	"bytes"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/pkg/logger" // 你的zap日志包
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// 自定义 ResponseWriter 用于捕获状态码
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger
// -----------------------------------
// 1. 全局日志中间件（记录请求、耗时、状态码）
// -----------------------------------
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 读取请求信息
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()
		uuid, _ := c.Get("RequestID")

		// 替换 writer 捕获返回内容
		w := &responseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = w

		// 执行后续逻辑
		c.Next()

		// 请求结束后记录
		cost := time.Since(start)
		status := c.Writer.Status()

		// 用 zap 记录结构化日志
		logger.Log.Info(path,
			zap.String("method", method),
			zap.String("ip", clientIP),
			zap.Any("request_id", uuid),
			zap.Int("status", status),
			zap.Duration("cost", cost),
		)
	}
}

// CORS
// -----------------------------------
// 2. 跨域中间件（CORS）
// -----------------------------------
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "*")
		c.Header("Access-Control-Expose-Headers", "Content-Length,Access-Control-Allow-Origin")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行OPTIONS
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestID
// -----------------------------------
// 3. 请求ID 中间件（追踪链路）
// -----------------------------------
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = time.Now().Format("20060102150405") + "-" + c.ClientIP()
		}
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// Recovery
// -----------------------------------
// 4. 自定义 Recovery（崩溃不宕机）
// -----------------------------------
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Log.Error("系统崩溃",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
				)
				c.JSON(500, gin.H{
					"code": 500,
					"msg":  "服务器内部错误",
					"err":  err,
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}
