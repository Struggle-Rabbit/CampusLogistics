package model

// UtilityPrice 水电费单价配置表
type UtilityPrice struct {
	BaseModel     `gorm:"embedded"`
	WaterPrice    float64 `gorm:"column:water_price;type:decimal(10,2);not null" json:"water_price"`        // 水价（元/吨）
	ElectricPrice float64 `gorm:"column:electric_price;type:decimal(10,2);not null" json:"electric_price"`    // 电价（元/度）
}

func (UtilityPrice) TableName() string {
	return "utility_price"
}
