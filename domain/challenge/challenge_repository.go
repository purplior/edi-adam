package challenge

import (
	"time"

	"github.com/purplior/podoroot/domain/shared/inner"
)

type (
	ChallengeRepository interface {
		InsertOne(
			ctx inner.Context,
			challenge Challenge,
		) (
			Challenge,
			error,
		)

		FindOne_ByID(
			ctx inner.Context,
			id string,
		) (
			Challenge,
			error,
		)

		FindOne_ByUserIDAndMissionID(
			ctx inner.Context,
			userID string,
			missionID string,
		) (
			Challenge,
			error,
		)

		FindList_ByUserIDAndMissionIDs(
			ctx inner.Context,
			userID string,
			missionIDs []string,
		) (
			[]Challenge,
			error,
		)

		UpdateOne_ReceivedStatus_ByID(
			ctx inner.Context,
			id string,
			isReceived bool,
			receivedAt time.Time,
		) error

		UpdateOne_AchievedStatus_ByID(
			ctx inner.Context,
			id string,
			isAchieved bool,
		) error
	}
)
