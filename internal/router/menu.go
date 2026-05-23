package router

import (
	menuCtl "github.com/Struggle-Rabbit/CampusLogistics/internal/controller/menu"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/middleware"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/menu"
	"github.com/gin-gonic/gin"
)

func LoadMenuRouter(api *gin.RouterGroup, menuSvc menu.MenuService) {
	mc := &menuCtl.MenuControllerProvider{MenuSvc: menuSvc}

	user := api.Group("/menu")
	{
		user.GET("/listPage", middleware.PermissionValidator("sys:menu:list"), mc.GetListByPage)
		user.GET("/detail", middleware.PermissionValidator("sys:menu:detail"), mc.QueryDetail)
		user.GET("/list", middleware.PermissionValidator("sys:menu:list"), mc.GetList)
		user.POST("/add", middleware.PermissionValidator("sys:menu:add"), mc.CreateMenu)
		user.POST("/del", middleware.PermissionValidator("sys:menu:del"), mc.DelMenu)
		user.POST("/update", middleware.PermissionValidator("sys:menu:update"), mc.UpdateMenu)
	}

}
