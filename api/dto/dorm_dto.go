package dto

type DormCreateReq struct {
	BuildingID string `json:"building_id" binding:"required"`             // 楼栋ID
	RoomNo     string `json:"room_no" binding:"required,max=50"`          // 宿舍编号
	Floor      int    `json:"floor" binding:"required"`                   // 楼层
	RoomType   int    `json:"room_type" binding:"required,oneof=1 2 3 4"` // 宿舍类型：1-4人间 2-6人间 3-8人间 4-其他
	MaxCount   int    `json:"max_count" binding:"required,min=1"`         // 最大入住人数
	Remark     string `json:"remark" binding:"max=500"`                   // 备注
}

type DormUpdateReq struct {
	ID         string `json:"id" binding:"required"`                      // 宿舍ID
	BuildingID string `json:"building_id" binding:"required"`             // 楼栋ID
	RoomNo     string `json:"room_no" binding:"required,max=50"`          // 宿舍编号
	Floor      int    `json:"floor" binding:"required"`                   // 楼层
	RoomType   int    `json:"room_type" binding:"required,oneof=1 2 3 4"` // 宿舍类型：1-4人间 2-6人间 3-8人间 4-其他
	MaxCount   int    `json:"max_count" binding:"required,min=1"`         // 最大入住人数
	Remark     string `json:"remark" binding:"max=500"`                   // 备注
}

type DormListPageReq struct {
	PageReq
	BuildingID string `form:"building_id"` // 楼栋ID
	CampusID   string `form:"campus_id"`   // 校区ID
	Floor      int    `form:"floor"`       // 楼层
	RoomType   int    `form:"room_type"`   // 宿舍类型
	Status     int    `form:"status"`      // 状态：1-可入住 2-已满 3-装修中
}

type DormResult struct {
	ID           string  `json:"id"`
	BuildingID   string  `json:"building_id"`
	BuildingName string  `json:"building_name"` // 楼栋名称
	CampusID     string  `json:"campus_id"`
	CampusName   string  `json:"campus_name"` // 校区名称
	RoomNo       string  `json:"room_no"`
	Floor        int     `json:"floor"`
	RoomType     int     `json:"room_type"` // 宿舍类型：1-4人间 2-6人间 3-8人间 4-其他
	RoomTypeName string  `json:"room_type_name"`
	MaxCount     int     `json:"max_count"`
	CurrentCount int     `json:"current_count"`
	AvailableBed int     `json:"available_bed"` // 剩余床位
	FillRate     float64 `json:"fill_rate"`     // 入住率
	Status       int     `json:"status"`        // 状态：1-可入住 2-已满 3-装修中
	Remark       string  `json:"remark"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
}

type DormAssignReq struct {
	RoomID string `json:"room_id" binding:"required"` // 宿舍ID
	UserID string `json:"user_id" binding:"required"` // 用户ID
}

type DormTransferReq struct {
	RoomID       string `json:"room_id" binding:"required"`        // 当前宿舍ID
	UserID       string `json:"user_id" binding:"required"`        // 用户ID
	TargetRoomID string `json:"target_room_id" binding:"required"` // 目标宿舍ID
}

type DormCheckOutReq struct {
	RoomID string `json:"room_id" binding:"required"` // 宿舍ID
	UserID string `json:"user_id" binding:"required"` // 用户ID
}

type DormUserListReq struct {
	PageReq
	RoomID   string `form:"room_id"`   // 宿舍ID
	UserID   string `form:"user_id"`   // 用户学号/工号
	UserName string `form:"user_name"` // 用户姓名
	Status   int    `form:"status"`    // 状态：1-在住 2-已迁出
}

type DormUserResult struct {
	ID           string `json:"id"`
	RoomID       string `json:"room_id"`
	RoomNo       string `json:"room_no"`       // 宿舍编号
	BuildingName string `json:"building_name"` // 楼栋名称
	CampusName   string `json:"campus_name"`   // 校区名称
	UserID       string `json:"user_id"`
	UserName     string `json:"user_name"`
	UserType     string `json:"user_type"`
	CheckInTime  string `json:"check_in_time"`
	CheckOutTime string `json:"check_out_time"`
	Status       int    `json:"status"` // 状态：1-在住 2-已迁出
}

const (
	DormRoomType4     = 1 // 4人间
	DormRoomType6     = 2 // 6人间
	DormRoomType8     = 3 // 8人间
	DormRoomTypeOther = 4 // 其他
)
