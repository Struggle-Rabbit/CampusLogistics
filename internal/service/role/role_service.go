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
	return s.app.DB.Create(&role).Error
}

func (s *RoleService) UpdateRole(req *dto.UpdateRoleReq) error {
	var role model.SysRole
	if err := mapstructure.Decode(req, &role); err != nil {
		return err
	}
	return s.app.DB.Save(&role).Error
}

func (s *RoleService) DelRole(id string) error {

	return s.app.DB.Delete(&model.SysRole{}, id).Error
}

func (s *RoleService) GetRoleList(req *dto.RoleListReq) ([]*dto.RoleResult, error) {
	var roleReq model.SysRole
	if err := mapstructure.Decode(req, &roleReq); err != nil {
		return nil, err
	}
	var roleSqlRes []*dto.RoleResult

	s.app.DB.Where(&roleReq).Find(&roleSqlRes)

	return roleSqlRes, nil
}

func (s *RoleService) GetRoleListByPage(req *dto.RoleListByPageReq) (*dto.PageResult, error) {
	var roleReq model.SysRole
	var total int64
	if err := s.app.DB.Model(&model.SysRole{}).Count(&total).Error; err != nil {
		return nil, err
	}
	if err := mapstructure.Decode(req, &roleReq); err != nil {
		return nil, err
	}
	var list []*dto.RoleResult

	if err := s.app.DB.Model(&model.SysRole{}).Where(&roleReq).Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Find(&list).Error; err != nil {
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
