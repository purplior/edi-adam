package assistant

import (
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/shared/pagination"
)

type (
	AssistantRepository interface {
		InsertOne(
			ctx inner.Context,
			assistant Assistant,
		) (
			Assistant,
			error,
		)

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

		FindList_ByAuthorID(
			ctx inner.Context,
			authorID string,
			joinOption AssistantJoinOption,
		) (
			[]Assistant,
			error,
		)

		FindList_ByCategoryAlias(
			ctx inner.Context,
			categoryAlias string,
			joinOption AssistantJoinOption,
		) (
			[]Assistant,
			error,
		)

		FindPaginatedList_ByAuthorID(
			ctx inner.Context,
			authorID string,
			page int,
			pageSize int,
		) (
			[]Assistant,
			pagination.PaginationMeta,
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
