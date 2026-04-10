package router

import (
	"fmt"

	"github.com/Struggle-Rabbit/CampusLogistics/internal/app"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/config"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/controller/common"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/middleware"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// InitRouter 初始化总路由
func initRouter(app *app.App) *gin.Engine {
	// 初始化 gin
	r := gin.Default()

	// 全局中间件
	r.Use(
		middleware.Recovery(),  // 崩溃恢复
		middleware.RequestID(), // 请求ID
		middleware.Logger(),    // 日志
		middleware.CORS(),      // 跨域
	)

	srv := service.NewServiceProvider(app)

	commonCtl := common.NewCommonController(srv)

	// ===================== 路由分组 =====================
	api := r.Group("/api/v1")
	{
		api.POST("/register", commonCtl.Register)
		api.POST("/login", commonCtl.Login)
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

func Run(app *app.App) error {
	// 注册路由
	fmt.Println("注册路由....")
	r := initRouter(app)

	logger.Info("服务启动",
		zap.String("env", config.GlobalConfig.App.Env),
		zap.Int("port", config.GlobalConfig.App.Port),
	)

	fmt.Println("服务启动中....")
	if err := r.Run(fmt.Sprintf(":%d", config.GlobalConfig.App.Port)); err != nil {
		return err
	}
	return nil
}
