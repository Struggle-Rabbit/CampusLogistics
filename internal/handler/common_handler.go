package handler

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/user"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
)

type commonHandler struct {
	userSvc *user.Service
}

func NewCommonHandler(userSvc *user.Service) *commonHandler {
	return &commonHandler{userSvc: userSvc}
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
func (h *commonHandler) Login(c *gin.Context) {
	var req dto.LoginReq

	if isValidate := utils.ShouldBind(c, &req); isValidate {
		res, err := h.userSvc.Login(&req)

		if err != nil {
			utils.Fail(c, err.Error())
			return
		}

		utils.Success(c, res, "登录成功")
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
func (h *commonHandler) Register(c *gin.Context) {
	var req dto.RegisterReq

	if isValidate := utils.ShouldBind(c, &req); isValidate {
		err := h.userSvc.Register(&req)

		if err != nil {
			utils.Fail(c, err.Error())
			return
		}

		utils.Success(c, "注册成功")
		return
	}
}
