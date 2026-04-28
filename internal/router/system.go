package router

import (
	"github.com/Struggle-Rabbit/CampusLogistics/internal/controller/system"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/middleware"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/gin-gonic/gin"
)

func LoadSystemRouter(api *gin.RouterGroup, srv *service.ServiceProvider) {

	sysCtl := system.NewSystemController(srv)

	api.POST("/OperationLogList", middleware.PermissionValidator("sys:optLog"), sysCtl.GetOperationLogListByPage)
}
