package assister

import (
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	AssisterRepository interface {
		Read(
			session inner.Session,
			queryOption QueryOption,
		) (
			model.Assister,
			error,
		)

		ReadPaginatedList(
			session inner.Session,
			query pagination.PaginationQuery[QueryOption],
		) (
			[]model.Assister,
			pagination.PaginationMeta,
			error,
		)

		Create(
			session inner.Session,
			m model.Assister,
		) (
			model.Assister,
			error,
		)

		Updates(
			session inner.Session,
			queryOption QueryOption,
			m model.Assister,
		) error

		Delete(
			session inner.Session,
			queryOption QueryOption,
		) error
	}
)
