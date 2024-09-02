package router

import (
	"github.com/labstack/echo/v4"
	auth "github.com/podossaem/podoroot/domain/auth/app"
	user "github.com/podossaem/podoroot/domain/user/app"
	verification "github.com/podossaem/podoroot/domain/verification/app"
)

type (
	Router interface {
		Attach(app *echo.Echo)
	}

	router struct {
		authRouter         auth.AuthRouter
		userRouter         user.UserRouter
		verificationRouter verification.VerificationRouter
	}
)

func (r *router) Attach(app *echo.Echo) {
	api := app.Group("/api/:version")

	r.authRouter.Attach(api)
	r.userRouter.Attach(api)
	r.verificationRouter.Attach(api)
}

func New(
	authRouter auth.AuthRouter,
	userRouter user.UserRouter,
	verificationRouter verification.VerificationRouter,
) Router {
	return &router{
		authRouter:         authRouter,
		userRouter:         userRouter,
		verificationRouter: verificationRouter,
	}
}
