package user

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	srv *service.ServiceProvider
}

func NewUserController(srv *service.ServiceProvider) *UserController {
	return &UserController{
		srv: srv,
	}
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
func (s *UserController) DelUser(c *gin.Context) {
	var data map[string]interface{}

	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	id, ok := data["id"].(string)
	if !ok || id == "" {
		utils.Fail(c, "请选择要删除的数据")
		return
	}
	if err := s.srv.UserService.DelUser(id); err != nil {
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
func (s *UserController) UpdateUser(c *gin.Context) {
	var userReq dto.UserUpdateReq
	_ = utils.ShouldBind(c, &userReq)

	if err := s.srv.UserService.UpdateUser(&userReq); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c)
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
func (s *UserController) GetListByPage(c *gin.Context) {
	var userReq dto.UserListPageReq
	_ = utils.ShouldBind(c, &userReq)

	res, err := s.srv.UserService.GetListByPage(&userReq)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// func (s *UserController) GetList(c *gin.Context) {
// 	var userReq dto.UserListReq
// 	_ = utils.ShouldBind(c, &userReq)

// 	res, err := s.srv.UserService.(&userReq)
// 	if err != nil {
// 		utils.Fail(c, err.Error())
// 		return
// 	}
// 	utils.Success(c, res, "获取成功")
// }

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
func (s *UserController) QueryDetail(c *gin.Context) {
	res, err := s.srv.UserService.GetUserInfo(c)
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
func (s *UserController) GetUserPermission(c *gin.Context) {
	userId, _ := c.Get("user_id")
	res, err := s.srv.UserService.GetUserPermission(userId.(string))
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}
