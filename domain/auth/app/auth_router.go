package app

import (
	"github.com/labstack/echo/v4"
	"github.com/podossaem/podoroot/application/api/controller"
)

type (
	AuthRouter interface {
		Attach(router *echo.Group)
	}
)

type (
	authRouter struct {
		authController AuthController
	}
)

func (r *authRouter) Attach(router *echo.Group) {
	authRouterGroup := router.Group("/auth")

	authRouterGroup.POST(
		"/sign-in-by-email",
		controller.Handler(
			r.authController.SignInByEmailVerification(),
		),
	)

	authRouterGroup.POST(
		"/sign-up-by-email",
		controller.Handler(
			r.authController.SignUpByEmailVerification(),
		),
	)
}

func NewAuthRouter(
	authController AuthController,
) AuthRouter {
	return &authRouter{
		authController: authController,
	}
}
