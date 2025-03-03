package verification

import (
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	VerificationRepository interface {
		Read(
			session inner.Session,
			queryOption QueryOption,
		) (
			model.Verification,
			error,
		)

		ReadCount(
			session inner.Session,
			queryOption QueryOption,
		) (
			int,
			error,
		)

		ReadPaginatedList(
			session inner.Session,
			query pagination.PaginationQuery[QueryOption],
		) (
			[]model.Verification,
			pagination.PaginationMeta,
			error,
		)

		Create(
			session inner.Session,
			m model.Verification,
		) (
			model.Verification,
			error,
		)

		Updates(
			session inner.Session,
			queryOption QueryOption,
			m model.Verification,
		) error

		Delete(
			session inner.Session,
			queryOption QueryOption,
		) error
	}
)
