package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务实例
func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		db,
	}
}

// Register 用户注册
func (s *UserService) Register(req *dto.RegisterReq) error {
	// 检查是否已存在
	var existUser model.SysUser
	err := s.db.Where("mobile = ?", req.Mobile).First(&existUser).Error
	if err == nil {
		return errors.New("手机号已注册")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 密码 bcrypt 加密
	hashedPassword, err := utils.HashedPasswordFunc(req.Password)
	if err != nil {
		return err
	}
	// 获取当前数据库中数据数量
	var userList []model.SysUser
	s.db.Find(&userList)

	// 生成根据时间的自增工号
	userCode := fmt.Sprintf("%s00%d", time.Now().Format("20060102"), len(userList)+1)
	// 创建用户（默认分配普通用户角色）
	user := &model.SysUser{
		UserCode: userCode,
		Name:     req.Name,
		Mobile:   req.Mobile,
		Password: hashedPassword,
		Status:   1, // 1-启用
	}
	return s.db.Create(user).Error
}

func (s *UserService) Login(req *dto.LoginReq) (*dto.LoginResult, error) {
	var sysUser model.SysUser
	err := s.db.Where("mobile = ? OR user_code = ?", req.Account).First(&sysUser).Error
	if err == nil {
		// 密码 bcrypt 加密
		hashedPassword, pwd_err := utils.HashedPasswordFunc(req.Password)
		if pwd_err != nil {
			return nil, pwd_err
		}
		if sysUser.Password != hashedPassword {
			return nil, errors.New("密码不正确")
		}
	} else {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		} else {
			return nil, errors.New("账号输入不正确")
		}
	}

	accessToken, refreshToken, err := utils.GenerateToken(sysUser.ID)

	if err != nil {
		return nil, err
	}

	return &dto.LoginResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserService) GetUserInfo(c *gin.Context) (*dto.UserInfoResult, error) {
	userId, exists := c.Get("userID")
	var sysUser dto.UserInfoResult
	if exists {
		err := s.db.First(&sysUser, userId).Error
		if err == nil {
			return &sysUser, nil
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("未查询到用户信息! ")
		} else {
			return nil, err
		}
	}

	return nil, errors.New("未获取到UserID")
}

func (s *UserService) GetListByPage(req *dto.UserListPageReq) (*dto.PageResult, error) {
	var list []model.SysUser
	var total int64
	if err := s.db.Model(&model.SysUser{}).Count(&total).Error; err != nil {
		return nil, err
	}

	if err := s.db.Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Find(&list).Error; err != nil {
		return nil, err
	}

	return &dto.PageResult{
		List:        list,
		Total:       total,
		PageSize:    req.PageSize,
		CurrentPage: req.CurrentPage,
	}, nil
}
