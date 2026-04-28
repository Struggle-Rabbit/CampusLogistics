package dto

type CampusCreateReq struct {
	CampusName string `json:"campus_name" binding:"required,max=100"` // 校区名称
	Address    string `json:"address" binding:"max=200"`              // 校区地址
	Contact    string `json:"contact" binding:"max=20"`               // 联系方式
	Phone      string `json:"phone" binding:"max=20"`                 // 联系电话
	Remark     string `json:"remark" binding:"max=500"`               // 备注
}

type CampusUpdateReq struct {
	ID         string `json:"id" binding:"required"`                  // 校区ID
	CampusName string `json:"campus_name" binding:"required,max=100"` // 校区名称
	Address    string `json:"address" binding:"max=200"`              // 校区地址
	Contact    string `json:"contact" binding:"max=20"`               // 联系方式
	Phone      string `json:"phone" binding:"max=20"`                 // 联系电话
	Remark     string `json:"remark" binding:"max=500"`               // 备注
}

type CampusListPageReq struct {
	PageReq
	CampusName string `form:"campus_name"` // 校区名称
}

type CampusResult struct {
	ID            string `json:"id"`
	CampusName    string `json:"campus_name"`
	Address       string `json:"address"`
	Contact       string `json:"contact"`
	Phone         string `json:"phone"`
	Remark        string `json:"remark"`
	BuildingCount int    `json:"building_count"` // 楼栋数量
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
