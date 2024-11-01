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

	challengeRouterGroup.GET(
		"/info-pagination",
		api.Handler(
			r.challengeController.GetPaginatedList(),
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
