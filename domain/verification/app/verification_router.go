package app

import (
	"github.com/labstack/echo/v4"
	"github.com/podossaem/podoroot/application/api"
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
		),
	)
	emailVerificationRouter.POST(
		"/email/verify-code",
		api.Handler(
			r.emailVerificationController.VerifyCode(),
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
