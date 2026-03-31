package middleware

import (
	"bytes"
	"time"

	"github.com/Struggle-Rabbit/CampusLogistics/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
// 全局日志中间件（记录请求、耗时、状态码）
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
		logger.Info(path,
			zap.String("method", method),
			zap.String("ip", clientIP),
			zap.Any("request_id", uuid),
			zap.Int("status", status),
			zap.Duration("cost", cost),
		)
	}
}

// RequestID
// -----------------------------------
// 请求ID 中间件（追踪链路）
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
