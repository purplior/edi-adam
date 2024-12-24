package repository

import (
	domain "github.com/purplior/podoroot/domain/category"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/infra/database"
	"github.com/purplior/podoroot/infra/database/podosql"
	"github.com/purplior/podoroot/infra/entity"
	"github.com/purplior/podoroot/lib/dt"
)

type (
	categoryRepository struct {
		client *podosql.Client
	}
)

func (r *categoryRepository) FindListByIDs(
	ctx inner.Context,
	ids []string,
) (
	[]domain.Category,
	error,
) {
	db := r.client.DBWithContext(ctx)
	eIDs := make([]int, len(ids))
	for i, id := range ids {
		eIDs[i] = dt.Int(id)
	}

	var entities []entity.Category
	if err := db.Find(&entities, eIDs).Error; err != nil {
		return nil, database.ToDomainError(err)
	}

	models := make([]domain.Category, len(entities))
	for i, entity := range entities {
		models[i] = entity.ToModel()
	}

	return models, nil
}

func NewCategoryRepository(
	client *podosql.Client,
) domain.CategoryRepository {
	return &categoryRepository{
		client: client,
	}
}
