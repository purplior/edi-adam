package repository

import (
	domain "github.com/purplior/edi-adam/domain/bookmark"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/infra/database/postgre"
)

type (
	bookmarkRepository struct {
		postgreRepository[model.Bookmark, domain.QueryOption]
	}
)

func NewBookmarkRepository(
	client *postgre.Client,
) domain.BookmarkRepository {
	var repo postgreRepository[model.Bookmark, domain.QueryOption]

	repo.client = client
	repo.applyQueryOption = func(db *postgre.DB, opt domain.QueryOption) (*postgre.DB, error) {
		// TODO:
		return db, nil
	}

	return &bookmarkRepository{
		postgreRepository: repo,
	}
}
