package service

import (
	"github.com/Struggle-Rabbit/CampusLogistics/internal/app"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/menu"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/role"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/system"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/user"
)

type ServiceProvider struct {
	UserService    *user.UserService
	MenuService    *menu.MenuService
	RoleService    *role.RoleService
	SysteamService *system.SystemService
}

func NewServiceProvider(app *app.App) *ServiceProvider {
	return &ServiceProvider{
		UserService:    user.NewUserService(app),
		MenuService:    menu.NewMenuService(app),
		RoleService:    role.NewRoleService(app),
		SysteamService: system.NewSystemService(app),
	}
}
