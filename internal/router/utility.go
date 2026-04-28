package router

import (
	"github.com/Struggle-Rabbit/CampusLogistics/internal/controller/utility"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/middleware"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/gin-gonic/gin"
)

func LoadUtilityRouter(api *gin.RouterGroup, srv *service.ServiceProvider) {
	utilityCtl := utility.NewUtilityController(srv)

	utilityGroup := api.Group("/utility")
	{
		utilityGroup.POST("/create", middleware.PermissionValidator("utility:create"), utilityCtl.Create)
		utilityGroup.POST("/update", middleware.PermissionValidator("utility:update"), utilityCtl.Update)
		utilityGroup.POST("/del", middleware.PermissionValidator("utility:del"), utilityCtl.Delete)
		utilityGroup.GET("/list", middleware.PermissionValidator("utility:list"), utilityCtl.GetListByPage)
		utilityGroup.GET("/detail", middleware.PermissionValidator("utility:detail"), utilityCtl.GetDetail)
		utilityGroup.POST("/pay", middleware.PermissionValidator("utility:pay"), utilityCtl.Pay)
		utilityGroup.POST("/batchPay", middleware.PermissionValidator("utility:batchPay"), utilityCtl.BatchPay)
		utilityGroup.GET("/price", middleware.PermissionValidator("utility:price"), utilityCtl.GetPrice)
		utilityGroup.POST("/price", middleware.PermissionValidator("utility:price"), utilityCtl.UpdatePrice)
		utilityGroup.GET("/statistics", middleware.PermissionValidator("utility:statistics"), utilityCtl.GetStatistics)
		utilityGroup.GET("/warning", middleware.PermissionValidator("utility:warning"), utilityCtl.GetUnpaidWarning)
		utilityGroup.GET("/myUtility", middleware.PermissionValidator(""), utilityCtl.GetUserDormUtility)
	}
}
