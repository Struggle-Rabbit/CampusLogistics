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
	menuSvc := menu.NewMenuService(app)
	roleSvc := role.NewRoleService(app)
	userSvc := user.NewUserService(app, menuSvc)
	systemSvc := system.NewSystemService(app)

	return &ServiceProvider{
		UserService:    userSvc,
		MenuService:    menuSvc,
		RoleService:    roleSvc,
		SysteamService: systemSvc,
	}
}
