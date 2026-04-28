package unittest

import (
	"testing"

	"github.com/Struggle-Rabbit/CampusLogistics/api/dto"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/menu"
	"github.com/stretchr/testify/assert"
)

func TestMenuService(t *testing.T) {
	_, appInstance := SetupTestDB()
	svc := menu.NewMenuService(appInstance)

	// 1. 创建菜单
	req1 := &dto.CreateMenuReq{
		ParentID: "0",
		Name:     "系统管理",
		Type:     1,
		Perms:    "sys:mgr",
		Status:   1,
	}
	err := svc.CreateMenu(req1)
	assert.NoError(t, err)

	list1, err := svc.GetMenuList(&dto.MenuListReq{})
	assert.NoError(t, err)
	assert.Len(t, list1, 1)
	menuID := list1[0].ID

	// 2. 创建子菜单
	req2 := &dto.CreateMenuReq{
		ParentID:  menuID,
		Name:      "用户管理",
		Type:      2,
		Path:      "/user",
		Component: "User",
		Perms:     "sys:user:list",
		Status:    1,
	}
	err = svc.CreateMenu(req2)
	assert.NoError(t, err)

	tree, err := svc.GetMenuList(&dto.MenuListReq{})
	assert.NoError(t, err)
	assert.Len(t, tree, 1)
	assert.Len(t, tree[0].Children, 1)
	assert.Equal(t, "用户管理", tree[0].Children[0].Name)

	// 3. 更新菜单
	req3 := &dto.UpdateMenuReq{
		ID:       menuID,
		Name:     "基础管理",
		ParentID: "0",
		Type:     1,
		Status:   1,
	}
	err = svc.UpdateMenu(req3)
	assert.NoError(t, err)

	detail, err := svc.MenuDetailById(menuID)
	assert.NoError(t, err)
	assert.Equal(t, "基础管理", detail.Name)

	// 4. 删除菜单
	err = svc.DelMenu([]string{menuID})
	assert.NoError(t, err)

	list4, _ := svc.GetMenuList(&dto.MenuListReq{})
	assert.Len(t, list4, 0)
}
