package model

// Campus 校区表
type Campus struct {
	BaseModel  `gorm:"embedded"`
	CampusName string `gorm:"column:campus_name;not null;size:100" json:"campus_name"`
	Address    string `gorm:"column:address;size:200" json:"address"`
	Contact    string `gorm:"column:contact;size:50" json:"contact"`
	Phone      string `gorm:"column:phone;size:20" json:"phone"`
	Remark     string `gorm:"column:remark;size:500" json:"remark"`
}

func (Campus) TableName() string {
	return "campus"
}

// Building 楼栋表
type Building struct {
	BaseModel    `gorm:"embedded"`
	CampusID     string `gorm:"column:campus_id;not null;index" json:"campus_id"`
	BuildingNo   string `gorm:"column:building_no;not null;size:50;index" json:"building_no"`
	BuildingName string `gorm:"column:building_name;not null;size:100" json:"building_name"`
	FloorCount   int    `gorm:"column:floor_count;not null" json:"floor_count"`
	RoomCount    int    `gorm:"column:room_count;not null" json:"room_count"`
	Remark       string `gorm:"column:remark;size:500" json:"remark"`
}

func (Building) TableName() string {
	return "building"
}

// DormRoom 宿舍表
type DormRoom struct {
	BaseModel    `gorm:"embedded"`
	BuildingID   string `gorm:"column:building_id;not null;index" json:"building_id"`
	RoomNo       string `gorm:"column:room_no;not null;size:50;index" json:"room_no"`
	Floor        int    `gorm:"column:floor;not null" json:"floor"`
	RoomType     int    `gorm:"column:room_type;not null" json:"room_type"`
	MaxCount     int    `gorm:"column:max_count;not null" json:"max_count"`
	CurrentCount int    `gorm:"column:current_count;not null;default:0" json:"current_count"`
	Remark       string `gorm:"column:remark;size:500" json:"remark"`
}

func (DormRoom) TableName() string {
	return "dorm_room"
}
