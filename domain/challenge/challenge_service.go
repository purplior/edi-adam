package challenge

import (
	"github.com/podossaem/podoroot/domain/shared/inner"
	"github.com/podossaem/podoroot/lib/mydate"
)

type (
	ChallengeService interface {
		PatchOne_ReceivedStatus(
			ctx inner.Context,
			id string,
		) error
	}
)

type (
	challengeService struct {
		challengeRepository ChallengeRepository
		cm                  inner.ContextManager
	}
)

func (s *challengeService) PatchOne_ReceivedStatus(
	ctx inner.Context,
	id string,
) error {
	return s.challengeRepository.UpdateOne_ReceivedStatus_ByID(
		ctx,
		id,
		true,
		mydate.Now(),
	)
}

func NewChallengeService(
	challengeRepository ChallengeRepository,
	cm inner.ContextManager,
) ChallengeService {
	return &challengeService{
		challengeRepository: challengeRepository,
		cm:                  cm,
	}
}
