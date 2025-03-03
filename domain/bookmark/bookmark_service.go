package bookmark

import (
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/exception"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	BookmarkService interface {
		Get(
			session inner.Session,
			queryOption QueryOption,
		) (
			model.Bookmark,
			error,
		)

		GetPaginatedList(
			session inner.Session,
			query pagination.PaginationQuery[QueryOption],
		) (
			[]model.Bookmark,
			pagination.PaginationMeta,
			error,
		)

		Toggle(
			session inner.Session,
			queryOption QueryOption,
			toggle bool,
		) error
	}
)

type (
	bookmarkService struct {
		bookmarkRepository BookmarkRepository
	}
)

func (s *bookmarkService) Get(
	session inner.Session,
	queryOption QueryOption,
) (
	model.Bookmark,
	error,
) {
	return s.bookmarkRepository.Read(
		session,
		queryOption,
	)
}

func (s *bookmarkService) GetPaginatedList(
	session inner.Session,
	query pagination.PaginationQuery[QueryOption],
) (
	mArr []model.Bookmark,
	meta pagination.PaginationMeta,
	err error,
) {
	return s.bookmarkRepository.ReadPaginatedList(
		session,
		query,
	)
}

func (s *bookmarkService) Toggle(
	session inner.Session,
	queryOption QueryOption,
	toggle bool,
) (err error) {
	if queryOption.UserID == 0 || queryOption.AssistantID == 0 {
		return exception.ErrBadRequest
	}

	var m model.Bookmark
	if m, err = s.bookmarkRepository.Read(
		session,
		queryOption,
	); err != nil && err != exception.ErrNoRecord {
		return err
	}

	isEmpty := m.ID == 0 || err == exception.ErrNoRecord
	err = nil

	if toggle {
		if isEmpty {
			m, err = s.bookmarkRepository.Create(
				session,
				model.Bookmark{
					UserID:      queryOption.UserID,
					AssistantID: queryOption.AssistantID,
				},
			)
		}
		return err
	}
	if isEmpty {
		return nil
	}

	return s.bookmarkRepository.Delete(
		session,
		queryOption,
	)
}

func NewBookmarkService(
	bookmarkRepository BookmarkRepository,
) BookmarkService {
	return &bookmarkService{
		bookmarkRepository: bookmarkRepository,
	}
}
