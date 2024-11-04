package repository

import (
	domain "github.com/podossaem/podoroot/domain/assistant"
	"github.com/podossaem/podoroot/domain/shared/inner"
	"github.com/podossaem/podoroot/domain/shared/pagination"
	"github.com/podossaem/podoroot/infra/database"
	"github.com/podossaem/podoroot/infra/database/podosql"
	"github.com/podossaem/podoroot/infra/entity"
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

func (r *assistantRepository) FindList_ByAuthorID(
	ctx inner.Context,
	authorID string,
	joinOption domain.AssistantJoinOption,
) (
	[]domain.Assistant,
	error,
) {
	var entities []entity.Assistant

	db := r.client.DBWithContext(ctx)
	query := db.
		Where("author_id = ? AND is_public = ?", authorID, true).
		Order("created_at asc")
	if joinOption.WithAuthor {
		query = query.Preload("Author")
	}

	result := query.Find(&entities)
	if result.Error != nil {
		return []domain.Assistant{}, database.ToDomainError(result.Error)
	}

	assistants := make([]domain.Assistant, len(entities))
	for i, entity := range entities {
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
	var totalCount int64
	db := r.client.DBWithContext(ctx)

	if err := db.Model(&entity.Assistant{}).
		Count(&totalCount).Error; err != nil {
		return nil, pagination.PaginationMeta{}, database.ToDomainError(err)
	}

	if len(authorID) > 0 {
		db = db.Where("author_id = ?", authorID)
	}

	offset := (page - 1) * pageSize
	if err := db.Offset(offset).
		Limit(pageSize).
		Find(&entities).Error; err != nil {
		return nil, pagination.PaginationMeta{}, database.ToDomainError(err)
	}

	meta := pagination.PaginationMeta{
		Page:      page,
		Size:      pageSize,
		TotalPage: int((totalCount + int64(pageSize) - 1) / int64(pageSize)),
	}

	assistants := make([]domain.Assistant, len(entities))
	for i, entity := range entities {
		assistants[i] = entity.ToModel()
	}

	return assistants, meta, nil
}

func NewAssistantRepository(
	client *podosql.Client,
) domain.AssistantRepository {
	return &assistantRepository{
		client: client,
	}
}
