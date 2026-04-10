package user

import (
	"github.com/Struggle-Rabbit/CampusLogistics/internal/service"
)

type UserController struct {
	srv *service.ServiceProvider
}

func NewUserController(srv *service.ServiceProvider) *UserController {
	return &UserController{
		srv: srv,
	}
}
