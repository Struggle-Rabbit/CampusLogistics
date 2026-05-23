package common

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/user"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
)

// CommonController 通用控制器接口
type CommonController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

// CommonControllerProvider 通用控制器实现
type CommonControllerProvider struct {
	UserSvc user.UserService
}

// Login 用户登录
// @Summary 用户登录接口
// @Description 登录系统
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body dto.LoginReq true "入参"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/login [post]
func (c *CommonControllerProvider) Login(ctx *gin.Context) {
	var req dto.LoginReq

	if isValidate := utils.ShouldBind(ctx, &req); isValidate {
		res, err := c.UserSvc.Login(&req)

		if err != nil {
			utils.Fail(ctx, err.Error())
			return
		}

		utils.Success(ctx, res, "登录成功")
		return
	}
}

// Register 用户注册
// @Summary 用户注册接口
// @Description 注册新用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body dto.RegisterReq true "注册参数"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/register [post]
func (c *CommonControllerProvider) Register(ctx *gin.Context) {
	var req dto.RegisterReq

	if isValidate := utils.ShouldBind(ctx, &req); isValidate {
		err := c.UserSvc.Register(&req)

		if err != nil {
			utils.Fail(ctx, err.Error())
			return
		}

		utils.Success(ctx, "注册成功")
		return

	}

}
