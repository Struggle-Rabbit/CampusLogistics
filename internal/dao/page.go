package dao

import "gorm.io/gorm"

// Paginate 分页通用方法
func Paginate(page int, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page < 1 {
			page = 1
		}
		if size < 1 {
			size = 10
		}
		if size > 100 {
			size = 100 // 最大限制100条
		}
		offset := (page - 1) * size
		return db.Offset(offset).Limit(size)
	}
}
