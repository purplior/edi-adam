package app

import (
	"github.com/podossaem/podoroot/domain/user"
)

type (
	UserController interface {
	}
)

type (
	userController struct {
		userService user.UserService
	}
)

func NewUserController() UserController {
	return &userController{}
}
