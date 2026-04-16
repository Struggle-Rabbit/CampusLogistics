package router

import (
	"github.com/Struggle-Rabbit/CampusLogistics/internal/controller/systeam"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/middleware"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/gin-gonic/gin"
)

func LoadSysteamRouter(api *gin.RouterGroup, srv *service.ServiceProvider) {

	sysCtl := systeam.NewSysteamController(srv)

	api.POST("/OperationLogList", middleware.PermissionValidator("sys:optLog"), sysCtl.GetOperationLogListByPage)
}
