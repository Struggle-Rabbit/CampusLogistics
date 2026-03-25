package model

// SysOperationLog 操作日志表
type SysOperationLog struct {
	BaseModel
	UserID    string `gorm:"column:user_id;size:50;index" json:"user_id"`         // 操作人学号/工号（未登录为空）
	Operation string `gorm:"column:operation;not null;size:100" json:"operation"` // 操作描述
	Path      string `gorm:"column:path;size:255" json:"path"`                    // 接口路径
	Method    string `gorm:"column:method;size:10" json:"method"`                 // 请求方法
	IP        string `gorm:"column:ip;size:50" json:"ip"`                         // 操作IP
	Params    string `gorm:"column:params;type:text" json:"params"`               // 请求参数（JSON字符串）
	Result    string `gorm:"column:result;type:text" json:"result"`               // 响应结果（JSON字符串）
}

func (SysOperationLog) TableName() string {
	return "sys_operation_log"
}
