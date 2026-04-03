package dao

import "gorm.io/gorm"

// Paginate 分页通用方法
func Paginate(page int, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page <= 0 {
			page = 1
		}
		switch {
		case size > 100:
			size = 100 // 最大限制100条
		case size <= 0:
			size = 10
		}
		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}
