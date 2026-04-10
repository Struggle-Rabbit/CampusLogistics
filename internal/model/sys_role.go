package model

// SysRole 角色表
type SysRole struct {
	BaseModel   `gorm:"embedded"`
	RoleName    string `gorm:"column:role_name;not null;size:50" json:"role_name"`             // 角色名称
	RoleCode    string `gorm:"column:role_code;uniqueIndex;not null;size:50" json:"role_code"` // 角色编码（唯一）
	Description string `gorm:"column:description;size:255" json:"description"`                 // 角色描述
	IsBuiltIn   int    `gorm:"column:is_built_in;not null;default:2" json:"is_built_in"`       // 是否内置：1-是 2-否（内置角色不可删除）
}

func (SysRole) TableName() string {
	return "sys_role"
}
