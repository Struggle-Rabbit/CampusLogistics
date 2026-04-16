package middleware

import (
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PermissionValidator(perms string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, _ := c.Get("user_id")

		if !CheckPermission(dao.DB, userId.(string), perms) {
			utils.NoPermission(c)
			c.Abort()
			return
		}

		c.Next()

	}
}

func CheckPermission(db *gorm.DB, userId string, perms string) bool {
	var count int64

	err := db.Table("sys_menu as m").
		Joins("JOIN sys_role_menu as rm ON rm.menu_id = m.id").
		Joins("JOIN sys_user_role as ur ON ur.role_id = rm.role_id").
		Where("ur.user_id = ?", userId).
		Where("m.perms = ?", perms).
		Count(&count).Error

	if err != nil {
		return false
	}

	return count > 0
}
