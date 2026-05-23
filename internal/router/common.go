package router

import (
	"github.com/Struggle-Rabbit/CampusLogistics/internal/controller/common"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/controller/notice"
	noticesvc "github.com/Struggle-Rabbit/CampusLogistics/internal/service/notice"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/user"
	"github.com/gin-gonic/gin"
)

func LoadCommonRouter(api *gin.RouterGroup, userSvc user.UserService, noticeSvc noticesvc.NoticeService) {
	commonCtl := &common.CommonControllerProvider{UserSvc: userSvc}
	noticeCtl := &notice.NoticeControllerProvider{NoticeSvc: noticeSvc}

	api.POST("/register", commonCtl.Register)
	api.POST("/login", commonCtl.Login)
	noticeGroup := api.Group("/notice")
	{
		noticeGroup.GET("/public", noticeCtl.GetPublicList)
	}
}
