package router

import (
	"github.com/labstack/echo/v4"
	assistant "github.com/podossaem/podoroot/domain/assistant/app"
	assister "github.com/podossaem/podoroot/domain/assister/app"
	assisterform "github.com/podossaem/podoroot/domain/assisterform/app"
	auth "github.com/podossaem/podoroot/domain/auth/app"
	challenge "github.com/podossaem/podoroot/domain/challenge/app"
	me "github.com/podossaem/podoroot/domain/me/app"
	mission "github.com/podossaem/podoroot/domain/mission/app"
	user "github.com/podossaem/podoroot/domain/user/app"
	verification "github.com/podossaem/podoroot/domain/verification/app"
)

type (
	Router interface {
		Attach(app *echo.Echo)
	}

	router struct {
		assistantRouter    assistant.AssistantRouter
		assisterRouter     assister.AssisterRouter
		assisterFormRouter assisterform.AssisterFormRouter
		authRouter         auth.AuthRouter
		challengeRouter    challenge.ChallengeRouter
		meRouter           me.MeRouter
		missionRouter      mission.MissionRouter
		userRouter         user.UserRouter
		verificationRouter verification.VerificationRouter
	}
)

func (r *router) Attach(app *echo.Echo) {
	api := app.Group("/api/:version")

	r.assistantRouter.Attach(api)
	r.assisterRouter.Attach(api)
	r.assisterFormRouter.Attach(api)
	r.authRouter.Attach(api)
	r.challengeRouter.Attach(api)
	r.meRouter.Attach(api)
	r.missionRouter.Attach(api)
	r.userRouter.Attach(api)
	r.verificationRouter.Attach(api)
}

func New(
	assistantRouter assistant.AssistantRouter,
	assisterRouter assister.AssisterRouter,
	assisterFormRouter assisterform.AssisterFormRouter,
	authRouter auth.AuthRouter,
	challengeRouter challenge.ChallengeRouter,
	meRouter me.MeRouter,
	missionRouter mission.MissionRouter,
	userRouter user.UserRouter,
	verificationRouter verification.VerificationRouter,
) Router {
	return &router{
		assistantRouter:    assistantRouter,
		assisterRouter:     assisterRouter,
		assisterFormRouter: assisterFormRouter,
		authRouter:         authRouter,
		challengeRouter:    challengeRouter,
		meRouter:           meRouter,
		missionRouter:      missionRouter,
		userRouter:         userRouter,
		verificationRouter: verificationRouter,
	}
}
