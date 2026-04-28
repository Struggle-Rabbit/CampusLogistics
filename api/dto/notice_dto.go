package dto

import "time"

type NoticeCreateReq struct {
	Title       string    `json:"title" binding:"required,max=200"`           // 公告标题
	Content     string    `json:"content" binding:"required"`                 // 公告内容
	NoticeType  int       `json:"notice_type" binding:"required,oneof=1 2 3"` // 公告类型：1-后勤通知 2-校园通知 3-紧急通知
	IsTop       int       `json:"is_top" binding:"oneof=1 2"`                 // 是否置顶：1-是 2-否
	PublishTime time.Time `json:"publish_time" binding:"required"`            // 发布时间
	Attachments []string  `json:"attachments"`                                // 附件URL列表
}

type NoticeUpdateReq struct {
	ID          string    `json:"id" binding:"required"`                      // 公告ID
	Title       string    `json:"title" binding:"required,max=200"`           // 公告标题
	Content     string    `json:"content" binding:"required"`                 // 公告内容
	NoticeType  int       `json:"notice_type" binding:"required,oneof=1 2 3"` // 公告类型：1-后勤通知 2-校园通知 3-紧急通知
	IsTop       int       `json:"is_top" binding:"oneof=1 2"`                 // 是否置顶：1-是 2-否
	PublishTime time.Time `json:"publish_time"`                               // 发布时间
	Attachments []string  `json:"attachments"`                                // 附件URL列表
}

type NoticeListPageReq struct {
	PageReq
	Title      string `form:"title"`       // 公告标题
	NoticeType int    `form:"notice_type"` // 公告类型
	IsTop      int    `form:"is_top"`      // 是否置顶
	StartTime  string `form:"start_time"`  // 开始时间
	EndTime    string `form:"end_time"`    // 结束时间
}

type NoticeResult struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	NoticeType     int       `json:"notice_type"`
	NoticeTypeName string    `json:"notice_type_name"`
	IsTop          int       `json:"is_top"`
	IsTopName      string    `json:"is_top_name"`
	PublishTime    time.Time `json:"publish_time"`
	ViewCount      int       `json:"view_count"`
	CreatorID      string    `json:"creator_id"`
	CreatorName    string    `json:"creator_name"`
	Attachments    []string  `json:"attachments"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type NoticePublicResult struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Content        string    `json:"content"`
	NoticeType     int       `json:"notice_type"`
	NoticeTypeName string    `json:"notice_type_name"`
	PublishTime    time.Time `json:"publish_time"`
	ViewCount      int       `json:"view_count"`
	Attachments    []string  `json:"attachments"`
}

type NoticeTopReq struct {
	ID    string `json:"id" binding:"required"`
	IsTop int    `json:"is_top" binding:"required,oneof=1 2"`
}

const (
	NoticeTypeLogistics = 1 // 后勤通知
	NoticeTypeCampus    = 2 // 校园通知
	NoticeTypeUrgent    = 3 // 紧急通知

	IsTopYes = 1 // 置顶
	IsTopNo  = 2 // 不置顶

	MaxTopNotice = 3 // 最多置顶数量
)

func GetNoticeTypeName(noticeType int) string {
	switch noticeType {
	case NoticeTypeLogistics:
		return "后勤通知"
	case NoticeTypeCampus:
		return "校园通知"
	case NoticeTypeUrgent:
		return "紧急通知"
	default:
		return "未知"
	}
}

func GetIsTopName(isTop int) string {
	switch isTop {
	case IsTopYes:
		return "已置顶"
	case IsTopNo:
		return "未置顶"
	default:
		return "未置顶"
	}
}
