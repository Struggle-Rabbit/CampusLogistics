package model

// SysRoleMenu 角色菜单关联表
type SysRoleMenu struct {
	BaseModel
	RoleID uint // 角色ID
	MenuID uint `gorm:"column:menu_id;not null;index:idx_role_menu" json:"menu_id"` // 菜单ID
}

func (SysRoleMenu) TableName() string {
	return "sys_role_menu"
}
