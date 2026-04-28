package menu

import (
	"errors"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/app"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
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
	menu := &model.SysMenu{
		ParentID:    req.ParentID,
		Name:        req.Name,
		Type:        req.Type,
		Perms:       req.Perms,
		Status:      req.Status,
		Sort:        req.Sort,
		Icon:        req.Icon,
		Description: req.Description,
	}

	if req.Type == 3 {
		if req.Path == "" {
			return errors.New("路由地址不能为空！")
		}
		if req.Component == "" {
			return errors.New("组件地址不能为空")
		}
		menu.Path = req.Path
		menu.Component = req.Component
	}

	return s.app.DB.Create(menu).Error
}

func (s *MenuService) UpdateMenu(req *dto.UpdateMenuReq) error {
	menu := &model.SysMenu{
		ParentID:    req.ParentID,
		Name:        req.Name,
		Type:        req.Type,
		Perms:       req.Perms,
		Status:      req.Status,
		Sort:        req.Sort,
		Icon:        req.Icon,
		Description: req.Description,
	}

	if req.Type == 3 {
		if req.Path == "" {
			return errors.New("路由地址不能为空！")
		}
		if req.Component == "" {
			return errors.New("组件地址不能为空")
		}
		menu.Path = req.Path
		menu.Component = req.Component
	}
	return s.app.DB.Where("id = ?", req.ID).Updates(&menu).Error
}

func (s *MenuService) DelMenu(id []string) error {

	return s.app.DB.Delete(&model.SysMenu{}, id).Error
}

func (s *MenuService) GetMenuList(req *dto.MenuListReq) ([]dto.MenuResult, error) {
	var menuSqlRes []model.SysMenu
	db := s.app.DB.Model(&model.SysMenu{})
	if req.Name != nil {
		db.Where("name LIKE ?", "%"+*req.Name+"%")
	}
	if req.Type != nil {
		db.Where("type = ?", *req.Type)
	}
	if req.Status != nil {
		db.Where("status = ?", *req.Status)
	}
	if req.ParentID != nil {
		db.Where("parent_id = ?", *req.ParentID)
	}
	if req.Perms != nil {
		db.Where("perms LIKE ?", "%"+*req.Perms+"%")
	}

	db.Find(&menuSqlRes)

	menuTree := s.BuildMenuTree(menuSqlRes)
	return menuTree, nil
}

func (s *MenuService) GetMenuListByPage(req *dto.MenuListByPageReq) (*dto.PageResult, error) {
	var total int64
	db := s.app.DB.Model(&model.SysMenu{})
	if err := db.Where("parent_id = ?", "0").Count(&total).Error; err != nil {
		return nil, err
	}
	if req.Name != nil {
		db.Where("name LIKE ?", "%"+*req.Name+"%")
	}
	if req.Type != nil {
		db.Where("type = ?", *req.Type)
	}
	if req.Status != nil {
		db.Where("status = ?", *req.Status)
	}
	if req.ParentID != nil {
		db.Where("parent_id = ?", *req.ParentID)
	}
	if req.Perms != nil {
		db.Where("perms LIKE ?", "%"+*req.Perms+"%")
	}
	var list []model.SysMenu

	if err := db.Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Find(&list).Error; err != nil {
		return nil, err
	}

	menuTree := s.BuildMenuTree(list)

	return &dto.PageResult{
		List:        menuTree,
		Total:       total,
		PageSize:    req.PageSize,
		CurrentPage: req.CurrentPage,
	}, nil
}

func (s *MenuService) MenuDetailById(id string) (*dto.MenuResult, error) {
	var menuResult dto.MenuResult

	if err := s.app.DB.Model(&model.SysMenu{}).Where("id = ?", id).Scan(&menuResult).Error; err != nil {
		return nil, err
	}

	return &menuResult, nil
}

// 处理树形菜单结构
func (s *MenuService) BuildMenuTree(allMenus []model.SysMenu) []dto.MenuResult {

	menuMap := make(map[string]*dto.MenuResult)

	for _, menu := range allMenus {
		menuMap[menu.ID] = &dto.MenuResult{
			ID:          menu.ID,
			ParentID:    menu.ParentID,
			Name:        menu.Name,
			Path:        menu.Path,
			Component:   menu.Component,
			Type:        menu.Type,
			Perms:       menu.Perms,
			Icon:        menu.Icon,
			Sort:        menu.Sort,
			Status:      menu.Status,
			Description: menu.Description,
			CreatedAt:   menu.CreatedAt,
			UpdatedAt:   menu.UpdatedAt,
			Children:    []*dto.MenuResult{}, // 初始化 Children
		}
	}

	var tree []dto.MenuResult

	for _, menu := range allMenus {
		// 从 Map 中获取当前节点的 DTO 指针
		currentNode := menuMap[menu.ID]

		// 判断是否为根节点（根据你的逻辑，ParentID 为 "0" 或空 是根节点）
		if menu.ParentID == "0" || menu.ParentID == "" {
			tree = append(tree, *currentNode)
		} else {
			// 如果不是根节点，尝试找到它的父节点
			if parentNode, ok := menuMap[menu.ParentID]; ok {
				// 将当前节点添加到父节点的 Children 中
				parentNode.Children = append(parentNode.Children, currentNode)
			}
		}
	}

	return tree
}
