package unittest

import (
	"os"

	"github.com/Struggle-Rabbit/CampusLogistics/internal/app"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/config"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/model"
	"github.com/Struggle-Rabbit/CampusLogistics/pkg/utils"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func SetupTestDB() (*gorm.DB, *app.App) {
	tmpFile := os.TempDir() + "/test_campus_" + "db"
	os.Remove(tmpFile)

	db, err := gorm.Open(sqlite.Open(tmpFile+"?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:        "test-secret",
			AccessExpire:  3600,
			RefreshExpire: 86400,
		},
	}
	config.GlobalConfig = cfg

	utils.InitSnowflake()

	db.AutoMigrate(
		&model.SysUser{},
		&model.SysRole{},
		&model.SysMenu{},
		&model.SysOperationLog{},
		&model.RepairOrder{},
		&model.RepairRecord{},
		&model.Campus{},
		&model.Building{},
		&model.DormRoom{},
		&model.DormUser{},
		&model.DormUtility{},
		&model.UtilityPrice{},
		&model.Notice{},
	)

	appInstance := app.NewApp(cfg, db)

	return db, appInstance
}
