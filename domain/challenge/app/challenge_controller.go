package app

import (
	domain "github.com/purplior/sbec/domain/challenge"
	"github.com/purplior/sbec/domain/shared/inner"
)

type (
	ChallengeController interface {
	}
)

type (
	challengeController struct {
		challengeService domain.ChallengeService
		cm               inner.ContextManager
	}
)

func NewChallengeController(
	challengeService domain.ChallengeService,
	cm inner.ContextManager,
) ChallengeController {
	return &challengeController{
		challengeService: challengeService,
		cm:               cm,
	}
}
