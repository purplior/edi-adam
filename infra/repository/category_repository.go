package repository

import (
	domain "github.com/purplior/edi-adam/domain/category"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/infra/database/postgre"
)

type (
	categoryRepository struct {
		postgreRepository[model.Category, domain.QueryOption]
	}
)

func NewCategoryRepository(
	client *postgre.Client,
) domain.CategoryRepository {
	var repo postgreRepository[model.Category, domain.QueryOption]

	repo.client = client
	repo.applyQueryOption = func(db *postgre.DB, opt domain.QueryOption) (*postgre.DB, error) {
		query := db
		if len(opt.ID) > 0 {
			query = query.Where("id = ?", opt.ID)
		}

		return query, nil
	}

	return &categoryRepository{
		postgreRepository: repo,
	}
}
