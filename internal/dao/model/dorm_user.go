package model

import "time"

// DormUser 宿舍人员关联表
type DormUser struct {
	BaseModel
	RoomID       uint       `gorm:"column:room_id;not null;index" json:"room_id"`         // 宿舍ID
	UserID       string     `gorm:"column:user_id;not null;size:50;index" json:"user_id"` // 用户学号/工号
	CheckInTime  *time.Time `gorm:"column:check_in_time" json:"check_in_time"`            // 入住时间
	CheckOutTime *time.Time `gorm:"column:check_out_time" json:"check_out_time"`          // 迁出时间
	Status       int        `gorm:"column:status;not null;default:1" json:"status"`       // 状态：1-在住 2-已迁出
}

func (DormUser) TableName() string {
	return "dorm_user"
}
