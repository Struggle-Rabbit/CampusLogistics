package router

import (
	"github.com/Struggle-Rabbit/CampusLogistics/internal/controller/user"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/middleware"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/gin-gonic/gin"
)

func LoadUserRouter(api *gin.RouterGroup, srv *service.ServiceProvider) {
	userCtl := user.NewUserController(srv)

	user := api.Group("/user")
	{
		user.GET("/listPage", middleware.PermissionValidator("sys:user:list"), userCtl.GetListByPage)
		user.GET("/detail", middleware.PermissionValidator("sys:user:detail"), userCtl.QueryDetail)
		user.GET("/del", middleware.PermissionValidator("sys:user:del"), userCtl.DelUser)
		user.GET("/update", middleware.PermissionValidator("sys:user:update"), userCtl.UpdateUser)
		user.GET("/getUserPermission", userCtl.GetUserPermission)
	}

}
