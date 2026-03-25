package dao

import (
	"log"
	"os"
	"time"

	"github.com/Struggle-Rabbit/CampusLogistics/configs"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/dao/model"
	"gorm.io/driver/mysql"
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

	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)

	if configs.GlobalConfig.App.Env == "dev" {
		if err := db.AutoMigrate(
			&model.SysUser{},
			&model.SysRole{},
			&model.SysPermission{},
			&model.RepairOrder{},
			// 其他模型...
		); err != nil {
			return err
		}
	}

	DB = db

	return nil

}
