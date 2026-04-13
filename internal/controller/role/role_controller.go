package role

import "github.com/Struggle-Rabbit/CampusLogistics/internal/service"

type RoleController struct {
	srv *service.ServiceProvider
}

func NewRoleController(srv *service.ServiceProvider) *RoleController {
	return &RoleController{
		srv: srv,
	}
}
