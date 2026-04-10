package dao

import "gorm.io/gorm"

func IsFieldUnique(db *gorm.DB, tableName string, field string, value string, id *string) bool {
	query := db.Table(tableName).Where(field+" = ? ", value)

	if id != nil {
		query.Where("id != ?", id)
	}
	var count int64
	query.Count(&count)

	return count == 0
}
