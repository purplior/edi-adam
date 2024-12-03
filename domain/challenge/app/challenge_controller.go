package app

import (
	domain "github.com/purplior/podoroot/domain/challenge"
	"github.com/purplior/podoroot/domain/shared/inner"
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
