package app

import (
	"github.com/labstack/echo/v4"
	"github.com/podossaem/podoroot/application/api/controller"
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
	userRouterGroup := router.Group("/users")

	userRouterGroup.POST(
		"/sign-up-by-email",
		controller.Handler(
			r.userController.SignUpByEmailVerification(),
		),
	)
}

func NewUserRouter(
	userController UserController,
) UserRouter {
	return &userRouter{
		userController: userController,
	}
}
