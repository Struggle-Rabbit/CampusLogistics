package unittest

import (
	"testing"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/menu"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/system"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/user"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/constant"
	"github.com/stretchr/testify/assert"
)

func TestSystemService(t *testing.T) {
	_, appInstance := SetupTestDB()
	svc := system.NewSystemService(appInstance)
	menuSvc := menu.NewMenuService(appInstance)
	userSvc := user.NewUserService(appInstance, menuSvc)

	// 先注册一个用户用于测试 Token 刷新
	userSvc.Register(&dto.RegisterReq{
		Name:     "系统用户",
		Mobile:   "13600136000",
		Password: "password123",
		UserType: constant.UserTypeAdmin,
	})

	loginRes, _ := userSvc.Login(&dto.LoginReq{
		Account:  "13600136000",
		Password: "password123",
	})

	t.Run("刷新 Token", func(t *testing.T) {
		res, err := svc.RefreshToken(loginRes.RefreshToken)
		assert.NoError(t, err)
		assert.NotEmpty(t, res.AccessToken)
		assert.NotEmpty(t, res.RefreshToken)
	})

	t.Run("分页查询操作日志", func(t *testing.T) {
		req := &dto.OperationLogByPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		}
		res, err := svc.GetOperationLogListByPage(req)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
}
