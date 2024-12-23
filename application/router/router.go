package router

import (
	"github.com/labstack/echo/v4"
	assistant "github.com/purplior/podoroot/domain/assistant/app"
	assister "github.com/purplior/podoroot/domain/assister/app"
	assisterform "github.com/purplior/podoroot/domain/assisterform/app"
	auth "github.com/purplior/podoroot/domain/auth/app"
	challenge "github.com/purplior/podoroot/domain/challenge/app"
	customervoice "github.com/purplior/podoroot/domain/customervoice/app"
	me "github.com/purplior/podoroot/domain/me/app"
	mission "github.com/purplior/podoroot/domain/mission/app"
	user "github.com/purplior/podoroot/domain/user/app"
	verification "github.com/purplior/podoroot/domain/verification/app"
)

type (
	Router interface {
		Attach(app *echo.Echo)
	}

	router struct {
		assistantRouter     assistant.AssistantRouter
		assisterRouter      assister.AssisterRouter
		assisterFormRouter  assisterform.AssisterFormRouter
		authRouter          auth.AuthRouter
		challengeRouter     challenge.ChallengeRouter
		customerVoiceRouter customervoice.CustomerVoiceRouter
		meRouter            me.MeRouter
		missionRouter       mission.MissionRouter
		userRouter          user.UserRouter
		verificationRouter  verification.VerificationRouter
	}
)

func (r *router) Attach(app *echo.Echo) {
	api := app.Group("/api/:version")

	r.assistantRouter.Attach(api)
	r.assisterRouter.Attach(api)
	r.assisterFormRouter.Attach(api)
	r.authRouter.Attach(api)
	r.challengeRouter.Attach(api)
	r.customerVoiceRouter.Attach(api)
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
	customerVoiceRouter customervoice.CustomerVoiceRouter,
	meRouter me.MeRouter,
	missionRouter mission.MissionRouter,
	userRouter user.UserRouter,
	verificationRouter verification.VerificationRouter,
) Router {
	return &router{
		assistantRouter:     assistantRouter,
		assisterRouter:      assisterRouter,
		assisterFormRouter:  assisterFormRouter,
		authRouter:          authRouter,
		challengeRouter:     challengeRouter,
		customerVoiceRouter: customerVoiceRouter,
		meRouter:            meRouter,
		missionRouter:       missionRouter,
		userRouter:          userRouter,
		verificationRouter:  verificationRouter,
	}
}
