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

func initRouter(app *app.App) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(
		middleware.Recovery(),
		middleware.RequestID(),
		middleware.Logger(),
		middleware.CORS(),
	)

	srv := service.NewServiceProvider(app)

	api := r.Group("/api/v1")
	api.Use(middleware.OperationLogMiddleware())
	{
		LoadCommonRouter(api, srv)

		api.Use(middleware.JWTAuth())
		{
			LoadUserRouter(api, srv)
			LoadSystemRouter(api, srv)
			LoadRoleRouter(api, srv)
			LoadMenuRouter(api, srv)
			LoadRepairRouter(api, srv)
			LoadCampusRouter(api, srv)
			LoadBuildingRouter(api, srv)
			LoadDormRouter(api, srv)
			LoadUtilityRouter(api, srv)
			LoadNoticeRouter(api, srv)
		}

	}

	return r
}

func Run(app *app.App) error {
	globalAppConfig := config.GlobalConfig.App
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
