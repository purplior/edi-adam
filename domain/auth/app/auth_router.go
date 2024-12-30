package app

import (
	"github.com/labstack/echo/v4"
	"github.com/purplior/podoroot/application/api"
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
		"/admin/email/sign-in",
		api.Handler(
			r.authController.SignIn_ByEmailVerification_ForAdmin(),
			api.HandlerFuncOption{AdminOnly: true},
		),
	)

	authRouterGroup.POST(
		"/email/sign-in",
		api.Handler(
			r.authController.SignIn_ByEmailVerification(),
			api.HandlerFuncOption{},
		),
	)

	authRouterGroup.POST(
		"/phone/sign-in",
		api.Handler(
			r.authController.SignIn_ByPhoneNumberVerification(),
			api.HandlerFuncOption{},
		),
	)

	authRouterGroup.POST(
		"/email/sign-up",
		api.Handler(
			r.authController.SignUp_ByEmailVerification(),
			api.HandlerFuncOption{},
		),
	)

	authRouterGroup.POST(
		"/phone/sign-up",
		api.Handler(
			r.authController.SignUp_ByPhoneNumberVerification(),
			api.HandlerFuncOption{},
		),
	)

	authRouterGroup.POST(
		"/reset-password",
		api.Handler(
			r.authController.ResetPassword_ByPhoneNumberVerification(),
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
