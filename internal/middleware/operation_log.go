package middleware

import (
	"bytes"
	"encoding/json"
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
		method := c.Request.Method
		path := c.Request.URL.Path
		ip := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// 2. 读取请求参数并脱敏
		params := getRequestParams(c)
		params = filterSensitiveData(params)

		// 3. 执行后续的接口逻辑
		c.Next()

		// 4. 接口执行完后，获取响应信息
		statusCode := c.Writer.Status()
		operationAt := time.Now() // 使用请求结束时间

		// 5. 从 Context 中获取当前登录用户信息
		userID, _ := c.Get("userID")
		userName, _ := c.Get("userName")

		// 如果没有登录，则记录为系统或游客
		uidStr := ""
		unameStr := "Guest"
		if userID != nil {
			uidStr = userID.(string)
		}
		if userName != nil {
			unameStr = userName.(string)
		}

		// 6. 组装日志对象
		log := model.SysOperationLog{
			UserID:      uidStr,
			UserName:    unameStr,
			Method:      method,
			Path:        path,
			Params:      params,
			StatusCode:  statusCode,
			IP:          ip,
			UserAgent:   userAgent,
			OperationAt: operationAt,
		}

		// 7. 异步写入数据库
		go func() {
			if err := dao.DB.Create(&log).Error; err != nil {
				logger.Error("操作日志写入失败: " + err.Error())
			}
		}()
	}
}

// filterSensitiveData 过滤敏感数据
func filterSensitiveData(params string) string {
	if params == "" {
		return ""
	}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(params), &data); err != nil {
		return params // 如果不是 JSON，直接返回
	}

	sensitiveFields := []string{"password", "old_password", "new_password", "token", "refresh_token"}
	for _, field := range sensitiveFields {
		if _, ok := data[field]; ok {
			data[field] = "******"
		}
	}

	filteredBytes, _ := json.Marshal(data)
	return string(filteredBytes)
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
