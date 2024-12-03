package repoutil

import (
	"github.com/purplior/podoroot/domain/shared/pagination"
	"github.com/purplior/podoroot/infra/database/podosql"
)

func FindPaginatedList(
	db *podosql.DB,
	model interface{},
	entities interface{},
	page int,
	pageSize int,
	appendConditions func(db *podosql.DB) *podosql.DB,
	appendCountConditions func(db *podosql.DB) *podosql.DB,
) (
	pagination.PaginationMeta,
	error,
) {
	var totalCount int64

	if err := appendCountConditions(db).
		Model(model).
		Count(&totalCount).Error; err != nil {
		return pagination.PaginationMeta{}, err
	}

	offset := (page - 1) * pageSize
	if err := appendConditions(db).
		Model(model).
		Offset(offset).
		Limit(pageSize).
		Find(entities).Error; err != nil {
		return pagination.PaginationMeta{}, err
	}

	meta := pagination.PaginationMeta{
		Page:      page,
		Size:      pageSize,
		TotalPage: int((totalCount + int64(pageSize) - 1) / int64(pageSize)),
	}

	return meta, nil
}
