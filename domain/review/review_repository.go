package review

import (
	"github.com/purplior/sbec/domain/shared/inner"
	"github.com/purplior/sbec/domain/shared/pagination"
)

type (
	ReviewRepository interface {
		FindOne_ByID(
			ctx inner.Context,
			id string,
			queryOption ReviewQueryOption,
		) (
			Review,
			error,
		)

		FindOne_ByAuthorAndAssistantID(
			ctx inner.Context,
			authorID string,
			assistantID string,
			queryOption ReviewQueryOption,
		) (
			Review,
			error,
		)

		FindPaginatedList_ByAssistantID(
			ctx inner.Context,
			assistantID string,
			pageRequest pagination.PaginationRequest,
		) (
			[]Review,
			pagination.PaginationMeta,
			error,
		)

		InsertOne(
			ctx inner.Context,
			review Review,
		) (
			Review,
			error,
		)

		UpdateOne_ByID(
			ctx inner.Context,
			id string,
			review Review,
		) error
	}
)
