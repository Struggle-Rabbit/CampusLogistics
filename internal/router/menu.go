package router

import (
	"github.com/Struggle-Rabbit/CampusLogistics/internal/controller/menu"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/middleware"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/gin-gonic/gin"
)

func LoadMenuRouter(api *gin.RouterGroup, srv *service.ServiceProvider) {
	menuCtl := menu.NewMenuController(srv)

	user := api.Group("/menu")
	{
		user.GET("/listPage", middleware.PermissionValidator("sys:menu:list"), menuCtl.GetListByPage)
		user.GET("/detail", middleware.PermissionValidator("sys:menu:detail"), menuCtl.QueryDetail)
		user.GET("/list", middleware.PermissionValidator("sys:menu:list"), menuCtl.GetList)
		user.POST("/add", middleware.PermissionValidator("sys:menu:add"), menuCtl.CreateMenu)
		user.POST("/del", middleware.PermissionValidator("sys:menu:del"), menuCtl.DelMenu)
		user.POST("/update", middleware.PermissionValidator("sys:menu:update"), menuCtl.UpdateMenu)
	}

}
