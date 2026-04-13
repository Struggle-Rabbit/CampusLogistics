package menu

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/app"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"github.com/go-viper/mapstructure/v2"
)

type MenuService struct {
	app *app.App
}

func NewMenuService(app *app.App) *MenuService {
	return &MenuService{
		app: app,
	}
}

func (s *MenuService) CreateMenu(req *dto.CreateMenuReq) error {
	var menu model.SysMenu
	if err := mapstructure.Decode(req, &menu); err != nil {
		return err
	}
	return s.app.DB.Create(&menu).Error
}

func (s *MenuService) UpdateMenu(req *dto.UpdateMenuReq) error {
	var menu model.SysMenu
	if err := mapstructure.Decode(req, &menu); err != nil {
		return err
	}
	return s.app.DB.Save(&menu).Error
}

func (s *MenuService) DelMenu(id string) error {

	return s.app.DB.Delete(&model.SysMenu{}, id).Error
}

func (s *MenuService) GetMenuList(req *dto.MenuListReq) ([]*dto.MenuResult, error) {
	var menuReq model.SysMenu
	if err := mapstructure.Decode(req, &menuReq); err != nil {
		return nil, err
	}
	var menuSqlRes []*dto.MenuResult

	s.app.DB.Where(&menuReq).Find(&menuSqlRes)

	menuTree := s.BuildMenuTree(menuSqlRes, "0")
	return menuTree, nil
}

func (s *MenuService) GetMenuListByPage(req *dto.MenuListByPageReq) (*dto.PageResult, error) {
	var menuReq model.SysMenu
	var total int64
	if err := s.app.DB.Model(&model.SysMenu{}).Where("parent_id = ?", "0").Count(&total).Error; err != nil {
		return nil, err
	}
	if err := mapstructure.Decode(req, &menuReq); err != nil {
		return nil, err
	}
	var list []*dto.MenuResult

	if err := s.app.DB.Model(&model.SysMenu{}).Where(&menuReq).Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Find(&list).Error; err != nil {
		return nil, err
	}

	menuTree := s.BuildMenuTree(list, "0")

	return &dto.PageResult{
		List:        menuTree,
		Total:       total,
		PageSize:    req.PageSize,
		CurrentPage: req.CurrentPage,
	}, nil
}

func (s *MenuService) MenuDetailById(id string) (*dto.MenuResult, error) {
	var menuResult dto.MenuResult

	if err := s.app.DB.Model(&model.SysMenu{}).Where("id = ?", id).First(&menuResult).Error; err != nil {
		return nil, err
	}

	return &menuResult, nil
}

// 处理树形菜单结构
func (s *MenuService) BuildMenuTree(allMenus []*dto.MenuResult, parentID string) []*dto.MenuResult {

	var tree []*dto.MenuResult

	for _, menu := range allMenus {
		if menu.ParentID == parentID {
			children := s.BuildMenuTree(allMenus, menu.ID)
			menu.Childen = children

			tree = append(tree, menu)
		}
	}
	return tree
}
