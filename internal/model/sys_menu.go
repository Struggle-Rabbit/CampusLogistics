package model

// SysMenu 菜单表
type SysMenu struct {
	BaseModel   `gorm:"embedded"`
	ParentID    string
	Name        string `gorm:"size:32;not null"`   // 菜单名称
	Path        string `gorm:"size:128"`           // 路由
	Component   string `gorm:"size:128"`           // 前端组件
	Type        int    `gorm:"not null;default:1"` // 1目录2菜单3按钮
	Perms       string `gorm:"size:64;unique"`     // 权限标识 sys:user:list
	Icon        string `gorm:"size:64"`
	Sort        int    `gorm:"default:0"`
	Status      int    `gorm:"default:1"`                                      // 状态
	Description string `gorm:"column:description;size:255" json:"description"` // 描述
}

func (SysMenu) TableName() string {
	return "sys_menu"
}
