package dto

type BuildingCreateReq struct {
	CampusID   string `json:"campus_id" binding:"required"` // 校区ID
	BuildingNo string `json:"building_no" binding:"required,max=50"` // 楼栋编号
	BuildingName string `json:"building_name" binding:"required,max=100"` // 楼栋名称
	FloorCount int    `json:"floor_count" binding:"required,min=1"`     // 楼层数
	RoomCount  int    `json:"room_count" binding:"required,min=0"`       // 房间数
	Remark     string `json:"remark" binding:"max=500"`                 // 备注
}

type BuildingUpdateReq struct {
	ID           string `json:"id" binding:"required"`                     // 楼栋ID
	CampusID     string `json:"campus_id" binding:"required"`              // 校区ID
	BuildingNo   string `json:"building_no" binding:"required,max=50"`    // 楼栋编号
	BuildingName string `json:"building_name" binding:"required,max=100"` // 楼栋名称
	FloorCount   int    `json:"floor_count" binding:"required,min=1"`     // 楼层数
	RoomCount    int    `json:"room_count" binding:"required,min=0"`      // 房间数
	Remark       string `json:"remark" binding:"max=500"`                 // 备注
}

type BuildingListPageReq struct {
	PageReq
	CampusID     string `form:"campus_id"`      // 校区ID
	BuildingNo   string `form:"building_no"`    // 楼栋编号
	BuildingName string `form:"building_name"`  // 楼栋名称
}

type BuildingResult struct {
	ID           string `json:"id"`
	CampusID     string `json:"campus_id"`
	CampusName   string `json:"campus_name"`     // 校区名称
	BuildingNo   string `json:"building_no"`
	BuildingName string `json:"building_name"`
	FloorCount   int    `json:"floor_count"`
	RoomCount    int    `json:"room_count"`
	Remark       string `json:"remark"`
	RoomUsedCount int   `json:"room_used_count"` // 已分配房间数
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type BuildingImportReq struct {
	CampusName   string `json:"campus_name" binding:"required"` // 校区名称（用于匹配校区ID）
	BuildingNo   string `json:"building_no" binding:"required"`
	BuildingName string `json:"building_name" binding:"required"`
	FloorCount   int    `json:"floor_count" binding:"required"`
	RoomCount    int    `json:"room_count" binding:"required"`
	Remark       string `json:"remark"`
}

type BuildingExportReq struct {
	CampusID   string `form:"campus_id"`
	BuildingNo string `form:"building_no"`
}