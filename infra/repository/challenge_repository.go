package repository

import (
	"time"

	domain "github.com/podossaem/podoroot/domain/challenge"
	"github.com/podossaem/podoroot/domain/shared/inner"
	"github.com/podossaem/podoroot/infra/database"
	"github.com/podossaem/podoroot/infra/database/podosql"
	"github.com/podossaem/podoroot/infra/entity"
)

type (
	challengeRepository struct {
		client *podosql.Client
	}
)

func (r *challengeRepository) InsertOne(
	ctx inner.Context,
	model domain.Challenge,
) (
	domain.Challenge,
	error,
) {
	e := entity.MakeChallenge(model)
	db := r.client.DBWithContext(ctx)
	result := db.Create(&e)

	if result.Error != nil {
		return domain.Challenge{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func (r *challengeRepository) FindOne_ByID(
	ctx inner.Context,
	id string,
) (
	domain.Challenge,
	error,
) {
	var e entity.Challenge

	db := r.client.DBWithContext(ctx)
	result := db.
		Model(&e).
		Where("id = ?", id).
		First(&e)

	if result.Error != nil {
		return domain.Challenge{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func (r *challengeRepository) FindOne_ByUserIDAndMissionID(
	ctx inner.Context,
	userID string,
	missionID string,
) (
	domain.Challenge,
	error,
) {
	var e entity.Challenge

	db := r.client.DBWithContext(ctx)
	result := db.
		Model(&e).
		Where("user_id = ?", userID).
		Where("mission_id = ?", missionID).
		First(&e)

	if result.Error != nil {
		return domain.Challenge{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func (r *challengeRepository) UpdateOne_ReceivedStatus_ByID(
	ctx inner.Context,
	id string,
	isReceived bool,
	receivedAt time.Time,
) error {
	db := r.client.DBWithContext(ctx)

	result := db.
		Model(&entity.Challenge{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_received": isReceived,
			"received_at": receivedAt,
		})

	if result.Error != nil {
		return database.ToDomainError(result.Error)
	}

	return nil
}

func (r *challengeRepository) UpdateOne_AchievedStatus_ByID(
	ctx inner.Context,
	id string,
	isAchieved bool,
) error {
	db := r.client.DBWithContext(ctx)

	result := db.
		Model(&entity.Challenge{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"is_achieved": isAchieved,
		})

	if result.Error != nil {
		return database.ToDomainError(result.Error)
	}

	return nil
}

func NewChallengeRepository(
	client *podosql.Client,
) domain.ChallengeRepository {
	return &challengeRepository{
		client: client,
	}
}
