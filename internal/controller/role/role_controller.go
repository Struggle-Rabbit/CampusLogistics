package role

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
)

type RoleController struct {
	srv *service.ServiceProvider
}

func NewRoleController(srv *service.ServiceProvider) *RoleController {
	return &RoleController{
		srv: srv,
	}
}

// CreateRole 创建角色
// @Summary 角色创建接口
// @Description 角色创建
// @Tags 系统模块
// @Accept json
// @Produce json
// @Param data body dto.CreateRoleReq true "入参"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/role/add [post]
func (s *RoleController) CreateRole(c *gin.Context) {
	var roleReq dto.CreateRoleReq
	_ = utils.ShouldBind(c, &roleReq)
	if err := s.srv.RoleService.CreateRole(&roleReq); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c)
}

// DelRole 删除角色
// @Summary 角色删除接口
// @Description 角色删除
// @Tags 系统模块
// @Accept json
// @Produce json
// @Param data body map[string]interface{} true "入参"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/role/del [post]
func (s *RoleController) DelRole(c *gin.Context) {
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
	if err := s.srv.RoleService.DelRole(id); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "删除成功")
}

// UpdateRole 更新角色
// @Summary 角色修改接口
// @Description 角色修改
// @Tags 系统模块
// @Accept json
// @Produce json
// @Param data body dto.UpdateRoleReq true "入参"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/role/update [post]
func (s *RoleController) UpdateRole(c *gin.Context) {
	var roleReq dto.UpdateRoleReq
	_ = utils.ShouldBind(c, &roleReq)

	if err := s.srv.RoleService.UpdateRole(&roleReq); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c)
}

// GetListByPage 角色分页
// @Summary 角色分页查询接口
// @Description 角色分页查询
// @Tags 系统模块
// @Accept json
// @Produce json
// @Param data query dto.RoleListByPageReq true "入参"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/role/listPage [get]
func (s *RoleController) GetListByPage(c *gin.Context) {
	var roleReq dto.RoleListByPageReq
	_ = utils.ShouldBind(c, &roleReq)

	res, err := s.srv.RoleService.GetRoleListByPage(&roleReq)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// GetList 角色列表
// @Summary 角色列表查询接口
// @Description 角色列表查询
// @Tags 系统模块
// @Accept json
// @Produce json
// @Param data query dto.RoleListByPageReq true "入参"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/role/list [get]
func (s *RoleController) GetList(c *gin.Context) {
	var roleReq dto.RoleListReq
	_ = utils.ShouldBind(c, &roleReq)

	res, err := s.srv.RoleService.GetRoleList(&roleReq)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// QueryDetail 角色详情
// @Summary 角色详情查询接口
// @Description 角色详情查询
// @Tags 系统模块
// @Accept json
// @Produce json
// @Param id query string true "入参"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/role/detail [get]
func (s *RoleController) QueryDetail(c *gin.Context) {
	var data map[string]interface{}
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	id, ok := data["id"].(string)
	if !ok || id == "" {
		utils.Fail(c, "数据ID不能为空")
		return
	}
	res, err := s.srv.RoleService.RoleDetailById(id)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}
