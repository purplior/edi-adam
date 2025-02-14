package repository

import (
	"fmt"
	"strings"
	"time"

	domain "github.com/purplior/sbec/domain/challenge"
	"github.com/purplior/sbec/domain/shared/inner"
	"github.com/purplior/sbec/infra/database"
	"github.com/purplior/sbec/infra/database/sqldb"
	"github.com/purplior/sbec/infra/entity"
	"github.com/purplior/sbec/lib/dt"
)

type (
	challengeRepository struct {
		client *sqldb.Client
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

func (r *challengeRepository) FindList_ByUserIDAndMissionIDs(
	ctx inner.Context,
	userID string,
	missionIDs []string,
) (
	[]domain.Challenge,
	error,
) {
	eMissionIDs := make([]int, len(missionIDs))
	for i, missionID := range missionIDs {
		eMissionIDs[i] = dt.Int(missionID)
	}

	var entities []entity.Challenge
	db := r.client.DBWithContext(ctx)
	orderStr := fmt.Sprintf("FIELD(mission_id, %s)", strings.Join(missionIDs, ","))
	if err := db.Where("mission_id IN ? AND user_id = ?", eMissionIDs, userID).
		Order(orderStr).
		Find(&entities).Error; err != nil {
		return nil, database.ToDomainError(err)
	}

	models := make([]domain.Challenge, len(entities))
	for i, entity := range entities {
		models[i] = entity.ToModel()
	}

	return models, nil
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
	client *sqldb.Client,
) domain.ChallengeRepository {
	return &challengeRepository{
		client: client,
	}
}
