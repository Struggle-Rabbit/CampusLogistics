package router

import (
	noticecontroller "github.com/Struggle-Rabbit/CampusLogistics/internal/controller/notice"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/middleware"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/notice"
	"github.com/gin-gonic/gin"
)

func LoadNoticeRouter(api *gin.RouterGroup, noticeSvc notice.NoticeService) {
	noticeCtl := &noticecontroller.NoticeControllerProvider{NoticeSvc: noticeSvc}

	noticeGroup := api.Group("/notice")
	{
		noticeGroup.POST("/create", middleware.PermissionValidator("notice:create"), noticeCtl.Create)
		noticeGroup.POST("/update", middleware.PermissionValidator("notice:update"), noticeCtl.Update)
		noticeGroup.POST("/del", middleware.PermissionValidator("notice:del"), noticeCtl.Delete)
		noticeGroup.GET("/list", middleware.PermissionValidator("notice:list"), noticeCtl.GetListByPage)
		noticeGroup.GET("/detail", middleware.PermissionValidator("notice:detail"), noticeCtl.GetDetail)
		noticeGroup.POST("/top", middleware.PermissionValidator("notice:top"), noticeCtl.SetTop)
	}
}
