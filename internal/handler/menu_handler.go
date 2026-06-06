package handler

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/menu"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
)

type menuHandler struct {
	svc *menu.Service
}

func NewMenuHandler(svc *menu.Service) *menuHandler {
	return &menuHandler{svc: svc}
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
func (h *menuHandler) CreateMenu(c *gin.Context) {
	var menuReq dto.CreateMenuReq
	if !utils.ShouldBind(c, &menuReq) {
		return
	}
	if err := h.svc.CreateMenu(&menuReq); err != nil {
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
func (h *menuHandler) DelMenu(c *gin.Context) {
	var data map[string]interface{}

	if err := c.ShouldBindJSON(&data); err != nil {
		utils.Fail(c, "参数错误")
		return
	}

	ids, ok := data["ids"].([]interface{})
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
	if err := h.svc.DelMenu(idStrs); err != nil {
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
func (h *menuHandler) UpdateMenu(c *gin.Context) {
	var menuReq dto.UpdateMenuReq
	if !utils.ShouldBind(c, &menuReq) {
		return
	}

	if err := h.svc.UpdateMenu(&menuReq); err != nil {
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
func (h *menuHandler) GetListByPage(c *gin.Context) {
	var menuReq dto.MenuListByPageReq
	if !utils.ShouldBind(c, &menuReq) {
		return
	}

	res, err := h.svc.GetMenuListByPage(&menuReq)
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
func (h *menuHandler) GetList(c *gin.Context) {
	var menuReq dto.MenuListReq
	if !utils.ShouldBind(c, &menuReq) {
		return
	}

	res, err := h.svc.GetMenuList(&menuReq)
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
func (h *menuHandler) QueryDetail(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		utils.Fail(c, "数据ID不能为空")
		return
	}
	res, err := h.svc.MenuDetailById(id)
	if err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c, res, "获取成功")
}
