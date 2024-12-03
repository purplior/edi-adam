package app

import (
	"github.com/labstack/echo/v4"
	"github.com/purplior/podoroot/application/api"
)

type (
	VerificationRouter interface {
		Attach(router *echo.Group)
	}
)

type (
	verificationRouter struct {
		emailVerificationController EmailVerificationController
	}
)

func (r *verificationRouter) Attach(router *echo.Group) {
	emailVerificationRouter := router.Group("/verifications")

	emailVerificationRouter.POST(
		"/email/request-code",
		api.Handler(
			r.emailVerificationController.RequestCode(),
			api.HandlerFuncOption{},
		),
	)
	emailVerificationRouter.POST(
		"/email/verify-code",
		api.Handler(
			r.emailVerificationController.VerifyCode(),
			api.HandlerFuncOption{},
		),
	)
}

func NewVerificationRouter(
	emailVerificationController EmailVerificationController,
) VerificationRouter {
	return &verificationRouter{
		emailVerificationController: emailVerificationController,
	}
}
