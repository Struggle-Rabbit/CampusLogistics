package user

import (
	"errors"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	DB *gorm.DB
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		DB: dao.DB,
	}
}

// Register 用户注册
func (s *UserService) Register(req *dto.RegisterReq) error {
	// 1. 检查学号/工号是否已存在
	var existUser model.SysUser
	err := s.DB.Where("user_id = ? OR phone = ?", req.UserID, req.Phone).First(&existUser).Error
	if err == nil {
		return errors.New("学号/工号或手机号已注册")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 2. 密码 bcrypt 加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 3. 创建用户（默认分配普通用户角色）
	user := &model.SysUser{
		UserID:   req.UserID,
		Name:     req.Name,
		Phone:    req.Phone,
		Password: string(hashedPassword),
		RoleIDs:  []uint{2}, // 假设 2 是普通用户角色ID
		Status:   1,         // 1-启用
	}
	return dao.DB.Create(user).Error
}
