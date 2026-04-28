package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/app"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/menu"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/constant"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	app  *app.App
	menu *menu.MenuService
}

// NewUserService 创建用户服务实例
func NewUserService(app *app.App, menuSvc *menu.MenuService) *UserService {
	return &UserService{
		app:  app,
		menu: menuSvc,
	}
}

// Register 用户注册
func (s *UserService) Register(req *dto.RegisterReq) error {
	return s.app.DB.Transaction(func(tx *gorm.DB) error {
		// 检查是否已存在
		var total int64
		if err := tx.Model(&model.SysUser{}).Where("mobile = ?", req.Mobile).Count(&total).Error; err != nil {
			return err
		}
		if total > 0 {
			return errors.New("手机号已注册")
		}

		hashedPassword, err := utils.HashedPasswordFunc(req.Password)
		if err != nil {
			return err
		}

		// 获取当前总数用于生成工号
		var count int64
		tx.Model(&model.SysUser{}).Count(&count)

		// 生成根据时间的自增工号
		return tx.Create(&model.SysUser{
			Name:     req.Name,
			Mobile:   req.Mobile,
			Password: hashedPassword,
			UserCode: fmt.Sprintf("%s00%d", time.Now().Format("20060102"), count+1),
			Status:   constant.UserStatusEnable,
			UserType: req.UserType,
		}).Error
	})
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
	var list []*model.SysUser
	var total int64
	db := s.app.DB.Model(&model.SysUser{})

	if req.Mobile != "" {
		db.Where("mobile = ?", req.Mobile)
	}

	if req.Name != "" {
		db.Where("name = ?", req.Name)
	}

	if req.UserType != "" {
		db.Where("user_type = ?", req.UserType)
	}

	if req.Status != "" {
		db.Where("status = ?", req.Status)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	if err := db.Scopes(dao.Paginate(req.CurrentPage, req.PageSize)).Preload("Roles").Find(&list).Error; err != nil {
		return nil, err
	}

	var dtoList []*dto.UserInfoResult

	for _, v := range list {
		var roleResults []*dto.RoleResult
		for _, r := range v.Roles {
			roleResults = append(roleResults, &dto.RoleResult{
				ID:       r.ID,
				RoleName: r.RoleName,
				RoleCode: r.RoleCode,
				Status:   r.Status,
			})
		}

		dtoList = append(dtoList, &dto.UserInfoResult{
			ID:        v.ID,
			Name:      v.Name,
			UserCode:  v.UserCode,
			UserType:  v.UserType,
			Avatar:    v.Avatar,
			Mobile:    v.Mobile,
			Status:    v.Status,
			Roles:     roleResults,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
		})
	}

	return &dto.PageResult{
		List:        dtoList,
		Total:       total,
		PageSize:    req.PageSize,
		CurrentPage: req.CurrentPage,
	}, nil
}

func (s *UserService) UpdateUser(req *dto.UserUpdateReq) error {
	if req.ID == "" {
		return errors.New("id不能为空")
	}
	if req.Mobile != "" {
		var user model.SysUser
		if err := s.app.DB.Where("mobile = ? AND id != ?", req.Mobile, req.ID).First(&user).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		} else {
			return errors.New("此手机号已被使用")
		}
	}
	return s.app.DB.Model(&model.SysUser{}).Where("id = ?", req.ID).Updates(model.SysUser{
		Name:     req.Name,
		Mobile:   req.Mobile,
		Status:   req.Status,
		Avatar:   req.Avatar,
		UserType: req.UserType,
	}).Error
}

func (s *UserService) DelUser(id []string) error {

	return s.app.DB.Delete(&model.SysUser{}, id).Error
}

// ResetPassword 重置密码
func (s *UserService) ResetPassword(req *dto.PasswordReset) error {
	var user model.SysUser
	// 查找用户
	err := s.app.DB.Where("mobile = ?", req.Mobile).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return err
	}

	// 校验原密码
	if err := utils.VerifyPasswordFunc(user.Password, req.OldPassword); err != nil {
		return errors.New("原密码错误")
	}

	// 加密新密码
	hashedPassword, err := utils.HashedPasswordFunc(req.NewPassword)
	if err != nil {
		return err
	}

	// 更新密码
	return s.app.DB.Model(&user).Update("password", hashedPassword).Error
}

func (s *UserService) GetUserPermission(user_id string) (*dto.UserPermissionResult, error) {
	var sysUser model.SysUser
	err := s.app.DB.Preload("Roles").Preload("Roles.Menus").First(&sysUser, "id = ?", user_id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	var roleRes []dto.RoleResult
	var menuList []model.SysMenu

	var roleIds []string
	var menuIds []string
	menuMap := make(map[string]model.SysMenu)
	for _, role := range sysUser.Roles {
		roleIds = append(roleIds, role.ID)

		roleRes = append(roleRes, dto.RoleResult{
			ID:          role.ID,
			RoleName:    role.RoleName,
			RoleCode:    role.RoleCode,
			Status:      role.Status,
			IsBuiltIn:   role.IsBuiltIn,
			Description: role.Description,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
		})

		for _, menu := range role.Menus {
			if _, exists := menuMap[menu.ID]; !exists {
				menuMap[menu.ID] = menu
			}
		}
	}

	for _, menu := range menuMap {
		menuIds = append(menuIds, menu.ID)
		menuList = append(menuList, menu)
	}

	return &dto.UserPermissionResult{
		UserId:  sysUser.ID,
		RoleIDs: roleIds,
		Roles:   roleRes,
		MenuIDs: menuIds,
		Menus:   s.menu.BuildMenuTree(menuList),
	}, nil
}
