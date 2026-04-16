package menu

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
)

type MenuController struct {
	srv *service.ServiceProvider
}

func NewMenuController(srv *service.ServiceProvider) *MenuController {
	return &MenuController{
		srv: srv,
	}
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
func (mc *MenuController) CreateMenu(c *gin.Context) {
	var menuReq dto.CreateMenuReq
	_ = utils.ShouldBind(c, &menuReq)
	if err := mc.srv.MenuService.CreateMenu(&menuReq); err != nil {
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
func (mc *MenuController) DelMenu(c *gin.Context) {
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
	if err := mc.srv.MenuService.DelMenu(id); err != nil {
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
func (mc *MenuController) UpdateMenu(c *gin.Context) {
	var menuReq dto.UpdateMenuReq
	_ = utils.ShouldBind(c, &menuReq)

	if err := mc.srv.MenuService.UpdateMenu(&menuReq); err != nil {
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
func (mc *MenuController) GetListByPage(c *gin.Context) {
	var menuReq dto.MenuListByPageReq
	_ = utils.ShouldBind(c, &menuReq)

	res, err := mc.srv.MenuService.GetMenuListByPage(&menuReq)
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
func (mc *MenuController) GetList(c *gin.Context) {
	var menuReq dto.MenuListReq
	_ = utils.ShouldBind(c, &menuReq)

	res, err := mc.srv.MenuService.GetMenuList(&menuReq)
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
func (mc *MenuController) QueryDetail(c *gin.Context) {
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
	res, err := mc.srv.MenuService.MenuDetailById(id)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}
