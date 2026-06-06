package menu

import (
	"errors"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewMenuService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateMenu(req *dto.CreateMenuReq) error {
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

	if req.Type == 2 {
		if req.Path == "" {
			return errors.New("路由地址不能为空！")
		}
		if req.Component == "" {
			return errors.New("组件地址不能为空")
		}
		menu.Path = req.Path
		menu.Component = req.Component
	}

	return s.db.Create(menu).Error
}

func (s *Service) UpdateMenu(req *dto.UpdateMenuReq) error {
	if req.Type == 2 {
		if req.Path == "" {
			return errors.New("菜单路由地址不能为空")
		}
		if req.Component == "" {
			return errors.New("组件地址不能为空")
		}
	}
	updateData := map[string]interface{}{
		"parent_id":   req.ParentID,
		"name":        req.Name,
		"type":        req.Type,
		"perms":       req.Perms,
		"status":      req.Status,
		"sort":        req.Sort,
		"icon":        req.Icon,
		"description": req.Description,
		"path":        req.Path,
		"component":   req.Component,
	}

	return s.db.Model(&model.SysMenu{}).Where("id = ?", req.ID).Updates(updateData).Error
}

func (s *Service) DelMenu(id []string) error {
	return s.db.Delete(&model.SysMenu{}, id).Error
}

func (s *Service) GetMenuList(req *dto.MenuListReq) ([]dto.MenuResult, error) {
	var menuSqlRes []model.SysMenu
	tx := s.db.Model(&model.SysMenu{})
	if req.Name != nil && *req.Name != "" {
		tx = tx.Where("name LIKE ?", "%"+*req.Name+"%")
	}
	if req.Type != nil {
		tx = tx.Where("type = ?", *req.Type)
	}
	if req.Status != nil {
		tx = tx.Where("status = ?", *req.Status)
	}
	if req.ParentID != nil {
		tx = tx.Where("parent_id = ?", *req.ParentID)
	}
	if req.Perms != nil && *req.Perms != "" {
		tx = tx.Where("perms LIKE ?", "%"+*req.Perms+"%")
	}

	if err := tx.Find(&menuSqlRes).Error; err != nil {
		return nil, err
	}

	return s.BuildMenuTree(menuSqlRes), nil
}

func (s *Service) GetMenuListByPage(req *dto.MenuListByPageReq) (*dto.PageResult, error) {
	var list []model.SysMenu
	var total int64
	db := s.db.Model(&model.SysMenu{})
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	if req.Name != nil && *req.Name != "" {
		db = db.Where("name LIKE ?", "%"+*req.Name+"%")
	}
	if req.Type != nil {
		db = db.Where("type = ?", *req.Type)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.ParentID != nil {
		db = db.Where("parent_id = ?", *req.ParentID)
	}
	err := db.Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Find(&list).Error
	if err != nil {
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

func (s *Service) MenuDetailById(id string) (*dto.MenuResult, error) {
	var m model.SysMenu
	err := s.db.Where("id = ?", id).First(&m).Error
	if err != nil {
		return nil, err
	}

	return &dto.MenuResult{
		ID:          m.ID,
		ParentID:    m.ParentID,
		Name:        m.Name,
		Path:        m.Path,
		Component:   m.Component,
		Type:        m.Type,
		Perms:       m.Perms,
		Icon:        m.Icon,
		Sort:        m.Sort,
		Status:      m.Status,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}, nil
}

// BuildMenuTree 处理树形菜单结构（接收者方法，内部调用私有实现）
func (s *Service) BuildMenuTree(allMenus []model.SysMenu) []dto.MenuResult {
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
			Children:    []*dto.MenuResult{},
		}
	}

	var tree []dto.MenuResult

	for _, menu := range allMenus {
		currentNode := menuMap[menu.ID]
		if menu.ParentID == "0" || menu.ParentID == "" {
			// skip, collect root nodes later
		} else {
			if parentNode, ok := menuMap[menu.ParentID]; ok {
				parentNode.Children = append(parentNode.Children, currentNode)
			}
		}
	}

	for _, menu := range allMenus {
		if menu.ParentID == "0" || menu.ParentID == "" {
			tree = append(tree, *menuMap[menu.ID])
		}
	}

	return tree
}
