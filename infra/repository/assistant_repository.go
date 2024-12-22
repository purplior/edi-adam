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

func (r *assistantRepository) InsertOne(
	ctx inner.Context,
	assistantForInsert domain.Assistant,
) (
	domain.Assistant,
	error,
) {
	e := entity.MakeAssistant(assistantForInsert)
	result := r.client.DBWithContext(ctx).
		Create(&e)

	if result.Error != nil {
		return domain.Assistant{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

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

	result := db.Model(&e).Where("id = ?", id).First(&e)
	if result.Error != nil {
		return domain.Assistant{}, database.ToDomainError(result.Error)
	}

	if joinOption.WithAuthor {
		if err := db.Model(&e).Association("Author").Find(&e.Author); err != nil {
			return domain.Assistant{}, database.ToDomainError(err)
		}
	}

	if joinOption.WithAssisters {
		if err := db.Preload("Assisters").Find(&e).Error; err != nil {
			return domain.Assistant{}, database.ToDomainError(err)
		}
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

	result := db.Model(&e).Where("view_id = ?", viewID).First(&e)
	if result.Error != nil {
		return domain.Assistant{}, database.ToDomainError(result.Error)
	}

	if joinOption.WithAuthor {
		if err := db.Model(&e).Association("Author").Find(&e.Author); err != nil {
			return domain.Assistant{}, database.ToDomainError(err)
		}
	}

	if joinOption.WithAssisters {
		if err := db.Preload("Assisters").Find(&e).Error; err != nil {
			return domain.Assistant{}, database.ToDomainError(err)
		}
	}

	return e.ToModel(), nil
}

func (r *assistantRepository) FindList_ByCategoryAlias(
	ctx inner.Context,
	categoryAlias string,
	joinOption domain.AssistantJoinOption,
) (
	[]domain.Assistant,
	error,
) {
	db := r.client.DBWithContext(ctx)

	var categoryEntity entity.Category
	if result := db.
		Model(&categoryEntity).
		Where("alias = ?", categoryAlias).
		First(&categoryEntity); result.Error != nil {
		return nil, database.ToDomainError(result.Error)
	}

	var entities []entity.Assistant
	query := db.
		Where("category_id = ? AND is_public = ?", categoryEntity.ID, true).
		Order("created_at asc")
	if joinOption.WithAuthor {
		query = query.Preload("Author")
	}

	if result := query.Find(&entities); result.Error != nil {
		return nil, database.ToDomainError(result.Error)
	}

	assistants := make([]domain.Assistant, len(entities))
	for i, entity := range entities {
		entity.Category = categoryEntity
		assistants[i] = entity.ToModel()
	}

	return assistants, nil
}

func (r *assistantRepository) FindPaginatedList_ByAuthorID(
	ctx inner.Context,
	authorID string,
	page int,
	pageSize int,
) (
	[]domain.Assistant,
	pagination.PaginationMeta,
	error,
) {
	var entities []entity.Assistant

	db := r.client.DBWithContext(ctx)
	meta, err := repoutil.FindPaginatedList(
		db,
		&entity.Assistant{},
		&entities,
		page,
		pageSize,
		func(db *podosql.DB) *podosql.DB {
			rdb := db
			if len(authorID) > 0 {
				rdb = db.Where("author_id = ?", authorID)
			}

			return rdb.Order("created_at DESC").Preload("Category")
		},
		func(db *podosql.DB) *podosql.DB {
			if len(authorID) > 0 {
				return db.Where("author_id = ?", authorID)
			}

			return db
		},
	)
	if err != nil {
		return nil, meta, database.ToDomainError(err)
	}

	assistants := make([]domain.Assistant, len(entities))
	for i, entity := range entities {
		assistants[i] = entity.ToModel()
	}

	return assistants, meta, nil
}

func (r *assistantRepository) UpdateOne(
	ctx inner.Context,
	assistant domain.Assistant,
) error {
	e := entity.MakeAssistant(assistant)
	db := r.client.DBWithContext(ctx)

	result := db.Save(e)
	if result.Error != nil {
		return database.ToDomainError(result.Error)
	}

	return nil
}

func (r *assistantRepository) DeleteOne_ByID(
	ctx inner.Context,
	id string,
) error {
	db := r.client.DBWithContext(ctx)

	result := db.Delete(&entity.Assistant{
		ID: dt.UInt(id),
	})
	if result.Error != nil {
		return database.ToDomainError(result.Error)
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
