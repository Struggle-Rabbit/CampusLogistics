package middleware

import (
	"strings"

	"github.com/Struggle-Rabbit/CampusLogistics/internal/pkg/jwt"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/pkg/response"
	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 登录校验
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 header 里的 Authorization
		auth := c.GetHeader("Authorization")
		if auth == "" {
			response.Unauth(c, "请先登录")
			c.Abort()
			return
		}

		// 格式：Bearer token
		parts := strings.SplitN(auth, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauth(c, "token 格式错误")
			c.Abort()
			return
		}

		// 2. 解析 token
		claims, err := jwt.ParseToken(parts[1])
		if err != nil {
			response.Unauth(c, "登录已过期或无效")
			c.Abort()
			return
		}

		// 3. 把用户信息存到 gin context
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		// 继续执行后续接口
		c.Next()
	}
}
