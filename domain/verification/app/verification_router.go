package app

import (
	"github.com/labstack/echo/v4"
	"github.com/podossaem/podoroot/application/api/controller"
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
		controller.Handler(
			r.emailVerificationController.RequestCode(),
		),
	)
	emailVerificationRouter.POST(
		"/verify-code",
		controller.Handler(
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
