package systeam

import "github.com/Struggle-Rabbit/CampusLogistics/internal/app"

type SysteamService struct {
	app *app.App
}

func NewSysteamService(app *app.App) *SysteamService {
	return &SysteamService{
		app: app,
	}
}
