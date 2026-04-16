package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/logger"
	"github.com/gin-gonic/gin"
)

// OperationLogMiddleware 操作日志中间件
func OperationLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 先获取请求信息（在 c.Next() 之前）
		startTime := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// 2. 读取请求参数（GET/POST）
		params := getRequestParams(c)

		// 3. 执行后续的接口逻辑
		c.Next()

		// 4. 接口执行完后，获取响应信息
		statusCode := c.Writer.Status()
		operationAt := startTime

		// 5. 从 Context 中获取当前登录用户信息（你登录时存进去的）
		userID, _ := c.Get("user_id")
		userName, _ := c.Get("user_name")

		// 6. 组装日志对象
		log := model.SysOperationLog{
			UserID:      userID.(string),
			UserName:    userName.(string),
			Method:      method,
			Path:        path,
			Params:      params,
			StatusCode:  statusCode,
			IP:          ip,
			UserAgent:   userAgent,
			OperationAt: operationAt,
		}

		// 7. 异步写入数据库（关键！不阻塞接口响应）
		go func() {
			if err := dao.DB.Create(&log).Error; err != nil {
				// 记录日志失败，只打印不影响接口
				logger.Error("操作日志写入失败: " + err.Error())
			}
		}()
	}
}

// getRequestParams 获取请求参数
func getRequestParams(c *gin.Context) string {
	// GET 请求：从 URL Query 获取
	if c.Request.Method == "GET" {
		return c.Request.URL.RawQuery
	}

	// POST/PUT 请求：从 Body 获取（注意：Body 读了要放回去）
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}
	// 把 Body 写回去，否则后续接口拿不到参数
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return string(bodyBytes)
}
