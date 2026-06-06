package router

import (
	"fmt"

	"github.com/Struggle-Rabbit/CampusLogistics/internal/config"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/handler"
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
	"gorm.io/gorm"

	_ "github.com/Struggle-Rabbit/CampusLogistics/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Use(
		middleware.Recovery(),
		middleware.RequestID(),
		middleware.Logger(),
		middleware.CORS(),
	)

	// 创建 service
	menuSvc := menu.NewMenuService(db)
	roleSvc := role.NewRoleService(db)
	userSvc := user.NewUserService(db, menuSvc)
	systemSvc := system.NewSystemService(db)
	repairSvc := repair.NewRepairService(db)
	campusSvc := campus.NewCampusService(db)
	buildingSvc := building.NewBuildingService(db)
	dormSvc := dorm.NewDormService(db)
	utilitySvc := utility.NewUtilityService(db)
	noticeSvc := notice.NewNoticeService(db)

	// 创建 handler
	commonH := handler.NewCommonHandler(userSvc)
	userH := handler.NewUserHandler(userSvc)
	menuH := handler.NewMenuHandler(menuSvc)
	roleH := handler.NewRoleHandler(roleSvc)
	systemH := handler.NewSystemHandler(systemSvc)
	repairH := handler.NewRepairHandler(repairSvc)
	campusH := handler.NewCampusHandler(campusSvc)
	buildingH := handler.NewBuildingHandler(buildingSvc)
	dormH := handler.NewDormHandler(dormSvc)
	utilityH := handler.NewUtilityHandler(utilitySvc)
	noticeH := handler.NewNoticeHandler(noticeSvc)

	api := r.Group("/api/v1")
	api.Use(middleware.OperationLogMiddleware())
	{
		// 公共路由（无需认证）
		api.POST("/register", commonH.Register)
		api.POST("/login", commonH.Login)
		noticeGroup := api.Group("/notice")
		{
			noticeGroup.GET("/public", noticeH.GetPublicList)
		}

		api.Use(middleware.JWTAuth())
		{
			// 用户路由
			userGroup := api.Group("/user")
			{
				userGroup.GET("/listPage", middleware.PermissionValidator("sys:user:list"), userH.GetListByPage)
				userGroup.GET("/detail", middleware.PermissionValidator("sys:user:detail"), userH.QueryDetail)
				userGroup.POST("/del", middleware.PermissionValidator("sys:user:del"), userH.DelUser)
				userGroup.POST("/update", middleware.PermissionValidator("sys:user:update"), userH.UpdateUser)
				userGroup.POST("/resetPassword", userH.ResetPassword)
				userGroup.GET("/getUserPermission", userH.GetUserPermission)
			}

			// 系统路由
			api.POST("/OperationLogList", middleware.PermissionValidator("sys:optLog"), systemH.GetOperationLogListByPage)

			// 角色路由
			roleGroup := api.Group("/role")
			{
				roleGroup.GET("/listPage", middleware.PermissionValidator("sys:role:list"), roleH.GetListByPage)
				roleGroup.GET("/detail", middleware.PermissionValidator("sys:role:detail"), roleH.QueryDetail)
				roleGroup.GET("/list", middleware.PermissionValidator("sys:role:list"), roleH.GetList)
				roleGroup.POST("/add", middleware.PermissionValidator("sys:role:add"), roleH.CreateRole)
				roleGroup.POST("/del", middleware.PermissionValidator("sys:role:del"), roleH.DelRole)
				roleGroup.POST("/update", middleware.PermissionValidator("sys:role:update"), roleH.UpdateRole)
			}

			// 菜单路由
			menuGroup := api.Group("/menu")
			{
				menuGroup.GET("/listPage", middleware.PermissionValidator("sys:menu:list"), menuH.GetListByPage)
				menuGroup.GET("/detail", middleware.PermissionValidator("sys:menu:detail"), menuH.QueryDetail)
				menuGroup.GET("/list", middleware.PermissionValidator("sys:menu:list"), menuH.GetList)
				menuGroup.POST("/add", middleware.PermissionValidator("sys:menu:add"), menuH.CreateMenu)
				menuGroup.POST("/del", middleware.PermissionValidator("sys:menu:del"), menuH.DelMenu)
				menuGroup.POST("/update", middleware.PermissionValidator("sys:menu:update"), menuH.UpdateMenu)
			}

			// 报修路由
			repairGroup := api.Group("/repair")
			{
				repairGroup.POST("/submit", middleware.PermissionValidator("repair:submit"), repairH.RepairOrderSubmit)
				repairGroup.GET("/list", middleware.PermissionValidator("repair:list"), repairH.GetListByPage)
				repairGroup.GET("/detail", middleware.PermissionValidator("repair:detail"), repairH.GetDetailById)
				repairGroup.POST("/update", middleware.PermissionValidator("repair:update"), repairH.UpdateRepairOrder)
				repairGroup.POST("/record", middleware.PermissionValidator("repair:record"), repairH.OrderRecord)
				repairGroup.POST("/del", middleware.PermissionValidator("repair:del"), repairH.DelRepairOrder)
			}

			// 校区路由
			campusGroup := api.Group("/campus")
			{
				campusGroup.POST("/create", middleware.PermissionValidator("campus:create"), campusH.Create)
				campusGroup.POST("/update", middleware.PermissionValidator("campus:update"), campusH.Update)
				campusGroup.POST("/del", middleware.PermissionValidator("campus:del"), campusH.Delete)
				campusGroup.GET("/list", middleware.PermissionValidator("campus:list"), campusH.GetListByPage)
				campusGroup.GET("/detail", middleware.PermissionValidator("campus:detail"), campusH.GetDetail)
				campusGroup.GET("/all", middleware.PermissionValidator("campus:list"), campusH.GetAll)
			}

			// 楼栋路由
			buildingGroup := api.Group("/building")
			{
				buildingGroup.POST("/create", middleware.PermissionValidator("building:create"), buildingH.Create)
				buildingGroup.POST("/update", middleware.PermissionValidator("building:update"), buildingH.Update)
				buildingGroup.POST("/del", middleware.PermissionValidator("building:del"), buildingH.Delete)
				buildingGroup.GET("/list", middleware.PermissionValidator("building:list"), buildingH.GetListByPage)
				buildingGroup.GET("/detail", middleware.PermissionValidator("building:detail"), buildingH.GetDetail)
				buildingGroup.GET("/byCampus", middleware.PermissionValidator("building:list"), buildingH.GetBuildingsByCampus)
				buildingGroup.POST("/import", middleware.PermissionValidator("building:import"), buildingH.ImportBuildings)
				buildingGroup.GET("/export", middleware.PermissionValidator("building:export"), buildingH.ExportBuildings)
			}

			// 宿舍路由
			dormGroup := api.Group("/dorm")
			{
				dormGroup.POST("/create", middleware.PermissionValidator("dorm:create"), dormH.Create)
				dormGroup.POST("/update", middleware.PermissionValidator("dorm:update"), dormH.Update)
				dormGroup.POST("/del", middleware.PermissionValidator("dorm:del"), dormH.Delete)
				dormGroup.GET("/list", middleware.PermissionValidator("dorm:list"), dormH.GetListByPage)
				dormGroup.GET("/detail", middleware.PermissionValidator("dorm:detail"), dormH.GetDetail)
				dormGroup.POST("/assign", middleware.PermissionValidator("dorm:assign"), dormH.AssignDorm)
				dormGroup.POST("/transfer", middleware.PermissionValidator("dorm:transfer"), dormH.TransferDorm)
				dormGroup.POST("/checkout", middleware.PermissionValidator("dorm:checkout"), dormH.CheckOut)
				dormGroup.GET("/users", middleware.PermissionValidator("dorm:users"), dormH.GetDormUsers)
				dormGroup.GET("/warning", middleware.PermissionValidator("dorm:warning"), dormH.GetCapacityWarning)
			}

			// 水电费路由
			utilityGroup := api.Group("/utility")
			{
				utilityGroup.POST("/create", middleware.PermissionValidator("utility:create"), utilityH.Create)
				utilityGroup.POST("/update", middleware.PermissionValidator("utility:update"), utilityH.Update)
				utilityGroup.POST("/del", middleware.PermissionValidator("utility:del"), utilityH.Delete)
				utilityGroup.GET("/list", middleware.PermissionValidator("utility:list"), utilityH.GetListByPage)
				utilityGroup.GET("/detail", middleware.PermissionValidator("utility:detail"), utilityH.GetDetail)
				utilityGroup.POST("/pay", middleware.PermissionValidator("utility:pay"), utilityH.Pay)
				utilityGroup.POST("/batchPay", middleware.PermissionValidator("utility:batchPay"), utilityH.BatchPay)
				utilityGroup.GET("/price", middleware.PermissionValidator("utility:price"), utilityH.GetPrice)
				utilityGroup.POST("/price", middleware.PermissionValidator("utility:price"), utilityH.UpdatePrice)
				utilityGroup.GET("/statistics", middleware.PermissionValidator("utility:statistics"), utilityH.GetStatistics)
				utilityGroup.GET("/warning", middleware.PermissionValidator("utility:warning"), utilityH.GetUnpaidWarning)
				utilityGroup.GET("/myUtility", middleware.PermissionValidator(""), utilityH.GetUserDormUtility)
			}

			// 公告路由
			noticeGroup := api.Group("/notice")
			{
				noticeGroup.POST("/create", middleware.PermissionValidator("notice:create"), noticeH.Create)
				noticeGroup.POST("/update", middleware.PermissionValidator("notice:update"), noticeH.Update)
				noticeGroup.POST("/del", middleware.PermissionValidator("notice:del"), noticeH.Delete)
				noticeGroup.GET("/list", middleware.PermissionValidator("notice:list"), noticeH.GetListByPage)
				noticeGroup.GET("/detail", middleware.PermissionValidator("notice:detail"), noticeH.GetDetail)
				noticeGroup.POST("/top", middleware.PermissionValidator("notice:top"), noticeH.SetTop)
			}
		}
	}

	return r
}

func Run(db *gorm.DB) error {
	globalAppConfig := config.GlobalConfig.App
	fmt.Println("注册路由....")
	r := initRouter(db)

	logger.Info("服务启动",
		zap.String("env", globalAppConfig.Env),
		zap.Int("port", globalAppConfig.Port),
	)

	if err := r.Run(fmt.Sprintf("%s:%d", globalAppConfig.Host, globalAppConfig.Port)); err != nil {
		return err
	}
	return nil
}
