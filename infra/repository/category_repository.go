package repository

import (
	"fmt"
	"strings"

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

func (r *categoryRepository) FindOne_ByAlias(
	ctx inner.Context,
	alias string,
) (
	domain.Category,
	error,
) {
	db := r.client.DBWithContext(ctx)

	var entity entity.Category
	if err := db.Where("alias = ?", alias).First(&entity).Error; err != nil {
		return domain.Category{}, database.ToDomainError(err)
	}

	return entity.ToModel(), nil
}

func (r *categoryRepository) FindList_ByIDs(
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
	orderStr := fmt.Sprintf("FIELD(id, %s)", strings.Join(ids, ","))
	if err := db.Where("id IN ?", eIDs).
		Order(orderStr).
		Find(&entities).Error; err != nil {
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
