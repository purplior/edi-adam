package assistant

import (
	"github.com/podossaem/podoroot/domain/shared/inner"
	"github.com/podossaem/podoroot/domain/shared/pagination"
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
	}
)
