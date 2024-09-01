package app

import (
	"github.com/labstack/echo/v4"
)

type (
	UserRouter interface {
		Attach(router *echo.Group)
	}

	userRouter struct {
		userController UserController
	}
)

func (r *userRouter) Attach(router *echo.Group) {
}

func NewUserRouter(
	userController UserController,
) UserRouter {
	return &userRouter{
		userController: userController,
	}
}
