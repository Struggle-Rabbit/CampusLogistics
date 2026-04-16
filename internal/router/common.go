package router

import (
	"github.com/Struggle-Rabbit/CampusLogistics/internal/controller/common"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/gin-gonic/gin"
)

func LoadCommonRouter(api *gin.RouterGroup, srv *service.ServiceProvider) {

	commonCtl := common.NewCommonController(srv)

	api.POST("/register", commonCtl.Register)
	api.POST("/login", commonCtl.Login)
}
