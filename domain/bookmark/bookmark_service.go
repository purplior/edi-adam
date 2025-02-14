package bookmark

import (
	"github.com/purplior/sbec/domain/shared/exception"
	"github.com/purplior/sbec/domain/shared/inner"
	"github.com/purplior/sbec/domain/shared/pagination"
)

type (
	BookmarkService interface {
		GetPaginatedList_ByUserID(
			ctx inner.Context,
			userID string,
			pageRequest pagination.PaginationRequest,
		) (
			[]Bookmark,
			pagination.PaginationMeta,
			error,
		)

		GetOne_ByUserIDAndAssistantID(
			ctx inner.Context,
			userID string,
			assistantID string,
		) (
			Bookmark,
			error,
		)

		ToggleOne(
			ctx inner.Context,
			userID string,
			assistantID string,
			toggle bool,
		) (
			err error,
		)
	}
)

type (
	bookmarkService struct {
		bookmarkRepository BookmarkRepository
	}
)

func (s *bookmarkService) GetOne_ByUserIDAndAssistantID(
	ctx inner.Context,
	userID string,
	assistantID string,
) (
	Bookmark,
	error,
) {
	return s.bookmarkRepository.FindOne_ByUserIDAndAssistantID(
		ctx,
		userID,
		assistantID,
	)
}

func (s *bookmarkService) GetPaginatedList_ByUserID(
	ctx inner.Context,
	userID string,
	pageRequest pagination.PaginationRequest,
) (
	[]Bookmark,
	pagination.PaginationMeta,
	error,
) {
	return s.bookmarkRepository.FindPaginatedList_ByUserID(
		ctx,
		userID,
		pageRequest,
	)
}

func (s *bookmarkService) ToggleOne(
	ctx inner.Context,
	userID string,
	assistantID string,
	toggle bool,
) error {
	bookmark, err := s.bookmarkRepository.FindOne_ByUserIDAndAssistantID(
		ctx,
		userID,
		assistantID,
	)
	if err != nil && err != exception.ErrNoRecord {
		return err
	}

	isEmpty := len(bookmark.ID) == 0 || err == exception.ErrNoRecord
	err = nil

	if toggle {
		if isEmpty {
			_, err = s.bookmarkRepository.InsertOne(
				ctx,
				Bookmark{
					UserID:      userID,
					AssistantID: assistantID,
				},
			)
		}
		return err
	}

	if isEmpty {
		return nil
	}

	return s.bookmarkRepository.DeleteOne_ByUserIDAndAssistantID(
		ctx,
		userID,
		assistantID,
	)
}

func NewBookmarkService(
	bookmarkRepository BookmarkRepository,
) BookmarkService {
	return &bookmarkService{
		bookmarkRepository: bookmarkRepository,
	}
}
