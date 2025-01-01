package router

import (
	"github.com/labstack/echo/v4"
	admin "github.com/purplior/podoroot/application/admin"
	assistant "github.com/purplior/podoroot/domain/assistant/app"
	assister "github.com/purplior/podoroot/domain/assister/app"
	auth "github.com/purplior/podoroot/domain/auth/app"
	bookmark "github.com/purplior/podoroot/domain/bookmark/app"
	category "github.com/purplior/podoroot/domain/category/app"
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
		adminRouter         admin.AdminRouter
		assistantRouter     assistant.AssistantRouter
		assisterRouter      assister.AssisterRouter
		authRouter          auth.AuthRouter
		bookmarkRouter      bookmark.BookmarkRouter
		categoryRouter      category.CategoryRouter
		challengeRouter     challenge.ChallengeRouter
		customerVoiceRouter customervoice.CustomerVoiceRouter
		meRouter            me.MeRouter
		missionRouter       mission.MissionRouter
		userRouter          user.UserRouter
		verificationRouter  verification.VerificationRouter
	}
)

func (r *router) Attach(app *echo.Echo) {
	adminGroup := app.Group("/_admin")
	r.adminRouter.Attach(adminGroup)

	apiGroup := app.Group("/api/:version")

	r.assistantRouter.Attach(apiGroup)
	r.assisterRouter.Attach(apiGroup)
	r.authRouter.Attach(apiGroup)
	r.bookmarkRouter.Attach(apiGroup)
	r.categoryRouter.Attach(apiGroup)
	r.challengeRouter.Attach(apiGroup)
	r.customerVoiceRouter.Attach(apiGroup)
	r.meRouter.Attach(apiGroup)
	r.missionRouter.Attach(apiGroup)
	r.userRouter.Attach(apiGroup)
	r.verificationRouter.Attach(apiGroup)
}

func New(
	adminRouter admin.AdminRouter,
	assistantRouter assistant.AssistantRouter,
	assisterRouter assister.AssisterRouter,
	authRouter auth.AuthRouter,
	bookmarkRouter bookmark.BookmarkRouter,
	categoryRouter category.CategoryRouter,
	challengeRouter challenge.ChallengeRouter,
	customerVoiceRouter customervoice.CustomerVoiceRouter,
	meRouter me.MeRouter,
	missionRouter mission.MissionRouter,
	userRouter user.UserRouter,
	verificationRouter verification.VerificationRouter,
) Router {
	return &router{
		adminRouter:         adminRouter,
		assistantRouter:     assistantRouter,
		assisterRouter:      assisterRouter,
		authRouter:          authRouter,
		bookmarkRouter:      bookmarkRouter,
		categoryRouter:      categoryRouter,
		challengeRouter:     challengeRouter,
		customerVoiceRouter: customerVoiceRouter,
		meRouter:            meRouter,
		missionRouter:       missionRouter,
		userRouter:          userRouter,
		verificationRouter:  verificationRouter,
	}
}
