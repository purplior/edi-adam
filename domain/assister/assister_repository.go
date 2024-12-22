package assister

import (
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/domain/shared/pagination"
)

type (
	AssisterRepository interface {
		FindOne_ByID(
			ctx inner.Context,
			id string,
		) (
			Assister,
			error,
		)

		FindPaginatedList_ByAssistantID(
			ctx inner.Context,
			assistantID string,
			page int,
			pageSize int,
		) (
			[]Assister,
			pagination.PaginationMeta,
			error,
		)

		UpdateOne(
			ctx inner.Context,
			assister Assister,
		) error

		InsertOne(
			ctx inner.Context,
			assister Assister,
		) (
			Assister,
			error,
		)

		DeleteAll_ByIDs(
			ctx inner.Context,
			ids []string,
		) error
	}
)
