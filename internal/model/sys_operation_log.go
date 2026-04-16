package model

import "time"

// SysOperationLog 操作日志表
type SysOperationLog struct {
	BaseModel
	UserID      string    `json:"user_id" gorm:"column:user_id;index;comment:'操作人ID'"`
	UserName    string    `json:"user_name" gorm:"column:user_name;comment:'操作人姓名'"`
	Method      string    `json:"method" gorm:"column:method;comment:'请求方法:GET/POST/PUT/DELETE'"`
	Path        string    `json:"path" gorm:"column:path;comment:'请求路径'"`
	Params      string    `json:"params" gorm:"column:params;type:text;comment:'请求参数'"`
	StatusCode  int       `json:"status_code" gorm:"column:status_code;comment:'响应状态码'"`
	IP          string    `json:"ip" gorm:"column:ip;comment:'操作人IP'"`
	UserAgent   string    `json:"user_agent" gorm:"column:user_agent;type:text;comment:'浏览器信息'"`
	OperationAt time.Time `json:"operation_at" gorm:"column:operation_at;index;comment:'操作时间'"`
}

// TableName 指定表名
func (SysOperationLog) TableName() string {
	return "operation_log"
}
