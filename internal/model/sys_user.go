package model

// SysUser 用户信息表
type SysUser struct {
	BaseModelWithDelete `gorm:"embedded"`
	UserCode            string    `gorm:"column:user_code;uniqueIndex;not null;size:50" json:"user_code"` // 学号/工号（唯一）
	Name                string    `gorm:"column:name;not null;size:50" json:"name"`                       // 姓名
	Mobile              string    `gorm:"column:mobile;uniqueIndex;not null;size:20" json:"mobile"`       // 手机号（唯一）
	Password            string    `gorm:"column:password;not null;size:255" json:"-"`                     // 密码（JSON隐藏）
	Status              int       `gorm:"column:status;not null;default:1;index" json:"status"`           // 状态：1-启用 2-禁用
	Avatar              string    `gorm:"column:avatar;size:255" json:"avatar"`                           // 头像URL
	UserType            string    `gorm:"column:user_type;not null;size:10" json:"user_type"`             // 用户类型: 00-管理员  01-职工  02-学生
	Roles               []SysRole `gorm:"many2many:sys_user_role;"`
	RefreshToken        string    `gorm:"column:refresh_token;uniqueIndex;not null;type:text" json:"refresh_token"`
}

// TableName 指定表名
func (SysUser) TableName() string {
	return "sys_user"
}
