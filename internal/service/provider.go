package service

import (
	"github.com/Struggle-Rabbit/CampusLogistics/internal/app"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/building"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/campus"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/dorm"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/menu"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/notice"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/repair"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/role"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/system"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/user"
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service/utility"
)

type ServiceProvider struct {
	UserService     *user.UserService
	MenuService     *menu.MenuService
	RoleService     *role.RoleService
	SystemService   *system.SystemService
	RepairService   *repair.RepairService
	CampusService   *campus.CampusService
	BuildingService *building.BuildingService
	DormService     *dorm.DormService
	UtilityService  *utility.UtilityService
	NoticeService   *notice.NoticeService
}

func NewServiceProvider(app *app.App) *ServiceProvider {
	menuSvc := menu.NewMenuService(app)
	roleSvc := role.NewRoleService(app)
	userSvc := user.NewUserService(app, menuSvc)
	systemSvc := system.NewSystemService(app)
	repairSvc := repair.NewRepairService(app)
	campusSvc := campus.NewCampusService(app)
	buildingSvc := building.NewBuildingService(app)
	dormSvc := dorm.NewDormService(app)
	utilitySvc := utility.NewUtilityService(app)
	noticeSvc := notice.NewNoticeService(app)

	return &ServiceProvider{
		UserService:     userSvc,
		MenuService:     menuSvc,
		RoleService:     roleSvc,
		SystemService:   systemSvc,
		RepairService:   repairSvc,
		CampusService:   campusSvc,
		BuildingService: buildingSvc,
		DormService:     dormSvc,
		UtilityService:  utilitySvc,
		NoticeService:   noticeSvc,
	}
}
