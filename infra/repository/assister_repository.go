package repository

import (
	domain "github.com/podossaem/podoroot/domain/assister"
	"github.com/podossaem/podoroot/domain/shared/inner"
	"github.com/podossaem/podoroot/infra/database"
	"github.com/podossaem/podoroot/infra/database/podosql"
	"github.com/podossaem/podoroot/infra/entity"
)

type (
	assisterRepository struct {
		client *podosql.Client
	}
)

func (r *assisterRepository) FindOneByID(
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

func NewAssisterRepository(
	client *podosql.Client,
) domain.AssisterRepository {
	return &assisterRepository{
		client: client,
	}
}
