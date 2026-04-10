package model

// RepairRecord 报修单流转记录表
type RepairRecord struct {
	BaseModel  `gorm:"embedded"`
	OrderID    string `gorm:"column:order_id;not null;index" json:"order_id"`         // 报修单ID
	OperatorID string `gorm:"column:operator_id;not null;size:50" json:"operator_id"` // 操作人学号/工号
	OldStatus  int    `gorm:"column:old_status;not null" json:"old_status"`           // 旧状态
	NewStatus  int    `gorm:"column:new_status;not null" json:"new_status"`           // 新状态
	Remark     string `gorm:"column:remark;type:text" json:"remark"`                  // 操作备注
}

func (RepairRecord) TableName() string {
	return "repair_record"
}
