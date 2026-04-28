package unittest

import (
	"testing"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/menu"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/user"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/constant"
	"github.com/stretchr/testify/assert"
)

func TestUserService_Register(t *testing.T) {
	_, appInstance := SetupTestDB()
	menuSvc := menu.NewMenuService(appInstance)
	svc := user.NewUserService(appInstance, menuSvc)

	t.Run("成功注册", func(t *testing.T) {
		req := &dto.RegisterReq{
			Name:     "测试用户",
			Mobile:   "13800138000",
			Password: "password123",
			UserType: constant.UserTypeStudent,
		}
		err := svc.Register(req)
		assert.NoError(t, err)

		// 重复注册
		err = svc.Register(req)
		assert.Error(t, err)
		assert.Equal(t, "手机号已注册", err.Error())
	})
}

func TestUserService_Login(t *testing.T) {
	_, appInstance := SetupTestDB()
	menuSvc := menu.NewMenuService(appInstance)
	svc := user.NewUserService(appInstance, menuSvc)

	// 先注册
	regReq := &dto.RegisterReq{
		Name:     "登录用户",
		Mobile:   "13900139000",
		Password: "password123",
		UserType: constant.UserTypeStudent,
	}
	svc.Register(regReq)

	t.Run("登录成功", func(t *testing.T) {
		loginReq := &dto.LoginReq{
			Account:  "13900139000",
			Password: "password123",
		}
		res, err := svc.Login(loginReq)
		assert.NoError(t, err)
		assert.NotEmpty(t, res.AccessToken)
		assert.NotEmpty(t, res.RefreshToken)
	})

	t.Run("登录失败-密码错误", func(t *testing.T) {
		loginReq := &dto.LoginReq{
			Account:  "13900139000",
			Password: "wrongpassword",
		}
		res, err := svc.Login(loginReq)
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, "账号密码不正确", err.Error())
	})
}

func TestUserService_ResetPassword(t *testing.T) {
	_, appInstance := SetupTestDB()
	menuSvc := menu.NewMenuService(appInstance)
	svc := user.NewUserService(appInstance, menuSvc)

	// 先注册
	mobile := "13700137000"
	svc.Register(&dto.RegisterReq{
		Name:     "重置用户",
		Mobile:   mobile,
		Password: "oldpassword",
		UserType: constant.UserTypeStudent,
	})

	t.Run("成功重置", func(t *testing.T) {
		req := &dto.PasswordReset{
			Mobile:      mobile,
			OldPassword: "oldpassword",
			NewPassword: "newpassword123",
		}
		err := svc.ResetPassword(req)
		assert.NoError(t, err)

		// 验证新密码登录
		loginRes, err := svc.Login(&dto.LoginReq{
			Account:  mobile,
			Password: "newpassword123",
		})
		assert.NoError(t, err)
		assert.NotNil(t, loginRes)
	})
}
