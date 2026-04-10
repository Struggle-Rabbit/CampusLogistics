package model

// RepairOrder 报修单表
type RepairOrder struct {
	BaseModel   `gorm:"embedded"`
	OrderNo     string   `gorm:"column:order_no;uniqueIndex;not null;size:32" json:"order_no"` // 报修单号（唯一）
	UserID      string   `gorm:"column:user_id;not null;size:50;index" json:"user_id"`         // 提交人学号/工号
	RepairType  int      `gorm:"column:repair_type;not null;index" json:"repair_type"`         // 报修类型：1-水电 2-家具 3-网络 4-其他
	Address     string   `gorm:"column:address;not null;size:255" json:"address"`              // 报修地点
	Description string   `gorm:"column:description;type:text" json:"description"`              // 问题描述
	Images      []string `gorm:"column:images;type:json" json:"images"`                        // 图片URL列表（JSON序列化）
	Contact     string   `gorm:"column:contact;not null;size:50" json:"contact"`               // 联系人
	Phone       string   `gorm:"column:phone;not null;size:20" json:"phone"`                   // 联系电话
	Status      int      `gorm:"column:status;not null;default:1;index" json:"status"`         // 状态：1-待分配 2-待处理 3-处理中 4-已完成 5-已驳回 6-已撤销
	HandlerID   *string  `gorm:"column:handler_id;size:50;index" json:"handler_id"`            // 处理人学号/工号（指针类型允许为空）
}

func (RepairOrder) TableName() string {
	return "repair_order"
}
