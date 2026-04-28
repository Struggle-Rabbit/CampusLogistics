package router

import (
	"github.com/Struggle-Rabbit/CampusLogistics/internal/controller/notice"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/middleware"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/gin-gonic/gin"
)

func LoadNoticeRouter(api *gin.RouterGroup, srv *service.ServiceProvider) {
	noticeCtl := notice.NewNoticeController(srv)

	noticeGroup := api.Group("/notice")
	{
		noticeGroup.POST("/create", middleware.PermissionValidator("notice:create"), noticeCtl.Create)
		noticeGroup.POST("/update", middleware.PermissionValidator("notice:update"), noticeCtl.Update)
		noticeGroup.POST("/del", middleware.PermissionValidator("notice:del"), noticeCtl.Delete)
		noticeGroup.GET("/list", middleware.PermissionValidator("notice:list"), noticeCtl.GetListByPage)
		noticeGroup.GET("/detail", middleware.PermissionValidator("notice:detail"), noticeCtl.GetDetail)
		noticeGroup.POST("/top", middleware.PermissionValidator("notice:top"), noticeCtl.SetTop)
		noticeGroup.GET("/public", noticeCtl.GetPublicList)
	}
}
