package repository

import (
	domain "github.com/purplior/edi-adam/domain/review"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/infra/database/postgre"
)

type (
	reviewRepository struct {
		postgreRepository[model.Review, domain.QueryOption]
	}
)

func NewReviewRepository(
	client *postgre.Client,
) domain.ReviewRepository {
	var repo postgreRepository[model.Review, domain.QueryOption]

	repo.client = client
	repo.applyQueryOption = func(db *postgre.DB, opt domain.QueryOption) (*postgre.DB, error) {
		return db, nil
	}

	return &reviewRepository{
		postgreRepository: repo,
	}
}
