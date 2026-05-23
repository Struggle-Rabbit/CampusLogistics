package system

import (
	"errors"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/app"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"gorm.io/gorm"
)

type SystemService interface {
	RefreshToken(token string) (*dto.RefreshTokenResult, error)
	GetOperationLogListByPage(req *dto.OperationLogByPageReq) (*dto.PageResult, error)
}

type SystemServiceProvider struct {
	App *app.App
}

func (s *SystemServiceProvider) RefreshToken(token string) (*dto.RefreshTokenResult, error) {
	info, err := utils.ParseToken(token)
	if err != nil {
		return nil, err
	}

	var sysUser model.SysUser

	if err := s.App.DB.Model(&model.SysUser{}).Where("id = ?", info.UserID).First(&sysUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("未查询到用户信息")
		} else {
			return nil, err
		}
	}

	if sysUser.RefreshToken != token {
		return nil, errors.New("RefreshToken无效")
	}

	accessToken, refreshToken, err := utils.GenerateToken(sysUser.ID, sysUser.Name)
	if err != nil {
		return nil, err
	}

	// 使用主键更新 refresh_token，避免 SQLite UPDATE...FROM 自引用歧义
	if err := s.App.DB.Model(&sysUser).Update("refresh_token", refreshToken).Error; err != nil {
		return nil, err
	}

	return &dto.RefreshTokenResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
