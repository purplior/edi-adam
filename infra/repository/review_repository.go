package repository

import (
	domain "github.com/purplior/podoroot/domain/review"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/shared/pagination"
	"github.com/purplior/podoroot/infra/database"
	"github.com/purplior/podoroot/infra/database/podosql"
	"github.com/purplior/podoroot/infra/entity"
	"github.com/purplior/podoroot/infra/repoutil"
	"github.com/purplior/podoroot/lib/dt"
	"gorm.io/gorm"
)

type (
	reviewRepository struct {
		client *podosql.Client
	}
)

func (r *reviewRepository) FindPaginatedList_ByAssistantID(
	ctx inner.Context,
	assistantID string,
	pageRequest pagination.PaginationRequest,
) (
	[]domain.Review,
	pagination.PaginationMeta,
	error,
) {
	var entities []entity.Review

	db := r.client.DBWithContext(ctx)
	pageMeta, err := repoutil.FindPaginatedList(
		db,
		&entity.Review{},
		&entities,
		pageRequest,
		repoutil.FindPaginatedListOption{
			Condition: func(db *podosql.DB) *podosql.DB {
				query := db.
					Preload("Author").
					Order(gorm.Expr("COALESCE(updated_at, created_at) DESC")).
					Where("assistant_id = ?", assistantID)

				return query
			},
		},
	)
	if err != nil {
		return nil, pagination.PaginationMeta{}, database.ToDomainError(err)
	}

	reviews := make([]domain.Review, len(entities))
	for i, entity := range entities {
		reviews[i] = entity.ToModel()
	}

	return reviews, pageMeta, nil
}

func (r *reviewRepository) FindOne_ByID(
	ctx inner.Context,
	id string,
	queryOption domain.ReviewQueryOption,
) (domain.Review, error) {
	db := r.client.DBWithContext(ctx)

	query := db.Where("id = ?", id)
	if queryOption.WithAuthor {
		query = query.Preload("Author")
	}

	var e entity.Review
	if err := query.First(&e).Error; err != nil {
		return domain.Review{}, database.ToDomainError(err)
	}

	return e.ToModel(), nil
}

func (r *reviewRepository) FindOne_ByAuthorAndAssistantID(
	ctx inner.Context,
	authorID string,
	assistantID string,
	queryOption domain.ReviewQueryOption,
) (
	domain.Review,
	error,
) {
	db := r.client.DBWithContext(ctx)

	query := db
	if queryOption.WithAuthor {
		query = query.Preload("Author")
	}

	var e entity.Review
	if err := query.
		Where("author_id = ? AND assistant_id = ?", authorID, assistantID).
		First(&e).Error; err != nil {
		return domain.Review{}, database.ToDomainError(err)
	}

	return e.ToModel(), nil
}

func (r *reviewRepository) InsertOne(
	ctx inner.Context,
	review domain.Review,
) (
	domain.Review,
	error,
) {
	db := r.client.DBWithContext(ctx)
	e := entity.MakeReview(review)

	if err := db.Create(&e).Error; err != nil {
		return domain.Review{}, database.ToDomainError(err)
	}

	return e.ToModel(), nil
}

func (r *reviewRepository) UpdateOne_ByID(
	ctx inner.Context,
	id string,
	review domain.Review,
) error {
	db := r.client.DBWithContext(ctx)
	e := entity.MakeReview(review)
	e.ID = 0

	if err := db.Where("id = ?", dt.UInt(id)).Updates(e).Error; err != nil {
		return database.ToDomainError(err)
	}

	return nil
}

func NewReviewRepository(
	client *podosql.Client,
) domain.ReviewRepository {
	return &reviewRepository{
		client: client,
	}
}
