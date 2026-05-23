package router

import (
	"github.com/Struggle-Rabbit/CampusLogistics/internal/controller/user"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/middleware"
	usersvc "github.com/Struggle-Rabbit/CampusLogistics/internal/service/user"
	"github.com/gin-gonic/gin"
)

func LoadUserRouter(api *gin.RouterGroup, userSvc usersvc.UserService) {
	userCtl := &user.UserControllerProvider{UserSvc: userSvc}

	user := api.Group("/user")
	{
		user.GET("/listPage", middleware.PermissionValidator("sys:user:list"), userCtl.GetListByPage)
		user.GET("/detail", middleware.PermissionValidator("sys:user:detail"), userCtl.QueryDetail)
		user.POST("/del", middleware.PermissionValidator("sys:user:del"), userCtl.DelUser)
		user.POST("/update", middleware.PermissionValidator("sys:user:update"), userCtl.UpdateUser)
		user.POST("/resetPassword", userCtl.ResetPassword)
		user.GET("/getUserPermission", userCtl.GetUserPermission)
	}

}
