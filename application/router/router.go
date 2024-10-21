package router

import (
	"github.com/labstack/echo/v4"
	assistant "github.com/podossaem/podoroot/domain/assistant/app"
	assisterform "github.com/podossaem/podoroot/domain/assisterform/app"
	auth "github.com/podossaem/podoroot/domain/auth/app"
	me "github.com/podossaem/podoroot/domain/me/app"
	user "github.com/podossaem/podoroot/domain/user/app"
	verification "github.com/podossaem/podoroot/domain/verification/app"
)

type (
	Router interface {
		Attach(app *echo.Echo)
	}

	router struct {
		assistantRouter    assistant.AssistantRouter
		assisterFormRouter assisterform.AssisterFormRouter
		authRouter         auth.AuthRouter
		meRouter           me.MeRouter
		userRouter         user.UserRouter
		verificationRouter verification.VerificationRouter
	}
)

func (r *router) Attach(app *echo.Echo) {
	api := app.Group("/api/:version")

	r.assistantRouter.Attach(api)
	r.assisterFormRouter.Attach(api)
	r.authRouter.Attach(api)
	r.meRouter.Attach(api)
	r.userRouter.Attach(api)
	r.verificationRouter.Attach(api)
}

func New(
	assistantRouter assistant.AssistantRouter,
	assisterFormRouter assisterform.AssisterFormRouter,
	authRouter auth.AuthRouter,
	meRouter me.MeRouter,
	userRouter user.UserRouter,
	verificationRouter verification.VerificationRouter,
) Router {
	return &router{
		assistantRouter:    assistantRouter,
		assisterFormRouter: assisterFormRouter,
		authRouter:         authRouter,
		meRouter:           meRouter,
		userRouter:         userRouter,
		verificationRouter: verificationRouter,
	}
}
