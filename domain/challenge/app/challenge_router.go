package app

import (
	"github.com/labstack/echo/v4"
	"github.com/podossaem/podoroot/application/api"
)

type (
	ChallengeRouter interface {
		Attach(router *echo.Group)
	}
)

type (
	challengeRouter struct {
		challengeController ChallengeController
	}
)

func (r *challengeRouter) Attach(router *echo.Group) {
	challengeRouterGroup := router.Group("/challenges")

	challengeRouterGroup.POST(
		"/receive",
		api.Handler(
			r.challengeController.ReceiveOne(),
		),
	)
}

func NewChallengeRouter(
	challengeController ChallengeController,
) ChallengeRouter {
	return &challengeRouter{
		challengeController: challengeController,
	}
}
