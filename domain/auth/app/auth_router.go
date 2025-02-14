package app

import (
	"github.com/labstack/echo/v4"
	"github.com/purplior/sbec/application/api"
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
			api.HandlerFuncOption{},
		),
	)

	authRouterGroup.POST(
		"/sign-in",
		api.Handler(
			r.authController.SignIn(),
			api.HandlerFuncOption{},
		),
	)

	authRouterGroup.POST(
		"/sign-up",
		api.Handler(
			r.authController.SignUp(),
			api.HandlerFuncOption{},
		),
	)

	authRouterGroup.POST(
		"/reset-password",
		api.Handler(
			r.authController.ResetPassword(),
			api.HandlerFuncOption{},
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
