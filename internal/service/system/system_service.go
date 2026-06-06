package system

import (
	"errors"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewSystemService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) RefreshToken(token string) (*dto.RefreshTokenResult, error) {
	info, err := utils.ParseToken(token)
	if err != nil {
		return nil, err
	}

	var sysUser model.SysUser

	if err := s.db.Model(&model.SysUser{}).Where("id = ?", info.UserID).First(&sysUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("未查询到用户信息")
		}
		return nil, err
	}

	if sysUser.RefreshToken != token {
		return nil, errors.New("RefreshToken无效")
	}

	accessToken, refreshToken, err := utils.GenerateToken(sysUser.ID, sysUser.Name)
	if err != nil {
		return nil, err
	}

	if err := s.db.Model(&sysUser).Update("refresh_token", refreshToken).Error; err != nil {
		return nil, err
	}

	return &dto.RefreshTokenResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) GetOperationLogListByPage(req *dto.OperationLogByPageReq) (*dto.PageResult, error) {
	var total int64
	db := s.db.Model(&model.SysOperationLog{})

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	if req.IP != "" {
		db.Where("ip = ?", req.IP)
	}

	if req.UserID != "" {
		db.Where("user_id = ?", req.UserID)
	}

	if !req.OperationTimeStart.IsZero() && !req.OperationTimeEnd.IsZero() {
		db.Where("operation_at >= ? AND operation_at <= ?", req.OperationTimeStart, req.OperationTimeEnd)
	}
	var list []model.SysOperationLog

	if err := db.Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Find(&list).Error; err != nil {
		return nil, err
	}

	return &dto.PageResult{
		List:        list,
		Total:       total,
		PageSize:    req.PageSize,
		CurrentPage: req.CurrentPage,
	}, nil
}
