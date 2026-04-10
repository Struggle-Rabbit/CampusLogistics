package common

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
)

type CommonController struct {
	srv *service.ServiceProvider
}

func NewCommonController(srv *service.ServiceProvider) *CommonController {
	return &CommonController{
		srv: srv,
	}
}

// Login 登录
func (uCtl *CommonController) Login(c *gin.Context) {
	var req dto.LoginReq

	if isValidate := utils.ShouldBind(c, &req); isValidate {
		res, err := uCtl.srv.UserService.Login(&req)

		if err != nil {
			utils.Fail(c, err.Error())
			return
		}

		utils.Success(c, res, "登录成功")
		return
	}
}

// Register 注册
func (uCtl *CommonController) Register(c *gin.Context) {
	var req dto.RegisterReq

	if isValidate := utils.ShouldBind(c, &req); isValidate {
		err := uCtl.srv.UserService.Register(&req)

		if err != nil {
			utils.Fail(c, err.Error())
			return
		}

		utils.Success(c, "注册成功")
		return

	}

}
