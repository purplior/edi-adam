package app

import (
	"github.com/labstack/echo/v4"
	"github.com/purplior/sbec/application/api"
)

type (
	VerificationRouter interface {
		Attach(router *echo.Group)
	}
)

type (
	verificationRouter struct {
		emailVerificationController EmailVerificationController
		phoneVerificationController PhoneVerificationController
	}
)

func (r *verificationRouter) Attach(router *echo.Group) {
	verificationRouterGroup := router.Group("/verifications")

	verificationRouterGroup.POST(
		"/email/request-code",
		api.Handler(
			r.emailVerificationController.RequestCode(),
			api.HandlerFuncOption{},
		),
	)

	verificationRouterGroup.POST(
		"/email/verify-code",
		api.Handler(
			r.emailVerificationController.VerifyCode(),
			api.HandlerFuncOption{},
		),
	)

	verificationRouterGroup.POST(
		"/phone/request-code",
		api.Handler(
			r.phoneVerificationController.RequestCode(),
			api.HandlerFuncOption{},
		),
	)

	verificationRouterGroup.POST(
		"/phone/verify-code",
		api.Handler(
			r.phoneVerificationController.VerifyCode(),
			api.HandlerFuncOption{},
		),
	)

	verificationRouterGroup.POST(
		"/joined-phone/request-code",
		api.Handler(
			r.phoneVerificationController.RequestCodeOfJoinedUser(),
			api.HandlerFuncOption{},
		),
	)
}

func NewVerificationRouter(
	emailVerificationController EmailVerificationController,
	phoneVerificationController PhoneVerificationController,
) VerificationRouter {
	return &verificationRouter{
		emailVerificationController: emailVerificationController,
		phoneVerificationController: phoneVerificationController,
	}
}
