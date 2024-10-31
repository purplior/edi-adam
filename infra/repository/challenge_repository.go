package repository

import (
	"fmt"

	domain "github.com/podossaem/podoroot/domain/challenge"
	"github.com/podossaem/podoroot/domain/shared/inner"
	"github.com/podossaem/podoroot/infra/database/podosql"
	"github.com/podossaem/podoroot/infra/entity"
)

type (
	challengeRepository struct {
		client *podosql.Client
	}
)

func (r *challengeRepository) FindPaginatedListByUserID(
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

	fmt.Println(len(entities))

	challenges := make([]domain.Challenge, len(entities))
	for i, entity := range entities {
		challenges[i] = entity.ToModel()
	}

	return challenges, nil
}

func NewChallengeRepository(
	client *podosql.Client,
) domain.ChallengeRepository {
	return &challengeRepository{
		client: client,
	}
}
