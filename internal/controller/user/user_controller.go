package user

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/user"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
)

// UserController 用户控制器接口
type UserController interface {
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	GetUserList(c *gin.Context)
	GetUserInfo(c *gin.Context)
	DeleteUser(c *gin.Context)
	ResetPassword(c *gin.Context)
	GetUserPermission(c *gin.Context)
}

// UserControllerProvider 用户控制器实现
type UserControllerProvider struct {
	UserSvc user.UserService
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
func (s *UserControllerProvider) DelUser(c *gin.Context) {
	var data map[string]interface{}

	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	id, ok := data["id"].([]string)
	if !ok || len(id) == 0 {
		utils.Fail(c, "请选择要删除的数据")
		return
	}
	if err := s.UserSvc.DelUser(id); err != nil {
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
func (s *UserControllerProvider) UpdateUser(c *gin.Context) {
	var userReq dto.UserUpdateReq
	if !utils.ShouldBind(c, &userReq) {
		return
	}

	if err := s.UserSvc.UpdateUser(&userReq); err != nil {
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
func (s *UserControllerProvider) ResetPassword(c *gin.Context) {
	var resetReq dto.PasswordReset
	if !utils.ShouldBind(c, &resetReq) {
		return
	}

	if err := s.UserSvc.ResetPassword(&resetReq); err != nil {
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
func (s *UserControllerProvider) GetListByPage(c *gin.Context) {
	var userReq dto.UserListPageReq
	if !utils.ShouldBind(c, &userReq) {
		return
	}

	res, err := s.UserSvc.GetListByPage(&userReq)
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
func (s *UserControllerProvider) QueryDetail(c *gin.Context) {
	res, err := s.UserSvc.GetUserInfo(c)
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
func (s *UserControllerProvider) GetUserPermission(c *gin.Context) {
	userId, _ := c.Get("user_id")
	res, err := s.UserSvc.GetUserPermission(userId.(string))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}
