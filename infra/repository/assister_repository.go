package repository

import (
	domain "github.com/purplior/podoroot/domain/assister"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/shared/pagination"
	"github.com/purplior/podoroot/infra/database"
	"github.com/purplior/podoroot/infra/database/podosql"
	"github.com/purplior/podoroot/infra/entity"
	"github.com/purplior/podoroot/infra/repoutil"
)

type (
	assisterRepository struct {
		client *podosql.Client
	}
)

func (r *assisterRepository) FindOne_ByID(
	ctx inner.Context,
	id string,
) (
	domain.Assister,
	error,
) {
	db := r.client.DBWithContext(ctx)

	var e entity.Assister
	result := db.Model(&e).Where("id = ?", id).First(&e)
	if result.Error != nil {
		return domain.Assister{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func (r *assisterRepository) FindPaginatedList_ByAssistantID(
	ctx inner.Context,
	assistantID string,
	page int,
	pageSize int,
) (
	[]domain.Assister,
	pagination.PaginationMeta,
	error,
) {
	var entities []entity.Assister

	db := r.client.DBWithContext(ctx)
	meta, err := repoutil.FindPaginatedList(
		db,
		&entity.Assister{},
		&entities,
		page,
		pageSize,
		func(db *podosql.DB) *podosql.DB {
			return db.Where("assistant_id = ?", assistantID).Order("created_at DESC")
		},
		func(db *podosql.DB) *podosql.DB {
			return db.Where("assistant_id = ?", assistantID)
		},
	)
	if err != nil {
		return nil, meta, database.ToDomainError(err)
	}

	assisters := make([]domain.Assister, len(entities))
	for i, entity := range entities {
		assisters[i] = entity.ToModel()
	}

	return assisters, meta, nil
}

func (r *assisterRepository) UpdateOne(
	ctx inner.Context,
	assister domain.Assister,
) error {
	e := entity.MakeAssister(assister)
	db := r.client.DBWithContext(ctx)

	result := db.Save(e)
	if result.Error != nil {
		return database.ToDomainError(result.Error)
	}

	return nil
}

func (r *assisterRepository) InsertOne(
	ctx inner.Context,
	assister domain.Assister,
) (
	domain.Assister,
	error,
) {
	db := r.client.DBWithContext(ctx)
	e := entity.MakeAssister(assister)

	result := db.Create(&e)

	if result.Error != nil {
		return domain.Assister{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func NewAssisterRepository(
	client *podosql.Client,
) domain.AssisterRepository {
	return &assisterRepository{
		client: client,
	}
}
