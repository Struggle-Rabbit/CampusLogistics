package app

import (
	"github.com/Struggle-Rabbit/CampusLogistics/internal/config"
	"gorm.io/gorm"
)

type App struct {
	Config *config.Config
	DB     *gorm.DB
}

func NewApp(cfg *config.Config, db *gorm.DB) *App {
	return &App{
		Config: cfg,
		DB:     db,
	}
}
