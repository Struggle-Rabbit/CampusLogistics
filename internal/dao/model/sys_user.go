package model

// SysUser 用户信息表
type SysUser struct {
	BaseModelWithDelete
	UserID   string `gorm:"column:user_id;uniqueIndex;not null;size:50" json:"user_id"` // 学号/工号（唯一）
	Name     string `gorm:"column:name;not null;size:50" json:"name"`                   // 姓名
	Phone    string `gorm:"column:phone;uniqueIndex;not null;size:20" json:"phone"`     // 手机号（唯一）
	Password string `gorm:"column:password;not null;size:255" json:"-"`                 // 密码（JSON隐藏）
	RoleIDs  []uint `gorm:"column:role_ids;type:json;not null" json:"role_ids"`         // 角色ID列表（JSON序列化）
	Status   int    `gorm:"column:status;not null;default:1;index" json:"status"`       // 状态：1-启用 2-禁用
	Avatar   string `gorm:"column:avatar;size:255" json:"avatar"`                       // 头像URL
}

// TableName 指定表名
func (SysUser) TableName() string {
	return "sys_user"
}
