package main

import (
	"fmt"

	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("scripts/sql/dev.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	utils.InitSnowflake()

	// AutoMigrate
	db.AutoMigrate(
		&model.SysUser{}, &model.SysRole{}, &model.SysMenu{},
		&model.SysOperationLog{}, &model.RepairOrder{}, &model.RepairRecord{},
		&model.Campus{}, &model.Building{}, &model.DormRoom{}, &model.DormUser{},
		&model.DormUtility{}, &model.Notice{}, &model.UtilityPrice{},
	)

	// 创建测试用户
	hashed, _ := utils.HashedPasswordFunc("Test12345678")
	user := model.SysUser{
		Name:     "测试用户",
		UserCode: "20260606001",
		Mobile:   "13800138001",
		Password: hashed,
		Status:   1,
		UserType: "02",
	}
	if err := db.Create(&user).Error; err != nil {
		panic(err)
	}
	fmt.Printf("测试用户: ID=%s\n", user.ID)

	// 创建管理员角色
	role := model.SysRole{
		Status:      "01",
		RoleName:    "超级管理员",
		RoleCode:    "admin",
		Description: "拥有所有权限",
		IsBuiltIn:   1,
	}
	if err := db.Create(&role).Error; err != nil {
		panic(err)
	}
	fmt.Printf("管理员角色: ID=%s\n", role.ID)

	// 所有权限标识
	perms := []string{
		"sys:user:list", "sys:user:detail", "sys:user:del", "sys:user:update",
		"sys:optLog",
		"sys:role:list", "sys:role:detail", "sys:role:add", "sys:role:del", "sys:role:update",
		"sys:menu:list", "sys:menu:detail", "sys:menu:add", "sys:menu:del", "sys:menu:update",
		"repair:submit", "repair:list", "repair:detail", "repair:update", "repair:record", "repair:del",
		"campus:create", "campus:update", "campus:del", "campus:list", "campus:detail",
		"building:create", "building:update", "building:del", "building:list", "building:detail", "building:import", "building:export",
		"dorm:create", "dorm:update", "dorm:del", "dorm:list", "dorm:detail", "dorm:assign", "dorm:transfer", "dorm:checkout", "dorm:users", "dorm:warning",
		"utility:create", "utility:update", "utility:del", "utility:list", "utility:detail", "utility:pay", "utility:batchPay", "utility:price", "utility:statistics", "utility:warning",
		"notice:create", "notice:update", "notice:del", "notice:list", "notice:detail", "notice:top",
	}

	var menuIDs []string
	for i, p := range perms {
		m := model.SysMenu{
			ParentID: "0", Name: p, Type: 3, Perms: p, Status: 1, Sort: i,
		}
		if err := db.Create(&m).Error; err != nil {
			panic(fmt.Sprintf("创建菜单[%s]失败: %v", p, err))
		}
		menuIDs = append(menuIDs, m.ID)
	}
	fmt.Printf("菜单: %d 个\n", len(menuIDs))

	// 用户-角色关联 (使用 GORM 中间表实际列名)
	db.Exec("INSERT INTO sys_user_role (sys_role_id, sys_user_id) VALUES (?, ?)", role.ID, user.ID)
	fmt.Println("用户-角色关联完成")

	// 角色-菜单关联
	for _, mid := range menuIDs {
		db.Exec("INSERT INTO sys_role_menu (sys_role_id, sys_menu_id) VALUES (?, ?)", role.ID, mid)
	}
	fmt.Println("角色-菜单关联完成")

	fmt.Println("\n种子数据注入完成！")
}
