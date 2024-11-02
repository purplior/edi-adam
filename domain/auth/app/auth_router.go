package app

import (
	"github.com/labstack/echo/v4"
	"github.com/podossaem/podoroot/application/api"
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
		"/identity/refresh",
		api.Handler(
			r.authController.RefreshIdentityToken(),
		),
	)

	authRouterGroup.POST(
		"/email/sign-in",
		api.Handler(
			r.authController.SignIn_ByEmailVerification(),
		),
	)

	authRouterGroup.POST(
		"/email/sign-up",
		api.Handler(
			r.authController.SignUp_ByEmailVerification(),
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
