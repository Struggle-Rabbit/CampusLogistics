package model

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel 基础模型（包含ID、创建时间、更新时间）
type BaseModel struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BaseModelWithDelete 带软删除的基础模型
type BaseModelWithDelete struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 软删除字段，JSON序列化时隐藏
}
