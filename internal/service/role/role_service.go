package role

import (
	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewRoleService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateRole(req *dto.CreateRoleReq) error {
	return s.db.Create(&model.SysRole{
		RoleName:    req.RoleName,
		RoleCode:    req.RoleCode,
		Status:      req.Status,
		IsBuiltIn:   2,
		Description: req.Description,
	}).Error
}

func (s *Service) UpdateRole(req *dto.UpdateRoleReq) error {
	return s.db.Model(&model.SysRole{}).Where("id = ?", req.ID).Updates(map[string]interface{}{
		"role_name":   req.RoleName,
		"role_code":   req.RoleCode,
		"status":      req.Status,
		"description": req.Description,
	}).Error
}

func (s *Service) DelRole(id []string) error {
	return s.db.Delete(&model.SysRole{}, id).Error
}

func (s *Service) GetRoleList(name string) ([]dto.RoleResult, error) {
	var roleSqlRes []dto.RoleResult

	s.db.Model(&model.SysRole{}).Where("role_name LIKE ?", "%"+name+"%").Scan(&roleSqlRes)

	return roleSqlRes, nil
}

func (s *Service) GetRoleListByPage(req *dto.RoleListByPageReq) (*dto.PageResult, error) {
	var total int64
	db := s.db.Model(&model.SysRole{})
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

func (s *Service) RoleDetailById(id string) (*dto.RoleResult, error) {
	var roleResult dto.RoleResult

	if err := s.db.Model(&model.SysRole{}).Where("id = ?", id).First(&roleResult).Error; err != nil {
		return nil, err
	}

	return &roleResult, nil
}
