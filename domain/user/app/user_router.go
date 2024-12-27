package app

import (
	"github.com/labstack/echo/v4"
	"github.com/purplior/podoroot/application/api"
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
		"/nickname-check",
		api.Handler(
			r.userController.CheckNicknameExistence(),
			api.HandlerFuncOption{},
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
