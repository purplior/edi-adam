package repository

import (
	domain "github.com/purplior/podoroot/domain/bookmark"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/infra/database"
	"github.com/purplior/podoroot/infra/database/podosql"
	"github.com/purplior/podoroot/infra/entity"
	"github.com/purplior/podoroot/lib/dt"
)

type (
	bookmarkRepository struct {
		client *podosql.Client
	}
)

func (r *bookmarkRepository) FindOne_ByUserIDAndAssistantID(
	ctx inner.Context,
	userID string,
	assistantID string,
) (
	domain.Bookmark,
	error,
) {
	db := r.client.DBWithContext(ctx)
	eUserID := dt.UInt(userID)
	eAssistantID := dt.UInt(assistantID)

	var e entity.Bookmark
	err := db.
		Where("user_id = ? AND assistant_id = ?", eUserID, eAssistantID).
		First(&e).
		Error

	if err != nil {
		return domain.Bookmark{}, database.ToDomainError(err)
	}

	return e.ToModel(), nil
}

func (r *bookmarkRepository) InsertOne(
	ctx inner.Context,
	target domain.Bookmark,
) (
	domain.Bookmark,
	error,
) {
	e := entity.MakeBookmark(target)
	if err := r.client.DBWithContext(ctx).Create(&e).Error; err != nil {
		return domain.Bookmark{}, database.ToDomainError(err)
	}

	return e.ToModel(), nil
}

func (r *bookmarkRepository) DeleteOne_ByUserIDAndAssistantID(
	ctx inner.Context,
	userID string,
	assistantID string,
) error {
	db := r.client.DBWithContext(ctx)
	eUserID := dt.UInt(userID)
	eAssistantID := dt.UInt(assistantID)

	err := db.
		Where("user_id = ? AND assistant_id = ?", eUserID, eAssistantID).
		Delete(&entity.Bookmark{}).
		Error

	if err != nil {
		return database.ToDomainError(err)
	}

	return nil
}

func NewBookmarkRepository(
	client *podosql.Client,
) domain.BookmarkRepository {
	return &bookmarkRepository{
		client: client,
	}
}
