package challenge

import "github.com/podossaem/podoroot/domain/shared/inner"

type (
	ChallengeService interface {
		GetPaginatedListByUserID(
			ctx inner.Context,
			userId string,
			limit int,
			offset int,
		) (
			[]Challenge,
			error,
		)
	}
)

type (
	challengeService struct {
		challengeRepository ChallengeRepository
	}
)

func (s *challengeService) GetPaginatedListByUserID(
	ctx inner.Context,
	userId string,
	limit int,
	offset int,
) (
	[]Challenge,
	error,
) {
	return s.challengeRepository.FindPaginatedListByUserID(
		ctx,
		userId,
		limit,
		offset,
	)
}

func NewChallengeService(
	challengeRepository ChallengeRepository,
) ChallengeService {
	return &challengeService{
		challengeRepository: challengeRepository,
	}
}
