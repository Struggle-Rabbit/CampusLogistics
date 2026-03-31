package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORS
// -----------------------------------
// 跨域中间件（CORS）
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
