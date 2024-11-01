package challenge

import "github.com/podossaem/podoroot/domain/shared/inner"

type (
	ChallengeService interface {
		GetPaginatedInfoListByUserID(
			ctx inner.Context,
			userId string,
			limit int,
			offset int,
		) (
			[]ChallengeInfo,
			error,
		)
	}
)

type (
	challengeService struct {
		challengeRepository ChallengeRepository
	}
)

func (s *challengeService) GetPaginatedInfoListByUserID(
	ctx inner.Context,
	userId string,
	limit int,
	offset int,
) (
	[]ChallengeInfo,
	error,
) {
	challenges, err := s.challengeRepository.FindPaginatedListByUserID(
		ctx,
		userId,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	challengeInfos := make([]ChallengeInfo, len(challenges))
	for i, challenge := range challenges {
		challengeInfos[i] = challenge.ToInfo()
	}

	return challengeInfos, nil
}

func NewChallengeService(
	challengeRepository ChallengeRepository,
) ChallengeService {
	return &challengeService{
		challengeRepository: challengeRepository,
	}
}
