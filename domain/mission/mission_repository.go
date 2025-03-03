package mission

import (
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	MissionRepository interface {
		Read(
			session inner.Session,
			queryOption QueryOption,
		) (
			model.Mission,
			error,
		)

		ReadPaginatedList(
			session inner.Session,
			query pagination.PaginationQuery[QueryOption],
		) (
			[]model.Mission,
			pagination.PaginationMeta,
			error,
		)

		Create(
			session inner.Session,
			m model.Mission,
		) (
			model.Mission,
			error,
		)

		Updates(
			session inner.Session,
			queryOption QueryOption,
			m model.Mission,
		) error

		Delete(
			session inner.Session,
			queryOption QueryOption,
		) error
	}
)
