package app

import (
	"github.com/labstack/echo/v4"
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
}

func NewChallengeRouter(
	challengeController ChallengeController,
) ChallengeRouter {
	return &challengeRouter{
		challengeController: challengeController,
	}
}
