package role

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/app"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"github.com/go-viper/mapstructure/v2"
)

type RoleService struct {
	app *app.App
}

func NewRoleService(app *app.App) *RoleService {
	return &RoleService{
		app: app,
	}
}

func (s *RoleService) CreateRole(req *dto.CreateRoleReq) error {
	var role model.SysRole
	if err := mapstructure.Decode(req, &role); err != nil {
		return err
	}
	return s.app.DB.Create(&model.SysRole{
		RoleName:    req.RoleName,
		RoleCode:    req.RoleCode,
		Status:      req.Status,
		IsBuiltIn:   2,
		Description: req.Description,
	}).Error
}

func (s *RoleService) UpdateRole(req *dto.UpdateRoleReq) error {
	return s.app.DB.Model(&model.SysRole{}).Where("id = ?", req.ID).Updates(model.SysRole{
		RoleName: req.RoleName,
		RoleCode: req.RoleCode,
		Status:   req.Status,
	}).Error
}

func (s *RoleService) DelRole(id []string) error {

	return s.app.DB.Delete(&model.SysRole{}, id).Error
}

func (s *RoleService) GetRoleList(name string) ([]dto.RoleResult, error) {
	var roleSqlRes []dto.RoleResult

	s.app.DB.Model(&model.SysRole{}).Where("role_name LIKE ?", "%"+name+"%").Scan(&roleSqlRes)

	return roleSqlRes, nil
}

func (s *RoleService) GetRoleListByPage(req *dto.RoleListByPageReq) (*dto.PageResult, error) {
	var total int64
	db := s.app.DB.Model(&model.SysRole{})
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	if req.RoleName != "" {
		db.Where("role_name = ?", req.RoleName)
	}
	if req.Status != "" {
		db.Where("status = ?", req.Status)
	}
	var list []*dto.RoleResult

	if err := db.Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Scan(&list).Error; err != nil {
		return nil, err
	}

	return &dto.PageResult{
		List:        list,
		Total:       total,
		PageSize:    req.PageSize,
		CurrentPage: req.CurrentPage,
	}, nil
}

func (s *RoleService) RoleDetailById(id string) (*dto.RoleResult, error) {
	var roleResult dto.RoleResult

	if err := s.app.DB.Model(&model.SysRole{}).Where("id = ?", id).First(&roleResult).Error; err != nil {
		return nil, err
	}

	return &roleResult, nil
}
