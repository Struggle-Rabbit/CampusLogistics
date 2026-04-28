package unittest

import (
	"testing"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/role"
	"github.com/stretchr/testify/assert"
)

func TestRoleService(t *testing.T) {
	_, appInstance := SetupTestDB()
	svc := role.NewRoleService(appInstance)

	var roleID string

	t.Run("创建角色", func(t *testing.T) {
		req := &dto.CreateRoleReq{
			RoleName:    "管理员",
			RoleCode:    "admin",
			Status:      "01",
			Description: "超级管理员",
		}
		err := svc.CreateRole(req)
		assert.NoError(t, err)

		// 获取列表确认
		list, err := svc.GetRoleList("管理员")
		assert.NoError(t, err)
		assert.Len(t, list, 1)
		roleID = list[0].ID
	})

	t.Run("更新角色", func(t *testing.T) {
		req := &dto.UpdateRoleReq{
			ID:          roleID,
			RoleName:    "超管",
			RoleCode:    "super_admin",
			Status:      "01",
			Description: "顶级管理员",
		}
		err := svc.UpdateRole(req)
		assert.NoError(t, err)

		// 验证详情
		detail, err := svc.RoleDetailById(roleID)
		assert.NoError(t, err)
		assert.Equal(t, "超管", detail.RoleName)
		assert.Equal(t, "super_admin", detail.RoleCode)
	})

	t.Run("分页查询", func(t *testing.T) {
		req := &dto.RoleListByPageReq{
			PageReq: dto.PageReq{CurrentPage: 1, PageSize: 10},
		}
		res, err := svc.GetRoleListByPage(req)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), res.Total)
	})

	t.Run("删除角色", func(t *testing.T) {
		err := svc.DelRole([]string{roleID})
		assert.NoError(t, err)

		list, _ := svc.GetRoleList("")
		assert.Len(t, list, 0)
	})
}
