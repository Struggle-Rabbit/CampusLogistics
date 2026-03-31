package model

// CampusBuilding 校区楼栋表
type CampusBuilding struct {
	BaseModel
	CampusName string `gorm:"column:campus_name;not null;size:100" json:"campus_name"` // 校区名称
	BuildingNo string `gorm:"column:building_no;not null;size:50" json:"building_no"`  // 楼栋号
	FloorCount int    `gorm:"column:floor_count;not null" json:"floor_count"`          // 楼层数
}

func (CampusBuilding) TableName() string {
	return "campus_building"
}
