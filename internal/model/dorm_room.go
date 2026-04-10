package model

// DormRoom 宿舍表
type DormRoom struct {
	BaseModel    `gorm:"embedded"`
	BuildingID   string `gorm:"column:building_id;not null;index" json:"building_id"`         // 楼栋ID
	RoomNo       string `gorm:"column:room_no;not null;size:50;index" json:"room_no"`         // 宿舍号
	RoomType     int    `gorm:"column:room_type;not null" json:"room_type"`                   // 宿舍类型：1-4人间 2-6人间 3-其他
	MaxCount     int    `gorm:"column:max_count;not null" json:"max_count"`                   // 可住人数
	CurrentCount int    `gorm:"column:current_count;not null;default:0" json:"current_count"` // 已住人数
	Floor        int    `gorm:"column:floor;not null" json:"floor"`                           // 楼层
}

func (DormRoom) TableName() string {
	return "dorm_room"
}
