package mission

import (
	"github.com/purplior/sbec/domain/shared/inner"
	"github.com/purplior/sbec/domain/shared/pagination"
)

type (
	MissionRepository interface {
		FindOne_ByIDAndUserID(
			ctx inner.Context,
			id string,
			userID string,
		) (
			Mission,
			error,
		)

		FindPaginatedList_OnlyPublic_ByUserID(
			ctx inner.Context,
			userID string,
			page int,
			pageSize int,
		) (
			[]Mission,
			pagination.PaginationMeta,
			error,
		)
	}
)
