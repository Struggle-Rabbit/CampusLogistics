package dto

type UtilityCreateReq struct {
	RoomID        string  `json:"room_id" binding:"required"`                    // 宿舍ID
	Year          int     `json:"year" binding:"required,min=2020,max=2100"`      // 年份
	Month         int     `json:"month" binding:"required,min=1,max=12"`          // 月份
	WaterUsage    float64 `json:"water_usage" binding:"min=0"`                    // 用水量（吨）
	ElectricUsage float64 `json:"electric_usage" binding:"min=0"`                 // 用电量（度）
}

type UtilityUpdateReq struct {
	ID            string  `json:"id" binding:"required"`                         // 记录ID
	WaterUsage    float64 `json:"water_usage" binding:"min=0"`                   // 用水量（吨）
	ElectricUsage float64 `json:"electric_usage" binding:"min=0"`                 // 用电量（度）
}

type UtilityListPageReq struct {
	PageReq
	RoomID    string `form:"room_id"`   // 宿舍ID
	CampusID  string `form:"campus_id"` // 校区ID
	BuildingID string `form:"building_id"` // 楼栋ID
	Year      int    `form:"year"`     // 年份
	Month     int    `form:"month"`    // 月份
	PayStatus int    `form:"pay_status"` // 缴费状态：1-未缴费 2-已缴费
}

type UtilityResult struct {
	ID            string `json:"id"`
	RoomID        string `json:"room_id"`
	RoomNo        string `json:"room_no"`       // 宿舍编号
	BuildingID    string `json:"building_id"`
	BuildingName  string `json:"building_name"` // 楼栋名称
	CampusID      string `json:"campus_id"`
	CampusName    string `json:"campus_name"`   // 校区名称
	Year          int    `json:"year"`
	Month         int    `json:"month"`
	WaterUsage    float64 `json:"water_usage"`
	ElectricUsage float64 `json:"electric_usage"`
	WaterPrice    float64 `json:"water_price"`    // 水价
	ElectricPrice float64 `json:"electric_price"` // 电价
	WaterAmount   float64 `json:"water_amount"`   // 水费
	ElectricAmount float64 `json:"electric_amount"` // 电费
	Amount        float64 `json:"amount"`
	PayStatus     int     `json:"pay_status"`    // 缴费状态：1-未缴费 2-已缴费
	PayStatusName string  `json:"pay_status_name"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type UtilityPayReq struct {
	ID     string `json:"id" binding:"required"` // 记录ID
}

type UtilityBatchPayReq struct {
	IDs []string `json:"ids" binding:"required"` // 记录ID列表
}

type UtilityImportReq struct {
	RoomID        string  `json:"room_id" binding:"required"`  // 宿舍ID
	Year          int     `json:"year" binding:"required"`
	Month         int     `json:"month" binding:"required"`
	WaterUsage    float64 `json:"water_usage"`
	ElectricUsage float64 `json:"electric_usage"`
}

type UtilityPriceReq struct {
	WaterPrice    float64 `json:"water_price" binding:"required,min=0"`    // 水价（元/吨）
	ElectricPrice float64 `json:"electric_price" binding:"required,min=0"` // 电价（元/度）
}

type UtilityPriceResult struct {
	WaterPrice    float64 `json:"water_price"`
	ElectricPrice float64 `json:"electric_price"`
}

type UtilityStatResult struct {
	TotalWaterUsage    float64 `json:"total_water_usage"`
	TotalElectricUsage float64 `json:"total_electric_usage"`
	TotalAmount        float64 `json:"total_amount"`
	UnpaidCount        int     `json:"unpaid_count"`
	UnpaidAmount       float64 `json:"unpaid_amount"`
}