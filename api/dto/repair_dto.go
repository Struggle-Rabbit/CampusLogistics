package dto

import (
	"time"
)

type RepairOrderSubmitReq struct {
	RepairType  int      `json:"repair_type"` // 报修类型：1-水电 2-家具 3-网络 4-其他
	Address     string   `json:"address"`     // 报修地点
	Description string   `json:"description"` // 问题描述
	Images      []string `json:"images"`      // 图片URL列表（JSON序列化）
	Contact     string   `json:"contact"`     // 联系人
	Phone       string   `json:"phone"`       // 联系电话
}

type UpdateRepairOrderSubmitReq struct {
	ID          string    `json:"id"`
	RepairType  *int      `json:"repair_type"` // 报修类型：1-水电 2-家具 3-网络 4-其他
	Address     *string   `json:"address"`     // 报修地点
	Description *string   `json:"description"` // 问题描述
	Images      []*string `json:"images"`      // 图片URL列表（JSON序列化）
	Contact     *string   `json:"contact"`     // 联系人
	Phone       *string   `json:"phone"`       // 联系电话
}

type RepairOrderListByPageReq struct {
	PageReq
	OrderNo    string `json:"order_no"`
	RepairType int    `json:"repair_type"` // 报修类型：1-水电 2-家具 3-网络 4-其他
	Contact    string `json:"contact"`     // 联系人
	Phone      string `json:"phone"`       // 联系电话
	Status     int    `json:"status"`
	HandlerID  string `json:"handler_id"`
}

type RepairOrderResult struct {
	ID          string    `json:"id"`
	OrderNo     string    `json:"order_no"`    // 报修单号（唯一）
	UserID      string    `json:"user_id"`     // 提交人学号/工号
	RepairType  int       `json:"repair_type"` // 报修类型：1-水电 2-家具 3-网络 4-其他
	Address     string    `json:"address"`     // 报修地点
	Description string    `json:"description"` // 问题描述
	Images      []string  `json:"images"`      // 图片URL列表（JSON序列化）
	Contact     string    `json:"contact"`     // 联系人
	Phone       string    `json:"phone"`       // 联系电话
	Status      int       `json:"status"`      // 状态：1-待分配 2-待处理 3-处理中 4-已完成 5-已驳回 6-已撤销
	HandlerID   *string   `json:"handler_id"`  // 处理人学号/工号（指针类型允许为空）
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
