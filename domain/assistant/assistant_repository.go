package assistant

import (
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/shared/pagination"
)

type (
	AssistantRepository interface {
		FindOne_ByID(
			ctx inner.Context,
			id string,
			joinOption AssistantJoinOption,
		) (
			Assistant,
			error,
		)

		FindOne_ByViewID(
			ctx inner.Context,
			viewID string,
			joinOption AssistantJoinOption,
		) (
			Assistant,
			error,
		)

		FindPaginatedList_ByCategoryID(
			ctx inner.Context,
			categoryID string,
			pageRequest pagination.PaginationRequest,
		) (
			[]Assistant,
			pagination.PaginationMeta,
			error,
		)

		FindPaginatedList_ByAuthorID(
			ctx inner.Context,
			authorID string,
			pageRequest pagination.PaginationRequest,
		) (
			[]Assistant,
			pagination.PaginationMeta,
			error,
		)

		InsertOne(
			ctx inner.Context,
			assistant Assistant,
		) (
			Assistant,
			error,
		)

		UpdateOne(
			ctx inner.Context,
			assistant Assistant,
		) error

		DeleteOne_ByID(
			ctx inner.Context,
			id string,
		) error
	}
)
