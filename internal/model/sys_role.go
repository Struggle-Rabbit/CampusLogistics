package model

// SysRole 角色表
type SysRole struct {
	BaseModel   `gorm:"embedded"`
	Status      string    `gorm:"column:status;not null;size:32;default:'01'" json:"status"`      // 角色状态：00-停用  01-启用
	RoleName    string    `gorm:"column:role_name;not null;size:50" json:"role_name"`             // 角色名称
	RoleCode    string    `gorm:"column:role_code;uniqueIndex;not null;size:50" json:"role_code"` // 角色编码（唯一）
	Description string    `gorm:"column:description;size:255" json:"description"`                 // 角色描述
	IsBuiltIn   int       `gorm:"column:is_built_in;not null;default:2" json:"is_built_in"`       // 是否内置：1-是 2-否（内置角色不可删除）
	Users       []SysUser `gorm:"many2many:sys_user_role;"`
	Menus       []SysMenu `gorm:"many2many:sys_role_menu;"`
}

func (SysRole) TableName() string {
	return "sys_role"
}
