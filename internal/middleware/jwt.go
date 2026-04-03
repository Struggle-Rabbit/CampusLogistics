package middleware

import (
	"strings"

	"github.com/Struggle-Rabbit/CampusLogistics/pkg/constant"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 登录校验
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 header 里的 Authorization
		auth := c.GetHeader("Authorization")
		if auth == "" {
			utils.Unauth(c, constant.MsgPleaseLoginFirst)
			c.Abort()
			return
		}

		// 格式：Bearer token
		parts := strings.SplitN(auth, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.Unauth(c, constant.MsgTokenFormatError)
			c.Abort()
			return
		}

		// 2. 解析 token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			utils.Unauth(c, constant.MsgLoginExpiredOrInvalid)
			c.Abort()
			return
		}

		// 3. 把用户信息存到 gin context
		c.Set("userID", claims.UserID)
		// c.Set("username", claims.Username)

		// 继续执行后续接口
		c.Next()
	}
}
