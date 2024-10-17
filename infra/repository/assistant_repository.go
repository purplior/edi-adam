package repository

import (
	"github.com/podossaem/podoroot/domain/assistant"
	"github.com/podossaem/podoroot/domain/shared/context"
	"github.com/podossaem/podoroot/infra/database"
	"github.com/podossaem/podoroot/infra/database/podosql"
	"github.com/podossaem/podoroot/infra/entity"
	"gorm.io/gorm"
)

type (
	assistantRepository struct {
		client *podosql.Client
	}
)

func (r *assistantRepository) InsertOne(
	ctx context.APIContext,
	assistantForInsert assistant.Assistant,
) (
	assistant.Assistant,
	error,
) {
	e := entity.MakeAssistant(assistantForInsert)
	result := r.client.DB.WithContext(ctx).
		Create(&e)

	if result.Error != nil {
		return assistant.Assistant{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func (r *assistantRepository) FindOneByID(
	ctx context.APIContext,
	id string,
	joinOption assistant.AssistantJoinOption,
) (
	assistant.Assistant,
	error,
) {
	var e entity.Assistant

	result := r.client.DB.Where("id = ?", id).First(&e)
	if result.Error != nil {
		return assistant.Assistant{}, database.ToDomainError(result.Error)
	}

	if joinOption.WithAuthor {
		if err := r.client.DB.Model(&e).Association("Author").Find(&e.Author); err != nil {
			return assistant.Assistant{}, database.ToDomainError(err)
		}
	}

	if joinOption.WithAssister {
		var assister entity.Assister

		result := r.client.DB.Where("assistant_id = ?", e.ID).
			Order("created_at DESC").
			Limit(1).
			First(&assister)

		switch result.Error {
		case nil:
			e.Assisters = []entity.Assister{assister}
		case gorm.ErrRecordNotFound:
			e.Assisters = []entity.Assister{}
		default:
			return assistant.Assistant{}, database.ToDomainError(result.Error)
		}

		e.Assisters = []entity.Assister{assister}
	}

	return e.ToModel(), nil
}

func (r *assistantRepository) FindListByAuthorID(
	ctx context.APIContext,
	authorID string,
	joinOption assistant.AssistantJoinOption,
) (
	[]assistant.Assistant,
	error,
) {
	var entities []entity.Assistant

	query := r.client.DB.
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
