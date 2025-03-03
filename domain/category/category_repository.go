package category

import (
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	CategoryRepository interface {
		Read(
			session inner.Session,
			queryOption QueryOption,
		) (
			model.Category,
			error,
		)

		ReadList(
			session inner.Session,
			queryOption QueryOption,
		) (
			[]model.Category,
			error,
		)

		ReadPaginatedList(
			session inner.Session,
			query pagination.PaginationQuery[QueryOption],
		) (
			[]model.Category,
			pagination.PaginationMeta,
			error,
		)

		Create(
			session inner.Session,
			m model.Category,
		) (
			model.Category,
			error,
		)

		Updates(
			session inner.Session,
			queryOption QueryOption,
			m model.Category,
		) error

		Delete(
			session inner.Session,
			queryOption QueryOption,
		) error
	}
)
