package repoutil

import (
	"github.com/purplior/podoroot/domain/shared/pagination"
	"github.com/purplior/podoroot/infra/database/podosql"
)

type (
	FindPaginatedListOption struct {
		Condition func(db *podosql.DB) *podosql.DB
	}
)

func FindPaginatedList(
	db *podosql.DB,
	entity interface{},
	entities interface{},
	pageRequest pagination.PaginationRequest,
	option FindPaginatedListOption,
) (
	pagination.PaginationMeta,
	error,
) {
	page := pageRequest.Page
	pageSize := pageRequest.Size
	useCountQuery := pageRequest.TotalCount == 0

	var totalCount int64 = int64(pageRequest.TotalCount)

	if useCountQuery {
		if err := option.
			Condition(db.Model(entity)).
			Count(&totalCount).Error; err != nil {
			return pagination.PaginationMeta{}, err
		}
	}

	offset := (page - 1) * pageSize
	query := option.Condition(db).
		Offset(offset).
		Limit(pageSize)

	var err error = query.Find(entities).Error

	if err != nil {
		return pagination.PaginationMeta{}, err
	}

	var totalPage int = 0
	if pageSize > 0 {
		totalPage = int((totalCount + int64(pageSize) - 1) / int64(pageSize))
	}

	meta := pagination.PaginationMeta{
		Page:       page,
		Size:       pageSize,
		TotalCount: int(totalCount),
		TotalPage:  totalPage,
	}

	return meta, nil
}
