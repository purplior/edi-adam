package repository

import (
	domain "github.com/purplior/podoroot/domain/assistant"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/shared/pagination"
	"github.com/purplior/podoroot/infra/database"
	"github.com/purplior/podoroot/infra/database/podosql"
	"github.com/purplior/podoroot/infra/entity"
	"github.com/purplior/podoroot/infra/repoutil"
	"github.com/purplior/podoroot/lib/dt"
)

type (
	assistantRepository struct {
		client *podosql.Client
	}
)

func (r *assistantRepository) FindOne_ByID(
	ctx inner.Context,
	id string,
	joinOption domain.AssistantJoinOption,
) (
	domain.Assistant,
	error,
) {
	var e entity.Assistant
	db := r.client.DBWithContext(ctx)
	eID := dt.UInt(id)

	query := db
	if joinOption.WithAuthor {
		query = query.Preload("Author")
	}
	if joinOption.WithCategory {
		query = query.Preload("Category")
	}

	if err := query.
		Where("id = ?", eID).
		First(&e).Error; err != nil {
		return domain.Assistant{}, database.ToDomainError(err)
	}

	return e.ToModel(), nil
}

func (r *assistantRepository) FindOne_ByViewID(
	ctx inner.Context,
	viewID string,
	joinOption domain.AssistantJoinOption,
) (
	domain.Assistant,
	error,
) {
	var e entity.Assistant
	db := r.client.DBWithContext(ctx)

	query := db
	if joinOption.WithAuthor {
		query = query.Preload("Author")
	}
	if joinOption.WithCategory {
		query = query.Preload("Category")
	}

	if err := query.
		Where("view_id = ?", viewID).
		First(&e).Error; err != nil {
		return domain.Assistant{}, database.ToDomainError(err)
	}

	return e.ToModel(), nil
}

func (r *assistantRepository) FindPaginatedList_ByCategoryID(
	ctx inner.Context,
	categoryID string,
	pageRequest pagination.PaginationRequest,
) (
	[]domain.Assistant,
	pagination.PaginationMeta,
	error,
) {
	var entities []entity.Assistant

	db := r.client.DBWithContext(ctx)
	pageMeta, err := repoutil.FindPaginatedList(
		db,
		&entity.Assistant{},
		&entities,
		pageRequest,
		repoutil.FindPaginatedListOption{
			Condition: func(db *podosql.DB) *podosql.DB {
				return db.
					Preload("Author").
					Order("created_at DESC").
					Where("category_id = ?", categoryID)
			},
		},
	)
	if err != nil {
		return nil, pagination.PaginationMeta{}, database.ToDomainError(err)
	}

	assistants := make([]domain.Assistant, len(entities))
	for i, entity := range entities {
		assistants[i] = entity.ToModel()
	}

	return assistants, pageMeta, nil
}

func (r *assistantRepository) FindPaginatedList_ByAuthorID(
	ctx inner.Context,
	authorID string,
	pageRequest pagination.PaginationRequest,
) (
	[]domain.Assistant,
	pagination.PaginationMeta,
	error,
) {
	var entities []entity.Assistant

	db := r.client.DBWithContext(ctx)
	pageMeta, err := repoutil.FindPaginatedList(
		db,
		&entity.Assistant{},
		&entities,
		pageRequest,
		repoutil.FindPaginatedListOption{
			Condition: func(db *podosql.DB) *podosql.DB {
				return db.
					Preload("Category").
					Order("created_at DESC").
					Where("author_id = ?", authorID)
			},
		},
	)
	if err != nil {
		return nil, pagination.PaginationMeta{}, database.ToDomainError(err)
	}

	assistants := make([]domain.Assistant, len(entities))
	for i, entity := range entities {
		assistants[i] = entity.ToModel()
	}

	return assistants, pageMeta, nil
}

func (r *assistantRepository) InsertOne(
	ctx inner.Context,
	assistantForInsert domain.Assistant,
) (
	domain.Assistant,
	error,
) {
	db := r.client.DBWithContext(ctx)
	e := entity.MakeAssistant(assistantForInsert)

	if err := db.Create(&e).Error; err != nil {
		return domain.Assistant{}, database.ToDomainError(err)
	}

	return e.ToModel(), nil
}

func (r *assistantRepository) UpdateOne(
	ctx inner.Context,
	assistant domain.Assistant,
) error {
	db := r.client.DBWithContext(ctx)
	e := entity.MakeAssistant(assistant)

	if err := db.Save(&e).Error; err != nil {
		return database.ToDomainError(err)
	}

	return nil
}

func (r *assistantRepository) DeleteOne_ByID(
	ctx inner.Context,
	id string,
) error {
	db := r.client.DBWithContext(ctx)
	eID := dt.UInt(id)

	if err := db.
		Where("id = ?", eID).
		Delete(&entity.Assistant{}).Error; err != nil {
		return database.ToDomainError(err)
	}

	return nil
}

func NewAssistantRepository(
	client *podosql.Client,
) domain.AssistantRepository {
	return &assistantRepository{
		client: client,
	}
}
