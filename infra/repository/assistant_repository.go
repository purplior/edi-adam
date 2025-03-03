package repository

import (
	domain "github.com/purplior/edi-adam/domain/assistant"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/infra/database/postgre"
	"gorm.io/gorm"
)

type (
	assistantRepository struct {
		postgreRepository[model.Assistant, domain.QueryOption]
	}
)

func NewAssistantRepository(
	client *postgre.Client,
) domain.AssistantRepository {
	var repo postgreRepository[model.Assistant, domain.QueryOption]

	repo.client = client
	repo.applyQueryOption = func(query *gorm.DB, opt domain.QueryOption) (*gorm.DB, error) {
		if opt.WithAuthor {
			query = query.Preload("Author")
		}
		if opt.WithCategory {
			query = query.Preload("Category")
		}
		if opt.ID > 0 {
			query = query.Where("id = ?", opt.ID)
		}
		if opt.AuthorID > 0 {
			query = query.Where("author_id = ?", opt.AuthorID)
		}
		if len(opt.CategoryID) > 0 {
			query = query.Where("category_id = ?", opt.CategoryID)
		}

		return query, nil
	}

	return &assistantRepository{
		postgreRepository: repo,
	}
}
