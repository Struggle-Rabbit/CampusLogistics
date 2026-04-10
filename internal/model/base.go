package model

import (
	"time"

	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"gorm.io/gorm"
)

// BaseModel 基础模型（包含ID、创建时间、更新时间）
type BaseModel struct {
	ID        string    `gorm:"primarykey;type:varchar(32);comment:主键ID" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BaseModelIntId struct {
	ID        uint      `gorm:"primarykey;comment:主键ID" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BaseModelWithDelete 带软删除的基础模型
type BaseModelWithDelete struct {
	ID        string         `gorm:"primarykey;type:varchar(32);comment:主键ID" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 软删除字段，JSON序列化时隐藏
}

func (bs *BaseModel) BeforeCreate() error {
	bs.ID = utils.GenStringID()
	return nil
}

func (bs *BaseModelWithDelete) BeforeCreate() error {
	bs.ID = utils.GenStringID()
	return nil
}
