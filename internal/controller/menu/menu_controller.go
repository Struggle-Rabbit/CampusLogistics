package menu

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/menu"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
)

type MenuController interface {
	CreateMenu(c *gin.Context)
	DelMenu(c *gin.Context)
	UpdateMenu(c *gin.Context)
	GetListByPage(c *gin.Context)
	GetList(c *gin.Context)
	QueryDetail(c *gin.Context)
}

type MenuControllerProvider struct {
	MenuSvc menu.MenuService
}

// CreateMenu 菜单创建
// @Summary 菜单创建接口
// @Description 菜单创建
// @Tags 系统模块
// @Accept json
// @Produce json
// @Param data body dto.CreateMenuReq true "入参"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/menu/add [post]
func (mc *MenuControllerProvider) CreateMenu(c *gin.Context) {
	var menuReq dto.CreateMenuReq
	_ = utils.ShouldBind(c, &menuReq)
	if err := mc.MenuSvc.CreateMenu(&menuReq); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c)
}

// DelMenu 菜单删除
// @Summary 菜单删除接口
// @Description 菜单删除
// @Tags 系统模块
// @Accept json
// @Produce json
// @Param id query string false "入参"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/menu/del [post]
func (mc *MenuControllerProvider) DelMenu(c *gin.Context) {
	var data map[string]interface{}

	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	id, ok := data["ids"].([]string)
	if !ok || len(id) == 0 {
		utils.Fail(c, "请选择要删除的数据")
		return
	}
	if err := mc.MenuSvc.DelMenu(id); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, "删除成功")
}

// UpdateMenu 菜单更新
// @Summary 菜单更新接口
// @Description 菜单更新
// @Tags 系统模块
// @Accept json
// @Produce json
// @Param data body dto.UpdateMenuReq false "入参"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/menu/update [post]
func (mc *MenuControllerProvider) UpdateMenu(c *gin.Context) {
	var menuReq dto.UpdateMenuReq
	_ = utils.ShouldBind(c, &menuReq)

	if err := mc.MenuSvc.UpdateMenu(&menuReq); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c)
}

// GetListByPage 菜单分页查询
// @Summary 菜单查询接口
// @Description 菜单分页查询
// @Tags 系统模块
// @Accept json
// @Produce json
// @Param data query dto.MenuListByPageReq false "入参"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/menu/listPage [get]
func (mc *MenuControllerProvider) GetListByPage(c *gin.Context) {
	var menuReq dto.MenuListByPageReq
	_ = utils.ShouldBind(c, &menuReq)

	res, err := mc.MenuSvc.GetMenuListByPage(&menuReq)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// GetList 菜单查询
// @Summary 菜单查询接口
// @Description 菜单查询
// @Tags 系统模块
// @Accept json
// @Produce json
// @Param data query dto.MenuListByPageReq false "入参"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/menu/listPage [get]
func (mc *MenuControllerProvider) GetList(c *gin.Context) {
	var menuReq dto.MenuListReq
	_ = utils.ShouldBind(c, &menuReq)

	res, err := mc.MenuSvc.GetMenuList(&menuReq)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}

// QueryDetail 菜单详情
// @Summary 菜单详情接口
// @Description 菜单详情查询
// @Tags 系统模块
// @Accept json
// @Produce json
// @Param id query string true "入参"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrResponse
// @Router /api/v1/menu/detail [get]
func (mc *MenuControllerProvider) QueryDetail(c *gin.Context) {
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
	res, err := mc.MenuSvc.MenuDetailById(id)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}
