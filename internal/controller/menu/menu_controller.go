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

func (mc *MenuController) CreateMenu(c *gin.Context) {
	var menuReq dto.CreateMenuReq
	_ = utils.ShouldBind(c, &menuReq)
	if err := mc.srv.MenuService.CreateMenu(&menuReq); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c)
}

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

func (mc *MenuController) UpdateMenu(c *gin.Context) {
	var menuReq dto.UpdateMenuReq
	_ = utils.ShouldBind(c, &menuReq)

	if err := mc.srv.MenuService.UpdateMenu(&menuReq); err != nil {
		utils.Fail(c, err.Error())
		return
	}
	utils.Success(c)
}

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
