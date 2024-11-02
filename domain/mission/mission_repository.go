package mission

import (
	"github.com/podossaem/podoroot/domain/shared/inner"
	"github.com/podossaem/podoroot/domain/shared/pagination"
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

		FindPaginatedList_ByUserID(
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
