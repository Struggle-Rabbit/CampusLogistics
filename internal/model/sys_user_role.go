package model

// SysUserRole 用户角色关联表
type SysUserRole struct {
	BaseModel
	UserID string // 用户ID
	RoleID string `gorm:"column:role_id;not null;index:idx_user_role" json:"role_id"` // 角色ID
}

func (SysUserRole) TableName() string {
	return "sys_user_role"
}
