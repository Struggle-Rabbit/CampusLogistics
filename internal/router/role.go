package router

import (
	"github.com/Struggle-Rabbit/CampusLogistics/internal/controller/role"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/middleware"
	rolesvc "github.com/Struggle-Rabbit/CampusLogistics/internal/service/role"
	"github.com/gin-gonic/gin"
)

func LoadRoleRouter(api *gin.RouterGroup, roleSvc rolesvc.RoleService) {
	roleCtl := &role.RoleControllerProvider{RoleSvc: roleSvc}

	user := api.Group("/role")
	{
		user.GET("/listPage", middleware.PermissionValidator("sys:role:list"), roleCtl.GetListByPage)
		user.GET("/detail", middleware.PermissionValidator("sys:role:detail"), roleCtl.QueryDetail)
		user.GET("/list", middleware.PermissionValidator("sys:role:list"), roleCtl.GetList)
		user.POST("/add", middleware.PermissionValidator("sys:role:add"), roleCtl.CreateRole)
		user.POST("/del", middleware.PermissionValidator("sys:role:del"), roleCtl.DelRole)
		user.POST("/update", middleware.PermissionValidator("sys:role:update"), roleCtl.UpdateRole)
	}

}
