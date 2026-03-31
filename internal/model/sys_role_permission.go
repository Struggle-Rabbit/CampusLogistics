package model

// SysRolePermission 角色权限关联表
type SysRolePermission struct {
	BaseModel
	RoleID       uint // 角色ID
	PermissionID uint `gorm:"column:permission_id;not null;index:idx_role_perm" json:"permission_id"` // 权限ID
}

func (SysRolePermission) TableName() string {
	return "sys_role_permission"
}
