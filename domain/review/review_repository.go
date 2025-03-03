package review

import (
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	ReviewRepository interface {
		Read(
			session inner.Session,
			queryOption QueryOption,
		) (
			model.Review,
			error,
		)

		ReadPaginatedList(
			session inner.Session,
			query pagination.PaginationQuery[QueryOption],
		) (
			[]model.Review,
			pagination.PaginationMeta,
			error,
		)

		Create(
			session inner.Session,
			m model.Review,
		) (
			model.Review,
			error,
		)

		Updates(
			session inner.Session,
			queryOption QueryOption,
			m model.Review,
		) error

		Delete(
			session inner.Session,
			queryOption QueryOption,
		) error
	}
)
