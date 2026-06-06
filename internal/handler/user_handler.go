package handler

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/user"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	svc *user.Service
}

func NewUserHandler(svc *user.Service) *userHandler {
	return &userHandler{svc: svc}
}

// DelUser 用户删除
// @Summary 用户删除接口
// @Description 删除用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id body string true "入参"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/user/del [post]
func (h *userHandler) DelUser(c *gin.Context) {
	var data map[string]interface{}

	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	ids, ok := data["id"].([]interface{})
	if !ok || len(ids) == 0 {
		utils.Fail(c, "请选择要删除的数据")
		return
	}

	idStrs := make([]string, 0, len(ids))
	for _, id := range ids {
		if idStr, ok := id.(string); ok {
			idStrs = append(idStrs, idStr)
		}
	}
	if err := h.svc.DelUser(idStrs); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "删除成功")
}

// UpdateUser 用户更新
// @Summary 用户更新接口
// @Description 用户信息更新
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body dto.UserUpdateReq true "入参"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/user/listPage [post]
func (h *userHandler) UpdateUser(c *gin.Context) {
	var userReq dto.UserUpdateReq
	if !utils.ShouldBind(c, &userReq) {
		return
	}

	if err := h.svc.UpdateUser(&userReq); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c)
}

// ResetPassword 重置密码
// @Summary 重置密码接口
// @Description 重置密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body dto.PasswordReset true "入参"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/user/resetPassword [post]
func (h *userHandler) ResetPassword(c *gin.Context) {
	var resetReq dto.PasswordReset
	if !utils.ShouldBind(c, &resetReq) {
		return
	}

	if err := h.svc.ResetPassword(&resetReq); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "密码重置成功")
}

// GetListByPage 用户分页
// @Summary 用户分页查询接口
// @Description 用户分页查询
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data query dto.UserListPageReq true "入参"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/user/listPage [get]
func (h *userHandler) GetListByPage(c *gin.Context) {
	var userReq dto.UserListPageReq
	if !utils.ShouldBind(c, &userReq) {
		return
	}

	res, err := h.svc.GetListByPage(&userReq)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// QueryDetail 用户详情
// @Summary 用户详情查询接口
// @Description 用户详情查询
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id query string true "入参"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/user/detail [get]
func (h *userHandler) QueryDetail(c *gin.Context) {
	userId, exists := c.Get("userID")
	if !exists {
		utils.Fail(c, "未获取到UserID")
		return
	}
	res, err := h.svc.GetUserInfo(userId.(string))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// GetUserPermission 用户菜单权限
// @Summary 用户菜单权限查询接口
// @Description 用户菜单权限查询
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/user/getUserPermission [get]
func (h *userHandler) GetUserPermission(c *gin.Context) {
	userId, _ := c.Get("userID")
	res, err := h.svc.GetUserPermission(userId.(string))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}
