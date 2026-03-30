package dao

import (
	"log"
	"os"
	"time"

	"github.com/Struggle-Rabbit/CampusLogistics/configs"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao/model"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() error {
	cfg := configs.GlobalConfig.MySQL

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	var (
		db    *gorm.DB
		dbErr error
	)
	if configs.IsDev() {
		db, dbErr = gorm.Open(sqlite.Open("dev.db"), &gorm.Config{
			Logger: newLogger,
		})

	} else {
		db, dbErr = gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
			Logger: newLogger,
		})
	}

	if dbErr != nil {
		return dbErr
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	if configs.IsDev() {
		if err := db.AutoMigrate(
			&model.SysUser{},
			&model.SysRole{},
			&model.SysPermission{},
			&model.SysRolePermission{},
			&model.RepairOrder{},
			&model.CampusBuilding{},
			&model.DormRoom{},
			&model.DormUser{},
			&model.DormUtility{},
			&model.Notice{},
			&model.RepairRecord{},
			&model.SysOperationLog{},
		); err != nil {
			return err
		}
	}

	DB = db

	return nil

}
