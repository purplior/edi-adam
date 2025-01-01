package bookmark

import (
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/shared/pagination"
)

type (
	BookmarkRepository interface {
		FindOne_ByUserIDAndAssistantID(
			ctx inner.Context,
			userID string,
			assistantID string,
		) (
			Bookmark,
			error,
		)

		FindPaginatedList_ByUserID(
			ctx inner.Context,
			userID string,
			pageRequest pagination.PaginationRequest,
		) (
			[]Bookmark,
			pagination.PaginationMeta,
			error,
		)

		InsertOne(
			ctx inner.Context,
			target Bookmark,
		) (
			registered Bookmark,
			err error,
		)

		DeleteOne_ByUserIDAndAssistantID(
			ctx inner.Context,
			userID string,
			assistantID string,
		) error
	}
)
