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
		Preload("Mission").
		First(&e)

	if result.Error != nil {
		return domain.Challenge{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func (r *challengeRepository) FindPaginatedList_ByUserID(
	ctx inner.Context,
	userId string,
	limit int,
	offset int,
) (
	[]domain.Challenge,
	error,
) {
	var entities []entity.Challenge

	db := r.client.DBWithContext(ctx)
	if err := db.Preload("Mission").
		Limit(limit).
		Offset(offset).
		Find(&entities).
		Error; err != nil {
		return nil, err
	}

	challenges := make([]domain.Challenge, len(entities))
	for i, entity := range entities {
		challenges[i] = entity.ToModel()
	}

	return challenges, nil
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

func NewChallengeRepository(
	client *podosql.Client,
) domain.ChallengeRepository {
	return &challengeRepository{
		client: client,
	}
}
