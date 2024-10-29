package repository

import (
	"github.com/podossaem/podoroot/domain/assistant"
	"github.com/podossaem/podoroot/domain/shared/inner"
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
	assistantForInsert assistant.Assistant,
) (
	assistant.Assistant,
	error,
) {
	e := entity.MakeAssistant(assistantForInsert)
	result := r.client.DBWithContext(ctx).
		Create(&e)

	if result.Error != nil {
		return assistant.Assistant{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func (r *assistantRepository) FindOneByViewID(
	ctx inner.Context,
	viewID string,
	joinOption assistant.AssistantJoinOption,
) (
	assistant.Assistant,
	error,
) {
	var e entity.Assistant
	db := r.client.DBWithContext(ctx)

	result := db.Model(&e).Where("view_id = ?", viewID).First(&e)
	if result.Error != nil {
		return assistant.Assistant{}, database.ToDomainError(result.Error)
	}

	if joinOption.WithAuthor {
		if err := db.Model(&e).Association("Author").Find(&e.Author); err != nil {
			return assistant.Assistant{}, database.ToDomainError(err)
		}
	}

	if joinOption.WithAssisters {
		if err := db.Preload("Assisters").Find(&e).Error; err != nil {
			return assistant.Assistant{}, database.ToDomainError(err)
		}
	}

	return e.ToModel(), nil
}

func (r *assistantRepository) FindListByAuthorID(
	ctx inner.Context,
	authorID string,
	joinOption assistant.AssistantJoinOption,
) (
	[]assistant.Assistant,
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
		return []assistant.Assistant{}, database.ToDomainError(result.Error)
	}

	assistants := make([]assistant.Assistant, len(entities))
	for i, entity := range entities {
		assistants[i] = entity.ToModel()
	}

	return assistants, nil
}

func NewAssistantRepository(
	client *podosql.Client,
) assistant.AssistantRepository {
	return &assistantRepository{
		client: client,
	}
}
