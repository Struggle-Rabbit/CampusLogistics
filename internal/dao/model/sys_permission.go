package model

// SysPermission 权限表
type SysPermission struct {
	BaseModel
	PermissionName string `gorm:"column:permission_name;not null;size:100" json:"permission_name"`             // 权限名称
	PermissionCode string `gorm:"column:permission_code;uniqueIndex;not null;size:100" json:"permission_code"` // 权限编码（唯一）
	Path           string `gorm:"column:path;size:255" json:"path"`                                            // 接口路径
	Method         string `gorm:"column:method;size:10" json:"method"`                                         // 请求方法（GET/POST/PUT/DELETE）
	Description    string `gorm:"column:description;size:255" json:"description"`                              // 权限描述
}

func (SysPermission) TableName() string {
	return "sys_permission"
}
