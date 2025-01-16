package repository

import (
	domain "github.com/purplior/podoroot/domain/review"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/infra/database"
	"github.com/purplior/podoroot/infra/database/podosql"
	"github.com/purplior/podoroot/infra/entity"
)

type (
	reviewRepository struct {
		client *podosql.Client
	}
)

func (r *reviewRepository) FindOne_ByAuthorAndAssistantID(
	ctx inner.Context,
	authorID string,
	assistantID string,
	joinOption domain.ReviewJoinOption,
) (
	domain.Review,
	error,
) {
	db := r.client.DBWithContext(ctx)

	query := db
	if joinOption.WithAuthor {
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

func (r *reviewRepository) UpdateOne(
	ctx inner.Context,
	review domain.Review,
) error {
	db := r.client.DBWithContext(ctx)
	e := entity.MakeReview(review)

	if err := db.Save(&e).Error; err != nil {
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
