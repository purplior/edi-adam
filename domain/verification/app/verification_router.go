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
	emailVerificationRouter := router.Group("/email-verifications")

	emailVerificationRouter.POST(
		"/request-code",
		api.Handler(
			r.emailVerificationController.RequestCode(),
		),
	)
	emailVerificationRouter.POST(
		"/verify-code",
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
