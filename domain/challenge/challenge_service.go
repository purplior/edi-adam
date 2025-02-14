package challenge

import (
	"github.com/purplior/sbec/domain/shared/exception"
	"github.com/purplior/sbec/domain/shared/inner"
	"github.com/purplior/sbec/lib/mydate"
)

type (
	ChallengeService interface {
		GetList_ByUserIDAndMissionIDs(
			ctx inner.Context,
			userID string,
			missionIDs []string,
		) (
			[]Challenge,
			error,
		)

		PatchOne_ReceivedStatus(
			ctx inner.Context,
			id string,
		) error

		AchieveOne_ByUserAndMission(
			ctx inner.Context,
			userID string,
			missionID string,
		) error
	}
)

type (
	challengeService struct {
		challengeRepository ChallengeRepository
		cm                  inner.ContextManager
	}
)

func (s *challengeService) GetList_ByUserIDAndMissionIDs(
	ctx inner.Context,
	userID string,
	missionIDs []string,
) (
	[]Challenge,
	error,
) {
	return s.challengeRepository.FindList_ByUserIDAndMissionIDs(
		ctx,
		userID,
		missionIDs,
	)
}

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

func (s *challengeService) AchieveOne_ByUserAndMission(
	ctx inner.Context,
	userID string,
	missionID string,
) error {
	challenge, err := s.challengeRepository.FindOne_ByUserIDAndMissionID(
		ctx,
		userID,
		missionID,
	)

	isNoRecord := err == exception.ErrNoRecord
	if err != nil && !isNoRecord {
		return err
	}

	if isNoRecord {
		_, err = s.challengeRepository.InsertOne(
			ctx,
			Challenge{
				UserID:     userID,
				MissionID:  missionID,
				IsAchieved: true,
				IsReceived: false,
			},
		)

		return err
	}

	// 이미 달성했거나 수령한 경우 무시함.
	if challenge.IsAchieved || challenge.IsReceived {
		return nil
	}

	return s.challengeRepository.UpdateOne_AchievedStatus_ByID(
		ctx,
		challenge.ID,
		true,
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
