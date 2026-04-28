package router

import (
	"github.com/Struggle-Rabbit/CampusLogistics/internal/controller/dorm"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/middleware"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/gin-gonic/gin"
)

func LoadDormRouter(api *gin.RouterGroup, srv *service.ServiceProvider) {
	dormCtl := dorm.NewDormController(srv)

	dormGroup := api.Group("/dorm")
	{
		dormGroup.POST("/create", middleware.PermissionValidator("dorm:create"), dormCtl.Create)
		dormGroup.POST("/update", middleware.PermissionValidator("dorm:update"), dormCtl.Update)
		dormGroup.POST("/del", middleware.PermissionValidator("dorm:del"), dormCtl.Delete)
		dormGroup.GET("/list", middleware.PermissionValidator("dorm:list"), dormCtl.GetListByPage)
		dormGroup.GET("/detail", middleware.PermissionValidator("dorm:detail"), dormCtl.GetDetail)
		dormGroup.POST("/assign", middleware.PermissionValidator("dorm:assign"), dormCtl.AssignDorm)
		dormGroup.POST("/transfer", middleware.PermissionValidator("dorm:transfer"), dormCtl.TransferDorm)
		dormGroup.POST("/checkout", middleware.PermissionValidator("dorm:checkout"), dormCtl.CheckOut)
		dormGroup.GET("/users", middleware.PermissionValidator("dorm:users"), dormCtl.GetDormUsers)
		dormGroup.GET("/warning", middleware.PermissionValidator("dorm:warning"), dormCtl.GetCapacityWarning)
	}
}
