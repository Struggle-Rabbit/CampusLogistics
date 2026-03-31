package model

import "time"

// Notice 公告表
type Notice struct {
	BaseModelWithDelete
	Title       string     `gorm:"column:title;not null;size:200" json:"title"`            // 公告标题
	Content     string     `gorm:"column:content;type:text;not null" json:"content"`       // 公告内容
	NoticeType  int        `gorm:"column:notice_type;not null;index" json:"notice_type"`   // 公告类型：1-后勤通知 2-校园通知 3-紧急通知
	IsTop       int        `gorm:"column:is_top;not null;default:2" json:"is_top"`         // 是否置顶：1-是 2-否
	PublishTime *time.Time `gorm:"column:publish_time;index" json:"publish_time"`          // 发布时间
	ViewCount   int        `gorm:"column:view_count;not null;default:0" json:"view_count"` // 浏览量
	CreatorID   string     `gorm:"column:creator_id;not null;size:50" json:"creator_id"`   // 创建人学号/工号
	Attachments []string   `gorm:"column:attachments;type:json" json:"attachments"`        // 附件URL列表（JSON序列化）
}

func (Notice) TableName() string {
	return "notice"
}
