package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/app"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-viper/mapstructure/v2"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	app *app.App
}

// NewUserService 创建用户服务实例
func NewUserService(app *app.App) *UserService {
	return &UserService{
		app: app,
	}
}

// Register 用户注册
func (s *UserService) Register(req *dto.RegisterReq) error {
	// 检查是否已存在
	var existUser model.SysUser
	var total int64
	s.app.DB.Count(&total)
	err := s.app.DB.Where("mobile = ?", req.Mobile).First(&existUser).Error
	if err == nil {
		return errors.New("手机号已注册")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	hashedPassword, err := utils.HashedPasswordFunc(req.Password)
	if err != nil {
		return err
	}

	var user model.SysUser
	if err := mapstructure.Decode(req, &user); err != nil {
		return err
	}
	// 生成根据时间的自增工号
	user.Status = 1
	user.UserCode = fmt.Sprintf("%s00%d", time.Now().Format("20060102"), total+1)
	user.Password = hashedPassword

	// 创建用户（默认分配普通用户角色）
	// user := &model.SysUser{
	// 	UserCode: userCode,
	// 	Name:     req.Name,
	// 	Mobile:   req.Mobile,
	// 	Password: hashedPassword,
	// 	Status:   1, // 1-启用
	// }
	return s.app.DB.Create(user).Error
}

func (s *UserService) Login(req *dto.LoginReq) (*dto.LoginResult, error) {
	var sysUser model.SysUser
	err := s.app.DB.Where("mobile = ? OR user_code = ?", req.Account, req.Account).First(&sysUser).Error
	if err == nil {
		// 密码校验
		if err := utils.VerifyPasswordFunc(sysUser.Password, req.Password); err != nil {
			return nil, errors.New("账号密码不正确")
		}
	} else {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		} else {
			return nil, errors.New("账号密码不正确")
		}
	}

	accessToken, refreshToken, err := utils.GenerateToken(sysUser.ID, sysUser.Name)

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
		err := s.app.DB.First(&sysUser, userId).Error
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
	if err := s.app.DB.Model(&model.SysUser{}).Count(&total).Error; err != nil {
		return nil, err
	}

	if err := s.app.DB.Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Find(&list).Error; err != nil {
		return nil, err
	}

	return &dto.PageResult{
		List:        list,
		Total:       total,
		PageSize:    req.PageSize,
		CurrentPage: req.CurrentPage,
	}, nil
}

func (s *UserService) UpdateUser(req *dto.UserUpdateReq) error {
	var user model.SysUser
	if err := mapstructure.Decode(req, &user); err != nil {
		return err
	}
	return s.app.DB.Save(&user).Error
}

func (s *UserService) DelUser(id string) error {

	return s.app.DB.Delete(&model.SysUser{}, id).Error
}

func (s *UserService) GetUserPermission(user_id string) (*dto.UserPermissionResult, error) {
	// 根据userID获取对应的角色ID
	var roleIDs []string
	if err := s.app.DB.Model(&model.SysUserRole{}).Where("user_id = ?", user_id).Pluck("role_id", &roleIDs).Error; err != nil {
		return nil, err
	}
	// 根据角色ID获取角色详情
	var userRole []dto.RoleResult
	if err := s.app.DB.Model(&model.SysRole{}).Where("idIN ?", roleIDs).Find(&userRole).Error; err != nil {
		return nil, err
	}

	// 根据角色ID查菜单ID（去重）
	var menuIDs []string
	if err := s.app.DB.Model(&model.SysRoleMenu{}).
		Where("role_id IN ?", roleIDs).
		Distinct(). // 去重：多个角色可能有相同菜单
		Pluck("menu_id", &menuIDs).Error; err != nil {
		return nil, err
	}

	// 根据菜单ID查菜单详情
	var menus []dto.MenuResult
	if len(menuIDs) > 0 {
		if err := s.app.DB.Model(&model.SysMenu{}).Where("id IN ?", menuIDs).Find(&menus).Error; err != nil {
			return nil, err
		}
	}

	return &dto.UserPermissionResult{
		UserId:  user_id,
		RoleIDs: roleIDs,
		Roles:   userRole,
		MenuIDs: menuIDs,
		Menus:   menus,
	}, nil
}
