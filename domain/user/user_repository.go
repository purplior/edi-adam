package user

import (
	"github.com/purplior/edi-adam/domain/shared/dto/pagination"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	UserRepository interface {
		Read(
			session inner.Session,
			queryOption QueryOption,
		) (
			model.User,
			error,
		)

		ReadPaginatedList(
			session inner.Session,
			query pagination.PaginationQuery[QueryOption],
		) (
			[]model.User,
			pagination.PaginationMeta,
			error,
		)

		Create(
			session inner.Session,
			m model.User,
		) (
			model.User,
			error,
		)

		Updates(
			session inner.Session,
			queryOption QueryOption,
			m model.User,
		) error

		Delete(
			session inner.Session,
			queryOption QueryOption,
		) error
	}
)
