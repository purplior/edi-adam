package api

import (
	"github.com/labstack/echo/v4"
	user "github.com/podossaem/podoroot/domain/user/app"
	verification "github.com/podossaem/podoroot/domain/verification/app"
)

type (
	Router interface {
		Attach(app *echo.Echo)
	}

	router struct {
		userRouter         user.UserRouter
		verificationRouter verification.VerificationRouter
	}
)

func (r *router) Attach(app *echo.Echo) {
	api := app.Group("/api/:version")

	r.userRouter.Attach(api)
	r.verificationRouter.Attach(api)
}

func NewRouter(
	userRouter user.UserRouter,
	verificationRouter verification.VerificationRouter,
) Router {
	return &router{
		userRouter:         userRouter,
		verificationRouter: verificationRouter,
	}
}
