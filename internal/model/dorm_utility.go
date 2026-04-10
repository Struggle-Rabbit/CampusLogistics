package model

// DormUtility 水电费表
type DormUtility struct {
	BaseModel     `gorm:"embedded"`
	RoomID        string  `gorm:"column:room_id;not null;index:idx_room_month" json:"room_id"`    // 宿舍ID
	Year          int     `gorm:"column:year;not null;index:idx_room_month" json:"year"`          // 年份
	Month         int     `gorm:"column:month;not null;index:idx_room_month" json:"month"`        // 月份
	WaterUsage    float64 `gorm:"column:water_usage;type:decimal(10,2)" json:"water_usage"`       // 用水量（吨）
	ElectricUsage float64 `gorm:"column:electric_usage;type:decimal(10,2)" json:"electric_usage"` // 用电量（度）
	Amount        float64 `gorm:"column:amount;type:decimal(10,2);not null" json:"amount"`        // 应缴金额
	PayStatus     int     `gorm:"column:pay_status;not null;default:1" json:"pay_status"`         // 缴费状态：1-未缴费 2-已缴费
}

func (DormUtility) TableName() string {
	return "dorm_utility"
}
