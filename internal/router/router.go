package router

import (
	"fmt"

	"github.com/Struggle-Rabbit/CampusLogistics/internal/app"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/config"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/middleware"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	_ "github.com/Struggle-Rabbit/CampusLogistics/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter 初始化总路由
func initRouter(app *app.App) *gin.Engine {
	// 初始化 gin
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 全局中间件
	r.Use(
		middleware.Recovery(),  // 崩溃恢复
		middleware.RequestID(), // 请求ID
		middleware.Logger(),    // 日志
		middleware.CORS(),      // 跨域
	)

	srv := service.NewServiceProvider(app)

	// ===================== 路由分组 =====================
	api := r.Group("/api/v1")
	api.Use(middleware.OperationLogMiddleware())
	{
		LoadCommonRouter(api, srv) // 公共模块
		// 需要Token的
		api.Use(middleware.JWTAuth())
		{
			LoadUserRouter(api, srv)    // 用户模块
			LoadSysteamRouter(api, srv) // 系统模块
			LoadRoleRouter(api, srv)    // 角色
			LoadMenuRouter(api, srv)    // 菜单权限
			//LoadDormRouter(api)   // 宿舍模块
			//LoadRepairRouter(api) // 报修模块
		}

	}

	return r
}

func Run(app *app.App) error {
	globalAppConfig := config.GlobalConfig.App
	// 注册路由
	fmt.Println("注册路由....")
	r := initRouter(app)

	logger.Info("服务启动",
		zap.String("env", globalAppConfig.Env),
		zap.Int("port", globalAppConfig.Port),
	)

	fmt.Println("服务启动中....")
	if err := r.Run(fmt.Sprintf("%s:%d", globalAppConfig.Host, globalAppConfig.Port)); err != nil {
		return err
	}
	return nil
}
