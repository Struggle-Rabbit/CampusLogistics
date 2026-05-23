package router

import (
	"fmt"

	"github.com/Struggle-Rabbit/CampusLogistics/internal/app"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/config"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/middleware"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/building"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/campus"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/dorm"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/menu"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/notice"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/repair"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/role"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/system"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/user"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/utility"
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

	menuSvc := &menu.MenuServiceProvider{App: app}
	roleSvc := &role.RoleServiceProvider{App: app}
	userSvc := &user.UserServiceProvider{App: app, Menu: menuSvc}
	systemSvc := &system.SystemServiceProvider{App: app}
	repairSvc := &repair.RepairServiceProvider{App: app}
	campusSvc := &campus.CampusServiceProvider{App: app}
	buildingSvc := &building.BuildingServiceProvider{App: app}
	dormSvc := &dorm.DormServiceProvider{App: app}
	utilitySvc := &utility.UtilityServiceProvider{App: app}
	noticeSvc := &notice.NoticeServiceProvider{App: app}

	api := r.Group("/api/v1")
	api.Use(middleware.OperationLogMiddleware())
	{
		LoadCommonRouter(api, userSvc, noticeSvc)

		api.Use(middleware.JWTAuth())
		{
			LoadUserRouter(api, userSvc)
			LoadSystemRouter(api, systemSvc)
			LoadRoleRouter(api, roleSvc)
			LoadMenuRouter(api, menuSvc)
			LoadRepairRouter(api, repairSvc)
			LoadCampusRouter(api, campusSvc)
			LoadBuildingRouter(api, buildingSvc)
			LoadDormRouter(api, dormSvc)
			LoadUtilityRouter(api, utilitySvc)
			LoadNoticeRouter(api, noticeSvc)
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
