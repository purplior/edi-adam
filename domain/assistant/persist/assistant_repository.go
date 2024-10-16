package persist

import (
	domain "github.com/podossaem/podoroot/domain/assistant"
	"github.com/podossaem/podoroot/domain/context"
	"github.com/podossaem/podoroot/infra/database"
	"github.com/podossaem/podoroot/infra/database/podosql"
)

type (
	assistantRepository struct {
		client *podosql.Client
	}
)

func (r *assistantRepository) InsertOne(
	ctx context.APIContext,
	assistant domain.Assistant,
) (
	domain.Assistant,
	error,
) {
	entity := MakeAssistant(assistant)
	result := r.client.DB.WithContext(ctx).
		Create(&entity)

	if result.Error != nil {
		return domain.Assistant{}, database.ToDomainError(result.Error)
	}

	return entity.ToModel(), nil
}

func (r *assistantRepository) FindListByAuthorID(
	ctx context.APIContext,
	authorID string,
) (
	[]domain.Assistant,
	error,
) {
	var entities []Assistant
	result := r.client.DB.
		Where("author_id = ? AND is_public = ?", authorID, true).
		Order("created_at asc").
		Find(&entities)
	if result.Error != nil {
		return []domain.Assistant{}, database.ToDomainError(result.Error)
	}

	assistants := make([]domain.Assistant, len(entities))
	for i, entity := range entities {
		assistants[i] = entity.ToModel()
	}

	return assistants, nil
}

func NewAssistantRepository(
	client *podosql.Client,
) domain.AssistantRepository {
	client.AddModel(&Assistant{})

	return &assistantRepository{
		client: client,
	}
}
