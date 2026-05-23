package router

import (
	systemctl "github.com/Struggle-Rabbit/CampusLogistics/internal/controller/system"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/middleware"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/system"
	"github.com/gin-gonic/gin"
)

func LoadSystemRouter(api *gin.RouterGroup, systemSvc system.SystemService) {

	sysCtl := &systemctl.SystemControllerProvider{SystemSvc: systemSvc}

	api.POST("/OperationLogList", middleware.PermissionValidator("sys:optLog"), sysCtl.GetOperationLogListByPage)
}
