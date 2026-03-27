package router

import (
	"github.com/Struggle-Rabbit/CampusLogistics/internal/middleware"
	"github.com/gin-gonic/gin"
)

// InitRouter 初始化总路由
func InitRouter() *gin.Engine {
	// 初始化 gin
	r := gin.Default()

	// 全局中间件
	r.Use(
		middleware.CORS(),      // 跨域
		middleware.Logger(),    // 日志
		middleware.RequestID(), // 请求ID
		middleware.Recovery(),  // 崩溃恢复
	)

	// ===================== 路由分组 =====================
	api := r.Group("/api")
	{
		// 需要Token的
		api.Use(middleware.JWTAuth())
		{
			//LoadDormRouter(api)   // 宿舍模块
			//LoadRepairRouter(api) // 报修模块
		}
		// 加载各个模块路由
		//LoadUserRouter(api)   // 用户模块

	}

	return r
}
