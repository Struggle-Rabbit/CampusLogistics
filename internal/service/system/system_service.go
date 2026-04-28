package system

import (
	"errors"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/app"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"gorm.io/gorm"
)

type SystemService struct {
	app *app.App
}

func NewSystemService(app *app.App) *SystemService {
	return &SystemService{
		app: app,
	}
}

func (s *SystemService) RefreshToken(token string) (*dto.RefreshTokenResult, error) {
	info, err := utils.ParseToken(token)
	if err != nil {
		return nil, err
	}

	var sysUser model.SysUser

	err = s.app.DB.Model(&model.SysUser{}).Where("id = ?", info.UserID).First(&sysUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("未查询到用户信息")
		} else {
			return nil, err
		}
	}

	accessToken, refreshToken, err := utils.GenerateToken(sysUser.ID, sysUser.Name)

	if err != nil {
		return nil, err
	}

	return &dto.RefreshTokenResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil

}
