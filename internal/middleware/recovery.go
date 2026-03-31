package middleware

import (
	"fmt"

	"github.com/Struggle-Rabbit/CampusLogistics/pkg/constant"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/logger"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recovery
// -----------------------------------
// 自定义 Recovery（崩溃不宕机）
// -----------------------------------
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if val := recover(); val != nil {
				var err error
				switch v := val.(type) {
				case error:
					err = v
				default:
					err = fmt.Errorf("%v", v)
				}
				logger.Error(
					"系统错误",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
				)
				utils.Error(c, constant.MsgInternalServerError, err)

				c.Abort()
			}
		}()
		c.Next()
	}
}
