package repoutil

import (
	"github.com/purplior/podoroot/domain/shared/pagination"
	"github.com/purplior/podoroot/infra/database/podosql"
)

type (
	FindPaginatedListOption struct {
		Condition   func(db *podosql.DB) *podosql.DB
		Association func(db *podosql.DB) *podosql.DB
	}
)

func FindPaginatedList(
	db *podosql.DB,
	entities interface{},
	pageRequest pagination.PaginationRequest,
	option FindPaginatedListOption,
) (
	pagination.PaginationMeta,
	error,
) {
	page := pageRequest.Page
	pageSize := pageRequest.Size
	useCountQuery := pageRequest.TotalCount != 0

	var totalCount int64 = int64(pageRequest.TotalCount)

	if useCountQuery {
		if err := option.Condition(db).
			Count(&totalCount).Error; err != nil {
			return pagination.PaginationMeta{}, err
		}
	}

	offset := (page - 1) * pageSize
	if err := option.Condition(option.Association(db)).
		Offset(offset).
		Limit(pageSize).
		Find(entities).Error; err != nil {
		return pagination.PaginationMeta{}, err
	}

	meta := pagination.PaginationMeta{
		Page:       page,
		Size:       pageSize,
		TotalCount: int(totalCount),
		TotalPage:  int((totalCount + int64(pageSize) - 1) / int64(pageSize)),
	}

	return meta, nil
}
