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
		print(result.Error)
		return domain.Assistant{}, database.ToDomainError(result.Error)
	}

	return entity.ToModel(), nil
}

func NewAssistantRepository(
	client *podosql.Client,
) domain.AssistantRepository {
	client.AddModel(&Assistant{})

	return &assistantRepository{
		client: client,
	}
}
